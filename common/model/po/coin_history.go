package po

// ----------------------------------------------------------------------
// 漫画基础信息模型
// ----------------------------------------------------------------------
// Link  : http://www.hlzblog.top/
// GITHUB: https://github.com/HaleyLeoZhang
// ----------------------------------------------------------------------

type CoinHistory struct {
	Id        int64   `gorm:"column:id;" json:"id"`
	CoinId    int     `gorm:"column:coin_id;default:'1'" json:"coinId"`
	Price     float64 `gorm:"column:price;default:'0.00'" json:"price"`
	ChannelId int     `gorm:"column:channel_id;default:'0'" json:"channelId"`
	Status    int     `gorm:"column:status;default:'200'" json:"status"`
	CreatedAt string  `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt string  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

//数据表---必需
func (c *CoinHistory) TableName() string {
	return "coin_history"
}
