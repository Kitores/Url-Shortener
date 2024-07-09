package storage

import "errors"

var (
	ErrAliasNotFound = errors.New("alias not found")
	ErrUrlNotFound   = errors.New("url not found")
	ErrUrlExists     = errors.New("url already exists")
)
