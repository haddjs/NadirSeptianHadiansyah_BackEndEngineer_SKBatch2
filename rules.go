package main

import "errors"

type Stage string
type Status string

const (
	StageCVScreening   Stage = "cv_screening"
	StageHRInterview   Stage = "hr_interview"
	StageAssignment    Stage = "assignment"
	StageUserInterview Stage = "user_interview"
	StageOffering      Stage = "offering"
	StageOnboarding    Stage = "onboarding"
)

const (
	StatusPending     Status = "pending"
	StatusPassed      Status = "passed"
	StatusRejected    Status = "rejected"
	StatusOffered     Status = "offered"
	StatusScheduled   Status = "scheduled"
	StatusRescheduled Status = "rescheduled"
	StatusAccepted    Status = "accepted"
	StatusDeclined    Status = "declined"
	StatusCompleted   Status = "completed"
)

type StageRule struct {
	Initial     Status
	Transitions map[Status][]Status
	AdvanceOn   Status
	Next        Stage
}

var rules = map[Stage]StageRule{
	StageCVScreening: {
		Initial: StatusPending,
		Transitions: map[Status][]Status{
			StatusPending: {StatusPassed, StatusRejected},
		},
		AdvanceOn: StatusPassed,
		Next:      StageHRInterview,
	},
	StageHRInterview: {
		Initial: StatusOffered,
		Transitions: map[Status][]Status{
			StatusOffered:     {StatusScheduled, StatusRejected},
			StatusScheduled:   {StatusRescheduled, StatusPassed, StatusRejected},
			StatusRescheduled: {StatusRescheduled, StatusPassed, StatusRejected},
		},
		AdvanceOn: StatusPassed,
		Next:      StageAssignment,
	},
	StageAssignment: {
		Initial: StatusOffered,
		Transitions: map[Status][]Status{
			StatusOffered: {StatusPassed, StatusRejected},
		},
		AdvanceOn: StatusPassed,
		Next:      StageUserInterview,
	},
	StageUserInterview: {
		Initial: StatusOffered,
		Transitions: map[Status][]Status{
			StatusOffered:     {StatusRejected, StatusScheduled},
			StatusScheduled:   {StatusRescheduled, StatusPassed, StatusRejected},
			StatusRescheduled: {StatusRescheduled, StatusPassed, StatusRejected},
		},
		AdvanceOn: StatusPassed,
		Next:      StageOffering,
	},
	StageOffering: {
		Initial: StatusOffered,
		Transitions: map[Status][]Status{
			StatusOffered: {StatusAccepted, StatusDeclined},
		},
		AdvanceOn: StatusAccepted,
		Next:      StageOnboarding,
	},
	StageOnboarding: {
		Initial: StatusScheduled,
		Transitions: map[Status][]Status{
			StatusScheduled: {StatusCompleted, StatusDeclined},
		},
	},
}

var (
	ErrUnknownStatus     = errors.New("unknown status")
	ErrIllegalTransition = errors.New("illegal transition")
	ErrTerminalStatus    = errors.New("applicant is in a terminal status")
	ErrUnknownPosition   = errors.New("position not found")
	ErrUnknownStage      = errors.New("stage not found")
)

var knownStatus = func() map[Status]bool {
	set := map[Status]bool{}
	for _, r := range rules {
		set[r.Initial] = true
		for from, tos := range r.Transitions {
			set[from] = true
			for _, to := range tos {
				set[to] = true
			}
		}
	}

	return set
}()

var knownStages = func() map[Stage]bool {
	set := make(map[Stage]bool, len(rules))
	for st := range rules {
		set[st] = true
	}

	return set
}()

func NextState(stage Stage, current, requested Status) (Stage, Status, error) {
	if !knownStatus[requested] {
		return stage, current, ErrUnknownStatus
	}

	rule, ok := rules[stage]
	if !ok {
		return stage, current, ErrUnknownStage
	}

	allowed, ok := rule.Transitions[current]
	if !ok {
		return stage, current, ErrTerminalStatus
	}

	if !contains(allowed, requested) {
		return stage, current, ErrIllegalTransition
	}

	if requested == rule.AdvanceOn && rule.Next != "" {
		next := rules[rule.Next]
		return rule.Next, next.Initial, nil
	}

	return stage, requested, nil
}

func contains(s []Status, v Status) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}

	return false
}
