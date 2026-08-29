package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ContentItem 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) ContentItem(ctx context.Context, req *v1.ContentItemReq) (res *v1.ContentItemRes, err error) {
	return svc.StudyPlanet().ContentItem(ctx, req)
}
