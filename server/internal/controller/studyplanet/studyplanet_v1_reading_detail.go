package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ReadingDetail 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) ReadingDetail(ctx context.Context, req *v1.ReadingDetailReq) (res *v1.ReadingDetailRes, err error) {
	return svc.StudyPlanet().ReadingDetail(ctx, req)
}
