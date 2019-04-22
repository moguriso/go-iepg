package iepg

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func PrintReserve(target string, ch string) {
	doc, err := goquery.NewDocument("https://tv.yahoo.co.jp/search/?q=" + title + "&a=23&oa=1&t=" + ch) // 0=tidigi, 2=cs
	if err != nil {
		fmt.Print("document not found. ")
		os.Exit(1)
	}

	d := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.leftarea > p.yjMS > em").Each(func(_ int, s *goquery.Selection) {
		d = s.Text()
	})
	month, date, sub_date := ParseDate(d)
	if (month == "") || (date == "") {
		return
	}

	tm_range := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.leftarea > p:nth-child(2) > em").Each(func(_ int, s *goquery.Selection) {
		tm_range = s.Text()
	})
	start_time, end_time := ParseTime(tm_range)

	station := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.rightarea > p:nth-child(2) > span.pr35").Each(func(_ int, s *goquery.Selection) {
		station = s.Text()
	})
	fmt.Println(station)

	title := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.rightarea > p.yjLS.pb5p > a").Each(func(_ int, s *goquery.Selection) {
		title = s.Text()
	})
	fmt.Println(title + sub_date)

	fmt.Println("Content-type: application/x-tv-program-info; charset=shift_jis")
	fmt.Println("version: 1")
	fmt.Println("station: " + ConvertStation(station))
	fmt.Println("year: " + fmt.Sprintf("%04d", time.Now().Year()))
	fmt.Println("month: " + month)
	fmt.Println("date: " + date)
	fmt.Println("start: " + start_time)
	fmt.Println("end: " + end_time)
	fmt.Println("program-title: " + title + sub_date)
}

func ReserveTidigi(target string) {
	doc, err := goquery.NewDocument("https://tv.yahoo.co.jp/search/?q=" + target + "&a=23&oa=1&t=0")
	if err != nil {
		fmt.Print("document not found. ")
		os.Exit(1)
	}

	d := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.leftarea > p.yjMS > em").Each(func(_ int, s *goquery.Selection) {
		d = s.Text()
	})
	month, date, sub_date := ParseDate(d)
	if (month == "") || (date == "") {
		return
	}

	tm_range := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.leftarea > p:nth-child(2) > em").Each(func(_ int, s *goquery.Selection) {
		tm_range = s.Text()
	})
	start_time, end_time := ParseTime(tm_range)

	station := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.rightarea > p:nth-child(2) > span.pr35").Each(func(_ int, s *goquery.Selection) {
		station = s.Text()
	})
	fmt.Println(station)

	title := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.rightarea > p.yjLS.pb5p > a").Each(func(_ int, s *goquery.Selection) {
		title = s.Text()
	})
	fmt.Println(title + sub_date)

	OutputIepg("test.tvpi", ConvertStation(station), fmt.Sprintf("%04d", time.Now().Year()),
		month, date, start_time, end_time, title+sub_date)
	exe, err := os.Executable()
	err = exec.Command(".\\PLUMAGE\\x64\\PLUMAGE.exe", filepath.Dir(exe)+"\\test.tvpi").Run()
	fmt.Println(err)
}

func ReserveCs(target string) {
	doc, err := goquery.NewDocument("https://tv.yahoo.co.jp/search/?q=" + target + "&a=23&oa=1&t=2")
	if err != nil {
		fmt.Print("document not found. ")
		os.Exit(1)
	}

	d := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.leftarea > p.yjMS > em").Each(func(_ int, s *goquery.Selection) {
		d = s.Text()
	})
	month, date, sub_date := ParseDate(d)
	if (month == "") || (date == "") {
		return
	}

	tm_range := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.leftarea > p:nth-child(2) > em").Each(func(_ int, s *goquery.Selection) {
		tm_range = s.Text()
	})
	start_time, end_time := ParseTime(tm_range)

	station := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.rightarea > p:nth-child(2) > span.pr35").Each(func(_ int, s *goquery.Selection) {
		station = s.Text()
	})
	fmt.Println(station)

	title := ""
	doc.Find("#main > div:nth-child(7) > ul > li > div.rightarea > p.yjLS.pb5p > a").Each(func(_ int, s *goquery.Selection) {
		title = s.Text()
	})
	fmt.Println(title + sub_date)

	OutputIepg("test.tvpi", ConvertStation(station), fmt.Sprintf("%04d", time.Now().Year()),
		month, date, start_time, end_time, title+sub_date)
	exe, err := os.Executable()
	err = exec.Command(".\\PLUMAGE\\x64\\PLUMAGE.exe", filepath.Dir(exe)+"\\test.tvpi").Run()
	fmt.Println(err)
}

func OutputIepg(file_name, station, year, month, date, start_time, end_time, title string) {
	fp, err := os.OpenFile(file_name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	sjisWriter := bufio.NewWriter(transform.NewWriter(fp, japanese.ShiftJIS.NewEncoder()))
	_, err = sjisWriter.WriteString("Content-type: application/x-tv-program-info; charset=shift_jis\n")
	_, err = sjisWriter.WriteString("version: 1\n")
	_, err = sjisWriter.WriteString("station: " + station + "\n")
	_, err = sjisWriter.WriteString("year: " + year + "\n")
	_, err = sjisWriter.WriteString("month: " + month + "\n")
	_, err = sjisWriter.WriteString("date: " + date + "\n")
	_, err = sjisWriter.WriteString("start: " + start_time + "\n")
	_, err = sjisWriter.WriteString("end: " + end_time + "\n")
	_, err = sjisWriter.WriteString("program-title: " + title + "\n")
	err = sjisWriter.Flush()
}

func ConvertStation(in string) string {
	if strings.Contains(in, "TOKYO　MX") {
		return "TOKYO MX"
	} else if strings.Contains(in, "ＴＢＳ") || strings.Contains(in, "TBS") {
		return "ＴＢＳテレビ"
	} else if strings.Contains(in, "テレビ東京") {
		return "テレビ東京"
	} else if strings.Contains(in, "フジテレビ") {
		return "フジテレビ"
	} else if strings.Contains(in, "日テレ") {
		return "日本テレビ"
	} else if strings.Contains(in, "テレビ朝日") {
		return "テレビ朝日"
	}
	return ""
}

func ParseDate(in string) (string, string, string) {
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

func ParseTime(in string) (string, string) {
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

//func ParseTitleNumber(in string) string {
//	r1 := regexp.MustCompile(`#[1-9] `)
//	r2 := regexp.MustCompile(`#[1-9]$`)
//}
