package main

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

var socialPrefixes = map[string]string{
	"linkedin":  "linkedin.com/in/",
	"instagram": "instagram.com/",
	"tiktok":    "tiktok.com/@",
	"threads":   "threads.com/@",
}

func validateCreate(req *CreateApplicantRequest) []FieldError {
	var errs []FieldError

	req.Name = strings.TrimSpace(req.Name)
	if n := utf8.RuneCountInString(req.Name); n < 2 || n > 100 {
		errs = append(errs, FieldError{"name", "must be 2 to 100 characters"})
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if msg := emailError(req.Email); msg != "" {
		errs = append(errs, FieldError{"email", msg})
	}

	req.Position = strings.ToLower(strings.TrimSpace(req.Position))
	if !knownPosition[Position(req.Position)] {
		errs = append(errs, FieldError{"position", "position not found. valid values: " + strings.Join(validPositionList(), ", ")})
	}

	req.Motivation = strings.TrimSpace(req.Motivation)
	if w := len(strings.Fields(req.Motivation)); w < 20 || w > 200 {
		errs = append(errs, FieldError{"motivation", fmt.Sprintf("must be 20 to 200 words (got %d)", w)})
	}

	req.CVURL = strings.TrimSpace(req.CVURL)
	if msg := urlError(req.CVURL); msg != "" {
		errs = append(errs, FieldError{"cv_url", msg})
	}

	req.PortfolioPDFURL = strings.TrimSpace(req.PortfolioPDFURL)
	if msg := urlError(req.PortfolioPDFURL); msg != "" {
		errs = append(errs, FieldError{"portfolio_pdf_url", msg})
	}

	req.PortfolioURL = strings.TrimSpace(req.PortfolioURL)
	if req.PortfolioURL != "" {
		if msg := urlError(req.PortfolioURL); msg != "" {
			errs = append(errs, FieldError{"portfolio_url", msg})
		}
	}

	return errs

}

func emailError(s string) string {
	if strings.ContainsAny(s, " \t\r\n") {
		return "must not contain spaces"
	}

	if s == "" {
		return "is required"
	}
	if len(s) > 254 {
		return "must be at most 254 characters"
	}
	if strings.Count(s, "@") != 1 {
		return "must contain exactly one @"
	}

	at := strings.IndexByte(s, '@')
	if at == 0 {
		return "email is invalid"
	}

	domain := s[at+1:]
	if domain == "" {
		return "email is invalid"
	}

	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return "domain must contain a dot"
	}
	for _, l := range labels {
		if l == "" {
			return "domain has an empty label"
		}
	}

	return ""
}

func urlError(s string) string {
	if s == "" {
		return "is required"
	}

	if len(s) > 2048 {
		return "must be at most 2048 characters"
	}

	u, err := url.ParseRequestURI(s)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "must be a valid http or https URL"
	}

	return ""
}

func socialMediaErros(m map[string]string) []FieldError {
	var errs []FieldError
	for key, val := range m {
		prefix, ok := socialPrefixes[key]
		if !ok {
			errs = append(errs, FieldError{
				"social_media" + key,
				"unknown platform",
			})
			continue
		}

		val = strings.TrimSpace(val)
		if val == "" {
			errs = append(errs, FieldError{
				"social_media" + key,
				"must start with https://",
			})
			continue
		}

		rest := strings.TrimPrefix(val, "https://")
		rest = strings.TrimPrefix(rest, "www.")

		if !strings.HasPrefix(rest, prefix) {
			errs = append(errs, FieldError{
				"social_media" + key,
				"must be a " + prefix + "URL",
			})
			continue
		}

		if len(rest) == len(prefix) {
			errs = append(errs, FieldError{
				"social_media" + key,
				"must include a username after " + prefix,
			})
		}
	}

	return errs
}
