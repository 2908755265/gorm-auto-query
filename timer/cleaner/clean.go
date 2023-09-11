package cleaner

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"gorm.io/gorm"
	"time"
)

type ScopeFunc func(*gorm.DB) *gorm.DB

type DataCleaner[M any] struct {
	logx.Logger
	db        *gorm.DB
	m         *M
	timeout   int
	period    int
	scopeFunc ScopeFunc
}

func (c *DataCleaner[M]) Start() {
	threading.GoSafe(func() {
		c.startCleanLoop()
	})
}

func (c *DataCleaner[M]) SetScopeFunc(fn ScopeFunc) {
	c.scopeFunc = fn
}

func (c *DataCleaner[M]) startCleanLoop() {
	for range time.Tick(time.Duration(c.period) * time.Second) {
		err := c.clean()
		if err != nil {
			c.Error("clean data error: ", err)
		}
	}
}

func (c *DataCleaner[M]) clean() error {
	return c.db.Model(c.m).Scopes(c.getScopeFunc()).Delete(c.m).Error
}

func (c *DataCleaner[M]) getScopeFunc() ScopeFunc {
	if c.scopeFunc != nil {
		return c.scopeFunc
	}
	return c.defaultScopeFunc
}

func (c *DataCleaner[M]) defaultScopeFunc(db *gorm.DB) *gorm.DB {
	return db.Where("created_at <= ?", time.Now().Add(time.Duration(-c.timeout)*time.Second))
}

func NewCleaner[M any](db *gorm.DB, timeout, period int) *DataCleaner[M] {
	var m M
	return &DataCleaner[M]{
		Logger:  logx.WithContext(context.Background()),
		db:      db,
		m:       &m,
		timeout: timeout,
		period:  period,
	}
}
