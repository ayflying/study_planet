package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// StudentAddTask 控制器：学生自建任务。
func (c *ControllerV1) StudentAddTask(ctx context.Context, req *v1.StudentAddTaskReq) (res *v1.StudentAddTaskRes, err error) {
	return svc.StudyPlanet().StudentAddTask(ctx, req)
}