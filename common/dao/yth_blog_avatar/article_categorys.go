package yth_blog_avatar

import (
	"context"
	dbTool "github.com/HaleyLeoZhang/go-component/driver/db"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	po "search_gateway/common/model/po/yth_blog_avatar"
)

// 通过ID查询
func (d *Dao) ArticleCategoryById(ctx context.Context, id int64) (res *po.ArticleCategory, err error) {
	res = &po.ArticleCategory{}

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
func (d *Dao) ArticleCategoryList(ctx context.Context, conditions *dbTool.DBConditions) (res []*po.ArticleCategory, err error) {
	oneModel := &po.ArticleCategory{}

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
func (d *Dao) ArticleCategoryGetOne(ctx context.Context, where map[string]interface{}) (res *po.ArticleCategory, err error) {
	res = &po.ArticleCategory{}

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
func (d *Dao) ArticleCategoryUpdate(ctx context.Context, where map[string]interface{}, update map[string]interface{}) (affected int64, err error) {
	oneModel := &po.ArticleCategory{}

	db := d.db.Table(oneModel.TableName()).Where(where).Updates(update)
	if db.Error != nil {
		err = errors.WithStack(db.Error)
		return
	}

	affected = db.RowsAffected
	return
}

// 条件更新
func (d *Dao) ArticleCategoryUpdateUseCondition(ctx context.Context, conditions dbTool.DBConditions, update map[string]interface{}) (affected int64, err error) {
	oneModel := &po.ArticleCategory{}

	db := d.db.Table(oneModel.TableName())
	db = conditions.Fill(db).Updates(update)
	if db.Error != nil {
		err = errors.WithStack(db.Error)
		return
	}

	affected = db.RowsAffected
	return
}
