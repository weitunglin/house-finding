package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

// ParseLayout is for parsing date time
const ParseLayout string = "2006 1 2"

var cases []Case

func init() {
	cases = make([]Case, 0)
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
		exist := false
		for _, v := range cases {
			if v.Name == c.Name {
				exist = true
				break
			}
		}

		if !exist {
			cases = append(cases, c)
			// announce to line notify
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
			<-time.After(5 * time.Minute)
			crawler.Visit("https://dof.gov.taipei/News.aspx?n=624B36D56FB63705&sms=148C417C1585EF00")
		}
	}()

	go func() {
		jt := NewJobTicker()
		for {
			<-jt.t.C

			// get cases list
			caseArray := make([]Case, 0)
			for _, v := range cases {
				y, e := strconv.Atoi(v.BidDate[0:3])
				if e != nil {
					fmt.Println("error occurred during parsing date time year")
				}

				m, e := strconv.Atoi(v.BidDate[4:6])
				if e != nil {
					fmt.Println("error occurred during parsing date time month")
				}

				d, e := strconv.Atoi(v.BidDate[7:9])
				if e != nil {
					fmt.Println("error occurred during parsing date time day")
				}

				ct, err := time.Parse(ParseLayout, strconv.Itoa(y+1911)+" "+strconv.Itoa(m)+" "+strconv.Itoa(d))
				if err != nil {
					fmt.Println("error occurred during paring datet time")
				}

				if time.Now().Before(ct) {
					caseArray = append(caseArray, v)
				}
			}

			if bytes, err := json.MarshalIndent(caseArray, "", "  "); err != nil {
				fmt.Println("error occurred during parsing list json")
			} else {
				if err = notify(string(bytes)); err != nil {
					fmt.Println("error occurred during sending line message (list)")

				}
			}

			jt.updateJobTicker()
		}
	}()

	ch := make(chan bool)
	<-ch
}
