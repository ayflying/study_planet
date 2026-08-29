// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	v1 "studyplanet/api/studyplanet/v1"
)

type (
	IStudyPlanet interface {
		// AuthMode 告知前端当前登录模式：casdoor（未配置时为 pin）。
		AuthMode(ctx context.Context, req *v1.AuthModeReq) (res *v1.AuthModeRes, err error)
		// CasdoorLogin 计算 Casdoor 授权页地址，控制器据此 302 跳转。
		// 回调地址按用户实际访问地址实时生成（req.RedirectURI 由控制器按请求推导后注入）。
		CasdoorLogin(ctx context.Context, req *v1.CasdoorLoginReq) (res *v1.CasdoorLoginRes, err error)
		// CasdoorCallback 处理授权码：换 token → 拉用户信息 → upsert parents → 签发本站 JWT。
		// 返回的 HTML 由控制器以 text/html 写出。
		CasdoorCallback(ctx context.Context, req *v1.CasdoorCallbackReq) (res *v1.CasdoorCallbackRes, err error)
		// ListSubjects 学科目录：GET /api/subjects?grade=5
		// 返回学科列表 + 每科题量，前端动态渲染学习地图。
		// 传 grade 时只返回该学段开设的学科（如 5 年级不出物理/化学）。
		ListSubjects(ctx context.Context, req *v1.ListSubjectsReq) (res *v1.ListSubjectsRes, err error)
		// PickQuestions 从内容库随机抽题：GET /api/content/pick?subject=math&grade=5&limit=5
		// 返回的题目不带 answer（不泄露答案给前端），前端作答后走 /api/content/answer 判分。
		PickQuestions(ctx context.Context, req *v1.PickQuestionsReq) (res *v1.PickQuestionsRes, err error)
		// ContentAnswer 统一判分：POST /api/content/answer
		// body: {id, answer, session_id?}；session_id 传入则复用连击+XP+错题本链路。
		ContentAnswer(ctx context.Context, req *v1.ContentAnswerReq) (res *v1.ContentAnswerRes, err error)
		// ContentItem 按 id 取单题（不含答案）：GET /api/content/item?id=
		// 错题本巩固复习时回取题目内容用。
		ContentItem(ctx context.Context, req *v1.ContentItemReq) (res *v1.ContentItemRes, err error)
		// ImportContent 通用题目导入（家长身份）：POST /api/parent/content/import
		// body: {"questions": [{subject,grade,topic,qtype,passage?,question,options[],answer,explanation?,difficulty?,source?}]}
		// 按 content_hash 去重，重复导入自动跳过——以后采集新资料直接调本接口，无需改源码。
		ImportContent(ctx context.Context, req *v1.ImportContentReq) (res *v1.ImportContentRes, err error)
		// SubjectStats 内容库统计（家长端）：GET /api/parent/content/stats
		SubjectStats(ctx context.Context, req *v1.SubjectStatsReq) (res *v1.SubjectStatsRes, err error)
		// Health 健康检查：返回运行状态与版本号。
		Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error)
		// ListWords 单词列表（可选按 level 过滤）。
		ListWords(ctx context.Context, req *v1.ListWordsReq) (res *v1.ListWordsRes, err error)
		// WordDetail 单词详情 + 当前学生掌握状态。
		WordDetail(ctx context.Context, req *v1.WordDetailReq) (res *v1.WordDetailRes, err error)
		// WordProgress 标记单词掌握状态；session_id 传入则走连击+场次计分。
		WordProgress(ctx context.Context, req *v1.WordProgressReq) (res *v1.WordProgressRes, err error)
		// ReadingDetail 阅读详情 + 题目列表。
		ReadingDetail(ctx context.Context, req *v1.ReadingDetailReq) (res *v1.ReadingDetailRes, err error)
		// ReadingAnswer 阅读题目作答。
		ReadingAnswer(ctx context.Context, req *v1.ReadingAnswerReq) (res *v1.ReadingAnswerRes, err error)
		// ListMath 数学题目列表。
		ListMath(ctx context.Context, req *v1.ListMathReq) (res *v1.ListMathRes, err error)
		// MathAnswer 数学题作答。
		MathAnswer(ctx context.Context, req *v1.MathAnswerReq) (res *v1.MathAnswerRes, err error)
		// ListTasks 学生任务列表（按 student_id，公开接口）。
		ListTasks(ctx context.Context, req *v1.ListTasksReq) (res *v1.ListTasksRes, err error)
		// CompleteTask 完成任务（学生操作，公开接口）。
		CompleteTask(ctx context.Context, req *v1.CompleteTaskReq) (res *v1.CompleteTaskRes, err error)
		// PointsSummary 积分汇总（公开接口）。
		PointsSummary(ctx context.Context, req *v1.PointsSummaryReq) (res *v1.PointsSummaryRes, err error)
		// PointsLog 积分流水（最近 100 条，公开接口）。
		PointsLog(ctx context.Context, req *v1.PointsLogReq) (res *v1.PointsLogRes, err error)
		// ListRewards 奖励列表（公开接口）。
		ListRewards(ctx context.Context, req *v1.ListRewardsReq) (res *v1.ListRewardsRes, err error)
		// Redeem 学生兑换奖励（公开接口，需家长确认）。
		Redeem(ctx context.Context, req *v1.RedeemReq) (res *v1.RedeemRes, err error)
		// ParentLogin 家长 PIN 登录（Casdoor 启用时禁用）。
		ParentLogin(ctx context.Context, req *v1.ParentLoginReq) (res *v1.ParentLoginRes, err error)
		// AddTask 发布任务（家长鉴权）。
		AddTask(ctx context.Context, req *v1.AddTaskReq) (res *v1.AddTaskRes, err error)
		// DeleteTask 删除任务（家长鉴权）。
		DeleteTask(ctx context.Context, req *v1.DeleteTaskReq) (res *v1.DeleteTaskRes, err error)
		// AddReward 添加奖励（家长鉴权）。
		AddReward(ctx context.Context, req *v1.AddRewardReq) (res *v1.AddRewardRes, err error)
		// ConfirmRedemption 家长确认兑换（家长鉴权）。
		ConfirmRedemption(ctx context.Context, req *v1.ConfirmRedemptionReq) (res *v1.ConfirmRedemptionRes, err error)
		// SetPin 修改家长 PIN（家长鉴权）。
		SetPin(ctx context.Context, req *v1.SetPinReq) (res *v1.SetPinRes, err error)
		// ListStudents 全部学生（家长切换用）。
		ListStudents(ctx context.Context, req *v1.ListStudentsReq) (res *v1.ListStudentsRes, err error)
		// CreateStudent 新建学生账号。
		CreateStudent(ctx context.Context, req *v1.CreateStudentReq) (res *v1.CreateStudentRes, err error)
		// UpdateStudent 修改学生信息（姓名/头像/年级/用户名）。
		UpdateStudent(ctx context.Context, req *v1.UpdateStudentReq) (res *v1.UpdateStudentRes, err error)
		// DeleteStudent 删除学生；至少保留一个，且清空其学习数据。
		DeleteStudent(ctx context.Context, req *v1.DeleteStudentReq) (res *v1.DeleteStudentRes, err error)
		// CreateSession 开启一关（多邻国式关卡）。
		CreateSession(ctx context.Context, req *v1.CreateSessionReq) (res *v1.CreateSessionRes, err error)
		// FinishSession 结算一关。
		// 星级规则：正确率 >=90% 三星、>=70% 两星、>=50% 一星、否则零星；
		// 三星额外奖励 10 分，两星 5 分。同一关只结算一次。
		FinishSession(ctx context.Context, req *v1.FinishSessionReq) (res *v1.FinishSessionRes, err error)
		// ListSessions 学生最近的练习记录（家长端统计用）。
		ListSessions(ctx context.Context, req *v1.ListSessionsReq) (res *v1.ListSessionsRes, err error)
		// ListWrongQuestions 学生错题本（可选按 subject 过滤）。
		ListWrongQuestions(ctx context.Context, req *v1.ListWrongQuestionsReq) (res *v1.ListWrongQuestionsRes, err error)
		// WeeklyLeaderboard 周榜：返回当前 ISO 周经验值最高的学生名单。
		WeeklyLeaderboard(ctx context.Context, req *v1.WeeklyLeaderboardReq) (res *v1.WeeklyLeaderboardRes, err error)
	}
)

var (
	localStudyPlanet IStudyPlanet
)

func StudyPlanet() IStudyPlanet {
	if localStudyPlanet == nil {
		panic("implement not found for interface IStudyPlanet, forgot register?")
	}
	return localStudyPlanet
}

func RegisterStudyPlanet(i IStudyPlanet) {
	localStudyPlanet = i
}
