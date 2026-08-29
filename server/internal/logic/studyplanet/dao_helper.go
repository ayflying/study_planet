package studyplanet

import "studyplanet/internal/dao"

// dao 快捷入口（GF 规范：业务层统一通过 dao 访问数据库）。
// dao 由 gf gen dao 生成（internal/dao），底层使用 g.DB() ORM 数据源。
var (
	daoChildren  = dao.Children
	daoTasks     = dao.Tasks
	daoRewards   = dao.Rewards
	daoPointsLog = dao.PointsLog
	daoParents   = dao.Parents
	daoWords     = dao.Words
	daoReadings  = dao.Readings
	daoReadingQ  = dao.ReadingQuestions
	daoMath      = dao.MathProblems
	daoWordProg  = dao.WordProgress
	daoRedempt   = dao.Redemptions
	daoSessions  = dao.PracticeSessions
	daoAnswers   = dao.SessionAnswers
	daoWrongQ    = dao.WrongQuestions
	daoSubjects  = dao.Subjects
	daoQuestions = dao.Questions
	daoWeekly    = dao.LeaderboardWeekly
)
