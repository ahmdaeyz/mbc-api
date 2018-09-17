package main

import (
	"encoding/json"
	"github.com/corpix/uarand"
	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	"github.com/metakeule/fmtdate"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type day struct {
	Shows []show
	Date  date
}
type date struct {
	DayLiteral string    `json:"day"`
	Date       time.Time `json:"date"`
}
type show struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImgURL      string `json:"imgUrl"`
	ShowTimes   map[string]time.Time
	currently   bool
}
type scrapper struct {
	URL *url.URL
}

var channels = map[string]string{"mbc2": "http://www.mbc.net/ar/mbc2/grid.channel-mbc1-me.programtype-all.html", "mbcmasr": "http://www.mbc.net/ar/mbc2/grid.channel-mbc-masr.programtype-all.html"}
var days = map[string]string{"Sat": "", "Sun": "", "Mon": "", "Tue": "", "Wed": "", "Thu": "", "Fri": ""}

func (scrapper *scrapper) Scrap() []day {
	var allShows [][]show
	var dates []date
	var days []day
	c := colly.NewCollector(
		colly.UserAgent(uarand.GetRandom()),
	)
	c.OnHTML("ol.text-box-toc li", func(element *colly.HTMLElement) {
		strDate := strings.TrimSpace(strings.Replace(strings.Replace(element.Text, "\n", "", -1), "  ", "", -1))
		finDate := date{}
		if strings.Contains(strDate, "السبت") {
			finDate.DayLiteral = "Sat"
			finDate.Date, _ = fmtdate.Parse("DD/MM/YYYY", strings.TrimSpace(strings.Replace(strDate, "السبت", "", -1)))
		} else if strings.Contains(strDate, "الأحد") {
			finDate.DayLiteral = "Sun"
			finDate.Date, _ = fmtdate.Parse("DD/MM/YYYY", strings.TrimSpace(strings.Replace(strDate, "الأحد", "", -1)))
		} else if strings.Contains(strDate, "الاثنين") {
			finDate.DayLiteral = "Mon"
			finDate.Date, _ = fmtdate.Parse("DD/MM/YYYY", strings.TrimSpace(strings.Replace(strDate, "الاثنين", "", -1)))
		} else if strings.Contains(strDate, "الثلاثاء") {
			finDate.DayLiteral = "Tue"
			finDate.Date, _ = fmtdate.Parse("DD/MM/YYYY", strings.TrimSpace(strings.Replace(strDate, "الثلاثاء", "", -1)))
		} else if strings.Contains(strDate, "الأربعاء") {
			finDate.DayLiteral = "Wed"
			finDate.Date, _ = fmtdate.Parse("DD/MM/YYYY", strings.TrimSpace(strings.Replace(strDate, "الأربعاء", "", -1)))
		} else if strings.Contains(strDate, "الخميس") {
			finDate.DayLiteral = "Thu"
			finDate.Date, _ = fmtdate.Parse("DD/MM/YYYY", strings.TrimSpace(strings.Replace(strDate, "الخميس", "", -1)))
		} else if strings.Contains(strDate, "الجمعة") {
			finDate.DayLiteral = "Fri"
			finDate.Date, _ = fmtdate.Parse("DD/MM/YYYY", strings.TrimSpace(strings.Replace(strDate, "الجمعة", "", -1)))
		}
		dates = append(dates, finDate)
	})
	c.Visit(scrapper.URL.String())
	for i := 0; i < 7; i++ {
		c := colly.NewCollector(
			colly.UserAgent(uarand.GetRandom()),
		)
		var imgUrls []string
		var playTimes []map[string]time.Time
		var titles []string
		var descriptions []string
		var states []bool
		var shows []show
		c.OnHTML("#tab-1-"+strconv.Itoa(i)+"-b6a2c45e-9bb2-427f-9f33-2d78a9efa4ca", func(element *colly.HTMLElement) {
			element.ForEach("div.archttl", func(i int, element *colly.HTMLElement) {
				element.ForEach("img", func(i int, element *colly.HTMLElement) {
					imgUrls = append(imgUrls, "http://www.mbc.net"+element.Attr("src"))
				})
				element.ForEach("div.img-box ul.info", func(i int, element *colly.HTMLElement) {
					times := make(map[string]time.Time)
					element.ForEach("ul.info li", func(i int, element *colly.HTMLElement) {
						if strings.Contains(element.Text, "توقيت") {
							if strings.Contains(element.Text, "توقيت مصر") {
								EgTime, _ := fmtdate.Parse("hh:mm", strings.TrimSpace(strings.Replace(strings.Replace(element.Text, "\n", "", -1), "توقيت مصر", "", -1)))
								times["EG"] = EgTime
							} else if strings.Contains(element.Text, "توقيت السعودية") {
								ksaTime, _ := fmtdate.Parse("hh:mm", strings.TrimSpace(strings.Replace(strings.Replace(element.Text, "\n", "", -1), "توقيت السعودية", "", -1)))
								times["KSA"] = ksaTime
							} else if strings.Contains(element.Text, "توقيت جرينتش") {
								GMTTime, _ := fmtdate.Parse("hh:mm", strings.TrimSpace(strings.Replace(strings.Replace(element.Text, "\n", "", -1), "توقيت جرينتش", "", -1)))
								times["GMT"] = GMTTime
							}
						}
					})
					playTimes = append(playTimes, times)
				})
				element.ForEach("div.archttl h3", func(i int, element *colly.HTMLElement) {
					titles = append(titles, strings.Replace(element.Text, "\n", "", -1))
				})
				element.ForEach("div.archttl p", func(i int, element *colly.HTMLElement) {
					descriptions = append(descriptions, strings.TrimSpace(strings.Replace(element.Text, "\n", "", -1)))
				})
				if strings.Contains(element.Attr("class"), "currently") {
					states = append(states, true)
				} else {
					states = append(states, false)
				}
			})
		})
		c.Visit(scrapper.URL.String())
		for i := 0; i < len(titles); i++ {
			shows = append(shows, show{Title: titles[i], Description: descriptions[i], ImgURL: imgUrls[i], currently: states[i], ShowTimes: playTimes[i]})
		}
		if len(shows) != 0 {
			allShows = append(allShows, shows)
		}
	}
	for i := 0; i < 7; i++ {
		days = append(days, day{allShows[i], dates[i]})
	}
	return days
}
func currentlyDisplaying(w http.ResponseWriter, r *http.Request) {
	scrper := &scrapper{}
	params := mux.Vars(r)
	value, ok := channels[params["channel"]]
	if ok {
		parsedURL, _ := url.Parse(value)
		scrper.URL = parsedURL
		w.Header().Set("Content-Type", "application/json")
		for _, day := range scrper.Scrap() {
			for _, aShow := range day.Shows {
				if aShow.currently {
					json.NewEncoder(w).Encode(aShow)
				}
			}
		}
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Can't find this channel"))
	}
}
func dayShows(w http.ResponseWriter, r *http.Request) {
	scrper := &scrapper{}
	params := mux.Vars(r)
	value, ok := channels[params["channel"]]
	if ok {
		parsedURL, _ := url.Parse(value)
		scrper.URL = parsedURL
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("encoding", "utf-8")
		_, dayExists := days[params["day"]]
		if dayExists {
			for _, day := range scrper.Scrap() {
				if day.Date.DayLiteral == params["day"] {
					json.NewEncoder(w).Encode(day.Shows)
				}
			}
		} else {
			w.WriteHeader(400)
		}
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Can't find this channel"))
	}
}
func wholeWeek(w http.ResponseWriter, r *http.Request) {
	scrper := &scrapper{}
	params := mux.Vars(r)
	value, ok := channels[params["channel"]]
	if ok {
		parsedURL, _ := url.Parse(value)
		scrper.URL = parsedURL
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("encoding", "utf-8")
		json.NewEncoder(w).Encode(scrper.Scrap())
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Can't find this channel"))
	}
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{channel}/currently", currentlyDisplaying).Methods("GET")
	r.HandleFunc("/{channel}/{day}", dayShows).Methods("GET")
	r.HandleFunc("/{channel}", wholeWeek).Methods("GET")
	http.ListenAndServe(":8080", r)
}
