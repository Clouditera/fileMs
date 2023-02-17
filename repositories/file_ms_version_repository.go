package repositories

import (
	"github.com/gannicus-w/yunqi_mysql/sqls"
	"github.com/gannicus-w/yunqi_mysql/web/params"
	"gorm.io/gorm"

	"fileMS/model"
)

var FileMsVersionRepository = newFileMsVersionRepository()

func newFileMsVersionRepository() *fileMsVersionRepository {
	return &fileMsVersionRepository{}
}

type fileMsVersionRepository struct {
}

func (r *fileMsVersionRepository) Get(db *gorm.DB, id int64) *model.FileVersion {
	ret := &model.FileVersion{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *fileMsVersionRepository) GetByName(db *gorm.DB, name string) *model.FileVersion {
	ret := &model.FileVersion{}
	if err := db.First(ret, "name = ?", name).Error; err != nil {
		return nil
	}
	return ret
}

func (r *fileMsVersionRepository) Take(db *gorm.DB, where ...interface{}) *model.FileVersion {
	ret := &model.FileVersion{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *fileMsVersionRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.FileVersion) {
	cnd.Find(db, &list)
	return
}

func (r *fileMsVersionRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.FileVersion {
	ret := &model.FileVersion{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *fileMsVersionRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.FileVersion, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *fileMsVersionRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.FileVersion, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.FileVersion{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *fileMsVersionRepository) Create(db *gorm.DB, t *model.FileVersion) (err error) {
	err = db.Create(t).Error
	return
}

func (r *fileMsVersionRepository) Update(db *gorm.DB, t *model.FileVersion) (err error) {
	err = db.Save(t).Error
	return
}

func (r *fileMsVersionRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.FileChunk{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *fileMsVersionRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.FileVersion{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *fileMsVersionRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.FileVersion{}).Delete("id", id)
}
