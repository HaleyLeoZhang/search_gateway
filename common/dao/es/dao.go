package es

import (
	"context"
	"github.com/HaleyLeoZhang/go-component/driver/xelastic"
	v7 "github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"search_gateway/common/model/es"
)

type Dao struct {
	clt *v7.Client
}

func New(cfg *xelastic.Config) *Dao {
	var err error

	d := &Dao{}
	d.clt, err = xelastic.NewV7(cfg)
	if err != nil {
		panic(err)
	}
	return d
}

// 基础公共方法封装

func (d *Dao) UpsertEs(ctx context.Context, doc es.BaseDoc) (err error) {
	_, err = d.clt.Update().Index(doc.GetIndex()).Id(doc.GetIdString()).Doc(doc).DocAsUpsert(true).Do(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

func (d *Dao) DeleteEs(ctx context.Context, doc es.BaseDoc) (err error) {
	_, err = d.clt.Delete().Index(doc.GetIndex()).Id(doc.GetIdString()).Do(ctx)
	if err != nil {
		if v7.IsNotFound(err) || v7.IsConflict(err) {
			err = nil
		} else {
			return errors.WithStack(err)
		}
		return
	}

	return
}
