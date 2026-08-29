package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ImportContent 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) ImportContent(ctx context.Context, req *v1.ImportContentReq) (res *v1.ImportContentRes, err error) {
	return svc.StudyPlanet().ImportContent(ctx, req)
}
