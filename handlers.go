package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func (api *API) updateStage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid applicant id")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req UpdateStageRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	requested := Status(strings.ToLower(strings.TrimSpace(req.Status)))
	if requested == "" {
		writeError(w, http.StatusBadRequest, "status is required",
			FieldError{"status", "is required"})
		return
	}

	updated, err := api.store.Transition(id, func(a Applicant) (Applicant, error) {
		nextStage, nextStatus, err := NextState(a.Stage, a.Status, requested)
		if err != nil {
			return Applicant{}, err
		}

		var dateErrs []FieldError
		var scheduledAt, startDate *string

		switch requested {
		case StatusScheduled, StatusRescheduled:
			if req.ScheduledAt == nil || strings.TrimSpace(*req.ScheduledAt) == "" {
				dateErrs = append(dateErrs, FieldError{"scheduled_at", "is required for " + string(requested)})
			} else if v, err := validateDate(*req.ScheduledAt); err != nil {
				dateErrs = append(dateErrs, FieldError{"scheduled_at", err.Error()})
			} else {
				scheduledAt = &v
			}
		case StatusAccepted:
			if req.StartDate == nil || strings.TrimSpace(*req.StartDate) == "" {
				dateErrs = append(dateErrs, FieldError{"start_date", "is required for accepted"})
			} else if v, err := validateDate(*req.StartDate); err != nil {
				dateErrs = append(dateErrs, FieldError{"start_date", err.Error()})
			} else {
				startDate = &v
			}
		}

		if len(dateErrs) > 0 {
			return Applicant{}, &ValidationError{Details: dateErrs}
		}

		stageChanged := a.Stage != nextStage
		a.Stage = nextStage
		a.Status = nextStatus

		if scheduledAt != nil {
			a.ScheduledAt = scheduledAt
		}
		if startDate != nil {
			a.StartDate = startDate
		}

		if stageChanged {
			a.ScheduledAt = nil
		}

		return a, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			writeError(w, http.StatusNotFound, "applicant not found")
		case errors.Is(err, ErrUnknownStatus):
			writeError(w, http.StatusBadRequest, err.Error(),
				FieldError{"status", "unknown status"})
		case errors.Is(err, ErrTerminalStatus):
			writeError(w, http.StatusLocked, err.Error())
		case errors.Is(err, ErrIllegalTransition):
			writeError(w, http.StatusUnprocessableEntity, err.Error(),
				FieldError{"status", "not allowed from current status"})
		default:
			var ve *ValidationError
			if errors.As(err, &ve) {
				writeError(w, http.StatusBadRequest, "validation failed", ve.Details...)
				return
			}
			log.Printf("updateStage: %v", err)
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	writeJSON(w, http.StatusOK, updated)
}
