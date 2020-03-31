package common

type CrowdSource struct {
	CasesTimeSeries []CasesTimeSeries `json:"cases_time_series"`
	KeyValues       []KeyValues       `json:"key_values"`
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
type KeyValues struct {
	Confirmeddelta           string `json:"confirmeddelta"`
	Counterforautotimeupdate string `json:"counterforautotimeupdate"`
	Deceaseddelta            string `json:"deceaseddelta"`
	Lastupdatedtime          string `json:"lastupdatedtime"`
	Recovereddelta           string `json:"recovereddelta"`
	Statesdelta              string `json:"statesdelta"`
}
type Delta struct {
	Active    int `json:"active"`
	Confirmed int `json:"confirmed"`
	Deaths    int `json:"deaths"`
	Recovered int `json:"recovered"`
}
type Statewise struct {
	Active          string `json:"active"`
	Confirmed       string `json:"confirmed"`
	Deaths          string `json:"deaths"`
	Delta           Delta  `json:"delta"`
	Lastupdatedtime string `json:"lastupdatedtime"`
	Recovered       string `json:"recovered"`
	State           string `json:"state"`
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

type CovidState struct {
	Name      string
	Id        string
	Districts []*CovidDistrict
}

type CovidDistrict struct {
	Name            string
	Confirmed       int32  `json:"confirmed"`
	Lastupdatedtime string `json:"lastupdatedtime"`
	Delta           Delta  `json:"delta"`
}
