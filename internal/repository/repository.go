package repository

import "errors"

// Общие ошибки для repository package

var ErrNotFound = errors.New("shorter url not found")
var ErrAlreadyExists = errors.New("shorter url already exists")
