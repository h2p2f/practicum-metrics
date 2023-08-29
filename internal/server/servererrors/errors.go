// package servererrors describes server errors.
package servererrors

import "errors"

// ErrNotFound - an error that occurs when the requested object is missing.
var ErrNotFound = errors.New("not found")

// ErrNotImplemented - an error that occurs when trying to call an unimplemented method.
var ErrNotImplemented = errors.New("not implemented")

// ErrEmptyKey - an error that occurs when trying to use an empty key.
var ErrEmptyKey = errors.New("empty key")
