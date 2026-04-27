package entity

type SysMailRow struct {
	Id string `gorm:"column:id;primaryKey"`

	RcvId   string `gorm:"column:receiver_id"`
	Title   string `gorm:"column:title"`
	Content string `gorm:"column:content"`

	Read bool `gorm:"column:read"`

	CreatedAt string `gorm:"column:created_at"`
}

type SysMailCreRow struct {
	Id string `gorm:"column:id;primaryKey"`

	RcvId   string `gorm:"column:receiver_id"`
	Title   string `gorm:"column:title"`
	Content string `gorm:"column:content"`

	CreatedAt string `gorm:"column:created_at"`
}
