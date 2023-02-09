package repositories

import (
	"github.com/gannicus-w/yunqi_mysql/sqls"
	"github.com/gannicus-w/yunqi_mysql/web/params"
	"gorm.io/gorm"

	"fileMS/model"
)

var FileMsRepository = newFileMsRepository()

func newFileMsRepository() *fileMsRepository {
	return &fileMsRepository{}
}

type fileMsRepository struct {
}

func (r *fileMsRepository) Get(db *gorm.DB, id int64) *model.FileChunk {
	ret := &model.FileChunk{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *fileMsRepository) GetByName(db *gorm.DB, name string) *model.FileChunk {
	ret := &model.FileChunk{}
	if err := db.First(ret, "name = ?", name).Error; err != nil {
		return nil
	}
	return ret
}

func (r *fileMsRepository) Take(db *gorm.DB, where ...interface{}) *model.FileChunk {
	ret := &model.FileChunk{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *fileMsRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.FileChunk) {
	cnd.Find(db, &list)
	return
}

func (r *fileMsRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.FileChunk {
	ret := &model.FileChunk{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *fileMsRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.FileChunk, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *fileMsRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.FileChunk, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.FileChunk{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *fileMsRepository) Create(db *gorm.DB, t *model.FileChunk) (err error) {
	err = db.Create(t).Error
	return
}

func (r *fileMsRepository) Update(db *gorm.DB, t *model.FileChunk) (err error) {
	err = db.Save(t).Error
	return
}

func (r *fileMsRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.FileChunk{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *fileMsRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.FileChunk{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *fileMsRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.FileChunk{}).Delete("id", id)
}
