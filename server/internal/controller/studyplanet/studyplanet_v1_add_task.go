package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// AddTask 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) AddTask(ctx context.Context, req *v1.AddTaskReq) (res *v1.AddTaskRes, err error) {
	return svc.StudyPlanet().AddTask(ctx, req)
}
