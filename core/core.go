package core

import (
	"database/sql"
	"messenger/core/models"
	"messenger/database"
	"messenger/web"
)

type Core struct {
	db  *sql.DB
	web *web.Web
}

func NewCore() (*Core, error) {
	db, err := database.Init()
	if err != nil {
		return nil, err
	}

	webServer := web.NewWeb(models.NewModelManager(db))
	return &Core{
		db:  db,
		web: webServer,
	}, nil
}

func (c *Core) Run() {
	c.web.Run()
}
