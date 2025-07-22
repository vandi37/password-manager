package repo

import (
	"database/sql"

	"github.com/vandi37/vanerrors"
)

const (
	ErrorPreparing = "error preparing"
	ErrorExecuting = "error executing"
	ErrorScanning  = "error scanning"
)

// ReturnByRes Check checks is the number of rows affected wright
func ReturnByRes(res sql.Result, check func(int64) bool) error {
	n, err := res.RowsAffected()
	if err != nil {
		return vanerrors.Wrap(ErrorExecuting, err)
	}

	if !check(n) {
		return vanerrors.Simple(ErrorExecuting)
	}

	return nil
}

func Equals(n int64) func(int64) bool {
	return func(i int64) bool {
		return i == n
	}
}

//func More(n int64) func(int64) bool {
//	return func(i int64) bool {
//		return i > n
//	}
//}
