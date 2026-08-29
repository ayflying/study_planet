package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// DeleteTask 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) DeleteTask(ctx context.Context, req *v1.DeleteTaskReq) (res *v1.DeleteTaskRes, err error) {
	return svc.StudyPlanet().DeleteTask(ctx, req)
}
