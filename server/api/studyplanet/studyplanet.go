// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package studyplanet

import (
	"context"

	"studyplanet/api/studyplanet/v1"
)

type IStudyplanetV1 interface {
	AuthMode(ctx context.Context, req *v1.AuthModeReq) (res *v1.AuthModeRes, err error)
	ParentLogin(ctx context.Context, req *v1.ParentLoginReq) (res *v1.ParentLoginRes, err error)
	CasdoorLogin(ctx context.Context, req *v1.CasdoorLoginReq) (res *v1.CasdoorLoginRes, err error)
	CasdoorCallback(ctx context.Context, req *v1.CasdoorCallbackReq) (res *v1.CasdoorCallbackRes, err error)
	SetPin(ctx context.Context, req *v1.SetPinReq) (res *v1.SetPinRes, err error)
	ListSubjects(ctx context.Context, req *v1.ListSubjectsReq) (res *v1.ListSubjectsRes, err error)
	PickQuestions(ctx context.Context, req *v1.PickQuestionsReq) (res *v1.PickQuestionsRes, err error)
	ContentItem(ctx context.Context, req *v1.ContentItemReq) (res *v1.ContentItemRes, err error)
	ContentAnswer(ctx context.Context, req *v1.ContentAnswerReq) (res *v1.ContentAnswerRes, err error)
	Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error)
	ImportContent(ctx context.Context, req *v1.ImportContentReq) (res *v1.ImportContentRes, err error)
	SubjectStats(ctx context.Context, req *v1.SubjectStatsReq) (res *v1.SubjectStatsRes, err error)
	ListWords(ctx context.Context, req *v1.ListWordsReq) (res *v1.ListWordsRes, err error)
	WordDetail(ctx context.Context, req *v1.WordDetailReq) (res *v1.WordDetailRes, err error)
	WordProgress(ctx context.Context, req *v1.WordProgressReq) (res *v1.WordProgressRes, err error)
	ReadingDetail(ctx context.Context, req *v1.ReadingDetailReq) (res *v1.ReadingDetailRes, err error)
	ReadingAnswer(ctx context.Context, req *v1.ReadingAnswerReq) (res *v1.ReadingAnswerRes, err error)
	ListMath(ctx context.Context, req *v1.ListMathReq) (res *v1.ListMathRes, err error)
	MathAnswer(ctx context.Context, req *v1.MathAnswerReq) (res *v1.MathAnswerRes, err error)
	PointsSummary(ctx context.Context, req *v1.PointsSummaryReq) (res *v1.PointsSummaryRes, err error)
	PointsLog(ctx context.Context, req *v1.PointsLogReq) (res *v1.PointsLogRes, err error)
	CreateSession(ctx context.Context, req *v1.CreateSessionReq) (res *v1.CreateSessionRes, err error)
	ListSessions(ctx context.Context, req *v1.ListSessionsReq) (res *v1.ListSessionsRes, err error)
	FinishSession(ctx context.Context, req *v1.FinishSessionReq) (res *v1.FinishSessionRes, err error)
	WeeklyLeaderboard(ctx context.Context, req *v1.WeeklyLeaderboardReq) (res *v1.WeeklyLeaderboardRes, err error)
	ListWrongQuestions(ctx context.Context, req *v1.ListWrongQuestionsReq) (res *v1.ListWrongQuestionsRes, err error)
	ListRewards(ctx context.Context, req *v1.ListRewardsReq) (res *v1.ListRewardsRes, err error)
	AddReward(ctx context.Context, req *v1.AddRewardReq) (res *v1.AddRewardRes, err error)
	Redeem(ctx context.Context, req *v1.RedeemReq) (res *v1.RedeemRes, err error)
	ConfirmRedemption(ctx context.Context, req *v1.ConfirmRedemptionReq) (res *v1.ConfirmRedemptionRes, err error)
	ListStudents(ctx context.Context, req *v1.ListStudentsReq) (res *v1.ListStudentsRes, err error)
	CreateStudent(ctx context.Context, req *v1.CreateStudentReq) (res *v1.CreateStudentRes, err error)
	UpdateStudent(ctx context.Context, req *v1.UpdateStudentReq) (res *v1.UpdateStudentRes, err error)
	DeleteStudent(ctx context.Context, req *v1.DeleteStudentReq) (res *v1.DeleteStudentRes, err error)
	ListTasks(ctx context.Context, req *v1.ListTasksReq) (res *v1.ListTasksRes, err error)
	AddTask(ctx context.Context, req *v1.AddTaskReq) (res *v1.AddTaskRes, err error)
	DeleteTask(ctx context.Context, req *v1.DeleteTaskReq) (res *v1.DeleteTaskRes, err error)
	CompleteTask(ctx context.Context, req *v1.CompleteTaskReq) (res *v1.CompleteTaskRes, err error)
}
