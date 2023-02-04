package main

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
)

type testCase struct {
	style   string
	explode bool
	empty   string
	str     string
	array   string
	object  string
}

var testCases = []testCase{
	{
		"matrix",
		false,
		";color",
		";color=blue",
		";color=blue,black,brown",
		";color=R,100,G,200,B,150",
	},
	{
		"matrix",
		true,
		";color",
		";color=blue",
		";color=blue;color=black;color=brown",
		";R=100;G=200;B=150",
	},
	{"label", false, ".", ".blue", ".blue.black.brown", ".R.100.G.200.B.150"},
	{"label", true, ".", ".blue", ".blue.black.brown", ".R=100.G=200.B=150"},
	{
		"form",
		false,
		"color=",
		"color=blue",
		"color=blue,black,brown",
		"color=R,100,G,200,B,150",
	},
	{
		"form",
		true,
		"color=",
		"color=blue",
		"color=blue&color=black&color=brown",
		"R=100&G=200&B=150",
	},
	{"simple", false, "n/a", "blue", "blue,black,brown", "R,100,G,200,B,150"},
	{"simple", true, "n/a", "blue", "blue,black,brown", "R=100,G=200,B=150"},
	{
		"spaceDelimited",
		false,
		"n/a",
		"n/a",
		"blue%20black%20brown",
		"R%20100%20G%20200%20B%20150",
	},
	{"pipeDelimited", false, "n/a", "n/a", "blue|black|brown", "R|100|G|200|B|150"},
	{"deepObject", true, "n/a", "n/a", "n/a", "color[R]=100&color[G]=200&color[B]=150"},
	{
		"matrix",
		false,
		";color;age",
		";color=blue;age=12",
		";color=blue,black,brown;age=16,12,44",
		";color=R,100,G,200,B,150;age=OLD,100,YOUNG,1",
	},
	{
		"matrix",
		true,
		";color;age",
		";color=blue;age=12",
		";color=blue;color=black;color=brown;age=33;age=22;age=12",
		";R=100;G=200;B=150;YOUNG=1;OLD=150",
	},
	{
		"form",
		false,
		"color=&age=",
		"color=blue&age=12",
		"color=blue,black,brown&age=11,22,33",
		"color=R,100,G,200,B,150&age=O,1Y,22",
	},
	{
		"form",
		true,
		"color=&age=",
		"color=blue&age=12",
		"color=blue&color=black&color=brown&age=1&age=2&age=3",
		"R=100&G=200&B=150&Y=1&O=2",
	},
	{
		"deepObject",
		true,
		"n/a",
		"n/a",
		"n/a",
		"color[R]=100&color[G]=200&color[B]=150&age[YOUNG]=123&bonzo[ok][a][b]=3",
	},
}

func TestParse(t *testing.T) {
	wants := []Params{
		{"color": interface{}(nil)},
		{"color": "blue"},
		{"color": []interface{}{"blue", "black", "brown"}},
		{"color": map[string]interface{}{"100": "G", "G": "200", "R": "100"}},
		{"color": interface{}(nil)},
		{"color": "blue"},
		{"color": []interface{}{"blue", "black", "brown"}},
		{"B": "150", "G": "200", "R": "100"},
		{},
		{"blue": ""},
		{"black": []interface{}{""}, "blue": []interface{}{""}, "brown": []interface{}{""}},
		{"100": map[string]interface{}{}, "150": map[string]interface{}{}, "200": map[string]interface{}{}, "B": map[string]interface{}{}, "G": map[string]interface{}{}, "R": map[string]interface{}{}},
		{},
		{"": "blue"},
		{"": []interface{}{"blue", "black", "brown"}},
		{"": map[string]interface{}{"B": "150", "G": "200", "R": "100"}},
		{"color": interface{}(nil)},
		{"color": "blue"},
		{"color": []interface{}{"blue", "black", "brown"}},
		{"color": map[string]interface{}{"100": "G", "G": "200", "R": "100"}},
		{"color": interface{}(nil)},
		{"color": "blue"},
		{"color": []interface{}{"blue", "black", "brown"}},
		{"B": "150", "G": "200", "R": "100"},
		{"n/a": interface{}(nil)},
		{"blue": ""},
		{"black": []interface{}{""}, "blue": []interface{}{""}, "brown": []interface{}{""}},
		{"100": map[string]interface{}{}, "150": map[string]interface{}{}, "200": map[string]interface{}{}, "B": map[string]interface{}{}, "G": map[string]interface{}{}, "R": map[string]interface{}{}},
		{"": interface{}(nil)},
		{"": "blue"},
		{"": []interface{}{"blue", "black", "brown"}},
		{"": map[string]interface{}{"B": "150", "G": "200", "R": "100"}},
		{},
		{},
		{},
		{},
		{"n/a": interface{}(nil)},
		{"n/a": ""},
		{"black": []interface{}{""}, "blue": []interface{}{""}, "brown": []interface{}{""}},
		{"100": "", "150": "", "200": "", "B": "", "G": "", "R": ""},
		{"n/a": interface{}(nil)},
		{"n/a": ""},
		{"n/a": []interface{}{""}},
		{"color": map[string]interface{}{"B": "150", "G": "200", "R": "100"}},
		{"age": interface{}(nil), "color": interface{}(nil)},
		{"age": "12", "color": "blue"},
		{"age": []interface{}{"16", "12", "44"}, "color": []interface{}{"blue", "black", "brown"}},
		{"age": map[string]interface{}{"100": "YOUNG", "OLD": "100"}, "color": map[string]interface{}{"100": "G", "G": "200", "R": "100"}},
		{"age": interface{}(nil), "color": interface{}(nil)},
		{"age": "12", "color": "blue"},
		{"age": []interface{}{"33", "22", "12"}, "color": []interface{}{"blue", "black", "brown"}},
		{"B": "150", "G": "200", "OLD": "150", "R": "100", "YOUNG": "1"},
		{"age": interface{}(nil), "color": interface{}(nil)},
		{"age": "12", "color": "blue"},
		{"age": []interface{}{"11", "22", "33"}, "color": []interface{}{"blue", "black", "brown"}},
		{"age": map[string]interface{}{"O": "1Y"}, "color": map[string]interface{}{"100": "G", "G": "200", "R": "100"}},
		{"age": interface{}(nil), "color": interface{}(nil)},
		{"age": "12", "color": "blue"},
		{"age": []interface{}{"1", "2", "3"}, "color": []interface{}{"blue", "black", "brown"}},
		{"B": "150", "G": "200", "O": "2", "R": "100", "Y": "1"},
		{"n/a": interface{}(nil)},
		{"n/a": ""},
		{"n/a": []interface{}{""}},
		{"age": map[string]interface{}{"YOUNG": "123"}, "bonzo": map[string]interface{}{"ok][a][b": "3"}, "color": map[string]interface{}{"B": "150", "G": "200", "R": "100"}},
	}
	i := 0
	for _, cs := range testCases {
		inputs := []string{cs.empty, cs.str, cs.array, cs.object}
		for typ, input := range inputs {
			tname := fmt.Sprintf("%v%v%v", cs.style, cs.explode, typ)
			t.Run(tname, func(t *testing.T) {
				got := make(Params)
				if err := Parse(input, cs.style, cs.explode, ParamType(typ), got); err != nil {
					t.Errorf("Parse() error = %v ", err)
				}
				if diff := deep.Equal(wants[i], got); diff != nil {
					t.Errorf("Parse() diff = %v", diff)
				}
				i++
			})
		}
	}
}

func BenchmarkParse(b *testing.B) {
	for _, cs := range testCases {
		for typ, input := range []string{cs.empty, cs.str, cs.array, cs.object} {
			b.Run(fmt.Sprintf("%v%v%v", cs.style, cs.explode, typ), func(b *testing.B) {
				benchmarkParse(b, ParamType(typ), input, cs)
			})
		}
	}
}

func benchmarkParse(b *testing.B, typ ParamType, input string, cs testCase) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			params := make(Params)
			if err := Parse(input, cs.style, cs.explode, typ, params); err != nil {
				b.Errorf("Parse() error = %v ", err)
			}
		}
	})
}

func FuzzParse(f *testing.F) {
	for _, cs := range testCases {
		for typ, input := range []string{cs.empty, cs.str, cs.array, cs.object} {
			f.Add(typ, input, cs.style, cs.explode)
		}
	}
	f.Fuzz(func(t *testing.T, typ int, input string, style string, explode bool) {
		params := make(Params)
		if err := Parse(input, style, explode, ParamType(typ), params); err != nil {
			t.Fatalf("Parse() error = %v ", err)
		}
	})
}
