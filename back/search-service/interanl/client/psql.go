package client

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/0x726f6f6b6965/my-blog/lib/config"
	_ "github.com/lib/pq"
)

// PSQLBuildQueryString builds a query string.
func PSQLBuildQueryString(user, pass, dbname, host string, port int, sslmode string) string {
	parts := []string{}
	if len(user) != 0 {
		parts = append(parts, fmt.Sprintf("user=%s", user))
	}
	if len(pass) != 0 {
		parts = append(parts, fmt.Sprintf("password=%s", pass))
	}
	if len(dbname) != 0 {
		parts = append(parts, fmt.Sprintf("dbname=%s", dbname))
	}
	if len(host) != 0 {
		parts = append(parts, fmt.Sprintf("host=%s", host))
	}
	if port != 0 {
		parts = append(parts, fmt.Sprintf("port=%d", port))
	}
	if len(sslmode) != 0 {
		parts = append(parts, fmt.Sprintf("sslmode=%s", sslmode))
	}

	return strings.Join(parts, " ")
}

func NewPostgres(cfg *config.DBConfig) (db *sql.DB, cleanup func(), err error) {
	psqlInfo := PSQLBuildQueryString(cfg.User, cfg.Password, cfg.DBName, cfg.Host, cfg.Port, cfg.SSLmode)
	db, err = sql.Open("postgres", psqlInfo)
	cleanup = func() {
		db.Close()
	}
	return
}
