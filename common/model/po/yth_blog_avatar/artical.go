package yth_blog_avatar

import "time"

type Article struct {
	Id         int64     `gorm:"column:id"`
	Title      string    `gorm:"column:title"`
	Type       int64     `gorm:"column:type"`
	Sticky     int64     `gorm:"column:sticky"`
	Sequence   int64     `gorm:"column:sequence"`
	Original   int64     `gorm:"column:original"`
	IsOnline   int64     `gorm:"column:is_online"`
	Content    string    `gorm:"column:content"`
	Descript   string    `gorm:"column:descript"`
	CoverUrl   string    `gorm:"column:cover_url"`
	CateId     int64     `gorm:"column:cate_id"`
	BgId       int64     `gorm:"column:bg_id"`
	RawContent string    `gorm:"column:raw_content"`
	Statistic  int64     `gorm:"column:statistic"`
	IsDeleted  int64     `gorm:"column:is_deleted"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (Article) TableName() string {
	return "articles"
}
