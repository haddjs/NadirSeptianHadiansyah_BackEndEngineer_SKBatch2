package main

import "time"

type Applicant struct {
	ID              int               `json:"id"`
	Name            string            `json:"name"`
	Email           string            `json:"email"`
	Position        Position          `json:"position"`
	Stage           Stage             `json:"stage"`
	Status          Status            `json:"status"`
	Motivation      string            `json:"motivation"`
	CVURL           string            `json:"cv_url"`
	PortfolioPDFURL string            `json:"portfolio_pdf_url"`
	PortfolioURL    string            `json:"portfolio_url,omitempty"`
	SocialMedia     map[string]string `json:"social_media,omitempty"`
	ScheduledAt     *string           `json:"scheduled_at"`
	StartDate       *string           `json:"start_date"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type CreateApplicantRequest struct {
	Name            string            `json:"name"`
	Email           string            `json:"email"`
	Position        string            `json:"position"`
	Motivation      string            `json:"motivation"`
	CVURL           string            `json:"cv_url"`
	PortfolioPDFURL string            `json:"portfolio_pdf_url"`
	PortfolioURL    string            `json:"portfolio_url"`
	SocialMedia     map[string]string `json:"social_media"`
}
