package cache

import (
	"gorm.io/gorm"
)

type CacheLogRepInterface interface {
	InsertCacheLog(*CacheLog) error
	InsertCacheLogs([]*CacheLog) error
	RetrieveCacheLogByPath(*CacheLog) error
	UpdateCacheLog(*CacheLog) error
	DeleteCacheLog(string) error
}

var CacheLogRepository *CacheLogRepo

type CacheLogRepo struct {
	*gorm.DB
}

func CreateCacheRepo(db *gorm.DB) {
	CacheLogRepository = &CacheLogRepo{
		DB: db,
	}
}

func (c *CacheLogRepo) InsertCacheLog(cacheLog *CacheLog) error {
	result := c.DB.Create(cacheLog)
	return result.Error
}

func (c *CacheLogRepo) InsertCacheLogs(cacheLogs []*CacheLog) error {
	result := c.DB.Create(cacheLogs)
	return result.Error
}

func (c *CacheLogRepo) RetrieveCacheLogByPath(cacheLog *CacheLog) error {
	result := c.DB.Last(cacheLog, "filepath = ?", cacheLog.FilePath)
	return result.Error
}

func (c *CacheLogRepo) UpdateCacheLog(cacheLog *CacheLog) error {
	result := c.DB.Save(cacheLog)
	return result.Error
}

func (c *CacheLogRepo) DeleteCacheLog(cacheLogPath string) error {
	result := c.DB.Where("filepath = ?", cacheLogPath).Delete(&CacheLog{})
	return result.Error
}
