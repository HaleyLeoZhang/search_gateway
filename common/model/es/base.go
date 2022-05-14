package es

type BaseDoc interface {
	GetIndex() string    // 获取索引名 - 别名
	GetIdString() string // 获取文档ID
}
