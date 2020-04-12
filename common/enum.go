package common

const (
	IK_HIGHEST_CONFIRMED_STATE     string = "hcs"
	IK_HIGHEST_CONFIRMED_NUMBER    string = "hcn"
	IK_HIGHEST_CURED_STATE         string = "hds"
	IK_HIGHEST_CURED_NUMBER        string = "hdn"
	IK_HIGHEST_FATAL_STATE         string = "hfs"
	IK_HIGHEST_FATAL_NUMBER        string = "hfn"
	IK_HIGHEST_PI_CONFIRMED_DATE   string = "hpicd"
	IK_HIGHEST_PI_CONFIRMED_NUMBER string = "hpicn"
	IK_HIGHEST_PI_CURED_DATE       string = "hpidd"
	IK_HIGHEST_PI_CURED_NUMBER     string = "hpidn"
	IK_HIGHEST_PI_FATAL_DATE       string = "hpifd"
	IK_HIGHEST_PI_FATAL_NUMBER     string = "hpifn"
)

type InsightType string

const (
	IT_STATE   InsightType = "state"
	IT_COUNTRY InsightType = "country"
)
