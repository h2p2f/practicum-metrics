// Package servererrors описывает ошибки сервера.
//
// package servererrors describes server errors.
package servererrors

import "errors"

// ErrNotFound - ошибка, возникающая при отсутствии запрашиваемого объекта.
//
// ErrNotFound - an error that occurs when the requested object is missing.
var ErrNotFound = errors.New("not found")

// ErrNotImplemented - ошибка, возникающая при попытке вызова не реализованного метода.
//
// ErrNotImplemented - an error that occurs when trying to call an unimplemented method.
var ErrNotImplemented = errors.New("not implemented")

// ErrEmptyKey - ошибка, возникающая при попытке использования пустого ключа.
//
// ErrEmptyKey - an error that occurs when trying to use an empty key.
var ErrEmptyKey = errors.New("empty key")
