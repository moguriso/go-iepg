package yahoo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/unicode/norm"

	"go-iepg/log"
	p "go-iepg/param"
	aw "go-iepg/wrapper"
	"time"
)

const YAHOO_ENDPOINT = "https://tv.yahoo.co.jp/search?q="

const APPEND_QUERY = "&g=&d=&ob=&oc=&dts=0&dte=0&a=23&t="

func Search(target string, ch string) *goquery.Document {
	title := strings.ReplaceAll(target, "%20", " ")
	title = strings.ReplaceAll(title, " ", "%20")

	driver, err := aw.GetWebDriver()
	defer driver.Stop()
	if err != nil {
		log.L.Fatal("Search: ", err)
		return nil
	}

	page, err := driver.NewPage()
	if err != nil {
		log.L.Fatal("Search: ", err)
		return nil
	}
	log.L.Info("uri = ", YAHOO_ENDPOINT+title+APPEND_QUERY+ch)
	page.Navigate(YAHOO_ENDPOINT + title + APPEND_QUERY + ch)
	time.Sleep(500 * time.Millisecond)

	ht, _ := page.HTML()
	red := strings.NewReader(ht)
	doc, err := goquery.NewDocumentFromReader(red)
	if err != nil {
		log.L.Error("document not found. ")
		return nil
	}
	return doc
}

func ParseSection(doc *goquery.Document, isCs bool, isRecReAir bool) []*p.ReadData {
	selection := doc.Find("#__next > div > main > div.inner > article > div.innerMain > section > ul")
	innserSelection := selection.Find("li.programListItem")

	var ret []*p.ReadData
	innserSelection.Each(func(_ int, s *goquery.Selection) {
		res := &p.ReadData{
			Station: "",
			Year:    0,
			Month:   0,
			Date:    0,
			Start_h: 0,
			Start_m: 0,
			End_h:   0,
			End_m:   0,
			Title:   "",
			Re:      false,
			IsCs:    false,
		}
		ht, _ := s.Html()
		log.L.Info("selection: ", ht)
		wd := parseWeekDay(s)
		if wd == "" {
			log.L.Error("weekday is not available")
			return
		}
		res.WeekDay = wd
		station := parseStation(s)
		if station == "" {
			log.L.Error("station is not available")
			return
		}
		res.Station = station
		month, date, _ := parseDate(s)
		if month == "" || date == "" {
			log.L.Error("month or date is not available")
			return
		}
		res.Month, _ = strconv.Atoi(month)
		res.Date, _ = strconv.Atoi(date)
		start_time, end_time := parseTime(s)
		if start_time == "" || end_time == "" {
			log.L.Error("start/end time is not available")
			return
		}
		res.Start_h, _ = strconv.Atoi(start_time[:2])
		res.Start_m, _ = strconv.Atoi(start_time[3:])
		res.End_h, _ = strconv.Atoi(end_time[:2])
		res.End_m, _ = strconv.Atoi(end_time[3:])
		title := parseTitle(s)
		if title == "" {
			log.L.Error("title is not available")
			return
		}
		res.Title = title
		re := parseRe(s)
		if strings.Contains(re, "再") && !isRecReAir {
			res.Re = true
		}
		if isCs {
			res.IsCs = true
		}
		log.L.Debug("res = ")
		log.L.Debug(res)
		ret = append(ret, res)
	})
	log.L.Info("ccc: ret = ", ret)
	return ret
}

func getDate(in string) (string, string, string) {
	tm := strings.ReplaceAll(in, " ", "")
	md := strings.Split(tm, "/")
	if (md[0] == "") || (md[1] == "") {
		return "", "", ""
	} else {
		month, _ := strconv.Atoi(md[0])
		date, _ := strconv.Atoi(md[1])
		return fmt.Sprintf("%02d", month), fmt.Sprintf("%02d", date), fmt.Sprintf(" %d月%d日", month, date)
	}
}

func parseDate(doc *goquery.Selection) (string, string, string) {
	d := ""
	count := 0
	tm := doc.Find("div.schedule")
	tm = tm.Find("time.scheduleText")
	tm.Find("span").Each(func(_ int, s *goquery.Selection) {
		ds := s.Text()
		count++
		log.L.Info("date(", count, ") = ", ds)
		if strings.Contains(ds, "/") {
			d = ds
		}
	})

	month, date, sub_date := getDate(d)
	if (month == "") || (date == "") {
		return "", "", ""
	}
	log.L.Info("date = ", month, " ", date, " ", sub_date)
	return month, date, sub_date
}

func parseWeekDay(doc *goquery.Selection) string {
	d := ""
	doc.Find("span.scheduleTextWeek").Each(func(_ int, s *goquery.Selection) {
		d = s.Text()
	})
	log.L.Info("week day = ", d)
	return d
}

func getTime(in string) (string, string) {
	tm := strings.ReplaceAll(in, " ", "")
	start_end := strings.Split(tm, "～")

	start := strings.Split(start_end[0], ":")
	end := strings.Split(start_end[1], ":")

	start_h, _ := strconv.Atoi(start[0])
	start_m, _ := strconv.Atoi(start[1])
	end_h, _ := strconv.Atoi(end[0])
	end_m, _ := strconv.Atoi(end[1])

	return fmt.Sprintf("%02d:%02d", start_h, start_m),
		fmt.Sprintf("%02d:%02d", end_h, end_m)
}

func parseTime(doc *goquery.Selection) (string, string) {
	d := ""
	count := 0
	tm := doc.Find("div.schedule")
	tm = tm.Find("time.scheduleText")
	tm.Find("span").Each(func(_ int, s *goquery.Selection) {
		ds := s.Text()
		count++
		log.L.Info("date(", count, ") = ", ds)
		if strings.Contains(ds, ":") && !s.HasClass("scheduleTextTimeEnd") {
			d = ds
		}
	})
	start_time := d

	end_time := ""
	doc.Find("span.scheduleTextTimeEnd").Each(func(_ int, s *goquery.Selection) {
		end_time = s.Text()
	})
	log.L.Info("start/end: ", start_time, "～", end_time)
	//return start_time, end_time
	return getTime(start_time + "～" + end_time)
}

func parseStation(doc *goquery.Selection) string {
	station := ""
	doc.Find("p.channelText").Each(func(_ int, s *goquery.Selection) {
		station = s.Text()
	})
	log.L.Info("station = ", station)
	return station
}

func parseTitle(doc *goquery.Selection) string {
	title := ""
	doc.Find("a.programListItemTitleLink").Each(func(_ int, s *goquery.Selection) {
		title = s.Text()
	})
	title = strings.ReplaceAll(title, "　", " ")
	title = strings.ReplaceAll(title, "＃", "#")
	title = strings.ReplaceAll(title, "♯", "#")
	s := title
	s = string(norm.NFKC.Bytes([]byte(s)))
	log.L.Info(s)
	title = s
	title = strings.ReplaceAll(title, "〜", "～")

	log.L.Info("before:", title)
	if r := convertEpisodeNumber(title); r != "" {
		title = r
		log.L.Info("after:", title)
	} else {
		log.L.Debug("title episode convert \"no need\" or \"error\"")
	}
	log.L.Info("title: ", title)
	return title
}

func convertEpisodeNumber(in string) string {
	ret := ""
	str := fmt.Sprintf("%s ", in)
	rep := regexp.MustCompile(`#[0-9][^0-9]`)

	n := rep.FindAllStringSubmatch(str, -1)
	for _, rr := range n {
		z := strings.Replace(rr[0], " ", "", -1)
		z = strings.Replace(z, "#", "", -1)
		p, _ := strconv.Atoi(z)
		if p > 0 && p < 10 {
			z = fmt.Sprintf("#%02d", p)
			ret = rep.ReplaceAllString(str, z)
		} else if p == 0 {
			rep2 := regexp.MustCompile(`#[0-9]`)
			n2 := rep2.FindAllStringSubmatch(str, -1)
			for _, rr2 := range n2 {
				z2 := strings.Replace(rr2[0], " ", "", -1)
				z2 = strings.Replace(z2, "#", "", -1)
				p2, _ := strconv.Atoi(z2)
				z2 = fmt.Sprintf("#%02d", p2)
				ret = rep2.ReplaceAllString(str, z2)
			}
		}
	}
	return ret
}

func parseRe(doc *goquery.Selection) string {
	re := ""
	tm := doc.Find("div.programListItemTitleGroup")
	tm = tm.Find("span")
	is_repeat := tm.HasClass("iconRepeat")
	if is_repeat {
		re = "再"
	}
	log.L.Info(re)
	return re
}
