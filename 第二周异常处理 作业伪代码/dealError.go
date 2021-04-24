package main

import (
	"database/sql"
	xerrors "github.com/pkg/errors"
)

func DealError() error {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:9090)/test?parseTime=true")
	if err != nil {
		return xerrors.Wrap(err, "msg: open \"mysql\" failed")
	}
	var s string
	rows := db.QueryRow("SELECT shuaibi FROM test.hello LIMIT 1")
	xerrors.Wrap(rows.Err(), "msg: can not find shuaibi")
	err = rows.Scan(&s)
	return xerrors.Wrap(err, "msg: Scan copies the columns from the matched row into the value failed")
}