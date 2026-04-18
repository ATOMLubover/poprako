package val

// UnitInfo 表示用于在应用层和接口层传递的翻译单元信息 VO
type UnitInfo struct {
	// ID 是翻译单元的唯一标识
	ID string `json:"id"`

	// PageID 是所属页面 ID
	PageID string `json:"page_id"`
	// Index 是翻译单元在页面中的序号
	Index int `json:"index"`

	// XCoord 是翻译单元的 X 坐标
	XCoord int `json:"x_coord"`
	// YCoord 是翻译单元的 Y 坐标
	YCoord int `json:"y_coord"`

	// IsBubble 表示该单元是否是气泡框
	IsBubble bool `json:"is_bubble"`

	// TranslatedText 是翻译后的文本
	TranslatedText *string `json:"translated_text"`
	// TranslatorID 是翻译者的用户 ID
	TranslatorID *string `json:"translator_id"`
	// TranslatorComment 是翻译者的备注
	TranslatorComment *string `json:"translator_comment"`

	// IsProofread 表示该单元是否已经校对
	IsProofread bool `json:"is_proofread"`
	// ProofreadText 是校对后的文本
	ProofreadText *string `json:"proofread_text"`
	// ProofreaderID 是校对者的用户 ID
	ProofreaderID *string `json:"proofreader_id"`
	// ProofreaderComment 是校对者的备注
	ProofreaderComment *string `json:"proofreader_comment"`
}

// SavePageUnitArgs 表示保存页面翻译单元的参数（diff 语义）
type SavePageUnitArgs struct {
	// PageID 是目标页面 ID
	PageID string `json:"page_id" validate:"required"`
	// UnitDiff 是翻译单元的变更信息
	UnitDiff UnitDiff `json:"unit_diff" validate:"required"`
}

// UnitDiff 表示翻译单元的变更信息
type UnitDiff struct {
	// Insert 是要新增的翻译单元列表
	Insert []UnitCreation `json:"insert"`
	// Patch 是要修改的翻译单元列表
	Patch []UnitPatch `json:"patch"`
	// Delete 是要删除的翻译单元 ID 列表
	Delete []string `json:"delete"`
}

// UnitCreation 表示新增翻译单元的数据
type UnitCreation struct {
	// ID 是翻译单元的唯一标识
	ID string `json:"id" validate:"required"`
	// Index 是翻译单元在页面中的序号
	Index int `json:"index"`

	// XCoord 是翻译单元的 X 坐标
	XCoord int `json:"x_coord"`
	// YCoord 是翻译单元的 Y 坐标
	YCoord int `json:"y_coord"`

	// IsBubble 表示该单元是否是气泡框
	IsBubble bool `json:"is_bubble"`

	// TranslatedText 是翻译后的文本
	TranslatedText *string `json:"translated_text"`
	// TranslatorID 是翻译者的用户 ID
	TranslatorID *string `json:"translator_id"`
	// TranslatorComment 是翻译者的备注
	TranslatorComment *string `json:"translator_comment"`

	// IsProofread 表示该单元是否已经校对
	IsProofread bool `json:"is_proofread"`
	// ProofreadText 是校对后的文本
	ProofreadText *string `json:"proofread_text"`
	// ProofreaderID 是校对者的用户 ID
	ProofreaderID *string `json:"proofreader_id"`
	// ProofreaderComment 是校对者的备注
	ProofreaderComment *string `json:"proofreader_comment"`
}

// UnitPatch 表示修改翻译单元的数据（PATCH 语义）
type UnitPatch struct {
	// ID 是要修改的翻译单元标识
	ID string `json:"id" validate:"required"`

	// Index 是翻译单元在页面中的序号
	Index *int `json:"index"`

	// XCoord 是翻译单元的 X 坐标
	XCoord *int `json:"x_coord"`
	// YCoord 是翻译单元的 Y 坐标
	YCoord *int `json:"y_coord"`

	// IsBubble 表示该单元是否是气泡框
	IsBubble *bool `json:"is_bubble"`

	// TranslatedText 是翻译后的文本
	TranslatedText **string `json:"translated_text"`
	// TranslatorID 是翻译者的用户 ID
	TranslatorID **string `json:"translator_id"`
	// TranslatorComment 是翻译者的备注
	TranslatorComment **string `json:"translator_comment"`

	// IsProofread 表示该单元是否已经校对
	IsProofread *bool `json:"is_proofread"`
	// ProofreadText 是校对后的文本
	ProofreadText **string `json:"proofread_text"`
	// ProofreaderID 是校对者的用户 ID
	ProofreaderID **string `json:"proofreader_id"`
	// ProofreaderComment 是校对者的备注
	ProofreaderComment **string `json:"proofreader_comment"`
}
