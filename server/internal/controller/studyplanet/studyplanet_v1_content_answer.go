package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ContentAnswer 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) ContentAnswer(ctx context.Context, req *v1.ContentAnswerReq) (res *v1.ContentAnswerRes, err error) {
	return svc.StudyPlanet().ContentAnswer(ctx, req)
}
