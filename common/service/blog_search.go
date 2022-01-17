package service

import (
	"context"
	"encoding/json"
	"fmt"
	dbTool "github.com/HaleyLeoZhang/go-component/driver/db"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"github.com/HaleyLeoZhang/go-component/errgroup"
	esSdk "github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"search_gateway/common/constant"
	"search_gateway/common/model/po/es"
	po "search_gateway/common/model/po/yth_blog_avatar"
	"search_gateway/common/model/vo"
)

// 博客主搜

// Part 刷全量数据脚本 - 后面有时间了，改成完整流程
// - 推送kafka消息去统一更新 TODO
func (s *Service) BlogFrontFlushAll() {
	var (
		err error
		ctx = context.Background()

		list         []*po.Article
		listCategory []*po.ArticleCategory

		minId = 0 //  起始 id
	)
	for {
		// Part 1  主表
		cond := &dbTool.DBConditions{}
		cond.And = make(map[string]interface{})
		cond.And["is_deleted = ?"] = constant.BASE_TABLE_DELETED_NO
		cond.And["id > ?"] = minId
		cond.Order = "id ASC"
		cond.Limit = constant.BUSINESS_FETCH // 每次拉取的数量
		list, err = s.BlogDao.ArticleList(ctx, cond)
		if err != nil {
			xlog.Errorf("BlogFrontFlushAll Err(%+v) ArticleList", err)
			return
		}
		lenList := len(list) // 考虑性能，后续对应关联该表达切片一开始就分配当前未当前的数量
		if lenList == 0 {
			xlog.Infof("BlogFrontFlushAll lastId(%v) done", minId)
			return
		}
		// Part 2  分类ID - 虽然数据库会自动去重，但是考虑后期每次取数量变化的问题，考虑本地手动去重
		mapCategory := make(map[int64]struct{})
		for _, one := range list {
			mapCategory[one.CateId] = struct{}{}
		}
		realCategoryIds := make([]int64, 0, len(mapCategory))
		for cateId, _ := range mapCategory {
			realCategoryIds = append(realCategoryIds, cateId)
		}
		condCategory := &dbTool.DBConditions{}
		condCategory.And = make(map[string]interface{})
		condCategory.And["id in (?)"] = realCategoryIds
		condCategory.Order = "id ASC"
		condCategory.Limit = constant.BUSINESS_FETCH // 每次拉取的数量
		listCategory, err = s.BlogDao.ArticleCategoryList(ctx, condCategory)
		if err != nil {
			xlog.Errorf("BlogFrontFlushAll Err(%+v) ArticleCategoryList", err)
			return
		}
		// Part 3  分类Map --- 仅供并发读
		mapCategoryItem := make(map[int64]*po.ArticleCategory)
		for _, one := range listCategory {
			mapCategoryItem[one.Id] = one
		}
		// Part 4  准备操作的数据
		// --- 并发写入
		eg := errgroup.Group{}
		eg.GOMAXPROCS(constant.BUSINESS_BLOG_SEARCH_COROUTINE)
		for _, tmpOne := range list {
			one := tmpOne
			eg.Go(func(context.Context) error {
				doc := &es.Blog{}
				doc.Id = int(one.Id)
				if one.IsDeleted != constant.BASE_TABLE_DELETED_NO && one.IsOnline != constant.BASE_TABLE_ONLINE {
					_, _ = s.EsDao.Delete().Index(doc.GetIndex()).Id(doc.GetIdString()).Do(ctx)
					return nil
				}
				doc.Title = one.Title
				doc.Describe = one.Descript
				category, ok := mapCategoryItem[one.CateId]
				if ok {
					doc.Category = category.Title
				}
				_, err := s.EsDao.Update().Index(doc.GetIndex()).Id(doc.GetIdString()).Doc(doc).DocAsUpsert(true).Do(ctx)
				if err != nil {
					xlog.Errorf("BlogFrontFlushAll Err(%+v) id(%v)", err, doc.Id)
					return nil
				}
				xlog.Infof("BlogFrontFlushAll Success(%+v) id(%v)", err, doc.Id)
				return nil
			})
		}
		_ = eg.Wait()
		// Part 计算下次的开头ID
		minId = int(list[lenList-1].Id)
		xlog.Infof("BlogFrontFlushAll lastId(%v) doing", minId)
	}
}

// Part 搜素
func (s *Service) BlogFrontSearch(ctx context.Context, req *vo.BlogFrontRequest) (res *vo.BlogFrontResponse, err error) {
	var (
		baseDoc = &es.Blog{}

		list  []*es.Blog
		total int64
	)
	res = &vo.BlogFrontResponse{}
	res.Ids = make([]int, 0)
	// 拉取数据
	// - 设置需要的字段
	fields := []string{"id", "title^2", "describe^2", "category"}
	search := esSdk.NewBoolQuery()
	if req.Id != nil {
		search.Must(esSdk.NewTermsQuery("id", *req.Id))
	}
	if req.Title != "" {
		search.Should(esSdk.NewMatchQuery("title", req.Title))
	}
	if req.Describe != "" {
		search.Should(esSdk.NewMatchQuery("describe", req.Describe))
	}
	if req.Category != "" {
		search.Should(esSdk.NewMatchQuery("category", req.Category))
	}
	// - 计算分页
	offset := req.GetOffset()
	limit := *req.PageSize

	result, err := s.EsDao.Search().Index(baseDoc.GetIndex()).
		From(offset).Size(limit). // 取数据区间
		Query(search).
		FetchSourceContext(esSdk.NewFetchSourceContext(true).Include(fields...)).
		Do(context.Background())
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	// 解析数据
	list, total = handleSearchData(result)
	// 循环包装
	res.Ids = make([]int, 0, len(list))
	res.Total = total
	for _, item := range list {
		res.Ids = append(res.Ids, item.Id)
	}
	return
}

// 处理返回结果
func handleSearchData(result *esSdk.SearchResult) (list []*es.Blog, totalInt64 int64) {
	list = make([]*es.Blog, 0)

	totalInt64 = result.TotalHits()
	if totalInt64 == 0 {
		return
	}
	for _, hit := range result.Hits.Hits {
		d := &es.Blog{}
		buf, _ := hit.Source.MarshalJSON()
		err := json.Unmarshal(buf, &d)
		if err != nil {
			err = errors.WithStack(err)
			xlog.Warnf("Warning value(%+v) err(%+v)", string(buf), err)
			continue
		}
		// 固定解析数据类型返回
		list = append(list, d)
	}

	return
}

// Part 初始化
func (s *Service) BlogFrontIni() (err error) {
	var (
		ctx = context.Background()
		doc = &es.Blog{}
	)
	b, err := s.EsDao.IndexExists(doc.GetIndex()).Do(ctx)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if !b {
		createIndex, errBusiness := s.EsDao.CreateIndex(doc.GetIndexOriginal()).Body(doc.GetMapping()).Do(ctx)
		if errBusiness != nil {
			err = errors.WithStack(fmt.Errorf("NewClient err(%+v)", errBusiness))
			return
		}
		if createIndex == nil {
			err = errors.WithStack(fmt.Errorf("Expected result to be != nil"))
			return
		}
	}
	return
}
