package common

type CrowdSource struct {
	CasesTimeSeries []CasesTimeSeries `json:"cases_time_series"`
	Statewise       []Statewise       `json:"statewise"`
	Tested          []Tested          `json:"tested"`
}
type CasesTimeSeries struct {
	Dailyconfirmed string `json:"dailyconfirmed"`
	Dailydeceased  string `json:"dailydeceased"`
	Dailyrecovered string `json:"dailyrecovered"`
	Date           string `json:"date"`
	Totalconfirmed string `json:"totalconfirmed"`
	Totaldeceased  string `json:"totaldeceased"`
	Totalrecovered string `json:"totalrecovered"`
}
type Delta struct {
	Active    string `json:"active"`
	Confirmed string `json:"confirmed"`
	Deaths    string `json:"deaths"`
	Recovered string `json:"recovered"`
}
type Statewise struct {
	Active          string `json:"active"`
	Confirmed       string `json:"confirmed"`
	Deaths          string `json:"deaths"`
	Deltaconfirmed  string `json:"deltaconfirmed"`
	Deltadeaths     string `json:"deltadeaths"`
	Deltarecovered  string `json:"deltarecovered"`
	Lastupdatedtime string `json:"lastupdatedtime"`
	Recovered       string `json:"recovered"`
	State           string `json:"state"`
	Statecode       string `json:"statecode"`
}
type Tested struct {
	Source                      string `json:"source"`
	Testsconductedbyprivatelabs string `json:"testsconductedbyprivatelabs"`
	Totalindividualstested      string `json:"totalindividualstested"`
	Totalpositivecases          string `json:"totalpositivecases"`
	Totalsamplestested          string `json:"totalsamplestested"`
	Updatetimestamp             string `json:"updatetimestamp"`
}

type StateDistrict struct {
	DistrictData map[string]PerDistrict `json:"districtData"`
}

type PerDistrict struct {
	Confirmed       int32  `json:"confirmed"`
	Lastupdatedtime string `json:"lastupdatedtime"`
	Delta           Delta  `json:"delta"`
}

type CovidDistrict struct {
	Name            string
	Confirmed       int32  `json:"confirmed"`
	Lastupdatedtime string `json:"lastupdatedtime"`
	Delta           Delta  `json:"delta"`
}
