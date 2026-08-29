// Package contentgen 内置题库生成：数学程序化生成（保证答案正确），
// 英语/语文/理科手写基础题库。输出 []contentlib.Question 供导入。
package contentgen

import (
	"fmt"
	"math/rand"

	"studyplanet/internal/contentlib"
)

// Generate 生成全部内置题目（确定性：固定随机种子，重复运行产生相同题库）。
func Generate() []contentlib.Question {
	var qs []contentlib.Question
	qs = append(qs, genMath()...)
	qs = append(qs, genFill()...)
	qs = append(qs, genEnglish()...)
	qs = append(qs, genChinese()...)
	qs = append(qs, genScience()...)
	return qs
}

// ---------- 数学：程序化生成 ----------

// mathSpec 各年级数学题型。
func genMath() []contentlib.Question {
	rng := rand.New(rand.NewSource(20260828)) // 固定种子保证可重复
	var qs []contentlib.Question
	add := func(grade int, topic string, question, answer string, options []string, explanation string) {
		qs = append(qs, contentlib.Question{
			Subject: "math", Grade: grade, Topic: topic, QType: "choice",
			Question: question, Options: options, Answer: answer, Explanation: explanation, Source: "builtin",
		})
	}
	// 1-2 年级：20/100 以内加减法
	for i := 0; i < 40; i++ {
		a, b := rng.Intn(20)+1, rng.Intn(20)+1
		if rng.Intn(2) == 0 && a >= b {
			add(pickGrade(rng, 1, 2), "加减法",
				fmt.Sprintf("计算：%d − %d = ?", a, b),
				itoa(a-b), numOptions(rng, a-b, 12), fmt.Sprintf("%d−%d=%d", a, b, a-b))
		} else {
			add(pickGrade(rng, 1, 2), "加减法",
				fmt.Sprintf("计算：%d + %d = ?", a, b),
				itoa(a+b), numOptions(rng, a+b, 12), fmt.Sprintf("%d+%d=%d", a, b, a+b))
		}
	}
	// 3 年级：乘法口诀 / 两位数乘法
	for i := 0; i < 30; i++ {
		a, b := rng.Intn(9)+2, rng.Intn(9)+2
		add(3, "乘法",
			fmt.Sprintf("计算：%d × %d = ?", a, b),
			itoa(a*b), numOptions(rng, a*b, 18), fmt.Sprintf("%d×%d=%d", a, b, a*b))
	}
	// 4 年级：大数加减 / 简单除法
	for i := 0; i < 25; i++ {
		a, b := (rng.Intn(90)+10)*10, rng.Intn(9)+2
		add(4, "除法",
			fmt.Sprintf("计算：%d ÷ %d = ?", a*b, b),
			itoa(a), numOptions(rng, a, 100), fmt.Sprintf("%d÷%d=%d", a*b, b, a))
	}
	// 5 年级：小数运算 / 长方形面积 / 简易方程
	for i := 0; i < 30; i++ {
		a := float64(rng.Intn(400)+100) / 100
		b := float64(rng.Intn(300)+50) / 100
		sum := fmt.Sprintf("%.2f", a+b)
		add(5, "小数运算",
			fmt.Sprintf("计算：%.2f + %.2f = ?", a, b),
			sum, numOptionsS(rng, sum, a+b), fmt.Sprintf("%.2f+%.2f=%.2f", a, b, a+b))
	}
	for i := 0; i < 15; i++ {
		l, w := rng.Intn(15)+3, rng.Intn(12)+2
		add(5, "图形面积",
			fmt.Sprintf("一个长方形长 %dcm、宽 %dcm，面积是多少平方厘米？", l, w),
			itoa(l*w), numOptions(rng, l*w, 40), fmt.Sprintf("面积=长×宽=%d×%d=%d", l, w, l*w))
	}
	for i := 0; i < 15; i++ {
		a, b := rng.Intn(8)+2, rng.Intn(20)+2
		x := rng.Intn(12)+2
		add(5, "简易方程",
			fmt.Sprintf("解方程：%dx + %d = %d，x = ?", a, b, a*x+b),
			itoa(x), numOptions(rng, x, 15), fmt.Sprintf("%dx=%d，x=%d", a, a*x, x))
	}
	// 6 年级：分数运算 / 百分数 / 比例
	fracs := [][2]int{{1, 2}, {1, 3}, {2, 3}, {1, 4}, {3, 4}, {1, 5}, {2, 5}, {3, 5}, {1, 6}, {5, 6}, {1, 8}, {3, 8}}
	for i := 0; i < 25; i++ {
		f1, f2 := fracs[rng.Intn(len(fracs))], fracs[rng.Intn(len(fracs))]
		num, den := f1[0]*f2[1]+f2[0]*f1[1], f1[1]*f2[1]
		g := gcd(num, den)
		num, den = num/g, den/g
		add(6, "分数运算",
			fmt.Sprintf("计算：%d/%d + %d/%d = ?", f1[0], f1[1], f2[0], f2[1]),
			fmt.Sprintf("%d/%d", num, den),
			fracOptions(rng, num, den),
			fmt.Sprintf("通分后相加：%d/%d", num, den))
	}
	for i := 0; i < 15; i++ {
		p := []int{10, 20, 25, 50, 75}[rng.Intn(5)]
		n := (rng.Intn(18) + 2) * 10
		add(6, "百分数",
			fmt.Sprintf("%d 的 %d%% 是多少？", n, p),
			itoa(n*p/100), numOptions(rng, n*p/100, 30), fmt.Sprintf("%d×%d%%=%d", n, p, n*p/100))
	}
	// 7 年级（初一）：有理数 / 整式 / 一元一次方程
	for i := 0; i < 25; i++ {
		a, b := rng.Intn(40)-20, rng.Intn(40)-20
		add(7, "有理数",
			fmt.Sprintf("计算：(%d) + (%d) = ?", a, b),
			itoa(a+b), signedOptions(rng, a+b), fmt.Sprintf("(%d)+(%d)=%d", a, b, a+b))
	}
	for i := 0; i < 15; i++ {
		a, b := rng.Intn(9)+2, rng.Intn(9)+2
		add(7, "整式运算",
			fmt.Sprintf("计算：%d² × %d = ?", a, b),
			itoa(a*a*b), numOptions(rng, a*a*b, 60), fmt.Sprintf("%d²=%d，%d×%d=%d", a, a*a, a*a, b, a*a*b))
	}
	for i := 0; i < 15; i++ {
		a := rng.Intn(6) + 2
		x := rng.Intn(15) + 2
		b := rng.Intn(30) - 15
		add(7, "一元一次方程",
			fmt.Sprintf("解方程：%dx %s = %d，x = ?", a, signedTerm(b), a*x-b),
			itoa(x), numOptions(rng, x, 20), fmt.Sprintf("%dx=%d，x=%d", a, a*x, x))
	}
	// 8 年级（初二）：一次函数 / 勾股定理 / 二次根式
	triples := [][3]int{{3, 4, 5}, {6, 8, 10}, {5, 12, 13}, {9, 12, 15}, {8, 15, 17}, {7, 24, 25}}
	for i := 0; i < 15; i++ {
		t := triples[rng.Intn(len(triples))]
		add(8, "勾股定理",
			fmt.Sprintf("直角三角形两直角边分别为 %d 和 %d，斜边长为？", t[0], t[1]),
			itoa(t[2]), numOptions(rng, t[2], 25), fmt.Sprintf("√(%d²+%d²)=√%d=%d", t[0], t[1], t[2]*t[2], t[2]))
	}
	for i := 0; i < 12; i++ {
		k, m := rng.Intn(6)+1, rng.Intn(20)-10
		xv := rng.Intn(6) + 1
		add(8, "一次函数",
			fmt.Sprintf("直线 y = %dx %s，当 x = %d 时 y = ?", k, signedConst(m), xv),
			itoa(k*xv+m), signedOptions(rng, k*xv+m), fmt.Sprintf("y=%d×%d%s=%d", k, xv, signedConst(m), k*xv+m))
	}
	for i := 0; i < 12; i++ {
		n := []int{4, 9, 16, 25, 36, 49, 64, 81, 100}[rng.Intn(9)]
		add(8, "二次根式",
			fmt.Sprintf("化简：√%d = ?", n * n),
			itoa(n), numOptions(rng, n, 15), fmt.Sprintf("√%d=%d", n*n, n))
	}
	// 9 年级（初三）：二次函数 / 一元二次方程 / 圆
	for i := 0; i < 15; i++ {
		r1, r2 := rng.Intn(8)+1, rng.Intn(8)+1
		b := -(r1 + r2)
		c := r1 * r2
		add(9, "一元二次方程",
			fmt.Sprintf("方程 x² %s %d = 0 的解是 x₁=%d，x₂=？", signedTerm(b), c, r1),
			itoa(r2), numOptions(rng, r2, 12), fmt.Sprintf("因式分解：(x−%d)(x−%d)=0", r1, r2))
	}
	for i := 0; i < 12; i++ {
		k := -rng.Intn(6) - 1
		xv := rng.Intn(5) + 1
		add(9, "二次函数",
			fmt.Sprintf("抛物线 y = x² %dx，当 x = %d 时 y = ?", k, xv),
			itoa(xv*xv+k*xv), signedOptions(rng, xv*xv+k*xv), fmt.Sprintf("y=%d²+(%d)×%d=%d", xv, k, xv, xv*xv+k*xv))
	}
	for i := 0; i < 10; i++ {
		r := rng.Intn(9) + 2
		area := fmt.Sprintf("%.2f", 3.14*float64(r*r))
		add(9, "圆",
			fmt.Sprintf("半径为 %d 的圆，面积是多少？（π≈3.14）", r),
			area, numOptionsS(rng, area, 3.14*float64(r*r)),
			fmt.Sprintf("S=πr²=3.14×%d²=%.2f", r, 3.14*float64(r*r)))
	}
	return qs
}

func pickGrade(rng *rand.Rand, a, b int) int { return a + rng.Intn(b-a+1) }
func itoa(v int) string                      { return fmt.Sprintf("%d", v) }
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}
func signedTerm(v int) string {
	if v >= 0 {
		return fmt.Sprintf("+ %d", v)
	}
	return fmt.Sprintf("− %d", -v)
}
func signedConst(v int) string {
	if v >= 0 {
		return fmt.Sprintf("+ %d", v)
	}
	return fmt.Sprintf("− %d", -v)
}

// numOptions 生成 4 个互不相同的数值选项（包含正确答案，干扰项为邻近数值）。
func numOptions(rng *rand.Rand, correct int, spread int) []string {
	set := map[string]bool{itoa(correct): true}
	opts := []string{itoa(correct)}
	for len(opts) < 4 {
		d := rng.Intn(2*spread+1) - spread
		if d == 0 || correct+d < 0 {
			continue
		}
		v := itoa(correct + d)
		if !set[v] {
			set[v] = true
			opts = append(opts, v)
		}
	}
	return opts
}

// numOptionsS 字符串版：正确答案 fixed，干扰项为数值扰动。
func numOptionsS(rng *rand.Rand, fixed string, correct float64) []string {
	set := map[string]bool{fixed: true}
	opts := []string{fixed}
	for len(opts) < 4 {
		d := float64(rng.Intn(21)-10) / 10
		if d == 0 {
			continue
		}
		v := fmt.Sprintf("%.2f", correct+d)
		if !set[v] {
			set[v] = true
			opts = append(opts, v)
		}
	}
	return opts
}

// signedOptions 含负数的选项。
func signedOptions(rng *rand.Rand, correct int) []string {
	set := map[string]bool{itoa(correct): true}
	opts := []string{itoa(correct)}
	for len(opts) < 4 {
		d := rng.Intn(25) - 12
		if d == 0 {
			continue
		}
		v := itoa(correct + d)
		if !set[v] {
			set[v] = true
			opts = append(opts, v)
		}
	}
	return opts
}

// fracOptions 分数干扰项（同分母邻近值）。
func fracOptions(rng *rand.Rand, num, den int) []string {
	ans := fmt.Sprintf("%d/%d", num, den)
	set := map[string]bool{ans: true}
	opts := []string{ans}
	for len(opts) < 4 {
		n := num + rng.Intn(7) - 3
		if n <= 0 || n == num {
			continue
		}
		v := fmt.Sprintf("%d/%d", n, den)
		if !set[v] {
			set[v] = true
			opts = append(opts, v)
		}
	}
	return opts
}
