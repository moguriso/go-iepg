package iepg

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"bufio"
	"fmt"
	"strings"

	"go-iepg/log"
	p "go-iepg/param"
	pl "go-iepg/search"
)

type ReserveType int32

const (
	TIDIJI ReserveType = iota
	CS
)

func getReadData(dp *p.DynamicParam) []*pl.ReadData {
	ch := "0"
	if dp.IsCs {
		ch = "1&t=2"
	}

	doc := pl.Search(dp.Title, ch)
	if doc == nil {
		return nil
	}

	r := pl.ParseSection(doc, dp.IsCs)
	return r
}

func PrintReserve(dp *p.DynamicParam) {
	r := getReadData(dp)
	for _, v := range r {
		if !confirmTitle(v.Title, dp.Title) {
			log.L.Error("failed to reserve [ " + v.Title + "]")
			log.L.Error("title unmatch src[" + v.Title + "]  [" + dp.Title + "]")
		} else if !confirmStartTime(v.Start_h, v.Start_m, dp.Start_time) {
			log.L.Error("failed to reserve [ " + v.Title + "]")
			log.L.Error("start time unmatch src[" + fmt.Sprintf("%02d:%02d", v.Start_h, v.Start_m) + "]  [" + dp.Start_time + "]")
		} else if !confirmWeekDay(v.WeekDay, dp.WeekDay) {
			log.L.Error("failed to reserve [ " + v.Title + "]")
			log.L.Error("weekday unmatch src[" + v.WeekDay + "]  [" + dp.WeekDay + "]")
		} else if v.Re {
			log.L.Error("failed to reserve [ " + v.Title + "]")
			log.L.Error("再放送番組は録画しません") /* TODO: if need reserve ... */
		} else {
			fmt.Println("Content-type: application/x-tv-program-info; charset=shift_jis")
			fmt.Println("version: 1")
			fmt.Println("station: " + convertStation(v.Station))
			/* TODO: 格好悪いので後で直す...が、そもそも年が入ってないほうが悪い... */
			if v.Month < 12 {
				fmt.Println("year: " + fmt.Sprintf("%04d", time.Now().Year()))
			} else {
				fmt.Println("year: " + fmt.Sprintf("%04d", time.Now().Year()+1))
			}
			fmt.Println("month: " + fmt.Sprintf("%02d", v.Month))
			fmt.Println("date: " + fmt.Sprintf("%02d", v.Date))
			fmt.Println("start: " + fmt.Sprintf("%02d:%02d", v.Start_h, v.Start_m))
			fmt.Println("end: " + fmt.Sprintf("%02d:%02d", v.End_h, v.End_m))
			fmt.Println("program-title: " + v.Title + fmt.Sprintf(" %d月%d日", v.Month, v.Date))
		}
	}
}

func confirmTitle(src, dst string) bool {
	s := strings.ReplaceAll(strings.ToLower(src), "　", " ")
	d := strings.ReplaceAll(strings.ToLower(dst), "　", " ")
	md := strings.Split(d, " ")
	for _, v := range md {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}

func confirmAlreadyStart(m, d, h, mm int) bool {
	t := time.Now()
	tm := int(t.Month())
	td := int(t.Day())
	th := int(t.Hour())
	tmm := int(t.Minute())
	log.L.Println("src = ", tm, td, th, tmm)
	log.L.Println("dst = ", m, d, h, mm)
	if tm >= m && td >= d {
		if (th > h) || (th == h && tmm > mm) {
			log.L.Println("this is time over.")
			return false
		}
	}
	return true
}

func confirmStartTime(h, mm int, dst string) bool {
	if dst == "" {
		return true
	}
	src_st := fmt.Sprintf("%02d:%02d", h, mm)
	if strings.Contains(src_st, dst) {
		return true
	}
	return false
}

func confirmStation(src, dst string) bool {
	log.L.Println("src =", src)
	log.L.Println("dst =", dst)
	if dst == "" {
		return true
	}
	if strings.Contains(src, dst) {
		return true
	}
	return false
}

func confirmWeekDay(src, dst string) bool {
	log.L.Println("src =", src)
	log.L.Println("dst =", dst)
	tmp := strings.ReplaceAll(src, "（", "")
	tmp = strings.ReplaceAll(tmp, "）", "")
	log.L.Println("tmp =", tmp)
	if dst == "" {
		return true
	}
	if strings.Contains(dst, tmp) {
		return true
	}
	return false
}

func Reserve(dp *p.DynamicParam) {
	s_conf := p.LoadStaticParam("config.json")
	r := getReadData(dp)
	for _, v := range r {
		if !confirmTitle(v.Title, dp.Title) {
			log.L.Error("failed to reserve [ " + v.Title + "]")
			log.L.Error("title unmatch src[" + v.Title + "]  [" + dp.Title + "]")
		} else if !confirmStartTime(v.Start_h, v.Start_m, dp.Start_time) {
			log.L.Error("failed to reserve [ " + v.Title + "]")
			log.L.Error("start time unmatch src[" + fmt.Sprintf("%02d:%02d", v.Start_h, v.Start_m) + "]  [" + dp.Start_time + "]")
		} else if !confirmAlreadyStart(v.Month, v.Date, v.Start_h, v.Start_m) {
			log.L.Error("failed to reserve [ " + v.Title + "]")
			log.L.Error("it may already started or done.")
		} else if !confirmWeekDay(v.WeekDay, dp.WeekDay) {
			log.L.Error("failed to reserve [ " + v.Title + "]")
			log.L.Error("weekday unmatch src[" + v.WeekDay + "]  [" + dp.WeekDay + "]")
		} else if !confirmStation(v.Station, dp.Station) {
			log.L.Error("failed to reserve [ " + v.Title + "]")
			log.L.Error("Station unmatch src[" + v.Station + "]  [" + dp.Station + "]")
		} else if v.Re {
			log.L.Error("failed to reserve [ " + v.Title + "]")
			log.L.Error("再放送番組は録画しません") /* TODO: if need reserve ... */
		} else {
			OutputIepg(s_conf.TempFileName, v)
			exe, err := os.Executable()
			err = exec.Command(s_conf.PlumagePath, filepath.Dir(exe)+"\\"+s_conf.TempFileName).Run()
			log.L.Error(err)
			log.L.Error(v.Title + "予約しました")
			//break
		}
	}
}

func OutputIepg(fileName string, in *pl.ReadData) {
	fp, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	sjisWriter := bufio.NewWriter(transform.NewWriter(fp, japanese.ShiftJIS.NewEncoder()))
	_, err = sjisWriter.WriteString("Content-type: application/x-tv-program-info; charset=shift_jis\n")
	_, err = sjisWriter.WriteString("version: 1\n")
	_, err = sjisWriter.WriteString("station: " + convertStation(in.Station) + "\n")
	/* TODO: 格好悪いので後で直す...が、そもそも年が入ってないほうが悪い... */
	if in.Month < 12 {
		_, err = sjisWriter.WriteString("year: " + fmt.Sprintf("%04d", time.Now().Year()+1) + "\n")
	} else {
		_, err = sjisWriter.WriteString("year: " + fmt.Sprintf("%04d", time.Now().Year()) + "\n")
	}
	_, err = sjisWriter.WriteString("month: " + fmt.Sprintf("%02d", in.Month) + "\n")
	_, err = sjisWriter.WriteString("date: " + fmt.Sprintf("%02d", in.Date) + "\n")
	_, err = sjisWriter.WriteString("start: " + fmt.Sprintf("%02d:%02d", in.Start_h, in.Start_m) + "\n")
	_, err = sjisWriter.WriteString("end: " + fmt.Sprintf("%02d:%02d", in.End_h, in.End_m) + "\n")
	tl := in.Title
	if in.IsCs {
		tl += "_cs"
	}
	tl += fmt.Sprintf(" %d月%d日", in.Month, in.Date) + "\n"
	_, err = sjisWriter.WriteString("program-title: " + tl)
	err = sjisWriter.Flush()
}

func convertStation(in string) string {
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
