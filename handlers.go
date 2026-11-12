package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type API struct {
	store *Store
}

func (api *API) createApplicant(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req CreateApplicantRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if errs := validateCreate(&req); len(errs) > 0 {
		writeError(w, http.StatusBadRequest, "validation failed", errs...)
		return
	}

	applicant := Applicant{
		Name:            req.Name,
		Email:           req.Email,
		Position:        Position(req.Position),
		Stage:           StageCVScreening,
		Status:          rules[StageCVScreening].Initial,
		Motivation:      req.Motivation,
		CVURL:           req.CVURL,
		PortfolioPDFURL: req.PortfolioPDFURL,
		PortfolioURL:    req.PortfolioURL,
		SocialMedia:     req.SocialMedia,
	}

	created, err := api.store.Create(applicant)
	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			writeError(w, http.StatusConflict, "email already exists")
			return
		}

		log.Printf("create applicant: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create applicant")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (api *API) listApplicants(w http.ResponseWriter, r *http.Request) {
	q, errs := parseListQuery(r)
	if len(errs) > 0 {
		writeError(w, http.StatusBadRequest, "invalid query parameters", errs...)
		return
	}

	page, total := applyList(api.store.List(), q)

	totalPages := (total + q.Limit - 1) / q.Limit

	writeJSON(w, http.StatusOK, ListResponse{
		Data: page,
		Meta: ListMeta{Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages},
	})
}
