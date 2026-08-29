package studyplanet

import (
	"bytes"
	"context"
	"os"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"

	"studyplanet/internal/config"
	"studyplanet/internal/leaderboard"
	"studyplanet/internal/middleware"
	"studyplanet/internal/service"
)

// AppVersion 由构建时注入（Dockerfile -ldflags -X 写入），未注入时读 VERSION 文件。
var AppVersion string

// CurrentVersion 返回构建注入的版本号（未注入时读 VERSION 文件，兜底 dev）。
func CurrentVersion() string {
	if AppVersion != "" {
		return AppVersion
	}
	if b, err := os.ReadFile("VERSION"); err == nil {
		return string(bytes.TrimSpace(b))
	}
	return "dev"
}

// sStudyPlanet 业务实现（命名符合 gf gen service 的 stPattern ^s([A-Z]\w+)$）。
// 按 GF 分层约束：logic 层不接触 ghttp.Request，数据访问统一走 dao（g.DB()），
// 配置与外部模块（排行榜）通过 SetDeps 注入。
type sStudyPlanet struct {
	Cfg   *config.Config
	Board *leaderboard.Board // 每周经验排行榜（独立模块，可为降级模式）
}

// SetDeps 注入运行依赖（由 cmd 启动时调用）。
// 注意：init() 注册时 local 必须已是可用实例，这里填充字段而非替换实例，
// 避免 service 层持有的接口指向旧指针。
func SetDeps(cfg *config.Config, board *leaderboard.Board) *sStudyPlanet {
	local.Cfg, local.Board = cfg, board
	return local
}

// local 当前业务实现实例：init 注册时创建（非 nil），SetDeps 填充依赖。
var local = &sStudyPlanet{}

// Study 供外部包获取当前业务实现实例。
func Study() *sStudyPlanet { return local }

// ExternalAddXP 供 battle 引擎等外部模块给学生加经验（联动周榜）。
// 返回闭包，避免 battle 直接依赖 logic 内部实现。
func ExternalAddXP(board *leaderboard.Board) func(childID int, delta int) {
	return func(childID int, delta int) {
		if delta == 0 || childID <= 0 {
			return
		}
		if _, err := daoChildren.Ctx(gctx.New()).Where("id", childID).Increment("xp", delta); err != nil {
			gLog("battle 外部加经验失败 child=%d: %v", childID, err)
		}
		if board != nil {
			board.AddXP(gctx.New(), childID, delta)
		}
	}
}

// ---------- 通用助手 ----------

// gLog 非关键路径日志（失败不影响主流程）。
func gLog(format string, args ...interface{}) {
	g.Log().Errorf(gctx.New(), format, args...)
}

// nowStr 当前时间（数据库 DATETIME 列格式）。
func nowStr() string { return time.Now().Format("2006-01-02 15:04:05") }

// todayStr 当前日期。
func todayStr() string { return time.Now().Format("2006-01-02") }

// issueToken 签发家长 JWT。
// parentID>0 时写入 claims（Casdoor 登录），数据隔离按它判定归属；
// 无家长账号的旧式登录（已废弃的 PIN）签 parentID=0，只能走只读兜底。
func issueToken(secret, name string, parentID int) (string, error) {
	// 会话没有固定过期时间：前端将 token 持久化到 localStorage，
	// 只有家长主动退出（或更换 JWT_SECRET）才会结束登录状态。
	claims := jwt.MapClaims{
		"sub":       "parent",
		"name":      name,
		"parent_id": parentID,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

// ---------- 业务错误（GF 中间件统一包装为 {code,message,data}） ----------

// errParam 参数/状态不满足（对应 HTTP 400 语义，业务 code=51）。
func errParam(msg string) error { return gerror.NewCode(gcode.CodeInvalidParameter, msg) }

// errNotFound 资源不存在（对应 HTTP 404 语义，业务 code=58）。
func errNotFound(msg string) error { return gerror.NewCode(gcode.CodeNotFound, msg) }

// errAuth 鉴权/密码错误（对应 HTTP 401 语义，业务 code=56）。
func errAuth(msg string) error { return gerror.NewCode(gcode.CodeNotAuthorized, msg) }

// errForbidden 无权操作他人数据（对应 HTTP 403 语义，业务 code=59）。
func errForbidden(msg string) error { return gerror.NewCode(gcode.CodeSecurityReason, msg) }

// ctxParentID 当前登录家长 id（ParentAuth 中间件注入；未携带时为 0）。
func ctxParentID(ctx context.Context) int { return middleware.ParentIDOf(ctx) }

// ---------- XP / 积分（多业务文件共用） ----------

// addXP 给学生累加经验值：children.xp 总量 + 周榜（Redis/降级数据库）。
func (s *sStudyPlanet) addXP(ctx context.Context, childID int, delta int) {
	if delta == 0 {
		return
	}
	if _, err := daoChildren.Ctx(ctx).Where("id", childID).Increment("xp", delta); err != nil {
		gLog("xp 累加失败 child=%d: %v", childID, err)
	}
	if s.Board != nil {
		s.Board.AddXP(ctx, childID, delta)
	}
}

// award 记录指定学生的积分变动（不阻塞主流程，失败仅记日志）。
func (s *sStudyPlanet) award(childID int, delta int, reason string) {
	if _, err := daoPointsLog.Ctx(gctx.New()).Data(doPointsLog{
		ChildId:   childID,
		Delta:     delta,
		Reason:    reason,
		CreatedAt: gtime.Now(),
	}).Insert(); err != nil {
		g.Log().Errorf(gctx.New(), "award 记录积分失败 child=%d delta=%d: %v", childID, delta, err)
	}
}

// init 业务实现注册：service 层接口 IStudyPlanet ← logic 实现绑定（gf gen service 规范）。
func init() {
	service.RegisterStudyPlanet(local)
}
