package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ListMath 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) ListMath(ctx context.Context, req *v1.ListMathReq) (res *v1.ListMathRes, err error) {
	return svc.StudyPlanet().ListMath(ctx, req)
}
