package matchString

// BuildTyp 构建入参类型
type BuildTyp interface {
	~string | ~[]string | ~[][]string
}

// ScanResTyp 匹配响应结果类型
type ScanResTyp interface {
	~int | ~[]string | ~[][]string
}

type Collision[BT BuildTyp, SRT ScanResTyp] interface {
	Build(BT) error  // 预构建结构
	Scan(string) SRT // 执行匹配响应
}
