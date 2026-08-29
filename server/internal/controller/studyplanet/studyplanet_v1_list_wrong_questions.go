package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ListWrongQuestions 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) ListWrongQuestions(ctx context.Context, req *v1.ListWrongQuestionsReq) (res *v1.ListWrongQuestionsRes, err error) {
	return svc.StudyPlanet().ListWrongQuestions(ctx, req)
}
