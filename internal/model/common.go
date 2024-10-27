package model

import "errors"

type TId = string

const (
	TSortAsc = "ASC"
	TSortDesc = "DESC"
)

var InternalErr = errors.New("Internal error")
var InvalidDataErr = errors.New("Bad request")
var NotFoundErr = errors.New("Not found")

type PaginationDto struct {
	TotalNumber int `json:"totalNumber"`
	PageSize int `json:"pageSize"`
	PageCount int `json:"pageCount"`
	PageNumber int `json:"pageNumber"`
}
