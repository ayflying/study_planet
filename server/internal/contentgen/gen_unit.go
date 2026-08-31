// Package contentgen 人教版单元测试题库：语文（1-6年级）、数学（1-6年级）、英语（3-6年级）。
// 按教材单元划分，每单元 3-4 道选择题，覆盖各单元核心知识点。
package contentgen

import "studyplanet/internal/contentlib"

// genUnitTests 生成人教版单元测试试题。
func genUnitTests() []contentlib.Question {
	var qs []contentlib.Question
	qs = append(qs, genChineseUnitTest()...)
	qs = append(qs, genMathUnitTest()...)
	qs = append(qs, genEnglishUnitTest()...)
	return qs
}

// ---------- 语文单元测试（人教版部编版 1-6年级上册）----------

func genChineseUnitTest() []contentlib.Question {
	var qs []contentlib.Question
	add := func(grade int, topic, question string, options []string, answer, explanation string) {
		qs = append(qs, contentlib.Question{
			Subject: "chinese", Grade: grade, Topic: topic, QType: "choice",
			Question: question, Options: options, Answer: answer, Explanation: explanation, Source: "builtin",
		})
	}

	// ======== 一年级 ========
	// 识字单元
	add(1, "识字", "下列哪个字的拼音是「mā」？", []string{"妈", "爸", "马", "大"}, "妈", "妈（mā）是母亲的称呼。")
	add(1, "识字", "「山水」中「山」的笔画数是？", []string{"3", "2", "4", "5"}, "3", "山字共3画：竖、竖折、竖。")
	add(1, "拼音", "下列哪个是整体认读音节？", []string{"zhi", "zi", "ci", "以上都是"}, "以上都是", "zhi、zi、ci 都是整体认读音节。")
	add(1, "课文", "《秋天》一文中，天气变凉了，树叶怎么了？", []string{"变黄了", "变红了", "落下来了", "变绿了"}, "落下来了", "秋天树叶从树上落下来。")
	add(1, "课文", "《小小的船》中「弯弯的月儿小小的船，小小的船儿两头尖」，月儿像什么？", []string{"小船", "香蕉", "镰刀", "圆盘"}, "小船", "弯弯的月儿像两头尖尖的小船。")
	add(1, "识字", "「日」加一笔是什么字？", []string{"目", "田", "申", "旧"}, "目", "日在内部加一横变为目。")
	add(1, "课文", "《四季》中草芽尖尖是什么季节？", []string{"春天", "夏天", "秋天", "冬天"}, "春天", "草芽尖尖是春天的特征。")
	add(1, "识字", "「大小」的反义词是？", []string{"多少", "上下", "左右", "前后"}, "多少", "大对少、小对多，大小对多少。")

	// ======== 二年级 ========
	add(2, "课文", "《小蝌蚪找妈妈》中，小蝌蚪长大后变成了什么？", []string{"青蛙", "蟾蜍", "鱼", "乌龟"}, "青蛙", "小蝌蚪长大后变成青蛙。")
	add(2, "识字", "「傍」字的部首是什么？", []string{"亻", "氵", "攵", "宀"}, "亻", "傍是单人旁，表示与人有关。")
	add(2, "课文", "《曹冲称象》中，曹冲用什么方法称象？", []string{"石头替换", "直接称", "切块称", "用水称"}, "石头替换", "曹冲把大象换成石头，分次称量。")
	add(2, "古诗", "《望庐山瀑布》中「飞流直下三千尺」的下一句是？", []string{"疑是银河落九天", "日照香炉生紫烟", "遥看瀑布挂前川", "两岸青山相对出"}, "疑是银河落九天", "李白《望庐山瀑布》：飞流直下三千尺，疑是银河落九天。")
	add(2, "识字", "「纸」字的最后一笔是什么？", []string{"点", "横", "竖", "撇"}, "点", "纸的笔顺最后是点。")
	add(2, "课文", "《黄山奇石》中「猴子观海」的石头像什么？", []string{"猴子", "桃子", "狮子", "大象"}, "猴子", "猴子观海是一块像猴子的奇石。")
	add(2, "识字", "「容」字的部首是？", []string{"宀", "穴", "口", "谷"}, "宀", "容是宝盖头，表示与房屋有关。")
	add(2, "课文", "《狐假虎威》中狐狸借用了谁的威风？", []string{"老虎", "狮子", "大象", "豹子"}, "老虎", "狐假虎威指狐狸借老虎的威风吓唬其他动物。")

	// ======== 三年级 ========
	add(3, "课文", "《大青树下的小学》中，学校在什么地方？", []string{"大青树下", "山顶上", "小河边", "草原上"}, "大青树下", "这是一所边疆小学，在大青树下面。")
	add(3, "词语", "下列词语中属于比喻义的是？", []string{"锦上添花", "一帆风顺", "一马当先", "一心一意"}, "锦上添花", "锦上添花比喻好上加好。")
	add(3, "古诗", "《山行》中「停车坐爱枫林晚」的「坐」是什么意思？", []string{"因为", "坐下", "乘坐", "坐落"}, "因为", "坐在这里是「因为」的意思。")
	add(3, "阅读", "一段话中概括主要意思的句子叫做？", []string{"中心句", "开头句", "过渡句", "结尾句"}, "中心句", "中心句概括段落的主要意思。")
	add(3, "词语", "「勇往直前」的反义词是？", []string{"畏缩不前", "一往无前", "奋勇向前", "所向披靡"}, "畏缩不前", "勇往直前反义是畏缩不前。")
	add(3, "课文", "《富饶的西沙群岛》中，西沙群岛位于哪个海？", []string{"南海", "东海", "黄海", "渤海"}, "南海", "西沙群岛位于中国南海。")
	add(3, "古诗", "《饮湖上初晴后雨》中「淡妆浓抹总相宜」描写的是？", []string{"西湖", "洞庭湖", "太湖", "鄱阳湖"}, "西湖", "苏轼用西施比喻西湖。")
	add(3, "修辞", "「春天像小姑娘，花枝招展的」运用了什么修辞？", []string{"拟人", "比喻", "夸张", "排比"}, "拟人", "把春天比作小姑娘，运用了拟人手法。")

	// ======== 四年级 ========
	add(4, "课文", "《观潮》描写的是哪里的潮水？", []string{"钱塘江", "长江", "黄河", "珠江"}, "钱塘江", "钱塘江大潮是世界奇观。")
	add(4, "词语", "下列词语中「鼎」字的意思与什么有关？", []string{"器物", "烹饪", "声音", "数量"}, "器物", "鼎是古代煮东西的器物。")
	add(4, "古诗", "《题西林壁》中「不识庐山真面目」的原因是？", []string{"身在此山中", "山太高了", "云雾遮住", "距离太远"}, "身在此山中", "只缘身在此山中，旁观者清。")
	add(4, "阅读", "说明文中常用的说明方法不包括？", []string{"比喻", "列数字", "举例子", "作比较"}, "比喻", "比喻是修辞手法，不是说明方法。")
	add(4, "课文", "《蟋蟀的住宅》中蟋蟀的住宅建造在？", []string{"草丛中", "洞穴里", "河岸边", "树干上"}, "草丛中", "蟋蟀的住宅建在草丛中。")
	add(4, "词语", "「随遇而安」的意思是？", []string{"顺应环境安于现状", "随朋友一起安家", "随意改变主意", "偶然遇到故人"}, "顺应环境安于现状", "随遇而安指能适应各种环境。")
	add(4, "古诗", "《出塞》中「但使龙城飞将在」的「飞将」指谁？", []string{"李广", "卫青", "霍去病", "岳飞"}, "李广", "飞将军李广，汉朝名将。")
	add(4, "修辞", "下列句子中运用了排比修辞的是？", []string{"有的…有的…有的…", "月亮像小船", "小鸟在唱歌", "风儿轻轻吹"}, "有的…有的…有的…", "排比是把三个或以上结构相似的句子并列。")

	// ======== 五年级 ========
	add(5, "课文", "《桂花雨》中作者最怀念的是？", []string{"故乡的桂花", "故乡的景色", "故乡的亲人", "故乡的房屋"}, "故乡的桂花", "作者通过桂花表达思乡之情。")
	add(5, "词语", "「呕心沥血」形容什么？", []string{"费尽心思", "吐了很多血", "非常辛苦", "用力呼吸"}, "费尽心思", "呕心沥血比喻费尽心思和精力。")
	add(5, "古诗", "《示儿》中「王师北定中原日」的下一句是？", []string{"家祭无忘告乃翁", "但悲不见九州同", "死去元知万事空", "每逢佳节倍思亲"}, "家祭无忘告乃翁", "陆游临终前盼望国家统一。")
	add(5, "阅读", "阅读时抓住关键词句有助于？", []string{"理解文章主旨", "加快阅读速度", "记住生字", "背诵全文"}, "理解文章主旨", "关键词句是理解文章的重要线索。")
	add(5, "课文", "《猎人海力布》中海力布变成了什么？", []string{"石头", "大树", "小鸟", "河流"}, "石头", "海力布为了救乡亲们变成了石头。")
	add(5, "词语", "下列词语中与「举世闻名」意思相近的是？", []string{"闻名遐迩", "默默无闻", "无人问津", "名不见经传"}, "闻名遐迩", "举世闻名和闻名遐迩都形容非常有名。")
	add(5, "古诗", "《山居秋暝》中「空山新雨后，天气晚来秋」描写了什么季节？", []string{"秋天", "春天", "夏天", "冬天"}, "秋天", "从「秋」字可知描写的是秋天。")
	add(5, "修辞", "「书籍是人类进步的阶梯」用了什么修辞？", []string{"比喻", "拟人", "夸张", "对偶"}, "比喻", "把书籍比作进步的阶梯，是比喻。")

	// ======== 六年级 ========
	add(6, "课文", "《草原》中老舍先生描写了哪里的草原？", []string{"内蒙古", "新疆", "西藏", "青海"}, "内蒙古", "文章描写了内蒙古草原的美丽风光。")
	add(6, "词语", "「别出心裁」的意思是？", []string{"独创一格与众不同", "心中另有打算", "离开原来的计划", "裁缝的新手艺"}, "独创一格与众不同", "别出心裁指想出的办法与众不同。")
	add(6, "古诗", "《七律·长征》中「金沙水拍云崖暖」的「暖」字写出了什么？", []string{"红军巧渡金沙江的喜悦", "天气非常炎热", "金沙江水温很高", "阳光照射很暖和"}, "红军巧渡金沙江的喜悦", "暖字表达了红军胜利渡江的喜悦心情。")
	add(6, "阅读", "文章的过渡句通常出现在？", []string{"段与段之间", "文章开头", "文章结尾", "标题中"}, "段与段之间", "过渡句起承上启下的作用。")
	add(6, "课文", "《少年闰土》中闰土给作者讲了哪些新鲜事？", []string{"雪地捕鸟看瓜刺猹", "打鱼砍柴", "放牛割草", "种地浇水"}, "雪地捕鸟看瓜刺猹", "闰土给作者讲了农村的新鲜事。")
	add(6, "词语", "「高山流水」比喻什么？", []string{"知音难觅", "山高水长", "风景优美", "水流湍急"}, "知音难觅", "高山流水比喻知音难遇或乐曲高妙。")
	add(6, "古诗", "《春夜喜雨》中「好雨知时节，当春乃发生」的「好」字表达了什么？", []string{"对春雨的喜爱", "雨下得很大", "雨来得及时", "雨很干净"}, "对春雨的喜爱", "好字表达了诗人对春雨的喜爱之情。")
	add(6, "文学常识", "中国四大名著中《红楼梦》的作者是？", []string{"曹雪芹", "罗贯中", "施耐庵", "吴承恩"}, "曹雪芹", "曹雪芹著《红楼梦》，是中国古典小说的巅峰。")

	return qs
}

// ---------- 数学单元测试（人教版 1-6年级上册）----------

func genMathUnitTest() []contentlib.Question {
	var qs []contentlib.Question
	add := func(grade int, topic, question string, options []string, answer, explanation string) {
		qs = append(qs, contentlib.Question{
			Subject: "math", Grade: grade, Topic: topic, QType: "choice",
			Question: question, Options: options, Answer: answer, Explanation: explanation, Source: "builtin",
		})
	}

	// ======== 一年级 ========
	add(1, "数一数", "数一数：下面哪个数比 5 大？", []string{"6", "3", "4", "2"}, "6", "6 比 5 大。")
	add(1, "比一比", "比一比：3 和 7 谁大？", []string{"7", "3", "一样大", "不知道"}, "7", "7 大于 3。")
	add(1, "认识图形", "下面哪个是正方体？", []string{"骰子", "皮球", "罐头", "铅笔"}, "骰子", "骰子是正方体，六个面都是正方形。")
	add(1, "加减法", "小明有 3 个苹果，妈妈又给了他 2 个，他一共有几个？", []string{"5", "4", "6", "3"}, "5", "3+2=5。")
	add(1, "位置", "你的右手边是？", []string{"右手", "左手", "前面", "后面"}, "右手", "右手边就是右手的方向。")
	add(1, "钟表", "分针指向 12，时针指向 3，是几点？", []string{"3点", "12点", "6点", "9点"}, "3点", "时针指向 3 就是 3 点。")
	add(1, "加减法", "10 − 4 = ?", []string{"6", "5", "7", "4"}, "6", "10−4=6。")
	add(1, "分类", "下面哪个是水果？", []string{"苹果", "铅笔", "橡皮", "尺子"}, "苹果", "苹果是水果，其他是文具。")

	// ======== 二年级 ========
	add(2, "长度单位", "1 米等于多少厘米？", []string{"100", "10", "1000", "50"}, "100", "1 米 = 100 厘米。")
	add(2, "加减法", "25 + 37 = ?", []string{"62", "52", "72", "60"}, "62", "25+37=62。")
	add(2, "乘法", "3 × 4 = ?", []string{"12", "7", "15", "9"}, "12", "3×4=12，三四十二。")
	add(2, "角的初步认识", "直角是几度？", []string{"90°", "45°", "180°", "360°"}, "90°", "直角等于 90 度。")
	add(2, "乘法", "5 × 6 = ?", []string{"30", "25", "35", "11"}, "30", "5×6=30，五六三十。")
	add(2, "观察物体", "从正面看一个长方体，看到的形状是？", []string{"长方形", "正方形", "圆形", "三角形"}, "长方形", "从正面看长方体通常是长方形。")
	add(2, "除法", "12 ÷ 3 = ?", []string{"4", "3", "6", "9"}, "4", "12÷3=4，三四十二。")
	add(2, "混合运算", "2 × 3 + 4 = ?", []string{"10", "14", "9", "12"}, "10", "先乘后加：2×3=6，6+4=10。")

	// ======== 三年级 ========
	add(3, "时分秒", "1 分钟等于多少秒？", []string{"60", "100", "30", "120"}, "60", "1 分钟 = 60 秒。")
	add(3, "万以内的加减法", "500 + 300 = ?", []string{"800", "700", "600", "900"}, "800", "500+300=800。")
	add(3, "测量", "1 千克等于多少克？", []string{"1000", "100", "10", "10000"}, "1000", "1 千克 = 1000 克。")
	add(3, "长方形和正方形", "长方形的周长公式是？", []string{"(长+宽)×2", "长×宽", "长+宽", "边长×4"}, "(长+宽)×2", "长方形周长 = 2×(长+宽)。")
	add(3, "分数初步", "把一个蛋糕平均分成 4 份，每份是几分之几？", []string{"1/4", "1/2", "1/3", "1/8"}, "1/4", "平均分成 4 份，每份是四分之一。")
	add(3, "乘法", "12 × 3 = ?", []string{"36", "32", "38", "40"}, "36", "12×3=36。")
	add(3, "集合", "「所有三角形的内角和都是」？", []string{"180°", "90°", "360°", "270°"}, "180°", "三角形内角和为 180 度。")
	add(3, "位置与方向", "太阳从哪个方向升起？", []string{"东方", "南方", "西方", "北方"}, "东方", "太阳从东方升起。")

	// ======== 四年级 ========
	add(4, "大数的认识", "10 个一万是多少？", []string{"十万", "一百万", "一千万", "一亿"}, "十万", "10×10000=100000，即十万。")
	add(4, "公顷和平方千米", "1 平方千米等于多少公顷？", []string{"100", "10", "1000", "10000"}, "100", "1 平方千米 = 100 公顷。")
	add(4, "角的度量", "平角等于多少度？", []string{"180°", "90°", "360°", "270°"}, "180°", "平角 = 180 度。")
	add(4, "三位数乘两位数", "20 × 30 = ?", []string{"600", "60", "6000", "60000"}, "600", "20×30=600。")
	add(4, "平行四边形和梯形", "平行四边形有几条边？", []string{"4", "3", "5", "6"}, "4", "平行四边形有 4 条边。")
	add(4, "除数是两位数的除法", "800 ÷ 20 = ?", []string{"40", "400", "4", "80"}, "40", "800÷20=40。")
	add(4, "条形统计图", "条形统计图的优点是什么？", []string{"直观比较数量", "便于计算", "显示比例", "美观好看"}, "直观比较数量", "条形统计图能直观比较不同数量的大小。")
	add(4, "数学广角", "烙一张饼需要 2 分钟（正反面各 1 分钟），同时烙 3 张饼最少需要几分钟？", []string{"3", "4", "5", "6"}, "3", "3 张饼交替烙，3 分钟即可。")

	// ======== 五年级 ========
	add(5, "小数乘法", "0.5 × 6 = ?", []string{"3", "0.3", "30", "0.03"}, "3", "0.5×6=3。")
	add(5, "位置", "数对(3,2)表示第几列第几行？", []string{"第3列第2行", "第2列第3行", "第3行第2列", "第2行第3列"}, "第3列第2行", "数对(列,行)，(3,2)表示第3列第2行。")
	add(5, "小数除法", "3.6 ÷ 0.6 = ?", []string{"6", "0.6", "60", "0.06"}, "6", "3.6÷0.6=6。")
	add(5, "可能性", "从装有 3 个红球和 1 个白球的盒子里摸一个球，摸到哪种球的可能性大？", []string{"红球", "白球", "一样大", "不确定"}, "红球", "红球数量多，摸到的可能性大。")
	add(5, "简易方程", "方程 x + 5 = 12 的解是 x = ?", []string{"7", "5", "12", "17"}, "7", "x=12−5=7。")
	add(5, "多边形的面积", "底为 6cm、高为 4cm 的三角形面积是多少？", []string{"12cm²", "24cm²", "10cm²", "20cm²"}, "12cm²", "三角形面积=底×高÷2=6×4÷2=12。")
	add(5, "植树问题", "在一条 100 米的小路一边每隔 5 米栽一棵树（两端都栽），需要多少棵树？", []string{"21", "20", "19", "22"}, "21", "100÷5+1=20+1=21 棵。")
	add(5, "数学广角", "鸡兔同笼，头 10 个，脚 28 只，鸡有几只？", []string{"6", "4", "5", "8"}, "6", "假设全是兔：10×4=40 脚，多 12 脚，每换一只鸡少 2 脚，12÷2=6 只鸡。")

	// ======== 六年级 ========
	add(6, "分数乘法", "1/2 × 1/3 = ?", []string{"1/6", "1/5", "2/5", "1/3"}, "1/6", "分子相乘，分母相乘：1×1/2×3=1/6。")
	add(6, "位置与方向", "东偏北 30° 是指从东方向向哪个方向偏转？", []string{"北", "南", "西", "东"}, "北", "东偏北即从东向北偏转。")
	add(6, "分数除法", "1/2 ÷ 1/4 = ?", []string{"2", "1/2", "1/8", "4"}, "2", "除以一个分数等于乘以它的倒数：1/2×4=2。")
	add(6, "比", "一个班级男生 20 人，女生 15 人，男女生比是？", []string{"4:3", "3:4", "5:3", "3:5"}, "4:3", "20:15=4:3。")
	add(6, "圆", "圆的周长公式是？", []string{"C=2πr", "C=πr", "C=πr²", "C=2r"}, "C=2πr", "周长=2×圆周率×半径。")
	add(6, "百分数", "1/4 化成百分数是？", []string{"25%", "40%", "75%", "50%"}, "25%", "1/4=0.25=25%。")
	add(6, "扇形统计图", "扇形统计图适合表示什么？", []string{"各部分与整体的比例", "数量变化趋势", "数据大小比较", "数据分布情况"}, "各部分与整体的比例", "扇形统计图直观显示各部分占比。")
	add(6, "数学广角", "用数字 1、2、3 可以组成几个不同的两位数？", []string{"6", "3", "9", "4"}, "6", "3×2=6 个：12,13,21,23,31,32。")

	return qs
}

// ---------- 英语单元测试（PEP 人教版 3-6年级上册）----------

func genEnglishUnitTest() []contentlib.Question {
	var qs []contentlib.Question
	add := func(grade int, topic, question string, options []string, answer, explanation string) {
		qs = append(qs, contentlib.Question{
			Subject: "english", Grade: grade, Topic: topic, QType: "choice",
			Question: question, Options: options, Answer: answer, Explanation: explanation, Source: "builtin",
		})
	}

	// ======== 三年级（PEP 上册）=======
	add(3, "Hello", "早上见到老师，你应该说？", []string{"Good morning.", "Goodbye.", "Good night.", "Good afternoon."}, "Good morning.", "早上好用 Good morning。")
	add(3, "Hello", "「How are you？」的答语是？", []string{"I'm fine, thank you.", "How are you?", "Goodbye.", "My name is Mike."}, "I'm fine, thank you.", "How are you？回答 I'm fine, thank you。")
	add(3, "Colours", "红色用英语怎么说？", []string{"red", "yellow", "blue", "green"}, "red", "红色是 red。")
	add(3, "Colours", "I see green. 的中文意思是？", []string{"我看到绿色。", "我喜欢绿色。", "这是绿色。", "绿色很美。"}, "我看到绿色。", "I see green = 我看到绿色。")
	add(3, "Animals", "What's this? It's a ____.（猫）", []string{"cat", "dog", "duck", "pig"}, "cat", "猫是 cat。")
	add(3, "Animals", "鸭子用英语怎么说？", []string{"duck", "dog", "cat", "bear"}, "duck", "鸭子是 duck。")
	add(3, "Food", "I'd like some ____.（牛奶）", []string{"milk", "bread", "cake", "juice"}, "milk", "牛奶是 milk。")
	add(3, "Food", "「吃一些面包」用英语怎么说？", []string{"Have some bread.", "Eat some rice.", "Drink some milk.", "Cut some cake."}, "Have some bread.", "bread 是面包。")
	add(3, "Toys", "Show me your ____.（铅笔）", []string{"pencil", "ruler", "eraser", "crayon"}, "pencil", "pencil 是铅笔。")
	add(3, "Toys", "I have a ____.（书包）", []string{"schoolbag", "book", "pen", "bag"}, "schoolbag", "schoolbag 是书包。")

	// ======== 四年级（PEP 上册）=======
	add(4, "School", "Where is the library? 是问什么？", []string{"图书馆在哪里？", "这是什么？", "你是谁？", "几点了？"}, "图书馆在哪里？", "Where is the library? 问图书馆在哪里。")
	add(4, "School", "Go to the ____. Read a book.（图书馆）", []string{"library", "classroom", "playground", "canteen"}, "library", "图书馆是读 book 的地方。")
	add(4, "Schoolbag", "What's in your schoolbag? 的正确回答是？", []string{"An English book.", "I'm fine.", "It's red.", "Yes, I do."}, "An English book.", "问书包里有什么，回答物品。")
	add(4, "Schoolbag", "「一包 candy」的英语是？", []string{"a bag of candy", "a candy bag", "some candy", "candy bag"}, "a bag of candy", "a bag of candy 一包糖果。")
	add(4, "Friends", "He has glasses and his hair is brown. 他长什么样？", []string{"戴眼镜棕头发", "很高很瘦", "黑头发", "戴帽子"}, "戴眼镜棕头发", "glasses 眼镜，brown hair 棕色头发。")
	add(4, "Friends", "My friend is tall and strong. 意思是？", []string{"我的朋友又高又壮。", "我的朋友很矮。", "我的朋友很瘦。", "我的朋友很可爱。"}, "我的朋友又高又壮。", "tall and strong 又高又壮。")
	add(4, "Home", "Is she in the kitchen? 否定回答是？", []string{"No, she isn't.", "Yes, she is.", "No, he isn't.", "Yes, it is."}, "No, she isn't.", "Is she...？否定回答 No, she isn't。")
	add(4, "Home", "Go to the living room. ____ TV.", []string{"Watch", "Read", "Play", "Listen"}, "Watch", "Watch TV 看电视，在客厅进行。")
	add(4, "Dinner", "What would you like for dinner? 回答？", []string{"I'd like some rice.", "I'm fine.", "I like apples.", "Yes, please."}, "I'd like some rice.", "What would you like? 回答 I'd like...")
	add(4, "Dinner", "牛肉 beef 用刀叉吃，放在____上？", []string{"plate", "bowl", "cup", "glass"}, "plate", "牛肉放在盘子里。")

	// ======== 五年级（PEP 上册）=======
	add(5, "Teachers", "Who is your English teacher? 的正确回答？", []string{"Miss White.", "She is tall.", "Yes, she is.", "No, she isn't."}, "Miss White.", "Who is...？回答人名。")
	add(5, "Teachers", "What's she like? 是问什么？", []string{"她长什么样？", "她是谁？", "她喜欢什么？", "她几岁了？"}, "她长什么样？", "What's...like？问外表或性格。")
	add(5, "Week", "What day is it today? 回答？", []string{"It's Monday.", "It's 8 o'clock.", "It's sunny.", "It's May 1st."}, "It's Monday.", "What day is it？问星期几。")
	add(5, "Week", "We have English ____ Mondays.", []string{"on", "in", "at", "for"}, "on", "星期几用介词 on。")
	add(5, "Food", "What would you like to eat? 回答？", []string{"I'd like a sandwich.", "I like cats.", "Yes, I do.", "No, thanks."}, "I'd like a sandwich.", "What would you like to eat? 回答想吃的食物。")
	add(5, "Food", "What's your favourite food? 的正确回答？", []string{"Ice cream.", "I like red.", "I'm ten.", "It's a book."}, "Ice cream.", "favourite food 问最喜欢的食物。")
	add(5, "Housework", "Can you wash the dishes? 肯定回答？", []string{"Yes, I can.", "Yes, I do.", "Yes, I am.", "Yes, I will."}, "Yes, I can.", "Can you...？回答 Yes, I can。")
	add(5, "Housework", "I can ____ the floor.（扫地）", []string{"clean", "sweep", "wash", "cook"}, "sweep", "sweep the floor 扫地。")
	add(5, "Nature", "There is a ____ in the park.（湖）", []string{"lake", "river", "mountain", "forest"}, "lake", "lake 是湖。")
	add(5, "Nature", "Are there any birds? 否定回答？", []string{"No, there aren't.", "No, there isn't.", "Yes, there are.", "Yes, there is."}, "No, there aren't.", "复数 birds 用 there aren't。")

	// ======== 六年级（PEP 上册）=======
	add(6, "Transportation", "How do you go to school? 回答？", []string{"By bus.", "I like school.", "On foot.", "Yes, I do."}, "By bus.", "How do you go...？询问交通方式。")
	add(6, "Transportation", "I go to school ____ foot.", []string{"on", "by", "in", "with"}, "on", "on foot 步行。")
	add(6, "Directions", "Where is the post office? 回答？", []string{"It's near the hospital.", "Yes, it is.", "No, it isn't.", "It's a post office."}, "It's near the hospital.", "Where is...？回答位置。")
	add(6, "Directions", "Turn left at the bookstore. 意思是？", []string{"在书店左转。", "在书店右转。", "直走。", "过马路。"}, "在书店左转。", "Turn left 左转。")
	add(6, "Hobbies", "What are your hobbies? 回答？", []string{"I like reading.", "I'm reading.", "Yes, I do.", "No, I don't."}, "I like reading.", "hobbies 爱好，回答喜欢的活动。")
	add(6, "Hobbies", "He likes ____ football.", []string{"playing", "play", "plays", "played"}, "playing", "like + doing 喜欢做某事。")
	add(6, "Pen pal", "Does she live in Sydney? 肯定回答？", []string{"Yes, she does.", "Yes, she is.", "Yes, she can.", "Yes, she has."}, "Yes, she does.", "Does 开头的一般疑问句，回答用 does。")
	add(6, "Pen pal", "My pen pal lives in ____（澳大利亚）.", []string{"Australia", "China", "Canada", "England"}, "Australia", "Australia 澳大利亚。")
	add(6, "Jobs", "What does your father do? 回答？", []string{"He is a doctor.", "He is tall.", "He likes reading.", "He is reading."}, "He is a doctor.", "What does...do？问职业。")
	add(6, "Jobs", "She works in a ____（医院）.", []string{"hospital", "school", "factory", "library"}, "hospital", "hospital 医院。")

	return qs
}

