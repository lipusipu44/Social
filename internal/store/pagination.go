package store

import (
	"net/http"
	"strconv"
)

type Pagination struct {
	Limit  int    `json:"limit" validate:"gte=0,lte=20"`
	Offset int    `json:"offset" validate:"gte=0,lte=100"`
	SortBy string `json:"sort_by" validate:"oneof=asc desc"`
}

//Parse
/*
Simple parse of the url to fetch all params.
and put it in Pagination struct.

I tried with chi.URLparam to get the qs and it did not work

This to be called in feed.go in cmd/api section to fetch the param and this struct
to be passed to storage class in Post for sql query
*/
func (pg Pagination) Parse(r *http.Request) (Pagination, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return pg, nil
		}
		pg.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return pg, nil
		}
		pg.Limit = o
	}

	sortBy := qs.Get("sort_by")
	if sortBy != "" {
		pg.SortBy = sortBy
	}
	return pg, nil
}
