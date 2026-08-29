// Package contentgen 填空题生成：数学计算填空（□ 形式），覆盖 1-9 年级。
// 与选择题同源程序化生成，保证答案正确；不含应用题。
package contentgen

import (
	"fmt"
	"math/rand"

	"studyplanet/internal/contentlib"
)

// genFill 生成数学填空题（qtype=fill），答案为纯数值/分数文本。
func genFill() []contentlib.Question {
	rng := rand.New(rand.NewSource(20260829)) // 固定种子保证可重复
	var qs []contentlib.Question
	add := func(grade int, topic, question, answer, explanation string) {
		qs = append(qs, contentlib.Question{
			Subject: "math", Grade: grade, Topic: topic, QType: "fill",
			Question: question, Options: nil, Answer: answer, Explanation: explanation, Source: "builtin-fill",
		})
	}
	// 1-2 年级：加减填空
	for i := 0; i < 20; i++ {
		a, b := rng.Intn(20)+1, rng.Intn(20)+1
		add(pickGrade(rng, 1, 2), "加减填空",
			fmt.Sprintf("填空：%d + □ = %d", a, a+b), itoa(b),
			fmt.Sprintf("%d+□=%d，□=%d−%d=%d", a, a+b, a+b, a, b))
	}
	for i := 0; i < 20; i++ {
		a, b := rng.Intn(15)+6, rng.Intn(15)+1
		if a <= b {
			a = b + rng.Intn(10) + 2
		}
		add(pickGrade(rng, 1, 2), "加减填空",
			fmt.Sprintf("填空：□ − %d = %d", b, a-b), itoa(a),
			fmt.Sprintf("□−%d=%d，□=%d+%d=%d", b, a-b, a-b, b, a))
	}
	// 3 年级：乘法填空
	for i := 0; i < 15; i++ {
		a, b := rng.Intn(9)+2, rng.Intn(9)+2
		add(3, "乘法填空",
			fmt.Sprintf("填空：%d × □ = %d", a, a*b), itoa(b),
			fmt.Sprintf("%d×□=%d，□=%d÷%d=%d", a, a*b, a*b, a, b))
	}
	// 4 年级：除法填空
	for i := 0; i < 12; i++ {
		a, b := (rng.Intn(80)+10)*10, rng.Intn(9)+2
		add(4, "除法填空",
			fmt.Sprintf("填空：□ ÷ %d = %d", b, a), itoa(a*b),
			fmt.Sprintf("□÷%d=%d，□=%d×%d=%d", b, a, a, b, a*b))
	}
	// 5 年级：小数填空
	for i := 0; i < 12; i++ {
		a := float64(rng.Intn(300)+100) / 100
		b := float64(rng.Intn(200)+50) / 100
		add(5, "小数填空",
			fmt.Sprintf("填空：%.2f + □ = %.2f", a, a+b), fmt.Sprintf("%.2f", b),
			fmt.Sprintf("□=%.2f−%.2f=%.2f", a+b, a, b))
	}
	// 6 年级：分数填空
	fracs := [][2]int{{1, 2}, {1, 3}, {2, 3}, {1, 4}, {3, 4}, {1, 5}, {2, 5}, {1, 6}, {1, 8}}
	for i := 0; i < 12; i++ {
		f := fracs[rng.Intn(len(fracs))]
		n := (rng.Intn(5) + 1) * f[1]
		add(6, "分数填空",
			fmt.Sprintf("填空：%d 的 □ 是 %d/%d", n, n/f[1]*f[0], f[1]),
			fmt.Sprintf("%d/%d", f[0], f[1]),
			fmt.Sprintf("%d×□=%d/%d，□=%d/%d", n, n/f[1]*f[0], f[1], f[0], f[1]))
	}
	// 7-8 年级：有理数/方程填空
	for i := 0; i < 12; i++ {
		a, x := rng.Intn(15)+2, rng.Intn(12)+2
		b := rng.Intn(20) - 10
		add(pickGrade(rng, 7, 8), "方程填空",
			fmt.Sprintf("填空：%dx %s = %d，x = □", a, signedTerm(b), a*x+b), itoa(x),
			fmt.Sprintf("%dx=%d，x=%d", a, a*x, x))
	}
	return qs
}
