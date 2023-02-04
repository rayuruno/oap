package main

import (
	"regexp"
	"strings"
)

type Params map[string]any

func (ps Params) Get(n string) any {
	return ps[n]
}

func (ps Params) Set(n string, v any) {
	ps[n] = v
}

func (ps Params) Add(n string, v any) {
	mm, ok := ps[n].([]any)
	if !ok {
		mm = make([]any, 0)
	}
	mm = append(mm, v)
	ps[n] = mm
}

var namesRe = regexp.MustCompile("\\[\\w+\\]")

var tmpls = map[string]map[bool][5]string{
	"matrix": {
		false: {";", "=", ",", ",", ""},
		true:  {";", "=", "", ";", "="},
	},
	"label": {
		false: {".", ".", ".", ".", ""},
		true:  {".", "", ".", "=", ""},
	},
	"form": {
		false: {"&", "=", ",", ",", ""},
		true:  {"&", "=", "", "&", "="},
	},
	"simple": {
		false: {",", ",", "", ",", ""},
		true:  {",", "", "", ",", "="},
	},
	"spacaeDelimited": {
		false: {"%20", "%20", "%20", "%20", "%20"},
	},
	"pipeDelimited": {
		false: {"|", "|", "|", "|", "|"},
	},
	"deepObject": {
		true: {"&", "=", "", "&", "="},
	},
}

type ParamType int

const (
	TypeEmpty ParamType = iota
	TypePrimitive
	TypeArray
	TypeObject
)

func Parse(qs string, style string, explode bool, typ ParamType, ps Params) (err error) {
	ts, ok := tmpls[style]
	if !ok {
		return nil
	}
	te, ok := ts[explode]
	if !ok {
		return nil
	}
	var param string
	for qs != "" {
		param, qs, _ = strings.Cut(qs, te[0])
		if param == "" {
			continue
		}
		name, value, _ := strings.Cut(param, te[1])
		switch typ {
		case 0:
			ps.Set(name, nil)
		case 1:
			ps.Set(name, value)
		case 2:
			if te[2] == "" {
				ps.Add(name, value)
				continue
			}
			for _, v := range strings.Split(value, te[2]) {
				ps.Add(name, v)
			}
		case 3:
			// deepObject
			if namesRe.MatchString(name) {
				n := namesRe.ReplaceAllString(name, "")
				k := strings.Replace(name, n, "", -1)[1:]
				if len(k) > 0 {
					k = k[:len(k)-1]
					value = k + "=" + value
					name = n
				}
			}

			tuples := strings.Split(value, te[3])
			if len(tuples) <= 1 && !strings.Contains(value, te[4]) {
				ps.Set(name, value)
				continue
			}

			mm, ok := ps.Get(name).(map[string]any)
			if !ok {
				mm = make(map[string]any)
			}
			if te[4] == "" {
				for i := 0; i < len(tuples)/2; i++ {
					mm[tuples[i]] = tuples[i+1]
				}
			} else {
				for _, tuple := range tuples {
					kv := strings.Split(tuple, te[4])
					mm[kv[0]] = kv[1]
				}
			}

			ps.Set(name, mm)
		}
	}

	return nil
}
