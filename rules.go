package main

import (
	"errors"
	"sort"
)

// --- Stage & Status Definitions ---

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

// --- Position Definitions ---

type Position string

const (
	PositionAIAutomationSpecialist      Position = "ai_automation_specialist"
	PositionAIContentSpecialist         Position = "ai_content_specialist"
	PositionAIEngineer                  Position = "ai_engineer"
	PositionAIProductSpecialist         Position = "ai_product_specialist"
	PositionAIResearchAndEvaluation     Position = "ai_research_and_evaluation"
	PositionAIWorkflowSpecialist        Position = "ai_workflow_specialist"
	PositionAdministrationSpecialist    Position = "administration_specialist"
	PositionBackendEngineer             Position = "backend_engineer"
	PositionCommunityAcquisition        Position = "community_acquisition"
	PositionContentCreator              Position = "content_creator"
	PositionContentWriter               Position = "content_writer"
	PositionDataAnalyst                 Position = "data_analyst"
	PositionEmployeeEngagement          Position = "employee_engagement"
	PositionEmployerBrandingSpecialist  Position = "employer_branding_specialist"
	PositionEventManagement             Position = "event_management"
	PositionFinanceAndOperationsSupport Position = "finance_and_operations_support"
	PositionFrontendEngineer            Position = "frontend_engineer"
	PositionGraphicDesigner             Position = "graphic_designer"
	PositionLearningAndDevelopment      Position = "learning_and_development"
	PositionPartnershipSpecialist       Position = "partnership_specialist"
	PositionPeopleOperations            Position = "people_operations"
	PositionPromptEngineer              Position = "prompt_engineer"
	PositionQAAndProductTester          Position = "qa_and_product_tester"
	PositionSocialMediaOfficer          Position = "social_media_officer"
	PositionTalentAcquisitionSpecialist Position = "talent_acquisition_specialist"
	PositionThreadSpecialist            Position = "thread_specialist"
	PositionVideoEditor                 Position = "video_editor"
	PositionWebDeveloper                Position = "web_developer"
)

var allPositions = []Position{
	PositionAIAutomationSpecialist,
	PositionAIContentSpecialist,
	PositionAIEngineer,
	PositionAIProductSpecialist,
	PositionAIResearchAndEvaluation,
	PositionAIWorkflowSpecialist,
	PositionAdministrationSpecialist,
	PositionBackendEngineer,
	PositionCommunityAcquisition,
	PositionContentCreator,
	PositionContentWriter,
	PositionDataAnalyst,
	PositionEmployeeEngagement,
	PositionEmployerBrandingSpecialist,
	PositionEventManagement,
	PositionFinanceAndOperationsSupport,
	PositionFrontendEngineer,
	PositionGraphicDesigner,
	PositionLearningAndDevelopment,
	PositionPartnershipSpecialist,
	PositionPeopleOperations,
	PositionPromptEngineer,
	PositionQAAndProductTester,
	PositionSocialMediaOfficer,
	PositionTalentAcquisitionSpecialist,
	PositionThreadSpecialist,
	PositionVideoEditor,
	PositionWebDeveloper,
}

var knownPosition = func() map[Position]bool {
	set := make(map[Position]bool, len(allPositions))
	for _, p := range allPositions {
		set[p] = true
	}
	return set
}()

func validPositionList() []string {
	out := make([]string, 0, len(allPositions))
	for _, p := range allPositions {
		out = append(out, string(p))
	}
	sort.Strings(out)
	return out
}

// --- Transition Rules ---

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

// --- Errors ---

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
