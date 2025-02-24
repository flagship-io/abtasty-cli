package feature_experimentation

type CampaignFE struct {
	Id              string           `json:"id,omitempty"`
	ProjectId       string           `json:"project_id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	Type            string           `json:"type,omitempty"`
	Status          string           `json:"status,omitempty"`
	State           string           `json:"state,omitempty"`
	Slug            string           `json:"slug,omitempty"`
	VariationGroups []VariationGroup `json:"variation_groups"`
	Scheduler       *Scheduler       `json:"scheduler,omitempty"`
	PrimaryGoal     *Goal            `json:"primary_goal,omitempty"`
	CreatedAt       string           `json:"created_at,omitempty"`
	UpdatedAt       string           `json:"updated_at,omitempty"`
	SecondaryGoals  []Goal           `json:"secondary_goals,omitempty"`
}

type Scheduler struct {
	StartDate string `json:"start_date"`
	StopDate  string `json:"stop_date"`
	TimeZone  string `json:"timezone"`
}

type CampaignFESwitchRequest struct {
	State string `json:"state"`
}
