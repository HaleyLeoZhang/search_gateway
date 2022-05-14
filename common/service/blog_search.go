package service

import (
	"context"
	"encoding/json"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"github.com/HaleyLeoZhang/go-component/errgroup"
	"github.com/Shopify/sarama"
	"search_gateway/common/constant"
	"search_gateway/common/model/es"
	"search_gateway/common/model/kafka"
	po "search_gateway/common/model/po/yth_blog_avatar"
	"search_gateway/common/model/vo"
)

// 博客主搜

// Part 刷全量数据脚本 - 后面有时间了，改成完整流程
// - 推送kafka消息去统一更新 TODO
//  ---------------- 搜索 -----------------
// Part 搜素
func (s *Service) BlogFrontSearch(ctx context.Context, req *vo.BlogFrontRequest) (res *vo.BlogFrontResponse, err error) {
	res = &vo.BlogFrontResponse{
		Ids: make([]int, 0),
	}
	var (
		list  []*es.Blog
		total int64
	)
	list, total, err = s.EsDao.BlogFrontSearch(ctx, req)
	if err != nil {
		return
	}
	// 循环包装
	res.Ids = make([]int, 0, len(list))
	res.Total = total
	for _, item := range list {
		res.Ids = append(res.Ids, int(item.Id))
	}
	return
}

//  ---------------- ES 处理 -----------------

//  消费端 处理ES数据
func (s *Service) KafkaBlogSearchEs(ctx context.Context, message *sarama.ConsumerMessage) (err error) {
	var (
		doc *es.Blog
		msg = new(kafka.BlogSearch)
	)
	err = json.Unmarshal(message.Value, msg)
	doc, err = s.blogSearchAssemble(ctx, msg.Id)
	if err == nil && doc == nil {
		msg.Action = constant.ES_DELETE
	}
	switch msg.Action {
	case constant.ES_CREATE, constant.ES_UPDATE:
		err = s.EsDao.UpsertEs(ctx, doc)
	case constant.ES_DELETE:
		err = s.EsDao.DeleteEs(ctx, &es.Blog{Id: msg.Id})
	}
	if err != nil {
		xlog.Errorf("KafkaBlogSearch Err(%+v) msg(%v)", err, msg)
		return
	}
	xlog.Infof("KafkaBlogSearch Success msg(%v)", msg)
	return
}

// ES组装数据 - 组装数据
func (s *Service) blogSearchAssemble(ctx context.Context, id int64) (doc *es.Blog, err error) {
	var (
		article         *po.Article
		articleCategory *po.ArticleCategory
	)
	article, err = s.BlogDao.ArticleById(ctx, id)
	if err != nil {
		return
	}
	// 数据不存在；任务已完成；隐藏  都不进ES
	if article == nil || article.IsDeleted == constant.IS_DELETED_YES { // 兼职的创建商品的，才用于搜索
		return
	}
	// 组装数据
	eg := &errgroup.Group{}
	eg.GOMAXPROCS(3)
	// - 获取类目信息
	eg.Go(func(context.Context) (errNil error) {
		articleCategory, err = s.BlogDao.ArticleCategoryById(ctx, article.CateId)
		if err != nil {
			return
		}
		if article == nil {
			articleCategory = &po.ArticleCategory{}
		}
		return
	})
	_ = eg.Wait()

	// 填充
	doc = &es.Blog{
		Id:       article.Id,
		Title:    article.Title,
		Describe: article.Descript,
		Category: articleCategory.Title,
	}
	return
}

// 定时全量有效数据刷一遍ES - 推送数据
func (s *Service) blogSearchEsSendAll(ctx context.Context) {
	var (
		minId       = int64(0)
		batchNumber = int64(constant.BUSINESS_FETCH)
	)
	xlog.Infof("blogSearchEsSendAll START")
	for {
		// 查询
		list, err := s.BlogDao.ArticleListByMinIdWithField(ctx, minId, batchNumber, "id")
		if err != nil {
			xlog.Infof("blogSearchEsSendAll minId(%v) Error(%+v)", minId, err)
			return
		}
		lenList := len(list)
		if lenList == 0 {
			xlog.Infof("blogSearchEsSendAll END minId(%v)", minId)
			return
		}
		for _, item := range list {
			// 通知任务ES更新 PV、UV等数据
			_ = s.producer.NotifyBlogSearch(ctx, item.Id, constant.ES_UPDATE, "blogSearchEsSendAll")
		}
		minId = list[len(list)-1].Id
	}
}
