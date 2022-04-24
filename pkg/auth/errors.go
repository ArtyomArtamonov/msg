package auth

import "errors"

var ErrAlreadyExists = errors.New("record already exists")
var ErrNotFound = errors.New("not found")
