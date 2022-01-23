package repository

import "errors"

var (
	// ErrRecordNotFound return an error if record is not found from database
	ErrRecordNotFound = errors.New("record not found")
	// ErrNotAnyRecordAffect return an error if record is not affect in database
	ErrNotAnyRecordAffect = errors.New("not any record(s) affect")
)
