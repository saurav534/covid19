package common

var StateCode map[string]string

func init() {
	StateCode = map[string]string{
		"maharashtra":                 "MH",
		"kerala":                      "KL",
		"tamil nadu":                  "TN",
		"delhi":                       "DL",
		"rajasthan":                   "RJ",
		"uttar pradesh":               "UP",
		"andhra pradesh":              "AP",
		"karnataka":                   "KA",
		"telangana":                   "TL",
		"telengana":                   "TL",
		"gujarat":                     "GJ",
		"madhya pradesh":              "MP",
		"jammu and kashmir":           "JK",
		"punjab":                      "PB",
		"haryana":                     "HR",
		"west bengal":                 "WB",
		"bihar":                       "BR",
		"chandigarh":                  "CH",
		"ladakh":                      "",
		"assam":                       "AS",
		"andaman and nicobar islands": "AN",
		"chhattisgarh":                "CT",
		"uttarakhand":                 "UT",
		"uttaranchal":                 "UT",
		"goa":                         "GA",
		"odisha":                      "OR",
		"himachal pradesh":            "HP",
		"puducherry":                  "PY",
		"jharkhand":                   "JH",
		"manipur":                     "MN",
		"mizoram":                     "MZ",
		"arunachal pradesh":           "AR",
		"tripura":                     "TR",
		"sikkim":                      "SK",
		"meghalaya":                   "ML",
		"nagaland":                    "NL",
		"dadra and nagar haveli":      "DN",
	}
}

func GetInfectColor(cnf int32) string {
	if cnf > 60000 {
		return "#8b0000"
	}
	if cnf > 50000 {
		return "#960f0f"
	}
	if cnf > 40000 {
		return "#a01e1e"
	}
	if cnf > 30000 {
		return "#a92929"
	}
	if cnf > 20000 {
		return "#b33c3c"
	}
	if cnf > 10000 {
		return "#bb4d4d"
	}
	if cnf > 5000 {
		return "#bf5b5b"
	}
	if cnf > 4000 {
		return "#c76b6b"
	}
	if cnf > 3000 {
		return "#d07c7c"
	}
	if cnf > 2000 {
		return "#d48787"
	}
	if cnf > 1000 {
		return "#d89696"
	}
	if cnf > 500 {
		return "#dca1a1"
	}
	if cnf > 0 {
		return "#e2b1b1"
	}
	return "#fff"
}
