package param

import (
	"encoding/json"
	"io/ioutil"

	"go-iepg/log"
)

type CHANNEL_MODE int

const (
	NO_MODE CHANNEL_MODE = iota
	TIDEGI_MODE
	CS_MODE
)

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

type StaticParam struct {
	PlumagePath  string
	TempFileName string
}

type DynamicParam struct {
	Title      string
	IgnorTitle string
	Start_time string
	Station    string
	WeekDay    string
	IsCs       bool
	IsRecReAir bool
}

func LoadStaticParam(fileName string) *StaticParam {
	var st map[string]interface{}
	ret := &StaticParam{
		PlumagePath:  "",
		TempFileName: "",
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.L.Error(fileName, " read error")
		return nil
	}
	err = json.Unmarshal(data, &st)

	ret.PlumagePath, _ = st["plumage_path"].(string)
	ret.TempFileName, _ = st["temp_name"].(string)
	log.L.Debug("plumage = " + ret.PlumagePath)
	log.L.Debug("temp = " + ret.TempFileName)

	return ret
}

func LoadDynamicParam(fileName string) []*DynamicParam {
	var dp map[string]interface{}
	var ret []*DynamicParam

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.L.Error(fileName, " read error")
		return nil
	}
	err = json.Unmarshal(data, &dp)

	d, _ := dp["reserve_info"].(map[string]interface{})
	for k, v := range d {
		res := &DynamicParam{
			Title:      "",
			IgnorTitle: "",
			Start_time: "",
			Station:    "",
			IsCs:       false,
		}
		res.Title = k
		ig, is_avail := v.(map[string]interface{})["ignore"].(string)
		if is_avail {
			res.IgnorTitle = ig
		}
		res.Start_time, _ = v.(map[string]interface{})["start_time"].(string)
		res.Station, _ = v.(map[string]interface{})["station"].(string)
		res.WeekDay, _ = v.(map[string]interface{})["weekday"].(string)
		res.IsCs, _ = v.(map[string]interface{})["is_cs"].(bool)
		re, is := v.(map[string]interface{})["is_rec_re_air"].(bool)
		if is {
			res.IsRecReAir = re
		} else {
			res.IsRecReAir = false
		}
		ret = append(ret, res)
	}
	return ret
}
