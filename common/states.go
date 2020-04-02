package common

var StateCode map[string]string

func init() {
	StateCode = map[string]string{
		"maharashtra" : "MH",
		"kerala" : "KL",
		"tamil nadu" : "TN",
		"delhi" : "DL",
		"rajasthan" : "RJ",
		"uttar pradesh" : "UP",
		"andhra pradesh" : "AP",
		"karnataka" : "KA",
		"telangana" : "TL",
		"telengana" : "TL",
		"gujarat" : "GJ",
		"madhya pradesh" : "MP",
		"jammu and kashmir" : "JK",
		"punjab" : "PB",
		"haryana" : "HR",
		"west bengal" : "WB",
		"bihar" : "BR",
		"chandigarh" : "CH",
		"ladakh" : "",
		"assam" : "AS",
		"andaman and nicobar islands" : "AN",
		"chhattisgarh" : "CT",
		"uttarakhand" : "UT",
		"uttaranchal" : "UT",
		"goa" : "GA",
		"odisha" : "OR",
		"himachal pradesh" : "HP",
		"puducherry" : "PY",
		"jharkhand" : "JH",
		"manipur" : "MN",
		"mizoram" : "MZ",
		"arunachal pradesh" : "AR",
	}
}

func GetInfectColor(cnf int32) string{
	if cnf >= 300 {
		return "#174ea6"
	}
	if cnf >= 270 {
		return "#1c52a9"
	}
	if cnf >= 250 {
		return "#2558a9"
	}
	if cnf >= 225 {
		return "#3160ab"
	}
	if cnf >= 200 {
		return "#406cb1"
	}
	if cnf >= 175 {
		return "#4d73af"
	}
	if cnf >= 150 {
		return "#5c7bad"
	}
	if cnf >= 125 {
		return "#718cb9"
	}
	if cnf >= 100 {
		return "#8ba2c7"
	}
	if cnf >= 75 {
		return "#a0b2d0"
	}
	if cnf >=50 {
		return "#bcc9de"
	}
	if cnf >=25 {
		return "#ced6e4"
	}
	if cnf >0 {
		return "#dee2ea"
	}
	return "#fff"
}
