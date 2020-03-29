package external

import (
	"covid19/common"
	"github.com/cenkalti/backoff/v4"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func LiveData() *common.CoronaUpdate {
	constantBackoff := backoff.NewConstantBackOff(500 * time.Millisecond)
	var resp *http.Response
	var err error
	err = backoff.Retry(func() error {
		resp, err = http.Get("https://www.mohfw.gov.in/")
		return err
	}, backoff.WithMaxRetries(constantBackoff, 4))

	stateData := make([]*common.StateData, 0)
	links := make([]string, 0)

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
	parseNode(node, update)
	update.Total = strconv.Itoa(toInt(update.Active) + toInt(update.Cured) +
		toInt(update.Migrated) + toInt(update.Death))

	return update
}

func parseNode(node *html.Node, update *common.CoronaUpdate) {
	if strings.Contains(node.Data, "foreign Nationals, as on") {
		split := strings.Split(node.Data, "foreign Nationals, as on")
		update.UpdateTime = strings.TrimSuffix(split[1], ")")
		return
	}

	if isElementNode(node, "a") {
		url := fetchUrl(node)
		if url != "" && strings.Contains(url,"mohfw") {
			update.Links = append(update.Links, url)
		}
		return
	}

	if rowHasClass(node, "div", "iblock_text") {
		setValues(node, update)
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
		parseNode(child, update)
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
