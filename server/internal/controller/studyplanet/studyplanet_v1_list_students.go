package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ListStudents 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) ListStudents(ctx context.Context, req *v1.ListStudentsReq) (res *v1.ListStudentsRes, err error) {
	return svc.StudyPlanet().ListStudents(ctx, req)
}
