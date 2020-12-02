# 学习笔记

## 作业题目

1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

## 作业内容

```go
// ErrNoRows is returned by Scan when QueryRow doesn't return a
// row. In such a case, QueryRow returns a placeholder *Row value that
// defers this error until a Scan.
var ErrNoRows = errors.New("sql: no rows in result set")

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (db *DB) QueryRow(query string, args ...interface{}) *Row {
	return db.QueryRowContext(context.Background(), query, args...)
}
```

以上为`database/sql/sql.go`中的部分源码，从源码注释中，可以知道`sql.ErrNoRows`仅在调用`QueryRow`方法时会出现。

而`QueryRow`方法最多返回一条数据，并且当查询的数据少于一条时，则会返回`sql.ErrNoRows`错误。

在对源码进行分析之后，结合实际场景，什么时候需要使用到`QueryRow`呢？很明显，大多数情况是在主观意识中认定存在该条数据，才会使用`QueryRow`去获取该数据，如：根据用户 ID 去获取该 ID 所在行的数据，所以，引发该错误的原因主要还是由于数据不完整，因此应该需要在`dao`层`Wrap`该错误并抛给上一层。

```go
package main

import (
	"database/sql"
	"errors"
	"fmt"
	_errors "github.com/pkg/errors"
	"log"
)

var ErrDataIncomplete = errors.New("dao: data incomplete")

type user struct {
}

func main() {
	_, err := service(1)
	if err != nil {
		log.Printf("error: %T %v\n", _errors.Cause(err), _errors.Cause(err))
		log.Printf("statck trace: \n%+v\n", err)
		return
	}
}

// service
func service(id int) (*user, error) {
	u, err := dao(id)
	if err != nil {
		return nil, _errors.WithMessage(err, "service: can't get user data")
	}

	// TODO service coding
	return u, nil
}

// dao
func dao(id int) (u *user, err error) {
	// TODO get data in db
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}

		return nil, _errors.Wrap(ErrDataIncomplete, fmt.Sprintf("can't found user that id is (%d)", id))
	}

	return
}
```

