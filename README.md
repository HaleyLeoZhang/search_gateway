## 搜索网关

云天河网站内的搜索服务

~~~bash
curl --location --request GET '服务地址/blog/front?title=安全&describe=消息&category=golang'
~~~

初始化ES数据--Blog  
~~~bash
_ = svr.BlogFrontIni()
svr.BlogFrontFlushAll()
~~~


