package models

import "time"

type RCA struct {
	ID                 int       `json:"id"`
	WorkItemID         int       `json:"work_item_id"`
	RootCauseCategory  string    `json:"root_cause_category" binding:"required"`
	FixApplied         string    `json:"fix_applied" binding:"required"`
	PreventionSteps    string    `json:"prevention_steps"`
	IncidentStart      time.Time `json:"incident_start" binding:"required"`
	IncidentEnd        time.Time `json:"incident_end" binding:"required"`
	MTTRMinutes        int       `json:"mttr_minutes"`
	CreatedAt          time.Time `json:"created_at"`
}
