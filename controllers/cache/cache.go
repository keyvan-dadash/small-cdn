package cache

import (
	"bufio"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/sod-lol/small-cdn/core/models/cache"
	"github.com/sod-lol/small-cdn/core/models/user"
)

type cacheFileReqBody struct {
	Type    string `form:"type" json:"type" xml:"type"  binding:"required"`
	Name    string `form:"name" json:"name" xml:"name"  binding:"required"`
	Content string `form:"content" json:"content" xml:"content"  binding:"required"`
	Minify  bool   `form:"minify" json:"minify" xml:"minify"`
}

func processCacheFile(cacheLog *cache.CacheLog, cacheFile cache.CacheFile) {
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	start := time.Now()
	err := cacheFile.Process()
	duration := time.Now().Sub(start)
	runtime.ReadMemStats(&m2)

	if err != nil {
		logrus.Errorf("[Error] Could not process cache file: %v", err)
		cacheLog = nil
		return
	}

	cacheLog.ConsumedMemory = m2.TotalAlloc - m1.TotalAlloc
	cacheLog.DurationOfMinifcation = duration

	dir, _ := filepath.Split(cacheLog.FilePath)
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		logrus.Errorf("[Error] Could not Mkdir folders: %v", err)
		cacheLog = nil
		return
	}

	f, err := os.Create(cacheLog.FilePath)
	if err != nil {
		logrus.Errorf("[Error] Could not create file: %v", err)
		cacheLog = nil
		return
	}

	w := bufio.NewWriter(f)
	err = cacheFile.Write(w)
	if err != nil {
		logrus.Errorf("[Error] Could write to the file: %v", err)
		cacheLog = nil
		return
	}

	w.Flush()
	f, err = os.Open(cacheLog.FilePath)
	stats, err := f.Stat()
	if err != nil {
		logrus.Errorf("[Error] Couldnt get file's stats: %v", err)
		cacheLog = nil
		return
	}

	cacheLog.FileSize = uint64(stats.Size())
}

func processCacheFilesInParrarel(cacheLogs []*cache.CacheLog, cacheFiles []cache.CacheFile) error {
	var wg sync.WaitGroup
	wg.Add(len(cacheLogs))

	for index, item := range cacheLogs {
		cacheFile := cacheFiles[index]
		go func(cacheLog *cache.CacheLog, cacheFile cache.CacheFile) {
			processCacheFile(cacheLog, cacheFile)
			wg.Done()
		}(item, cacheFile)
	}

	wg.Wait()

	for _, item := range cacheLogs {
		if item == nil {
			return errors.New("Could not process all cache files")
		}
	}

	return cache.CacheLogRepository.InsertCacheLogs(cacheLogs)
}

func HandleAddCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cacheFilesReq []cacheFileReqBody

		if err := c.ShouldBind(&cacheFilesReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		username := c.GetString("username") //come from middleware
		userID := c.GetUint("userID")       //come from middleware

		var cacheLogs []*cache.CacheLog
		var cacheFiles []cache.CacheFile
		for _, item := range cacheFilesReq {
			cacheLog := cache.CreateCacheLog(item.Name, username)
			cacheLogs = append(cacheLogs, cacheLog)
			cacheLog.UserID = userID
			var cacheLogType string
			if item.Minify {
				cacheLogType = strings.ToLower(item.Type)
			} else {
				cacheLogType = "simple"
			}

			err, cacheFile := cache.CreateCacheFileFactory(cacheLog, cacheLogType)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			cacheFiles = append(cacheFiles, cacheFile)
		}

		for index, item := range cacheFilesReq {
			cacheFiles[index].Read(strings.NewReader(item.Content))
		}

		err := processCacheFilesInParrarel(cacheLogs, cacheFiles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusCreated, gin.H{})
	}
}

type cacheLogResBody struct {
	FileName              string
	FileSize              uint64
	FileType              string
	DurationOfMinifcation time.Duration
	ConsumedMemory        uint64
	CreatedAt             time.Time
}

func convertCacheLogToCacheLogResBody(cacheFile *cache.CacheLog) cacheLogResBody {
	err, cacheFileType := cache.ConvertFileTypeToString(cacheFile.FileType)
	if err != nil {
		return cacheLogResBody{}
	}

	return cacheLogResBody{
		FileName:              cacheFile.FileName,
		FileSize:              cacheFile.FileSize,
		FileType:              cacheFileType,
		DurationOfMinifcation: cacheFile.DurationOfMinifcation,
		ConsumedMemory:        cacheFile.ConsumedMemory,
		CreatedAt:             cacheFile.CreatedAt,
	}
}

func HandleListOfCacheFiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username") //come from middleware

		err, cacheFiles := user.UserRepository.RetrieveUserPreloadCacheLogs(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		var cacheLogsRes []cacheLogResBody
		for _, item := range cacheFiles {
			cacheLogsRes = append(cacheLogsRes, convertCacheLogToCacheLogResBody(&item))
		}

		c.JSON(http.StatusOK, cacheLogsRes)
	}
}
