package main

import "sort"

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
