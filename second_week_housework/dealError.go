package main

import (
	"database/sql"
	xerrors "github.com/pkg/errors"
)

func DealError(db *sql.DB) error {
	var s string
	sqlQueryRow := "SELECT shuaibi FROM test.hello LIMIT 1"
	rows := db.QueryRow(sqlQueryRow)
	if rows.Err() != nil {
		// 查询不到往上抛，因为可能是username、pw校验等重要信息,让上层对rows.Err()做sentinel errors类型断言处理
		return xerrors.Wrapf(rows.Err(), "msg: Excute sql QueryRow [%v] failed", sqlQueryRow)
	}
	return xerrors.Wrap(rows.Scan(&s), "msg: Scan copies the columns from the matched row into the value failed")
}