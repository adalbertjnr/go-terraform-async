package main

import "errors"

var (
	ErrEmptyTask    = errors.New("empty list of tasks")
	ErrVerbNotFound = errors.New("the verb should be plan, apply or destroy")
)
