package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// StudentDeleteTask 控制器：学生删除自建任务。
func (c *ControllerV1) StudentDeleteTask(ctx context.Context, req *v1.StudentDeleteTaskReq) (res *v1.StudentDeleteTaskRes, err error) {
	return svc.StudyPlanet().StudentDeleteTask(ctx, req)
}