package client

import (
	"database/sql"

	"github.com/0x726f6f6b6965/my-blog/blog-service/inernal/helper"
	"github.com/0x726f6f6b6965/my-blog/lib/config"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg *config.DBConfig) (db *sql.DB, cleanup func(), err error) {
	psqlInfo := helper.PSQLBuildQueryString(cfg.User, cfg.Password, cfg.DBName, cfg.Host, cfg.Port, cfg.SSLmode)
	db, err = sql.Open("postgres", psqlInfo)
	cleanup = func() {
		db.Close()
	}
	return
}
