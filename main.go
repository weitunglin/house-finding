package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

var cases map[string]Case

func init() {
	cases = make(map[string]Case)
	cases["公開標售臺北市中正區水源路37之2號3樓之2等3戶市有不動產（均含車位），請踴躍參加投標。"] = Case{
		Name:         "公開標售臺北市中正區水源路37之2號3樓之2等3戶市有不動產（均含車位），請踴躍參加投標。",
		AnnounceDate: "109-07-29",
		BidDate:      "109-08-21",
		Progress:     "公告中",
	}
	cases["公開徵求臺北市士林區建業路7號等3處市有房地活化提案，請踴躍參加。"] = Case{
		Name:         "公開徵求臺北市士林區建業路7號等3處市有房地活化提案，請踴躍參加。",
		AnnounceDate: "109-07-22",
		BidDate:      "109-08-25",
		Progress:     "公告中",
	}
	cases["公開標租新北市板橋區江子翠段第二崁小段189-5地號市有不動產案"] = Case{
		Name:         "公開標租新北市板橋區江子翠段第二崁小段189-5地號市有不動產案",
		AnnounceDate: "109-08-10",
		BidDate:      "109-09-01",
		Progress:     "公告中",
	}
}

func main() {
	crawler := colly.NewCollector()

	crawler.OnHTML("tbody > tr", func(e *colly.HTMLElement) {
		var c Case
		e.ForEach("td", func(idx int, el *colly.HTMLElement) {
			fmt.Printf("Col %d", idx)
			column := el.Attr("data-title")
			fmt.Printf("column: %s", column)
			value := el.Text
			fmt.Printf("value: %s", value)
			switch column {
			case "標案名稱":
				c.Name = value
				break
			case "公告日期":
				c.AnnounceDate = value
				break
			case "開標日期":
				c.BidDate = value
				break
			case "標案進度":
				c.Progress = value
				break
			default:
			}
		})
		fmt.Printf(", row: %+v\n", c)
		_, exist := cases[c.Name]
		if exist {
			// update
			cases[c.Name] = c
		} else {
			// add to list
			// announce to line notify
			cases[c.Name] = c
			if bytes, err := json.MarshalIndent(c, "", "  "); err != nil {
				fmt.Println("error occurred during parsing json")
			} else {
				if err = notify(string(bytes)); err != nil {
					fmt.Println("error occurred during sending line message")
				}
			}
		}
	})

	crawler.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	crawler.OnError(func(r *colly.Response, err error) {
		fmt.Println("error occurred during visiting webste", "\nError:", err)
		if err = notify(fmt.Sprintf("request failed, err: %s", err.Error())); err != nil {
			fmt.Println("error occurred during sending line error message")
		}
	})

	go func() {
		for {
			<-time.After(1 * time.Minute)
			crawler.Visit("https://dof.gov.taipei/News.aspx?n=624B36D56FB63705&sms=148C417C1585EF00")
		}
	}()

	<-time.After(1 * time.Second)
	tick := time.Tick(1 * time.Minute)
	for range tick {
		caseArray := make([]Case, 0)
		for _, v := range cases {
			caseArray = append(caseArray, v)
		}

		if bytes, err := json.MarshalIndent(caseArray, "", "  "); err != nil {
			fmt.Println("error occurred during parsing list json")
		} else {
			if err = notify(string(bytes)); err != nil {
				fmt.Println("error occurred during sending line message (list)")

			}
		}
	}
}
