package manager

import (
	"errors"
	"strings"
)

/*
  // https://github.com/libwww-perl/uri-db/

  db:engine://username:password@example.com:8042/widgets.db?tz=utc&charset=utf8#users
  \/ \____/   \_______________/ \_________/ \__/ \________/ \/ \__/ \____/ \__/\____/
   |    |             |              |       |        |     |    |     |    |     |
   |    |         userinfo        hostname  port      |    key   |    key   |     |
   |    |     \________________________________/      |          |          |     |
   |    |                      |                      |        value      value   |
   |  engine                   |                      |     \_________________/   |
scheme  |                  authority         db name or path         |            |
 name   |     \___________________________________________/        query       fragment
   |    |                           |
   |    |                   hierarchical part
   |    |
   |    |      db name or path       query    fragment
   |  __|_   ________|________    _____|____  ____|____
  /\ /    \ /                 \  /          \/         \
  db:engine:my_big_fat_database?encoding=big5#log.animals
*/

// Params 数据库参数
// Engine & Database can't be empty
type Params struct {
	Engine   string
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Query    string
	Fragment string
}

// 从URI字符串中解析数据库参数
func ParseDbUri(uri string) (params *Params, err error) {
	if !strings.HasPrefix(uri, `db:`) {
		return nil, errors.New(`DB uri must be started with "db:"`)
	}
	uri = uri[3:]

	pos := strings.Index(uri, `:`)
	if pos < 1 {
		return nil, errors.New(`DB engine not found`)
	}
	params = &Params{
		Engine: uri[0:pos],
	}
	uri = uri[pos+1:]

	if strings.HasPrefix(uri, `//`) {
		uri = uri[2:]
		pos = strings.LastIndex(uri, `@`)
		if pos != -1 {
			params.Username = uri[0:pos]
			pos2 := strings.Index(params.Username, `:`)
			if pos2 != -1 {
				params.Password = params.Username[pos2+1:]
				params.Username = params.Username[0:pos2]
			}
			uri = uri[pos+1:]
		}

		pos = strings.Index(uri, `/`)
		if pos != -1 {
			params.Host = uri[0:pos]
			pos2 := strings.Index(params.Host, `:`)
			if pos2 != -1 {
				params.Port = params.Host[pos2+1:]
				params.Host = params.Host[0:pos2]
			}
			uri = uri[pos+1:]
		}
	}

	pos = strings.Index(uri, `?`)
	if pos != -1 {
		params.Database = uri[0:pos]
		params.Query = uri[pos+1:]
		pos = strings.Index(params.Query, `#`)
		if pos != -1 {
			params.Fragment = params.Query[pos+1:]
			params.Query = params.Query[0:pos]
		}
	} else {
		pos = strings.Index(uri, `#`)
		if pos != -1 {
			params.Database = uri[0:pos]
			params.Fragment = uri[pos+1:]
		} else {
			params.Database = uri
		}
	}

	if params.Database == `` {
		return nil, errors.New(`DB name not found.`)
	}

	return params, nil
}

// DbUri2Dsn 转换URI为DSN
func DbUri2Dsn(params *Params) string {
	dsn := ""
	if params.Username != `` {
		dsn += params.Username
		if params.Password != `` {
			dsn += `:` + params.Password
		}
		dsn += `@`
	}
	if params.Host != `` {
		dsn += `tcp(` + params.Host
		if params.Port != `` {
			dsn += `:` + params.Port
		}
		dsn += `)`
	}
	dsn += `/` + params.Database
	if params.Query != `` {
		dsn += `?` + params.Query
	}
	if params.Fragment != `` {
		dsn += `#` + params.Fragment
	}
	return dsn
}
