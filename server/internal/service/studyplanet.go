// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IStudyPlanet interface {
		// AuthMode 告知前端当前登录模式：casdoor（未配置时为 pin）。
		AuthMode(r *ghttp.Request)
		// CasdoorLogin 302 跳转到 Casdoor 授权页。回调地址按用户实际访问地址实时生成。
		CasdoorLogin(r *ghttp.Request)
		// CasdoorCallback 处理授权码：换 token → 拉用户信息 → upsert parents → 签发本站 JWT。
		CasdoorCallback(r *ghttp.Request)
		// ListSubjects 学科目录：GET /api/subjects?grade=5
		// 返回学科列表 + 每科题量，前端动态渲染学习地图。
		// 传 grade 时只返回该学段开设的学科（如 5 年级不出物理/化学）。
		ListSubjects(r *ghttp.Request)
		// PickQuestions 从内容库随机抽题：GET /api/content/pick?subject=math&grade=5&limit=5
		// 返回的题目不带 answer（不泄露答案给前端），前端作答后走 /api/content/answer 判分。
		PickQuestions(r *ghttp.Request)
		// ContentAnswer 统一判分：POST /api/content/answer
		// body: {id, answer, session_id?}；session_id 传入则复用连击+XP+错题本链路。
		ContentAnswer(r *ghttp.Request)
		// ContentItem 按 id 取单题（不含答案）：GET /api/content/item?id=
		// 错题本巩固复习时回取题目内容用。
		ContentItem(r *ghttp.Request)
		// ImportContent 通用题目导入（家长身份）：POST /api/parent/content/import
		// body: {"questions": [{subject,grade,topic,qtype,passage?,question,options[],answer,explanation?,difficulty?,source?}]}
		// 按 content_hash 去重，重复导入自动跳过——以后采集新资料直接调本接口，无需改源码。
		ImportContent(r *ghttp.Request)
		// SubjectStats 内容库统计（家长端）：GET /api/parent/content/stats
		SubjectStats(r *ghttp.Request)
		// ---------- 健康检查 ----------
		Health(r *ghttp.Request)
		// ---------- 单词卡片 ----------
		ListWords(r *ghttp.Request)
		WordDetail(r *ghttp.Request)
		WordProgress(r *ghttp.Request)
		// ---------- 语文阅读 ----------
		ReadingDetail(r *ghttp.Request)
		ReadingAnswer(r *ghttp.Request)
		// ---------- 数学题目 ----------
		ListMath(r *ghttp.Request)
		MathAnswer(r *ghttp.Request)
		// ---------- 每日任务 ----------
		ListTasks(r *ghttp.Request)
		CompleteTask(r *ghttp.Request)
		// ---------- 积分 ----------
		PointsSummary(r *ghttp.Request)
		PointsLog(r *ghttp.Request)
		// ---------- 奖励 / 兑换 ----------
		ListRewards(r *ghttp.Request)
		Redeem(r *ghttp.Request)
		// ---------- 家长端 ----------
		ParentLogin(r *ghttp.Request)
		// ---------- 以下为家长鉴权后接口 ----------
		AddTask(r *ghttp.Request)
		DeleteTask(r *ghttp.Request)
		AddReward(r *ghttp.Request)
		ConfirmRedemption(r *ghttp.Request)
		SetPin(r *ghttp.Request)
		// ---------- 学生账号管理（家长鉴权后） ----------
		// ListStudents 全部学生（家长切换用）。
		ListStudents(r *ghttp.Request)
		// CreateStudent 新建学生账号。
		CreateStudent(r *ghttp.Request)
		// UpdateStudent 修改学生信息（姓名/头像/年级/用户名）。
		UpdateStudent(r *ghttp.Request)
		// DeleteStudent 删除学生；至少保留一个，且清空其学习数据。
		DeleteStudent(r *ghttp.Request)
		// CreateSession 开启一关：POST /api/sessions {subject, level, total, student_id}
		CreateSession(r *ghttp.Request)
		// FinishSession 结算一关：POST /api/sessions/:id/finish
		// 星级规则：正确率 >=90% 三星、>=70% 两星、>=50% 一星、否则零星；
		// 三星额外奖励 10 分，两星 5 分。同一关只结算一次。
		FinishSession(r *ghttp.Request)
		// ListSessions 学生最近的练习记录（家长端统计用）：GET /api/sessions?student_id=
		ListSessions(r *ghttp.Request)
		// WeeklyLeaderboard 周榜：GET /api/leaderboard/weekly?limit=20
		// 返回当前 ISO 周经验值最高的学生名单。
		WeeklyLeaderboard(r *ghttp.Request)
		// ListWrongQuestions 学生错题本：GET /api/wrong-questions?subject=
		ListWrongQuestions(r *ghttp.Request)
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
