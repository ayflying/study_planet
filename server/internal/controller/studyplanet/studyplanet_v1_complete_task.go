package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// CompleteTask 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) CompleteTask(ctx context.Context, req *v1.CompleteTaskReq) (res *v1.CompleteTaskRes, err error) {
	return svc.StudyPlanet().CompleteTask(ctx, req)
}
