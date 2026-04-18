package model

// UnitInfo 表示一个页面中的翻译单元（如一个气泡框或注释）
type UnitInfo struct {
	ID string

	PageID string
	Index  int

	XCoord int
	YCoord int

	IsBubble bool

	TranslatedText    *string
	TranslatorID      *string
	TranslatorComment *string

	IsProofread        bool
	ProofreadText      *string
	ProofreaderID      *string
	ProofreaderComment *string
}

// UnitPatch 是 PATCH 语义的载荷，仅修改非 nil 的字段
// 对于指针类型的 Option 字段：nil 表示不修改；非 nil 表示必须修改（其值可为 nil 以清空）
type UnitPatch struct {
	// ID 是必须的，标识要修改的单元
	ID string

	Index *int

	XCoord *int
	YCoord *int

	IsBubble *bool

	TranslatedText    **string
	TranslatorID      **string
	TranslatorComment **string

	IsProofread        *bool
	ProofreadText      **string
	ProofreaderID      **string
	ProofreaderComment **string
}

// UnitCreation 与 UnitInfo 语义相同，创建时直接复用
type UnitCreation = UnitInfo

// UnitDiff 表示一次对页面翻译单元的修改，包括新增、修改和删除
type UnitDiff struct {
	Insert []UnitCreation // 新增的单元列表
	Patch  []UnitPatch    // 修改的单元列表
	Delete []string       // 删除的单元 ID 列表
}

// UnitQueryOpt 指定翻译单元查询的可选筛选条件
// 所有字段均为可空，nil 表示不参与筛选
type UnitQueryOpt struct {
	// PageID 按所属页面 ID 筛选
	PageID string
}
