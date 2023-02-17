package services

import (
	"fileMS/model"
	"fileMS/repositories"
	"github.com/gannicus-w/yunqi_mysql/sqls"
	"github.com/gannicus-w/yunqi_mysql/web/params"
)

var FileMsVersionService = newFileMsVersionService()

func newFileMsVersionService() *fileMsVersionService {
	return &fileMsVersionService{}
}

type fileMsVersionService struct {
}

func (s *fileMsVersionService) Get(id int64) *model.FileVersion {
	return repositories.FileMsVersionRepository.Get(sqls.DB(), id)
}

func (s *fileMsVersionService) Take(where ...interface{}) *model.FileVersion {
	return repositories.FileMsVersionRepository.Take(sqls.DB(), where...)
}

func (s *fileMsVersionService) Find(cnd *sqls.Cnd) []model.FileVersion {
	return repositories.FileMsVersionRepository.Find(sqls.DB(), cnd)
}

func (s *fileMsVersionService) FindOne(cnd *sqls.Cnd) *model.FileVersion {
	return repositories.FileMsVersionRepository.FindOne(sqls.DB(), cnd)
}

func (s *fileMsVersionService) FindPageByParams(params *params.QueryParams) (list []model.FileVersion, paging *sqls.Paging) {
	return repositories.FileMsVersionRepository.FindPageByParams(sqls.DB(), params)
}

func (s *fileMsVersionService) FindPageByCnd(cnd *sqls.Cnd) (list []model.FileVersion, paging *sqls.Paging) {
	return repositories.FileMsVersionRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *fileMsVersionService) Create(t *model.FileVersion) error {
	err := repositories.FileMsVersionRepository.Create(sqls.DB(), t)
	return err
}

func (s *fileMsVersionService) Update(t *model.FileVersion) error {
	err := repositories.FileMsVersionRepository.Update(sqls.DB(), t)
	return err
}

func (s *fileMsVersionService) Updates(id int64, columns map[string]interface{}) error {
	err := repositories.FileMsVersionRepository.Updates(sqls.DB(), id, columns)
	return err
}

func (s *fileMsVersionService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.FileMsVersionRepository.UpdateColumn(sqls.DB(), id, name, value)
	return err
}

func (s *fileMsVersionService) Delete(id int64) {
	repositories.FileMsVersionRepository.Delete(sqls.DB(), id)
}

// Scan 扫描
func (s *fileMsVersionService) Scan(callback func(tasks []model.FileVersion)) {
	var cursor int64
	for {
		list := repositories.FileMsVersionRepository.Find(sqls.DB(), sqls.NewCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}
