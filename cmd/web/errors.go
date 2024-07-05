package main

import "errors"

var (
	ErrNoRecord           = errors.New(`no matching record found`)
	ErrInvalidCredentials = errors.New(`invalid credentials`)
	ErrDuplicateEmail     = errors.New(`email already taken`)
)

// vim: ts=4 sw=4 fdm=indent
