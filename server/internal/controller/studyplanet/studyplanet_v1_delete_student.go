package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// DeleteStudent 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) DeleteStudent(ctx context.Context, req *v1.DeleteStudentReq) (res *v1.DeleteStudentRes, err error) {
	return svc.StudyPlanet().DeleteStudent(ctx, req)
}
