// Package judge 统一判分：文本/数值/分数答案比较（内容判分与对战判分共用）。
// 填空答案兼容：全角/半角、空格、前导零、"1/2" 分数与 "0.5" 小数等价（按值比较）。
package judge

import (
	"strconv"
	"strings"
)

// Judge 判定给定答案与标准答案是否一致。
func Judge(right, given string) bool {
	right, given = strings.TrimSpace(right), strings.TrimSpace(given)
	if right == "" || given == "" {
		return false
	}
	if strings.EqualFold(right, given) {
		return true
	}
	// 数值比较（含小数）
	if isNum(right) && isNum(given) {
		return numEq(toF(right), toF(given))
	}
	// 分数比较 "a/b"
	rn, rd, ok1 := splitFrac(right)
	gn, gd, ok2 := splitFrac(given)
	if ok1 && ok2 && rd != 0 && gd != 0 {
		return rn*gd == gn*rd
	}
	// 归一化文本兜底（全角→半角、去空格）
	return Norm(right) == Norm(given)
}

// toF 字符串转 float（失败返回 0，配合 isNum 使用）。
func toF(s string) float64 { f, _ := strconv.ParseFloat(Norm(s), 64); return f }

// isNum 判断是否为纯数值（含负号/小数点）。
func isNum(s string) bool {
	_, err := strconv.ParseFloat(Norm(s), 64)
	return err == nil
}

// numEq 数值相等（容差 1e-6）。
func numEq(a, b float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d < 1e-6
}

// splitFrac 拆分 "a/b" → (a, b, true)。
func splitFrac(s string) (int, int, bool) {
	parts := strings.SplitN(Norm(s), "/", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}
	a, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	b, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return a, b, true
}

// Norm 归一化：全角转半角、去空格；纯数值时消除浮点表示差异（0.50 → 0.5）。
func Norm(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 0xFF01 && r <= 0xFF5E: // 全角 ASCII
			r -= 0xFEE0
		case r == 0x3000: // 全角空格
			r = ' '
		}
		if r != ' ' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if n, err := strconv.ParseFloat(out, 64); err == nil {
		return strconv.FormatFloat(n, 'f', -1, 64)
	}
	return out
}
