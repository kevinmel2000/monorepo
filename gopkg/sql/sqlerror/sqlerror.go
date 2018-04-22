package sqlerror

import (
	"github.com/lib/pq"
)

// list of sql error
const (
	errQuestionMark       pq.ErrorCode = "??"
	PQErrCodeDuplicateKey pq.ErrorCode = "23505"
)

// PQGetErrCode function
func PQGetErrCode(err error) pq.ErrorCode {
	if err, ok := err.(*pq.Error); ok {
		return err.Code
	}
	return errQuestionMark
}
