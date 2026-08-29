package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ListSubjects 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) ListSubjects(ctx context.Context, req *v1.ListSubjectsReq) (res *v1.ListSubjectsRes, err error) {
	return svc.StudyPlanet().ListSubjects(ctx, req)
}
