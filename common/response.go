package common

type CoronaUpdate struct {
	UpdateTime        string
	Total             string
	Active            string
	Cured             string
	Death             string
	Delta             *Delta
	Closed            string
	FatalPercent      string
	LivePercent       string
	Migrated          string
	ScreenedAtAirport string
	StateWise         []*StateData
	DistrictWise      string
	SampleTested      string
	HelpLine          string
	Faq               string
	Youtube           string
	Facebook          string
	Twitter           string
	Links             []string
	SeriesDelta       *CovidDelta
	Tested            *CoronaTest
}

type CoronaTest struct {
	Date        []string
	TotalSample []int
}

type StateData struct {
	Id           string
	Name         string
	Code         string
	Color        string
	Display      string
	Total        string
	Active       string
	LiveExit     string
	Death        string
	Closed       string
	FatalPercent string
	LivePercent  string
	Delta        *Delta
	District     []*CovidDistrict
}

type CovidDelta struct {
	Dates     []string
	Confirmed []string
	Cured     []string
	Death     []string
}
