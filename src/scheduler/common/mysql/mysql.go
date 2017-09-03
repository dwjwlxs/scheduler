package mysql

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	MaxOpenConns = 0                //max conn, 0 for no limit
	MaxIdleConns = 500              //max number of idle connections
	MaxLifetime  = 20 * time.Second //timeout, 0 for alive forever
)

type Mysql struct {
	host       string
	port       string
	user       string
	password   string
	dbname     string
	enablePool bool
	timeout    time.Duration
	DB         *sql.DB
}

func NewMysql(dbInfo map[string]interface{}, enablePool bool, timeout time.Duration) (*Mysql, error) {
	if dbInfo["host"] == "" || dbInfo["port"] == 0 {
		return nil, errors.New("mysql host:port should not be empty")
	}

	m := &Mysql{
		host:       dbInfo["host"].(string),
		port:       dbInfo["port"].(string),
		dbname:     dbInfo["dbname"].(string),
		enablePool: enablePool,
		timeout:    timeout,
	}

	if dbInfo["user"] != nil {
		m.user = dbInfo["user"].(string)
	}
	if dbInfo["password"] != nil {
		m.password = dbInfo["password"].(string)
	}

	db, err := sql.Open("mysql", m.user+":"+m.password+"@tcp("+m.host+":"+m.port+")/"+m.dbname)
	if err != nil {
		return nil, err
	}
	m.DB = db
	return m, nil
}
