package common

type CoronaUpdate struct {
	UpdateTime        string       `json:"update_time"`
	Total             string       `json:"total"`
	Active            string       `json:"active"`
	Cured             string       `json:"cured"`
	Death             string       `json:"death"`
	Migrated          string       `json:"migrated"`
	ScreenedAtAirport string       `json:"screened_at_airport"`
	StateWise         []*StateData `json:"state_wise"`
	Youtube           string       `json:"youtube"`
	DistrictWise      string       `json:"district_wise"`
	HelpLine          string       `json:"help_line"`
	Faq               string       `json:"faq"`
	Links             []string     `json:"links"`
}

type StateData struct {
	Name         string `json:"name"`
	Total        string `json:"total"`
	TotalIndian  string `json:"total_indian"`
	TotalForeign string `json:"total_foreign"`
	LiveExit     string `json:"exit"`
	Death        string `json:"death"`
}
