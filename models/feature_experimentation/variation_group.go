package feature_experimentation

type VariationGroup struct {
	Id               string            `json:"id,omitempty"`
	Name             string            `json:"name"`
	Variations       []*VariationFE    `json:"variations"`
	Targeting        *Targeting        `json:"targeting"`
	AllocationConfig *AllocationConfig `json:"allocation_config,omitempty"`
}

type Targeting struct {
	TargetingGroups []*TargetingGroup `json:"targeting_groups"`
}

type TargetingGroup struct {
	Targetings []*InnerTargeting `json:"targetings"`
}

type InnerTargeting struct {
	Key      string      `json:"key"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type AllocationConfig struct {
	StartDate       string         `json:"start_date"`
	Timezone        string         `json:"timezone"`
	StartAllocation float64        `json:"start_allocation"`
	PeriodicSteps   *PeriodicSteps `json:"periodic_steps"`
}

type PeriodicSteps struct {
	Allocation float64 `json:"allocation"`
	Step       int     `json:"step"`
	Step_type  string  `json:"step_type"`
}
