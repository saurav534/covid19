package external

import (
	"covid19/common"
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

var crowdData *common.CoronaUpdate

func GetCrowdData() *common.CoronaUpdate {
	return crowdData
}

func UpdateCrowdData() {
	log.Printf("updating crowd Data starting at: %v", time.Now().Format("02 Jan, 03:04:05 PM"))
	dataChannel := make(chan map[string]*common.StateDistrict)
	go func() {
		crowdDistrictData(dataChannel)
	}()

	constantBackoff := backoff.NewConstantBackOff(500 * time.Millisecond)
	var resp *http.Response
	var err error
	err = backoff.Retry(func() error {
		resp, err = http.Get("https://api.covid19india.org/data.json")
		log.Printf("Complete crowd data response recieved")
		return err
	}, backoff.WithMaxRetries(constantBackoff, 4))

	// init insight
	insights := make(map[string]string)
	insights[common.IK_HIGHEST_CONFIRMED_NUMBER] = "0"
	insights[common.IK_HIGHEST_CURED_NUMBER] = "0"
	insights[common.IK_HIGHEST_FATAL_NUMBER] = "0"
	insights[common.IK_HIGHEST_PI_CONFIRMED_NUMBER] = "0"
	insights[common.IK_HIGHEST_PI_CURED_NUMBER] = "0"
	insights[common.IK_HIGHEST_PI_FATAL_NUMBER] = "0"

	update := &common.CoronaUpdate{
		StateWise: make([]*common.StateData, 0),
		Links:     make([]string, 0),
		Insights:  insights,
	}
	if err != nil {
		log.Printf("Error while updating crowd data, %v", err)
		return
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	cs := &common.CrowdSource{}
	_ = json.Unmarshal(bytes, cs)

	newsChan := make(chan *common.News)
	go func() {
		latestNews(newsChan)
	}()

	// set time series data
	covidDeltaChan := make(chan *common.CovidDelta)
	go func() {
		seriesDelta(covidDeltaChan, cs.CasesTimeSeries)
	}()

	for _, cts := range cs.CasesTimeSeries {
		updateInsights(common.IT_COUNTRY, nil, cts, update.Insights)
	}

	// set covid test data
	tests := cs.Tested
	coronaTestChan := make(chan *common.CoronaTest)
	go func() {
		coronaTested(coronaTestChan, tests)
	}()

	for i := len(cs.Tested) - 1; i >= 0; i-- {
		if strings.Trim(cs.Tested[i].Totalsamplestested, " ") != "" {
			update.SampleTested = cs.Tested[i].Totalsamplestested
			break
		}
	}

	// parse district data first
	stateDistrictData := <-dataChannel
	stateToDistrictList := toStateToDistrictList(stateDistrictData)

	for _, data := range cs.Statewise {
		if data.State == "Total" {
			setPanIndiaData(update, data)
		} else if data.Confirmed != "0" {
			st := getPerStateData(data, stateToDistrictList)
			updateInsights(common.IT_STATE, data, nil, update.Insights)
			update.StateWise = append(update.StateWise, st)
		}
	}
	update.SeriesDelta = <-covidDeltaChan
	tested, ok := <-coronaTestChan
	if ok {
		update.Tested = tested
		update.Tested.ConfirmRate = getConfirmRate(update)
	}

	if mohData != nil {
		update.Facebook = mohData.Facebook
		update.Youtube = mohData.Youtube
		update.Twitter = mohData.Twitter
		update.HelpLine = mohData.HelpLine
		update.Faq = mohData.Faq
		update.Links = mohData.Links
	}
	update.News = <-newsChan
	sortStates(update.StateWise)
	crowdData = update
	log.Printf("updating crowd data done at: %v", time.Now().Format("02 Jan, 03:04:05 PM"))
}

func toStateToDistrictList(stateDistrictData map[string]*common.StateDistrict) map[string][]*common.CovidDistrict {
	stateToDistrictList := make(map[string][]*common.CovidDistrict)
	for _, value := range reflect.ValueOf(stateDistrictData).MapKeys() {
		if value.String() == "Unkown" {
			continue
		}
		stateDistricts := stateDistrictData[value.String()]
		districts := stateDistricts.DistrictData
		covidDistrict := make([]*common.CovidDistrict, 0)
		for _, val := range reflect.ValueOf(districts).MapKeys() {
			districtData := districts[val.String()]
			covidDistrict = append(covidDistrict, &common.CovidDistrict{
				Name:            val.String(),
				Confirmed:       districtData.Confirmed,
				Lastupdatedtime: districtData.Lastupdatedtime,
				Delta: common.Delta{
					Confirmed: strconv.Itoa(districtData.Delta.Confirmed),
				},
			})
		}
		sort.Slice(covidDistrict, func(i, j int) bool {
			return covidDistrict[i].Confirmed > covidDistrict[j].Confirmed
		})
		stateToDistrictList[value.String()] = covidDistrict
	}
	return stateToDistrictList
}

func getPerStateData(data *common.Statewise, stateToDistrictList map[string][]*common.CovidDistrict) *common.StateData {
	st := &common.StateData{}
	st.Death = data.Deaths
	st.Total = data.Confirmed
	st.Active = data.Active
	st.LiveExit = data.Recovered
	st.Name = data.State
	st.Delta = &common.Delta{
		Confirmed: data.Deltaconfirmed,
		Deaths:    data.Deltadeaths,
		Recovered: data.Deltarecovered,
	}
	st.Code = common.StateCode[strings.ToLower(st.Name)]
	st.Color = common.GetInfectColor(int32(toInt(st.Total)))
	st.Display = st.Name + " - " + st.Total
	st.Id = strings.ReplaceAll(st.Name, " ", "-")
	st.District = stateToDistrictList[st.Name]
	cured := toInt(data.Recovered)
	death := toInt(data.Deaths)
	closed := cured + death
	st.Closed = strconv.Itoa(closed)
	if closed != 0 {
		st.FatalPercent = fmt.Sprint(math.Round((float64(death*100))/float64(closed))) + "%"
		st.LivePercent = fmt.Sprint(math.Round((float64((cured)*100))/float64(closed))) + "%"
	} else {
		st.FatalPercent = "NA"
		st.LivePercent = "NA"
	}
	return st
}

func updateInsights(it common.InsightType, data *common.Statewise, cdata *common.CasesTimeSeries, insights map[string]string) {
	if common.IT_STATE == it {
		hc, _ := strconv.Atoi(insights[common.IK_HIGHEST_CONFIRMED_NUMBER])
		hd, _ := strconv.Atoi(insights[common.IK_HIGHEST_CURED_NUMBER])
		hf, _ := strconv.Atoi(insights[common.IK_HIGHEST_FATAL_NUMBER])

		nhc, _ := strconv.Atoi(data.Confirmed)
		nhd, _ := strconv.Atoi(data.Recovered)
		nhf, _ := strconv.Atoi(data.Deaths)

		if nhc > hc {
			insights[common.IK_HIGHEST_CONFIRMED_STATE] = data.State
			insights[common.IK_HIGHEST_CONFIRMED_NUMBER] = data.Confirmed
		}
		if nhd > hd {
			insights[common.IK_HIGHEST_CURED_STATE] = data.State
			insights[common.IK_HIGHEST_CURED_NUMBER] = data.Recovered
		}
		if nhf > hf {
			insights[common.IK_HIGHEST_FATAL_STATE] = data.State
			insights[common.IK_HIGHEST_FATAL_NUMBER] = data.Deaths
		}
	}
	if common.IT_COUNTRY == it {
		hc, _ := strconv.Atoi(insights[common.IK_HIGHEST_PI_CONFIRMED_NUMBER])
		hd, _ := strconv.Atoi(insights[common.IK_HIGHEST_PI_CURED_NUMBER])
		hf, _ := strconv.Atoi(insights[common.IK_HIGHEST_PI_FATAL_NUMBER])

		nhc, _ := strconv.Atoi(cdata.Dailyconfirmed)
		nhd, _ := strconv.Atoi(cdata.Dailyrecovered)
		nhf, _ := strconv.Atoi(cdata.Dailydeceased)

		if nhc > hc {
			insights[common.IK_HIGHEST_PI_CONFIRMED_DATE] = cdata.Date[:6]
			insights[common.IK_HIGHEST_PI_CONFIRMED_NUMBER] = cdata.Dailyconfirmed
		}
		if nhd > hd {
			insights[common.IK_HIGHEST_PI_CURED_DATE] = cdata.Date[:6]
			insights[common.IK_HIGHEST_PI_CURED_NUMBER] = cdata.Dailyrecovered
		}
		if nhf > hf {
			insights[common.IK_HIGHEST_PI_FATAL_DATE] = cdata.Date[:6]
			insights[common.IK_HIGHEST_PI_FATAL_NUMBER] = cdata.Dailydeceased
		}
	}
}

func setPanIndiaData(update *common.CoronaUpdate, data *common.Statewise) {
	update.Total = data.Confirmed
	update.Active = data.Active
	update.Cured = data.Recovered
	update.Death = data.Deaths
	update.Delta = &common.Delta{
		Confirmed: data.Deltaconfirmed,
		Deaths:    data.Deltadeaths,
		Recovered: data.Deltarecovered,
	}
	updateTime, _ := time.Parse("02/01/2006 15:04:05", data.Lastupdatedtime)
	update.UpdateTime = updateTime.Format("02 Jan, 03:04 PM")
	cured := toInt(data.Recovered)
	death := toInt(data.Deaths)
	closed := cured + death
	update.Closed = strconv.Itoa(closed)
	if closed != 0 {
		update.FatalPercent = fmt.Sprintf("%.2f", (float32(death*100))/float32(closed))
		update.LivePercent = fmt.Sprintf("%.2f", (float32((cured)*100))/float32(closed))
	} else {
		update.FatalPercent = "0"
		update.LivePercent = "0"
	}
}

func getConfirmRate(update *common.CoronaUpdate) []float32 {
	dc := update.SeriesDelta.DateToConfirmed
	confirmRate := make([]float32, 0)
	for i, d := range update.Tested.Date {
		c := dc[d]
		if i == len(update.Tested.Date)-1 && c == 0 {
			newConfirmed, _ := strconv.Atoi(update.Delta.Confirmed)
			c = newConfirmed
		}
		rate := (float32(c) * 100.0) / float32(update.Tested.TotalSample[i])
		confirmRate = append(confirmRate, float32(math.Round(float64(rate)*100)/100))
	}
	return confirmRate
}

func latestNews(newsChannel chan<- *common.News) {
	defer close(newsChannel)
	constantBackoff := backoff.NewConstantBackOff(500 * time.Millisecond)
	var resp *http.Response
	var err error
	err = backoff.Retry(func() error {
		resp, err = http.Get("http://newsapi.org/v2/top-headlines?country=in&q=corona&apiKey=cd3cace91bf849509f00076f46c7f62e")
		log.Printf("News headline response recieved")
		return err
	}, backoff.WithMaxRetries(constantBackoff, 4))

	news := common.News{}
	if err != nil {
		log.Printf("error while calling to news headline service, %v", err)
		newsChannel <- &news
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bytes, &news)
	if err != nil {
		log.Printf("error happend %v", err)
	}
	newsChannel <- &news
}

func crowdDistrictData(dataChannel chan<- map[string]*common.StateDistrict) {
	defer close(dataChannel)
	constantBackoff := backoff.NewConstantBackOff(500 * time.Millisecond)
	var resp *http.Response
	var err error
	err = backoff.Retry(func() error {
		resp, err = http.Get("https://api.covid19india.org/state_district_wise.json")
		log.Printf("District Wise response recieved")
		return err
	}, backoff.WithMaxRetries(constantBackoff, 4))

	dd := make(map[string]*common.StateDistrict)
	if err != nil {
		log.Printf("error while calling to crowd district data service, %v", err)
		dataChannel <- dd
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bytes, &dd)
	if err != nil {
		log.Printf("error happend %v", err)
	}
	dataChannel <- dd
}

func seriesDelta(seriesChan chan<- *common.CovidDelta, seriesData []*common.CasesTimeSeries) {
	defer close(seriesChan)
	l := len(seriesData)
	dates := make([]string, 0)
	confirmed := make([]string, 0)
	cured := make([]string, 0)
	death := make([]string, 0)
	dateToConfirmed := make(map[string]int)
	dateToday := time.Now().Format("02 Jan")
	if strings.Contains(seriesData[l-1].Date, dateToday) {
		l = l - 1
	}
	i := l - 21
	if i < 0 {
		i = 0
	}
	for ; i < l; i++ {
		data := seriesData[i]
		dateString := data.Date[:6]
		dates = append(dates, dateString)
		confirmed = append(confirmed, data.Dailyconfirmed)
		cured = append(cured, data.Dailyrecovered)
		death = append(death, data.Dailydeceased)
		dc, _ := strconv.Atoi(data.Dailyconfirmed)
		dateToConfirmed[dateString] = dc
	}

	seriesChan <- &common.CovidDelta{
		Dates:           dates,
		Confirmed:       confirmed,
		Cured:           cured,
		Death:           death,
		DateToConfirmed: dateToConfirmed,
	}
}

func coronaTested(covidTestChan chan<- *common.CoronaTest, tested []*common.Tested) {
	defer close(covidTestChan)
	testMap := make(map[string]int)
	testMapRev := make(map[int]string)
	cumTest := make([]int, 0)
	for _, test := range tested {
		st := strings.Trim(test.Totalsamplestested, " ")
		if st != "" {
			date := strings.Split(test.Updatetimestamp, " ")[0]
			i, err := strconv.Atoi(st)
			if err == nil {
				testMap[date] = i
			}
		}
	}

	for k, v := range testMap {
		cumTest = append(cumTest, v)
		testMapRev[v] = k
	}

	sort.Ints(cumTest)
	testLen := len(cumTest)
	dates := make([]string, 0)
	sampleTested := make([]int, 0)
	i := testLen - 8
	if i < 0 {
		i = 0
	}
	for ; i < testLen; i++ {
		updateTime, _ := time.Parse("2/1/2006", testMapRev[cumTest[i]])
		dates = append(dates, updateTime.Format("02 Jan"))
		sampleTested = append(sampleTested, cumTest[i]-cumTest[i-1])
		updateTime.Format("02 Jan")
	}
	covidTestChan <- &common.CoronaTest{
		Date:        dates,
		TotalSample: sampleTested,
	}
}
