package main

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type ListQuery struct {
	Stage    Stage
	Status   Status
	Position Position
	Page     int
	Limit    int
	Sort     string
	Order    string
}

type ListMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type ListResponse struct {
	Data []Applicant `json:"data"`
	Meta ListMeta    `json:"meta"`
}

func parseListQuery(r *http.Request) (ListQuery, []FieldError) {
	q := ListQuery{Page: 1, Limit: 20, Sort: "created_at", Order: "desc"}
	var errs []FieldError
	v := r.URL.Query()

	if s := strings.TrimSpace(v.Get("stage")); s != "" {
		if !knownStages[Stage(s)] {
			errs = append(errs, FieldError{"stage", "must be a known stage"})
		} else {
			q.Stage = Stage(s)
		}
	}

	if s := strings.TrimSpace(v.Get("status")); s != "" {
		if !knownStatus[Status(s)] {
			errs = append(errs, FieldError{"status", "must be a known status"})
		} else {
			q.Status = Status(s)
		}
	}

	if s := strings.ToLower(strings.TrimSpace(v.Get("position"))); s != "" {
		if !knownPosition[Position(s)] {
			errs = append(errs, FieldError{"position", "must be a known position"})
		} else {
			q.Position = Position(s)
		}
	}

	if s := v.Get("page"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 {
			errs = append(errs, FieldError{"page", "must be an integer >= 1"})
		} else {
			q.Page = n
		}
	}

	if s := v.Get("limit"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 100 {
			errs = append(errs, FieldError{"limit", "must be an integer between 1 and 100"})
		} else {
			q.Limit = n
		}
	}

	if s := strings.TrimSpace(v.Get("sort")); s != "" {
		switch s {
		case "created_at", "name":
			q.Sort = s
		default:
			errs = append(errs, FieldError{"sort", "must be created_at or name"})
		}
	}

	if s := strings.ToLower(strings.TrimSpace(v.Get("order"))); s != "" {
		switch s {
		case "asc", "desc":
			q.Order = s
		default:
			errs = append(errs, FieldError{"order", "must be asc or desc"})
		}
	}

	return q, errs
}

func applyList(all []Applicant, q ListQuery) ([]Applicant, int) {
	filtered := make([]Applicant, 0, len(all))
	for _, a := range all {
		if q.Stage != "" && a.Stage != q.Stage {
			continue
		}
		if q.Status != "" && a.Status != q.Status {
			continue
		}
		if q.Position != "" && a.Position != q.Position {
			continue
		}

		filtered = append(filtered, a)
	}

	sort.SliceStable(filtered, func(i, j int) bool {
		a, b := filtered[i], filtered[j]
		var less bool
		switch q.Sort {
		case "name":
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		default:
			less = a.CreatedAt.Before(b.CreatedAt)
		}

		if primaryEqual(a, b, q.Sort) {
			return a.ID < b.ID
		}
		if q.Order == "desc" {
			return !less
		}

		return less
	})

	total := len(filtered)

	start := (q.Page - 1) * q.Limit
	if start < 0 {
		start = 0
	}
	if start >= total {
		return []Applicant{}, total
	}

	end := start + q.Limit
	if end > total {
		end = total
	}

	page := make([]Applicant, end-start)
	copy(page, filtered[start:end])
	return page, total
}

func primaryEqual(a, b Applicant, sortKey string) bool {
	switch sortKey {
	case "name":
		return strings.ToLower(a.Name) == strings.ToLower(b.Name)
	default:
		return a.CreatedAt.Equal(b.CreatedAt)
	}
}
