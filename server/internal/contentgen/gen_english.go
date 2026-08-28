package contentgen

import "studyplanet/internal/contentlib"

// genEnglish 分级词表：小学 1-6 年级 / 初中 7-9 年级常用词。
// 每个词条生成一道「选释义」选择题；部分年级附带例句填空题。
func genEnglish() []contentlib.Question {
	var qs []contentlib.Question
	add := func(grade int, word, meaning, phonetic, example string, distract []string) {
		opts := shuffleStr(append(append([]string{}, distract...), meaning))
		qs = append(qs, contentlib.Question{
			Subject: "english", Grade: grade, Topic: "词汇", QType: "choice",
			Question:  "单词「" + word + "」的中文意思是？",
			Options:   opts,
			Answer:    meaning,
			Explanation: word + " " + phonetic + " " + example,
			Source:    "builtin",
		})
		_ = phonetic
	}
	// 1-2 年级基础词
	g1 := [][2]string{
		{"apple", "苹果"}, {"cat", "猫"}, {"dog", "狗"}, {"egg", "鸡蛋"}, {"fish", "鱼"},
		{"girl", "女孩"}, {"hand", "手"}, {"ice", "冰"}, {"juice", "果汁"}, {"kite", "风筝"},
		{"lion", "狮子"}, {"milk", "牛奶"}, {"nose", "鼻子"}, {"orange", "橙子"}, {"pig", "猪"},
		{"queen", "女王"}, {"rain", "雨"}, {"sun", "太阳"}, {"tree", "树"}, {"water", "水"},
	}
	g1d := [][]string{{"香蕉", "桌子"}, {"鸟", "帽子"}, {"鸭子", "猪"}, {"眼睛", "大象"}, {"鸟", "牛奶"},
		{"男孩", "书"}, {"头", "脚"}, {"火", "米饭"}, {"茶", "果汁"}, {"鸟", "包"},
		{"老虎", "熊"}, {"米饭", "面包"}, {"耳朵", "嘴巴"}, {"香蕉", "葡萄"}, {"鸭子", "母鸡"},
		{"国王", "王子"}, {"雪", "风"}, {"月亮", "星星"}, {"花", "草"}, {"火", "冰"}}
	for i, w := range g1 {
		add(1+i%2, w[0], w[1], "", "", g1d[i])
	}
	// 3-4 年级
	g3 := [][3]string{
		{"animal", "动物", "An animal is a living thing."}, {"birthday", "生日", "Happy birthday to you!"},
		{"children", "孩子们", "The children are playing."}, {"different", "不同的", "We have different ideas."},
		{"evening", "晚上", "Good evening, everyone."}, {"family", "家庭", "I love my family."},
		{"garden", "花园", "Flowers grow in the garden."}, {"holiday", "假期", "We travel during the holiday."},
		{"important", "重要的", "This is an important question."}, {"kitchen", "厨房", "Mom cooks in the kitchen."},
		{"language", "语言", "English is a world language."}, {"mountain", "山", "The mountain is very high."},
		{"number", "数字", "Do you know my phone number?"}, {"office", "办公室", "Dad works in an office."},
		{"picture", "图画", "Draw a picture of your family."}, {"question", "问题", "May I ask a question?"},
		{"restaurant", "餐馆", "We eat at a restaurant."}, {"school", "学校", "I go to school by bus."},
		{"teacher", "老师", "My teacher is very kind."}, {"window", "窗户", "Open the window, please."},
	}
	g3d := [][]string{{"植物", "机器人"}, {"星期", "节日"}, {"同学", "父母"}, {"困难的", "相似的"}, {"上午", "中午"},
		{"学校", "班级"}, {"森林", "公园"}, {"周末", "星期天"}, {"有趣的", "简单的"}, {"客厅", "卧室"},
		{"方言", "历史"}, {"河流", "湖泊"}, {"字母", "颜色"}, {"工厂", "医院"}, {"照片", "文章"},
		{"答案", "作业"}, {"旅馆", "厨房"}, {"公园", "教室"}, {"学生", "医生"}, {"门", "墙"}}
	for i, w := range g3 {
		add(3+i%2, w[0], w[1], "", w[2], g3d[i])
	}
	// 5-6 年级
	g5 := [][3]string{
		{"because", "因为", "I stayed home because it rained."}, {"beautiful", "美丽的", "She has a beautiful smile."},
		{"favorite", "最喜欢的", "This is my favorite book."}, {"usually", "通常", "We usually eat dinner at seven."},
		{"expensive", "昂贵的", "The phone is too expensive."}, {"dangerous", "危险的", "Be careful, it is dangerous."},
		{"library", "图书馆", "I read books in the library."}, {"Wednesday", "星期三", "We have art on Wednesday."},
		{"dictionary", "字典", "Use a dictionary to learn words."}, {"exercise", "练习", "Do exercise every morning."},
		{"knowledge", "知识", "Knowledge is power."}, {"necessary", "必要的", "Water is necessary for life."},
		{"opinion", "观点", "In my opinion, he is right."}, {"possible", "可能的", "It is possible to finish today."},
		{"restaurant", "餐馆", "The restaurant is very famous."}, {"surprise", "惊喜", "What a nice surprise!"},
		{"together", "一起", "Let's work together."}, {"vegetable", "蔬菜", "Eat more vegetables."},
		{"weather", "天气", "The weather is nice today."}, {"yesterday", "昨天", "It rained yesterday."},
	}
	g5d := [][]string{{"所以", "但是"}, {"聪明的", "可爱的"}, {"讨厌的", "特别的"}, {"从不", "总是"}, {"便宜的", "结实的"},
		{"安全的", "兴奋的"}, {"医院", "市场"}, {"星期二", "星期四"}, {"小说", "杂志"}, {"休息", "睡眠"},
		{"力量", "财富"}, {"有用的", "免费的"}, {"错误", "事实"}, {"不可能的", "困难的"}, {"旅馆", "银行"},
		{"礼物", "派对"}, {"分开", "前面"}, {"水果", "肉类"}, {"季节", "气候"}, {"明天", "今天"}}
	for i, w := range g5 {
		add(5+i%2, w[0], w[1], "", w[2], g5d[i])
	}
	// 7-9 年级（初中）
	g7 := [][3]string{
		{"ability", "能力", "She has the ability to sing well."}, {"accept", "接受", "Please accept my gift."},
		{"achieve", "实现", "You can achieve your dream."}, {"advantage", "优势", "Height is his advantage."},
		{"behavior", "行为", "His behavior was very rude."}, {"celebrate", "庆祝", "We celebrate the New Year."},
		{"community", "社区", "Our community is friendly."}, {"confident", "自信的", "Be confident in yourself."},
		{"decision", "决定", "Make a decision right now."}, {"education", "教育", "Education is important."},
		{"environment", "环境", "Protect the environment."}, {"experience", "经验", "Travel gives you experience."},
		{"government", "政府", "The government made a new rule."}, {"improve", "改进", "I want to improve my English."},
		{"influence", "影响", "Parents influence their kids."}, {"knowledge", "知识", "Knowledge changes fate."},
		{"memory", "记忆", "This photo brings back memories."}, {"opportunity", "机会", "Don't miss this opportunity."},
		{"patient", "耐心的", "Please be patient with kids."}, {"purpose", "目的", "The purpose of the meeting is clear."},
	}
	g7d := [][]string{{"态度", "年龄"}, {"拒绝", "期待"}, {"放弃", "失败"}, {"劣势", "高度"}, {"语言", "性格"},
		{"参加", "祝福"}, {"国家", "村庄"}, {"谦虚的", "勇敢的"}, {"建议", "结论"}, {"娱乐", "健康"},
		{"能量", "资源"}, {"旅行", "考试"}, {"规则", "法律"}, {"翻译", "提高"}, {"信任", "影响"},
		{"力量", "朋友"}, {"记忆", "照片"}, {"困难", "挑战"}, {"生病的", "聪明的"}, {"原因", "结果"}}
	for i, w := range g7 {
		add(7+i%3, w[0], w[1], "", w[2], g7d[i])
	}
	return qs
}
