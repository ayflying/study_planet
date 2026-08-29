package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// MathAnswer 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) MathAnswer(ctx context.Context, req *v1.MathAnswerReq) (res *v1.MathAnswerRes, err error) {
	return svc.StudyPlanet().MathAnswer(ctx, req)
}
