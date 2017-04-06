package public

import (
	"errors"
	"strings"
)

type UrlParse struct {
	keyPos map[string]int
}

func (this *UrlParse) InitModle(pattern string) {
	this.keyPos = nil
	modles := parseStr(pattern)
	for i, modle := range modles {
		if modle != "" && modle != "*" {
			this.keyPos[modle] = i
		}
	}
}

func (this *UrlParse) Parse(path string) (map[string]string, error) {
	var keyValue map[string]string
	dest := parseStr(path)
	for key, pos := range this.keyPos {
		if pos > len(dest) {
			return nil, errors.New("path is not correct")
		}
		keyValue[key] = dest[pos]
	}
	return keyValue, nil
}

func parseStr(path string) []string {
	var dest []string
	src := strings.Split(path, "/")
	for _, elem := range src {
		dest = append(dest, elem)
	}
	return dest
}
