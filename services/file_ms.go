package services

import (
	"fileMS/model"
	"fileMS/repositories"
	"github.com/gannicus-w/yunqi_mysql/sqls"
	"github.com/gannicus-w/yunqi_mysql/web/params"
)

var FileMsService = newFileMsService()

func newFileMsService() *fileMsService {
	return &fileMsService{}
}

type fileMsService struct {
}

func (s *fileMsService) Get(id int64) *model.FileChunk {
	return repositories.FileMsRepository.Get(sqls.DB(), id)
}

func (s *fileMsService) Take(where ...interface{}) *model.FileChunk {
	return repositories.FileMsRepository.Take(sqls.DB(), where...)
}

func (s *fileMsService) Find(cnd *sqls.Cnd) []model.FileChunk {
	return repositories.FileMsRepository.Find(sqls.DB(), cnd)
}

func (s *fileMsService) FindOne(cnd *sqls.Cnd) *model.FileChunk {
	return repositories.FileMsRepository.FindOne(sqls.DB(), cnd)
}

func (s *fileMsService) FindPageByParams(params *params.QueryParams) (list []model.FileChunk, paging *sqls.Paging) {
	return repositories.FileMsRepository.FindPageByParams(sqls.DB(), params)
}

func (s *fileMsService) FindPageByCnd(cnd *sqls.Cnd) (list []model.FileChunk, paging *sqls.Paging) {
	return repositories.FileMsRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *fileMsService) Create(t *model.FileChunk) error {
	err := repositories.FileMsRepository.Create(sqls.DB(), t)
	return err
}

func (s *fileMsService) Update(t *model.FileChunk) error {
	err := repositories.FileMsRepository.Update(sqls.DB(), t)
	return err
}

func (s *fileMsService) Updates(id int64, columns map[string]interface{}) error {
	err := repositories.FileMsRepository.Updates(sqls.DB(), id, columns)
	return err
}

func (s *fileMsService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.FileMsRepository.UpdateColumn(sqls.DB(), id, name, value)
	return err
}

func (s *fileMsService) Delete(id int64) {
	repositories.FileMsRepository.Delete(sqls.DB(), id)
}

// Scan 扫描
func (s *fileMsService) Scan(callback func(tasks []model.FileChunk)) {
	var cursor int64
	for {
		list := repositories.FileMsRepository.Find(sqls.DB(), sqls.NewCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}
