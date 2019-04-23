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
	pl "go-iepg/search"
)

func PrintReserve(target string, ch string) {
	doc := pl.Search(target, ch)
	if doc == nil {
		return
	}
	r := pl.ParseSection(doc)
	for _, v := range r {
		src := strings.ToLower(v.Title)
		dst := strings.ToLower(target)
		if strings.Contains(src, dst) {
			fmt.Println("Content-type: application/x-tv-program-info; charset=shift_jis")
			fmt.Println("version: 1")
			fmt.Println("station: " + convertStation(v.Station))
			fmt.Println("year: " + fmt.Sprintf("%04d", time.Now().Year()))
			fmt.Println("month: " + fmt.Sprintf("%02d", v.Month))
			fmt.Println("date: " + fmt.Sprintf("%02d", v.Date))
			fmt.Println("start: " + fmt.Sprintf("%02d:%02d", v.Start_h, v.Start_m))
			fmt.Println("end: " + fmt.Sprintf("%02d:%02d", v.End_h, v.End_m))
			fmt.Println("program-title: " + v.Title)
		}
	}
}

func reserveCommon(target, ch string) {
	doc := pl.Search(target, ch)
	if doc == nil {
		return
	}
	r := pl.ParseSection(doc)
	for _, v := range r {
		src := strings.ToLower(v.Title)
		dst := strings.ToLower(target)
		if strings.Contains(src, dst) {
			OutputIepg(v)
			exe, err := os.Executable()
			err = exec.Command(".\\PLUMAGE\\x64\\PLUMAGE.exe", filepath.Dir(exe)+"\\test.tvpi").Run()
			log.L.Debug(err)
			break
		}
	}
}

func ReserveTidigi(target string) {
	reserveCommon(target, "0")
}

func ReserveCs(target string) {
	reserveCommon(target, "2")
}

func OutputIepg(in *pl.ReadData) {
	fp, err := os.OpenFile("tmp.tvpi", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	sjisWriter := bufio.NewWriter(transform.NewWriter(fp, japanese.ShiftJIS.NewEncoder()))
	_, err = sjisWriter.WriteString("Content-type: application/x-tv-program-info; charset=shift_jis\n")
	_, err = sjisWriter.WriteString("version: 1\n")
	_, err = sjisWriter.WriteString("station: " + convertStation(in.Station) + "\n")
	_, err = sjisWriter.WriteString("year: " + fmt.Sprintf("%04d", time.Now().Year()) + "\n")
	_, err = sjisWriter.WriteString("month: " + fmt.Sprintf("%02d", in.Month) + "\n")
	_, err = sjisWriter.WriteString("date: " + fmt.Sprintf("%02d", in.Date) + "\n")
	_, err = sjisWriter.WriteString("start: " + fmt.Sprintf("%02d:%02d", in.Start_h, in.Start_m) + "\n")
	_, err = sjisWriter.WriteString("end: " + fmt.Sprintf("%02d:%02d", in.End_h, in.End_m) + "\n")
	_, err = sjisWriter.WriteString("program-title: " + in.Title + "\n")
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
