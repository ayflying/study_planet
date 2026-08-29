package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// PickQuestions 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) PickQuestions(ctx context.Context, req *v1.PickQuestionsReq) (res *v1.PickQuestionsRes, err error) {
	return svc.StudyPlanet().PickQuestions(ctx, req)
}
