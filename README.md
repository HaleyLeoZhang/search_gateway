## 搜索网关

云天河网站内的搜索服务

~~~bash
curl --location --request GET '服务地址/blog/front?title=安全&describe=消息&category=golang'
~~~

> ES mapping结构

[点此查看](common/model/es/blog.go) 

> 初始化数据

`Step 1` 推送全量博文ID

~~~bash
svr.blogSearchEsSendAll(ctx)
~~~

`Step 2` 消费端写入或者删除 `ES` 数据  
`job` 消费端代码 [点此查看](job/service/consumer.go)  




