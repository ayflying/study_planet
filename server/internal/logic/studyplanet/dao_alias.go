package studyplanet

import "studyplanet/internal/model/do"

// doXxx 是 do 结构体的简写别名（GF 规范：写库统一用 do 结构体，字段全 interface{} 可空）。
type (
	doChildren  = do.Children
	doTasks     = do.Tasks
	doRewards   = do.Rewards
	doPointsLog = do.PointsLog
	doSettings  = do.Settings
	doParents   = do.Parents
	doWordProg  = do.WordProgress
	doRedempt   = do.Redemptions
	doSessions  = do.PracticeSessions
	doAnswers   = do.SessionAnswers
	doWrongQ    = do.WrongQuestions
)
