package main

import "testing"

func TestNextState(t *testing.T) {
	tests := []struct {
		name      string
		stage     Stage
		current   Status
		requested Status
		wantStage Stage
		wantState Status
		wantErr   error
	}{
		{"cv passed advances to HR", StageCVScreening, StatusPending, StatusPassed, StageHRInterview, StatusOffered, nil},
		{"hr passed advances to assignment", StageHRInterview, StatusScheduled, StatusPassed, StageAssignment, StatusOffered, nil},
		{"assignment passed advances to user interview", StageAssignment, StatusOffered, StatusPassed, StageUserInterview, StatusOffered, nil},
		{"user interview passed advances to offering", StageUserInterview, StatusScheduled, StatusPassed, StageOffering, StatusOffered, nil},
		{"offering accepted advances to onboarding", StageOffering, StatusOffered, StatusAccepted, StageOnboarding, StatusScheduled, nil},

		{"cv stays on rejected", StageCVScreening, StatusPending, StatusRejected, StageCVScreening, StatusRejected, nil},
		{"hr offered to scheduled stays in stage", StageHRInterview, StatusOffered, StatusScheduled, StageHRInterview, StatusScheduled, nil},
		{"hr scheduled to rescheduled stays in stage", StageHRInterview, StatusScheduled, StatusRescheduled, StageHRInterview, StatusRescheduled, nil},
		{"hr rescheduled to rescheduled stays in stage", StageHRInterview, StatusRescheduled, StatusRescheduled, StageHRInterview, StatusRescheduled, nil},
		{"hr scheduled rejected stays in stage", StageHRInterview, StatusScheduled, StatusRejected, StageHRInterview, StatusRejected, nil},
		{"assignment rejected stays in stage", StageAssignment, StatusOffered, StatusRejected, StageAssignment, StatusRejected, nil},
		{"user interview offered to scheduled stays in stage", StageUserInterview, StatusOffered, StatusScheduled, StageUserInterview, StatusScheduled, nil},
		{"offering declined stays in stage", StageOffering, StatusOffered, StatusDeclined, StageOffering, StatusDeclined, nil},
		{"onboarding scheduled to completed (terminal stay)", StageOnboarding, StatusScheduled, StatusCompleted, StageOnboarding, StatusCompleted, nil},
		{"onboarding scheduled to declined (terminal stay)", StageOnboarding, StatusScheduled, StatusDeclined, StageOnboarding, StatusDeclined, nil},

		{"unknown requested status", StageCVScreening, StatusPending, Status("bogus"), StageCVScreening, StatusPending, ErrUnknownStatus},
		{"unknown status even on terminal current", StageOnboarding, StatusCompleted, Status("bogus"), StageOnboarding, StatusCompleted, ErrUnknownStatus},
		{"unknown stage returns ErrUnknownStage", Stage("does_not_exist"), StatusPending, StatusPassed, Stage("does_not_exist"), StatusPending, ErrUnknownStage},

		{"cv rejected is terminal", StageCVScreening, StatusRejected, StatusPassed, StageCVScreening, StatusRejected, ErrTerminalStatus},
		{"hr rejected is terminal", StageHRInterview, StatusRejected, StatusPassed, StageHRInterview, StatusRejected, ErrTerminalStatus},
		{"assignment rejected is terminal", StageAssignment, StatusRejected, StatusPassed, StageAssignment, StatusRejected, ErrTerminalStatus},
		{"offering declined is terminal", StageOffering, StatusDeclined, StatusAccepted, StageOffering, StatusDeclined, ErrTerminalStatus},
		{"onboarding completed is terminal", StageOnboarding, StatusCompleted, StatusDeclined, StageOnboarding, StatusCompleted, ErrTerminalStatus},
		{"onboarding declined is terminal", StageOnboarding, StatusDeclined, StatusCompleted, StageOnboarding, StatusDeclined, ErrTerminalStatus},
		{"unknown stage treated as terminal", Stage("does_not_exist"), StatusPending, StatusPassed, Stage("does_not_exist"), StatusPending, ErrTerminalStatus},

		{"cv pending cannot jump to offered", StageCVScreening, StatusPending, StatusOffered, StageCVScreening, StatusPending, ErrIllegalTransition},
		{"cv pending cannot go to scheduled", StageCVScreening, StatusPending, StatusScheduled, StageCVScreening, StatusPending, ErrIllegalTransition},
		{"hr offered cannot jump to passed", StageHRInterview, StatusOffered, StatusPassed, StageHRInterview, StatusOffered, ErrIllegalTransition},
		{"hr offered cannot go to rescheduled", StageHRInterview, StatusOffered, StatusRescheduled, StageHRInterview, StatusOffered, ErrIllegalTransition},
		{"hr scheduled cannot go back to offered", StageHRInterview, StatusScheduled, StatusOffered, StageHRInterview, StatusScheduled, ErrIllegalTransition},
		{"assignment offered cannot go to accepted", StageAssignment, StatusOffered, StatusAccepted, StageAssignment, StatusOffered, ErrIllegalTransition},
		{"offering offered cannot go to passed", StageOffering, StatusOffered, StatusPassed, StageOffering, StatusOffered, ErrIllegalTransition},
		{"onboarding scheduled cannot go to passed", StageOnboarding, StatusScheduled, StatusPassed, StageOnboarding, StatusScheduled, ErrIllegalTransition},

		{"illegal transition returns input state unchanged", StageHRInterview, StatusOffered, StatusPassed, StageHRInterview, StatusOffered, ErrIllegalTransition},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStage, gotStatus, err := NextState(tt.stage, tt.current, tt.requested)
			if err != tt.wantErr {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if gotStage != tt.wantStage || gotStatus != tt.wantState {
				t.Errorf("got (%s, %s), want (%s, %s)", gotStage, gotStatus, tt.wantStage, tt.wantState)
			}
		})
	}
}
