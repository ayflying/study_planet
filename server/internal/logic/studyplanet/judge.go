// Package studyplanet 判分入口：转发到独立 judge 包（内容判分与对战判分共用）。
package studyplanet

import "studyplanet/internal/judge"

// judgeAnswer 统一判分：支持选择题文本比对 + 填空题数值/分数比对。
func judgeAnswer(right, given string) bool {
	return judge.Judge(right, given)
}
