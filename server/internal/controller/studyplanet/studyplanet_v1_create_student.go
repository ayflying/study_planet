package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// CreateStudent 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) CreateStudent(ctx context.Context, req *v1.CreateStudentReq) (res *v1.CreateStudentRes, err error) {
	return svc.StudyPlanet().CreateStudent(ctx, req)
}
