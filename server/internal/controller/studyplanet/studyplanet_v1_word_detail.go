package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// WordDetail 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) WordDetail(ctx context.Context, req *v1.WordDetailReq) (res *v1.WordDetailRes, err error) {
	return svc.StudyPlanet().WordDetail(ctx, req)
}
