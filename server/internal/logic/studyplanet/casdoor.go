package studyplanet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	"studyplanet/internal/model"
)

// ---------- Casdoor SSO（OIDC 授权码流程） ----------

type casdoorTokenResp struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
}

// AuthMode 告知前端当前登录模式：casdoor（未配置时为 pin）。
func (s *sStudyPlanet) AuthMode(r *ghttp.Request) {
	if s.Cfg.Casdoor.Enabled() {
		s.ok(r, map[string]interface{}{
			"mode":        "casdoor",
			"login_url":   "/api/parent/casdoor/login",
			"endpoint":    s.Cfg.Casdoor.Endpoint,
			"application": s.Cfg.Casdoor.Application,
		})
		return
	}
	s.ok(r, map[string]interface{}{"mode": "pin"})
}

// RedirectURIOf 按本次请求的访问地址实时推导回调地址（不写死、不落配置）。
// 依次识别 X-Forwarded-Proto / X-Forwarded-Host（反代场景），否则用请求自身 scheme+host。
func RedirectURIOf(r *ghttp.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	} else if p := r.Header.Get("X-Forwarded-Proto"); p != "" {
		scheme = p
	}
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.GetHost()
	}
	return fmt.Sprintf("%s://%s/api/parent/casdoor/callback", scheme, host)
}

// CasdoorLogin 302 跳转到 Casdoor 授权页。回调地址按用户实际访问地址实时生成。
func (s *sStudyPlanet) CasdoorLogin(r *ghttp.Request) {
	c := &s.Cfg.Casdoor
	if !c.Enabled() {
		s.fail(r, http.StatusBadRequest, "Casdoor 未配置")
		return
	}
	redirect := RedirectURIOf(r)
	r.Response.Header().Set("Location", c.AuthURL(redirect))
	r.Response.WriteStatus(http.StatusFound)
}

// casdoorUser Casdoor /api/userinfo 返回的用户信息（按需取子集）。
type casdoorUser struct {
	Sub               string `json:"sub"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	Avatar            string `json:"avatar"`
}

// CasdoorCallback 处理授权码：换 token → 拉用户信息 → upsert parents → 签发本站 JWT。
func (s *sStudyPlanet) CasdoorCallback(r *ghttp.Request) {
	c := &s.Cfg.Casdoor
	if !c.Enabled() {
		s.fail(r, http.StatusBadRequest, "Casdoor 未配置")
		return
	}
	code := r.GetQuery("code").String()
	if code == "" {
		s.fail(r, http.StatusBadRequest, "缺少授权码")
		return
	}

	// 1. 用授权码换 access_token（redirect_uri 必须与授权时一致，同样按请求推导）
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", c.ClientID)
	form.Set("client_secret", c.ClientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", RedirectURIOf(r))
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, c.TokenURL(), strings.NewReader(form.Encode()))
	if err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		s.fail(r, 502, "请求 Casdoor 失败: "+err.Error())
		return
	}
	defer resp.Body.Close()
	var tok casdoorTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil || tok.AccessToken == "" {
		s.fail(r, 502, "换取 token 失败")
		return
	}

	// 2. 拉 OAuth 用户信息
	req2, err := http.NewRequestWithContext(r.Context(), http.MethodGet, c.UserURL(), nil)
	if err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	req2.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	resp2, err := client.Do(req2)
	if err != nil {
		s.fail(r, 502, "拉取用户信息失败: "+err.Error())
		return
	}
	defer resp2.Body.Close()
	var cu casdoorUser
	if err := json.NewDecoder(resp2.Body).Decode(&cu); err != nil || cu.Sub == "" {
		s.fail(r, 502, "解析用户信息失败")
		return
	}

	// 3. upsert 家长账号（按 casdoor_sub 幂等，存在则覆盖展示信息与登录时间）
	now := time.Now().Format("2006-01-02 15:04:05")
	if _, err := s.DB.Exec(
		`INSERT INTO parents(casdoor_sub, display_name, avatar, last_login_at) VALUES(?,?,?,?)
		 ON DUPLICATE KEY UPDATE display_name=VALUES(display_name), avatar=VALUES(avatar), last_login_at=VALUES(last_login_at)`,
		cu.Sub, displayNameOf(&cu), cu.Avatar, now,
	); err != nil {
		s.fail(r, 500, err.Error())
		return
	}

	// 4. 签发本站 JWT，前端凭它调用家长接口
	jw, err := issueToken(s.Cfg.Parent.JWTSecret, displayNameOf(&cu))
	if err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	// 简单回跳页：把 token 写入 localStorage 后关闭/跳首页由前端接管
	html := `<script>
try{localStorage.setItem('sp_parent_jwt',` + jsQuote(jw) + `);localStorage.setItem('sp_parent_name',` + jsQuote(displayNameOf(&cu)) + `);}catch(e){}
location.href='/';
</script>`
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.Response.Write([]byte(html))
}

func displayNameOf(u *casdoorUser) string {
	if u.Name != "" {
		return u.Name
	}
	if u.PreferredUsername != "" {
		return u.PreferredUsername
	}
	return u.Sub
}

func jsQuote(v string) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// parentNameFromSub 由 casdoor_sub 取家长显示名（可选）。
func (s *sStudyPlanet) parentNameFromSub(sub string) string {
	var p model.Parent
	if err := s.DB.Get(&p, "SELECT id,casdoor_sub,display_name,avatar,created_at,COALESCE(last_login_at,'') AS last_login_at FROM parents WHERE casdoor_sub=?", sub); err != nil {
		return ""
	}
	return p.DisplayName
}
