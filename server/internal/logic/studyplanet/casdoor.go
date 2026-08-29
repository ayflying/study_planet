package studyplanet

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "studyplanet/api/studyplanet/v1"
)

// ---------- Casdoor SSO（OIDC 授权码流程） ----------
// 说明：302 跳转与回调页 HTML 输出属于 HTTP 边缘语义，由控制器层处理；
// 本文件只做「换 token / 拉用户信息 / 落库家长 / 签 JWT」的纯业务逻辑。

// casdoorTokenResp Casdoor token 端点的返回结构。
type casdoorTokenResp struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
}

// casdoorUser Casdoor /api/userinfo 返回的用户信息（按需取子集）。
type casdoorUser struct {
	Sub               string `json:"sub"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	Avatar            string `json:"avatar"`
}

// AuthMode 告知前端当前登录模式：casdoor（未配置时为 pin）。
func (s *sStudyPlanet) AuthMode(ctx context.Context, req *v1.AuthModeReq) (res *v1.AuthModeRes, err error) {
	if s.Cfg != nil && s.Cfg.Casdoor.Enabled() {
		return &v1.AuthModeRes{
			Mode:        "casdoor",
			LoginURL:    "/api/parent/casdoor/login",
			Endpoint:    s.Cfg.Casdoor.Endpoint,
			Application: s.Cfg.Casdoor.Application,
		}, nil
	}
	return &v1.AuthModeRes{Mode: "pin"}, nil
}

// CasdoorLogin 计算 Casdoor 授权页地址，控制器据此 302 跳转。
// 回调地址按用户实际访问地址实时生成（req.RedirectURI 由控制器按请求推导后注入）。
func (s *sStudyPlanet) CasdoorLogin(ctx context.Context, req *v1.CasdoorLoginReq) (res *v1.CasdoorLoginRes, err error) {
	if s.Cfg == nil || !s.Cfg.Casdoor.Enabled() {
		return nil, errAuth("Casdoor 未配置")
	}
	return &v1.CasdoorLoginRes{Location: s.Cfg.Casdoor.AuthURL(req.RedirectURI)}, nil
}

// CasdoorCallback 处理授权码：换 token → 拉用户信息 → upsert parents → 签发本站 JWT。
// 返回的 HTML 由控制器以 text/html 写出。
func (s *sStudyPlanet) CasdoorCallback(ctx context.Context, req *v1.CasdoorCallbackReq) (res *v1.CasdoorCallbackRes, err error) {
	if s.Cfg == nil || !s.Cfg.Casdoor.Enabled() {
		return nil, errAuth("Casdoor 未配置")
	}
	if req.Code == "" {
		return nil, errAuth("缺少授权码")
	}
	c := &s.Cfg.Casdoor

	// 1. 用授权码换 access_token（redirect_uri 必须与授权时一致）
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", c.ClientID)
	form.Set("client_secret", c.ClientSecret)
	form.Set("code", req.Code)
	form.Set("redirect_uri", req.RedirectURI)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.TokenURL(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, gerror.Wrap(err, "构造 token 请求失败")
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, gerror.Wrap(err, "请求 Casdoor 失败")
	}
	defer resp.Body.Close()
	var tok casdoorTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil || tok.AccessToken == "" {
		return nil, gerror.New("换取 token 失败")
	}

	// 2. 拉 OAuth 用户信息
	httpReq2, err := http.NewRequestWithContext(ctx, http.MethodGet, c.UserURL(), nil)
	if err != nil {
		return nil, gerror.Wrap(err, "构造用户信息请求失败")
	}
	httpReq2.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	resp2, err := client.Do(httpReq2)
	if err != nil {
		return nil, gerror.Wrap(err, "拉取用户信息失败")
	}
	defer resp2.Body.Close()
	var cu casdoorUser
	if err := json.NewDecoder(resp2.Body).Decode(&cu); err != nil || cu.Sub == "" {
		return nil, gerror.New("解析用户信息失败")
	}

	// 3. upsert 家长账号（按 casdoor_sub 幂等，存在则覆盖展示信息与登录时间）
	name := displayNameOf(&cu)
	if _, err := daoParents.Ctx(ctx).Data(doParents{
		CasdoorSub:  cu.Sub,
		DisplayName: name,
		Avatar:      cu.Avatar,
		LastLoginAt: gtime.Now(),
	}).Save(); err != nil {
		return nil, gerror.Wrap(err, "保存家长账号失败")
	}
	prow, err := daoParents.Ctx(ctx).Where("casdoor_sub", cu.Sub).One()
	if err != nil || prow.IsEmpty() {
		return nil, gerror.New("查询家长账号失败")
	}
	parentID := prow["id"].Int()

	// 4. 老数据接管：历史上无家长归属（parent_id 为 NULL）的孩子与奖励，
	// 归属到第一个登录的家长，此后各家长数据完全隔离。
	if err := s.claimOrphans(ctx, parentID); err != nil {
		return nil, gerror.Wrap(err, "接管历史数据失败")
	}

	// 5. 签发本站 JWT，前端凭它调用家长接口
	secret := ""
	if s.Cfg != nil {
		secret = s.Cfg.Parent.JWTSecret
	}
	jw, err := issueToken(secret, name, parentID)
	if err != nil {
		return nil, gerror.Wrap(err, "签发令牌失败")
	}
	// 回跳页：把 token 写入 localStorage 后跳首页（由控制器以 text/html 输出）
	return &v1.CasdoorCallbackRes{HTML: `<script>
try{localStorage.setItem('sp_parent_jwt',` + jsQuote(jw) + `);localStorage.setItem('sp_parent_name',` + jsQuote(name) + `);}catch(e){}
location.href='/';
</script>`}, nil
}

// claimOrphans 历史无归属数据接管：把 parent_id 为 NULL 的孩子与奖励划给指定家长。
// 仅在家长登录时执行；先判空避免每次登录都跑无谓的 UPDATE。
func (s *sStudyPlanet) claimOrphans(ctx context.Context, parentID int) error {
	if parentID <= 0 {
		return nil
	}
	// 孩子无归属 → 全部接管
	orphanChildren, err := daoChildren.Ctx(ctx).Where("parent_id IS NULL").Count()
	if err != nil {
		return err
	}
	if orphanChildren > 0 {
		if _, err := daoChildren.Ctx(ctx).Where("parent_id IS NULL").Data(doChildren{ParentId: parentID}).Update(); err != nil {
			return err
		}
		g.Log().Infof(ctx, "家长 %d 接管了 %d 个无归属孩子", parentID, orphanChildren)
	}
	// 奖励无归属 → 全部接管
	orphanRewards, err := daoRewards.Ctx(ctx).Where("parent_id IS NULL").Count()
	if err != nil {
		return err
	}
	if orphanRewards > 0 {
		if _, err := daoRewards.Ctx(ctx).Where("parent_id IS NULL").Data(doRewards{ParentId: parentID}).Update(); err != nil {
			return err
		}
		g.Log().Infof(ctx, "家长 %d 接管了 %d 条无归属奖励", parentID, orphanRewards)
	}
	return nil
}

// displayNameOf 家长展示名：优先 name，其次 preferred_username，最后 sub。
func displayNameOf(u *casdoorUser) string {
	if u.Name != "" {
		return u.Name
	}
	if u.PreferredUsername != "" {
		return u.PreferredUsername
	}
	return u.Sub
}

// jsQuote 安全地把字符串嵌入 JS 字面量。
func jsQuote(v string) string {
	b, _ := json.Marshal(v)
	return string(b)
}
