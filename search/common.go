package search

import (
	"github.com/PuerkitoBio/goquery"
	"go-iepg/kingdom"
	p "go-iepg/param"
	"go-iepg/yahoo"
)

func Search(target string, ch p.CHANNEL_MODE) *goquery.Document {
	if ch == p.TIDEGI_MODE {
		return yahoo.Search(target, "3") /* 3:tideji, 1:bs...*/
	} else if ch == p.CS_MODE {
		return kingdom.Search(target)
	}
	return nil
}

func ParseSection(doc *goquery.Document, isCs bool, isRecReAir bool) []*p.ReadData {
	if isCs {
		return kingdom.ParseSection(doc, isCs, isRecReAir)
	} else {
		return yahoo.ParseSection(doc, isCs, isRecReAir)
	}
	return nil
}
