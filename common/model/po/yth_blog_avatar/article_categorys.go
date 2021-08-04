package yth_blog_avatar

import "time"

type ArticleCategory struct {
	Id        int64     `gorm:"column:id"`
	Title     string    `gorm:"column:title"`
	IsDeleted int64     `gorm:"column:is_deleted"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (ArticleCategory) TableName() string {
	return "article_categorys"
}
