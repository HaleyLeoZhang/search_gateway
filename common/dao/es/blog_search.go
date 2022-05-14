package es

import (
	"context"
	"encoding/json"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	es7Sdk "github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"search_gateway/common/model/es"
	"search_gateway/common/model/vo"
)

// Part 博文搜索
func (d *Dao) BlogFrontSearch(ctx context.Context, req *vo.BlogFrontRequest) (list []*es.Blog, total int64, err error) {
	var (
		baseDoc = &es.Blog{}
	)
	// 拉取数据
	// - 设置需要的字段
	fields := []string{"id", "title^2", "describe^2", "category"}
	search := es7Sdk.NewBoolQuery()
	if req.Id != nil {
		search.Must(es7Sdk.NewTermsQuery("id", *req.Id))
	}
	if req.Title != "" {
		search.Should(es7Sdk.NewMatchQuery("title", req.Title))
	}
	if req.Describe != "" {
		search.Should(es7Sdk.NewMatchQuery("describe", req.Describe))
	}
	if req.Category != "" {
		search.Should(es7Sdk.NewMatchQuery("category", req.Category))
	}
	// - 计算分页
	offset := req.GetOffset()
	limit := *req.PageSize

	result, err := d.clt.Search().Index(baseDoc.GetIndex()).
		From(offset).Size(limit). // 取数据区间
		Query(search).
		FetchSourceContext(es7Sdk.NewFetchSourceContext(true).Include(fields...)).
		Do(context.Background())
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	// 解析数据
	list, total = blogFrontSearchHandler(result)
	return
}

// 处理返回结果
func blogFrontSearchHandler(result *es7Sdk.SearchResult) (list []*es.Blog, totalInt64 int64) {
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
