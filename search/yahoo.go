package search

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"go-iepg/log"
)

const YAHOO_ENDPOINT = "https://tv.yahoo.co.jp/search/?q="
const APPEND_QUERY = "&a=23&oa=1&t="

type ReadData struct {
	Station string
	Year    int
	Month   int
	Date    int
	Start_h int
	Start_m int
	End_h   int
	End_m   int
	Title   string
	Re      bool
	WeekDay string
	IsCs    bool
}

func Search(target string, ch string) *goquery.Document {
	title := strings.ReplaceAll(target, "%20", " ")
	title = strings.ReplaceAll(title, " ", "%20")
	doc, err := goquery.NewDocument(YAHOO_ENDPOINT + title + APPEND_QUERY + ch) // 0=tidigi, 2=cs
	if err != nil {
		log.L.Error("document not found. ")
		return nil
	}
	return doc
}

func ParseFindCount(doc *goquery.Document) int {
	count := ""
	doc.Find("#main > div.yjMS.search_number.mb10 > p.floatl > em:nth-child(1)").Each(func(_ int, s *goquery.Selection) {
		count = s.Text()
	})

	c, _ := strconv.Atoi(count)

	return c
}

func ParseSection(doc *goquery.Document, isCs bool) []*ReadData {
	selection := doc.Find("#main > div:nth-child(7) > ul")
	innserSelection := selection.Find("li")

	var ret []*ReadData
	innserSelection.Each(func(_ int, s *goquery.Selection) {
		res := &ReadData{
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
		wd := ParseWeekDay(s)
		log.L.Error("weekday = " + wd)
		res.WeekDay = wd
		station := ParseStation(s)
		res.Station = station
		month, date, _ := ParseDate(s)
		res.Month, _ = strconv.Atoi(month)
		res.Date, _ = strconv.Atoi(date)
		start_time, end_time := ParseTime(s)
		res.Start_h, _ = strconv.Atoi(start_time[:2])
		res.Start_m, _ = strconv.Atoi(start_time[3:])
		res.End_h, _ = strconv.Atoi(end_time[:2])
		res.End_m, _ = strconv.Atoi(end_time[3:])
		title := ParseTitle(s)
		res.Title = title
		re := ParseRe(s)
		if strings.Contains(re, "再") {
			res.Re = true
		}
		if isCs {
			res.IsCs = true
		}
		log.L.Debug("res = ")
		log.L.Debug(res)
		ret = append(ret, res)
	})
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
func ParseDate(doc *goquery.Selection) (string, string, string) {
	d := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.leftarea > p.yjMS > em").Each(func(_ int, s *goquery.Selection) {
		d = s.Text()
	})

	month, date, sub_date := getDate(d)
	if (month == "") || (date == "") {
		return "", "", ""
	}
	return month, date, sub_date
}

func ParseWeekDay(doc *goquery.Selection) string {
	d := ""
	//doc.Find("#main > div:nth-child(7) > ul > li:nth-child(1) > div.leftarea > p.yjMS").Each(func(_ int, s *goquery.Selection) {
	doc.Find("#main > div > ul > li > div.leftarea > p.yjMS").Each(func(_ int, s *goquery.Selection) {
		d = s.Text()
	})
	b := ""
	//doc.Find("#main > div:nth-child(7) > ul > li > div.leftarea > p.yjMS > em").Each(func(_ int, s *goquery.Selection) {
	doc.Find("#main > div > ul > li > div.leftarea > p.yjMS > em").Each(func(_ int, s *goquery.Selection) {
		b = s.Text()
	})
	//d = d[strings.Index(d, "（")+1 : strings.Index(d, "）")]
	d = strings.Trim(d, b)
	//d = strings.TrimRight(d, "）")
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

func ParseTime(doc *goquery.Selection) (string, string) {
	tm_range := ""
	doc.Find("#main > div > ul > li > div.leftarea > p > em").Each(func(_ int, s *goquery.Selection) {
		tm_range = s.Text()
	})
	start_time, end_time := getTime(tm_range)
	return start_time, end_time
}

func ParseStation(doc *goquery.Selection) string {
	station := ""
	doc.Find("#main > div > ul > li > div.rightarea > p:nth-child(2) > span.pr35").Each(func(_ int, s *goquery.Selection) {
		station = s.Text()
	})
	log.L.Debug(station)
	return station
}

func ParseTitle(doc *goquery.Selection) string {
	title := ""
	doc.Find("#main > div > ul > li > div.rightarea > p.yjLS.pb5p > a").Each(func(_ int, s *goquery.Selection) {
		title = s.Text()
	})
	title = strings.ReplaceAll(title, "　", " ")
	log.L.Info("before:", title)
	if r := convertEpisodeNumber(title); r != "" {
		title = r
		log.L.Info("after:", title)
	} else {
		log.L.Error("title episode convert error")
	}
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

func ParseRe(doc *goquery.Selection) string {
	re := ""
	doc.Find("#main > div > ul > li > div.rightarea > p.yjLS.pb5p > span").Each(func(_ int, s *goquery.Selection) {
		re = s.Text()
	})
	re = strings.ReplaceAll(re, "　", " ")
	log.L.Error(re)
	return re
}
