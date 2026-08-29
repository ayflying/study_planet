package middleware

import (
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

// HandlerResponse 统一响应包装（GF 标准）：所有接口返回 {code, message, data}。
// 约定：
//   - code = 0 表示成功，data 为业务数据；code != 0 表示业务/系统错误，message 为错误描述。
//   - 已自行写入响应内容的接口（如 Casdoor 回调 HTML、302 跳转）不再二次包装。
func HandlerResponse(r *ghttp.Request) {
	r.Middleware.Next()

	// 控制器已写出自定义内容（如 Casdoor 回调 HTML）或已自行设置非 2xx 状态，保持原样。
	if r.Response.BufferLength() > 0 || r.Response.Status >= http.StatusMultipleChoices {
		return
	}

	var (
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code gcode.Code
		msg  string
	)
	if err != nil {
		code = gerror.Code(err)
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		msg = err.Error()
	} else {
		code = gcode.CodeOK
	}
	r.Response.WriteHeader(http.StatusOK)
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    code.Code(),
		Message: msg,
		Data:    res,
	})
}
