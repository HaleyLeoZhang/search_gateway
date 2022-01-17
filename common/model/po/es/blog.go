package es

import "fmt"

type Blog struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Describe string `json:"describe"`
	Category string `json:"category"`
}

// 索引 - 生产一定要用别名，方便故障时切换 - 生产上用这个进行CURD
func (Blog) GetIndex() string {
	return "blog_front_search"
}

// 原始索引
func (Blog) GetIndexOriginal() string {
	return "blog_front_search_v1"
}

// 方便查询/删除 用格式化后的ID
func (b *Blog) GetIdString() string {
	return fmt.Sprintf("%v", b.Id)
}

// 初始化mapping
func (b *Blog) GetMapping() (realMapping string) {
	return `
{
    "settings":{
        "number_of_shards":1,
        "number_of_replicas":0
    },
    "aliases":{
        "blog_front_search":{

        }
    },
    "mappings":{
        "properties":{
            "id":{
                "type":"integer"
            },
            "title":{
                "type":"text",
                "analyzer":"ik_smart",
                "search_analyzer":"ik_smart"
            },
            "describe":{
                "type":"text",
                "analyzer":"ik_smart",
                "search_analyzer":"ik_smart"
            },
            "category":{
                "type":"text",
                "analyzer":"ik_smart",
                "search_analyzer":"ik_smart"
            }
        }
    }
}
	`
}
