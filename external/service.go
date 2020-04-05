package external

import (
	"covid19/common"
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"golang.org/x/net/html"
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

var updateTime int64
var crowdUpdateTime int64
var data *common.CoronaUpdate
var crowdData *common.CoronaUpdate

func LiveData() *common.CoronaUpdate {
	millisNow := time.Now().UnixNano() / 1000000
	diff := (millisNow - updateTime) / (60 * 1000)

	if diff <= 5 && data != nil {
		log.Printf("MOH data from cache found %v", diff)
		return data
	}

	constantBackoff := backoff.NewConstantBackOff(500 * time.Millisecond)
	var resp *http.Response
	var err error
	err = backoff.Retry(func() error {
		resp, err = http.Get("https://www.mohfw.gov.in/")
		return err
	}, backoff.WithMaxRetries(constantBackoff, 4))

	stateData := make([]*common.StateData, 0)
	links := make([]string, 0)
	linkMap := make(map[string]bool)

	update := &common.CoronaUpdate{
		StateWise: stateData,
		Links:     links,
	}
	if err != nil {
		log.Printf("error while calling to gov service, %v", err)
		return update
	}

	node, _ := html.Parse(resp.Body)
	defer resp.Body.Close()
	parseNode(node, update, linkMap)
	for _, value := range reflect.ValueOf(linkMap).MapKeys() {
		update.Links = append(update.Links, value.String())
	}
	cured := toInt(update.Cured)
	death := toInt(update.Death)
	closed := cured + death
	update.Total = strconv.Itoa(toInt(update.Active) + cured +
		toInt(update.Migrated) + death)

	update.Closed = strconv.Itoa(cured + death)
	if closed != 0 {
		update.FatalPercent = fmt.Sprintf("%.2f", (float32(death*100))/float32(closed))
		update.LivePercent = fmt.Sprintf("%.2f", (float32((cured)*100))/float32(closed))
	} else {
		update.FatalPercent = "0"
		update.LivePercent = "0"
	}

	// updating cache time
	updateTime = millisNow
	data = update

	return update
}

func CrowdData() *common.CoronaUpdate {
	millisNow := time.Now().UnixNano() / 1000000
	diff := (millisNow - crowdUpdateTime) / (60 * 1000)

	if diff <= 5 && crowdData != nil {
		log.Printf("Crowd data from cache found %v", diff)
		return crowdData
	}

	dataChannel := make(chan map[string]*common.StateDistrict)
	go func() {
		crowdDistrictData(dataChannel)
	}()

	constantBackoff := backoff.NewConstantBackOff(500 * time.Millisecond)
	var resp *http.Response
	var err error
	err = backoff.Retry(func() error {
		resp, err = http.Get("https://api.covid19india.org/data.json")
		log.Printf("Complete data response recieved")
		return err
	}, backoff.WithMaxRetries(constantBackoff, 4))

	stateData := make([]*common.StateData, 0)
	links := make([]string, 0)

	update := &common.CoronaUpdate{
		StateWise: stateData,
		Links:     links,
	}
	if err != nil {
		log.Printf("error while calling to crowd data service, %v", err)
		return update
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	cs := &common.CrowdSource{}
	_ = json.Unmarshal(bytes, cs)

	// set time series data
	covidDeltaChan := make(chan *common.CovidDelta)
	go func() {
		seriesDelta(covidDeltaChan, cs.CasesTimeSeries)
	}()

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
	stataDistrictData := <-dataChannel
	stateToDistrictList := make(map[string][]*common.CovidDistrict)

	for _, value := range reflect.ValueOf(stataDistrictData).MapKeys() {
		if value.String() == "Unkown" {
			continue
		}
		stateDistricts := stataDistrictData[value.String()]
		districts := stateDistricts.DistrictData
		covidDistrict := make([]*common.CovidDistrict, 0)
		for _, val := range reflect.ValueOf(districts).MapKeys() {
			districtData := districts[val.String()]
			covidDistrict = append(covidDistrict, &common.CovidDistrict{
				Name:            val.String(),
				Confirmed:       districtData.Confirmed,
				Lastupdatedtime: districtData.Lastupdatedtime,
				Delta:           districtData.Delta,
			})
		}
		sort.Slice(covidDistrict, func(i, j int) bool {
			return covidDistrict[i].Confirmed > covidDistrict[j].Confirmed
		})
		stateToDistrictList[value.String()] = covidDistrict
	}

	for _, data := range cs.Statewise {
		if data.State == "Total" {
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
		} else if data.Confirmed != "0" {
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
			update.StateWise = append(update.StateWise, st)
		}
	}
	update.SeriesDelta = <-covidDeltaChan
	update.Tested = <-coronaTestChan
	if data != nil {
		update.Facebook = data.Facebook
		update.Youtube = data.Youtube
		update.Twitter = data.Twitter
		update.HelpLine = data.HelpLine
		update.Faq = data.Faq
		update.Links = data.Links
	}

	// updating cache time
	crowdUpdateTime = millisNow
	crowdData = update

	return update
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

	if err != nil {
		log.Printf("error while calling to crowd district data service, %v", err)
		dataChannel <- nil
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	dd := make(map[string]*common.StateDistrict)
	err = json.Unmarshal(bytes, &dd)
	if err != nil {
		log.Printf("error happend %v", err)
	}
	dataChannel <- dd
}

func seriesDelta(seriesChan chan<- *common.CovidDelta, seriesData []common.CasesTimeSeries) {
	defer close(seriesChan)
	l := len(seriesData)
	dates := make([]string, 0)
	confirmed := make([]string, 0)
	cured := make([]string, 0)
	death := make([]string, 0)

	for i := l - 21; i < l; i++ {
		data := seriesData[i]
		dates = append(dates, data.Date[:6])
		confirmed = append(confirmed, data.Dailyconfirmed)
		cured = append(cured, data.Dailyrecovered)
		death = append(death, data.Dailydeceased)
	}

	seriesChan <- &common.CovidDelta{
		Dates:     dates,
		Confirmed: confirmed,
		Cured:     cured,
		Death:     death,
	}
}

func coronaTested(covidTestChan chan<- *common.CoronaTest, tested []common.Tested) {
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
	for i := testLen - 6; i < testLen; i++ {
		updateTime, _ := time.Parse("2/1/2006", testMapRev[cumTest[i]])
		dates = append(dates, updateTime.Format("02 Jan"))
		sampleTested = append(sampleTested, cumTest[i] - cumTest[i-1])
		updateTime.Format("02 Jan")
	}
	covidTestChan <- &common.CoronaTest{
		Date:        dates,
		TotalSample: sampleTested,
	}
}

func parseNode(node *html.Node, update *common.CoronaUpdate, linkMap map[string]bool) {
	if strings.Contains(node.Data, "as on :") {
		split := strings.Split(node.Data, " GMT")
		updateTime := strings.TrimPrefix(split[0], "as on : ")

		ut, _ := time.Parse("02 January 2006, 15:04", updateTime)
		update.UpdateTime = ut.Format("02 Jan, 03:04 PM")
		return
	}

	if isElementNode(node, "a") {
		url := fetchUrl(node)
		if url != "" {
			if strings.Contains(url, "youtube") {
				update.Youtube = url
			} else if strings.Contains(url, "twitter") {
				update.Twitter = url
			} else if strings.Contains(url, "facebook") {
				update.Facebook = url
			} else if strings.Contains(url, "District") {
				update.DistrictWise = url
			} else if strings.Contains(url, "FAQ") {
				update.Faq = url
			} else if strings.Contains(url, "helpline") {
				update.HelpLine = url
			} else {
				linkMap[url] = true
			}
		}
		return
	}

	if rowHasClass(node, "li", "bg-blue") {
		setStatValues(node, &update.Active)
		return
	}
	if rowHasClass(node, "li", "bg-green") {
		setStatValues(node, &update.Cured)
		return
	}
	if rowHasClass(node, "li", "bg-red") {
		setStatValues(node, &update.Death)
		return
	}
	if rowHasClass(node, "li", "bg-orange") {
		setStatValues(node, &update.Migrated)
		return
	}

	if isElementNode(node, "table") {
		for c1 := node.FirstChild; c1 != nil; c1 = c1.NextSibling {
			if isElementNode(c1, "thead") && !containText(c1, "Name of State") {
				return
			}
			if isElementNode(c1, "tbody") {
				setStateWise(c1, update)
				return
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		parseNode(child, update, linkMap)
	}

}

func setStateWise(tbody *html.Node, update *common.CoronaUpdate) {
	for tr := tbody.FirstChild; tr != nil; tr = tr.NextSibling {
		if isElementNode(tr, "tr") {
			st := &common.StateData{}
			index := 0
			ignore := false
			for td := tr.FirstChild; td != nil; td = td.NextSibling {
				if isElementNode(td, "td") {
					if td.FirstChild.Data == "strong" {
						ignore = true
						break
					}
					switch index {
					case 0:
					case 1:
						st.Name = td.FirstChild.Data
					case 2:
						st.Total = td.FirstChild.Data
					case 3:
						st.LiveExit = td.FirstChild.Data
					case 4:
						st.Death = td.FirstChild.Data
					}
					index = index + 1
				}
			}
			if !ignore {
				live := toInt(st.LiveExit)
				death := toInt(st.Death)
				closed := live + death
				st.Active = strconv.Itoa(toInt(st.Total) - closed)
				st.Closed = strconv.Itoa(closed)
				if closed != 0 {
					st.FatalPercent = fmt.Sprint(math.Round((float64(death*100))/float64(closed))) + "%"
					st.LivePercent = fmt.Sprint(math.Round((float64((live)*100))/float64(closed))) + "%"
				} else {
					st.FatalPercent = "NA"
					st.LivePercent = "NA"
				}
				st.Code = common.StateCode[strings.ToLower(st.Name)]
				st.Color = common.GetInfectColor(int32(toInt(st.Total)))
				st.Display = st.Name + " - " + st.Total
				update.StateWise = append(update.StateWise, st)
			}
		}
	}
}

func containText(node *html.Node, text string) bool {
	if strings.Contains(node.Data, text) {
		return true
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if containText(child, text) {
			return true
		}
	}
	return false;
}

func isElementNode(node *html.Node, nt string) bool {
	return node.Type == html.ElementNode && node.Data == nt;
}

func setValues(node *html.Node, update *common.CoronaUpdate) {
	var key, value string
	for infoChild := node.FirstChild; infoChild != nil; infoChild = infoChild.NextSibling {
		if rowHasClass(infoChild, "span", "icount") {
			value = infoChild.FirstChild.Data
		}
		if rowHasClass(infoChild, "div", "info_label") {
			key = infoChild.FirstChild.Data
		}
	}
	if strings.Index(key, "screened") >= 0 {
		update.ScreenedAtAirport = value
	}
	if strings.Index(key, "Active") >= 0 {
		update.Active = value
	}
	if strings.Index(key, "Cured") >= 0 {
		update.Cured = value
	}
	if strings.Index(key, "Death") >= 0 {
		update.Death = value
	}
	if strings.Index(key, "Migrated") >= 0 {
		update.Migrated = value
	}
}

func setStatValues(node *html.Node, val *string) {
	for infoChild := node.FirstChild; infoChild != nil; infoChild = infoChild.NextSibling {
		if isElementNode(infoChild, "strong") {
			*val = infoChild.FirstChild.Data
		}
	}
}

func rowHasClass(node *html.Node, nodeType string, class string) bool {
	return node.Type == html.ElementNode && node.Data == nodeType && func(n *html.Node) bool {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == class {
				return true
			}
		}
		return false
	}(node)
}

func fetchUrl(node *html.Node) string {
	for _, a := range node.Attr {
		if a.Key == "href" {
			if strings.Index(a.Val, "http") == 0 {
				return a.Val
			}
			return ""
		}
	}
	return ""
}

func toInt(input string) int {
	output, _ := strconv.Atoi(strings.ReplaceAll(input, ",", ""))
	return output
}
