package common

type CoronaUpdate struct {
	UpdateTime        string       `json:"update_time"`
	Total             string       `json:"total"`
	Active            string       `json:"active"`
	Cured             string       `json:"cured"`
	Death             string       `json:"death"`
	Delta             KeyValues    `json:"delta"`
	Closed            string       `json:"closed_case"`
	FatalPercent      string       `json:"fatal_percent"`
	LivePercent       string       `json:"live_percent"`
	Migrated          string       `json:"migrated"`
	ScreenedAtAirport string       `json:"screened_at_airport"`
	StateWise         []*StateData `json:"state_wise"`
	DistrictWise      string       `json:"district_wise"`
	HelpLine          string       `json:"help_line"`
	Faq               string       `json:"faq"`
	Youtube           string       `json:"youtube"`
	Facebook          string       `json:"facebook"`
	Twitter           string       `json:"twitter"`
	Links             []string     `json:"links"`
}

type StateData struct {
	Id           string           `json:"id"`
	Name         string           `json:"name"`
	Total        string           `json:"total"`
	Active       string           `json:"active"`
	LiveExit     string           `json:"exit"`
	Death        string           `json:"death"`
	Closed       string           `json:"closed_case"`
	FatalPercent string           `json:"fatal_percent"`
	LivePercent  string           `json:"live_percent"`
	District     []*CovidDistrict `json:"district"`
}
