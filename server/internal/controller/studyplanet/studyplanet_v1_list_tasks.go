package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ListTasks 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) ListTasks(ctx context.Context, req *v1.ListTasksReq) (res *v1.ListTasksRes, err error) {
	return svc.StudyPlanet().ListTasks(ctx, req)
}
