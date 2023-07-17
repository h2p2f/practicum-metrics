package servererrors

import "errors"

var ErrNotFound = errors.New("not found")
var ErrNotImplemented = errors.New("not implemented")

var ErrEmptyKey = errors.New("empty key")
