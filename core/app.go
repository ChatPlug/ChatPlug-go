package core

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type App struct {
	router *gin.Engine
	db     *gorm.DB
	sm     *ServiceManager
	sl     *ServiceLoader
	ch     *ConfigurationHandler
}

func NewApp() *App {
	app := &App{
		sm: &ServiceManager{},
		ch: &ConfigurationHandler{},
	}
	app.sl = &ServiceLoader{App: app}
	return app
}

func (app *App) Init() {
	var err error

	app.sm.LoadAvailableServices()
	app.sl.Initialize()
	app.ch.WatchForConfiguration()

	app.db, err = gorm.Open("sqlite3", "cp.db")
	if err != nil {
		panic("failed to connect to database")
	}

	app.db.AutoMigrate(&Thread{})
	app.db.AutoMigrate(&ThreadGroup{})
	app.db.AutoMigrate(&Message{})
	app.db.AutoMigrate(&MessageAuthor{})
	app.db.AutoMigrate(&ServiceInstance{})
	app.db.AutoMigrate(&Attachment{})
}
