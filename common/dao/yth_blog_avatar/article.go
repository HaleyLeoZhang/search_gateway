package yth_blog_avatar

import (
	"context"
	"fmt"
	dbTool "github.com/HaleyLeoZhang/go-component/driver/db"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"search_gateway/common/constant"
	po "search_gateway/common/model/po/yth_blog_avatar"
)

// 通过ID查询
func (d *Dao) ArticleById(ctx context.Context, id int64) (res *po.Article, err error) {
	res = &po.Article{}

	err = d.db.Table(res.TableName()).Where("id = ?", id).Take(&res).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return
}

// 条件查询
func (d *Dao) articleList(ctx context.Context, conditions *dbTool.DBConditions) (res []*po.Article, err error) {
	oneModel := &po.Article{}

	db := d.db.Table(oneModel.TableName())
	db = conditions.Fill(db)
	err = db.Find(&res).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// 查询获取单条
func (d *Dao) ArticleGetOne(ctx context.Context, where map[string]interface{}) (res *po.Article, err error) {
	res = &po.Article{}

	err = d.db.Table(res.TableName()).Where(where).First(&res).Error
	if gorm.IsRecordNotFoundError(err) {
		res = nil
		return
	}
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// 更新
func (d *Dao) ArticleUpdate(ctx context.Context, where map[string]interface{}, update map[string]interface{}) (affected int64, err error) {
	oneModel := &po.Article{}

	db := d.db.Table(oneModel.TableName()).Where(where).Updates(update)
	if db.Error != nil {
		err = errors.WithStack(db.Error)
		return
	}

	affected = db.RowsAffected
	return
}

// 条件更新
func (d *Dao) ArticleUpdateUseCondition(ctx context.Context, conditions dbTool.DBConditions, update map[string]interface{}) (affected int64, err error) {
	oneModel := &po.Article{}

	db := d.db.Table(oneModel.TableName())
	db = conditions.Fill(db).Updates(update)
	if db.Error != nil {
		err = errors.WithStack(db.Error)
		return
	}

	affected = db.RowsAffected
	return
}

// 通过表ID，查询数据，并限制返回的字段
func (d *Dao) ArticleListByMinIdWithField(ctx context.Context, minId, limit int64, fields string) (list []*po.Article, err error) {
	if fields == "" {
		err = errors.WithStack(fmt.Errorf("fields不能为空"))
		return
	}

	// Part 1  主表
	cond := &dbTool.DBConditions{}
	cond.Select = fields
	cond.And = make(map[string]interface{})
	cond.And["is_deleted = ?"] = constant.BASE_TABLE_DELETED_NO
	cond.And["id > ?"] = minId
	cond.Order = "id ASC"
	cond.Limit = limit
	list, err = d.articleList(ctx, cond)
	if err != nil {
		return
	}
	return
}
