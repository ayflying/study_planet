package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// UpdateStudent 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) UpdateStudent(ctx context.Context, req *v1.UpdateStudentReq) (res *v1.UpdateStudentRes, err error) {
	return svc.StudyPlanet().UpdateStudent(ctx, req)
}
