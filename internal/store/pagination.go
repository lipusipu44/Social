package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Pagination struct {
	Limit  int    `json:"limit" validate:"gte=0,lte=20"`
	Offset int    `json:"offset" validate:"gte=0,lte=100"`
	SortBy string `json:"sort_by" validate:"oneof=asc desc"`
	//adding new fields for filtering data from URL
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
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

	tags := qs.Get("tags")
	if tags != "" {
		pg.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		pg.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		pg.Since = parseTime(since)
	}

	until := qs.Get("until")
	if until != "" {
		pg.Until = parseTime(until)
	}

	return pg, nil
}

/*
our parse Time string in yyyy-mm-dd hh:mm:ss format
*/
func parseTime(since string) string {
	t, err := time.Parse(time.DateTime, since)
	if err != nil {
		return ""
	}
	return t.Format(time.DateTime)

}
