package alflag

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestFlagSet_parseName(t *testing.T) {
	type testcase struct {
		input     string
		short     string
		long      string
		wantPanic bool
	}

	tests := []testcase{
		{"a", "a", "", false},
		{"all", "", "all", false},
		{"a,all", "a", "all", false},
		{"a,all-tests", "a", "all-tests", false},

		{"1", "1", "", true},
		{"1all", "", "1all", true},
		{"vila*mq", "", "vila*mq", true},
		{"a,", "a", "", true},
		{",all", "", "all", true},
		{"all,a", "a", "all", true},
		{"1,all", "1", "all", true},
		{"a,1all", "a", "1all", true},
		{"a,all-", "a", "all-", true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil && !test.wantPanic {
					t.Errorf("parseName(%q) should not panic, but panicked: %v", test.input, r)
				}
				if r == nil && test.wantPanic {
					t.Errorf("parseName(%q) should panic, but did not panic", test.input)
				}
			}()

			fs := NewFlagSet("test")

			short, long := fs.parseName(test.input)
			if short != test.short || long != test.long {
				t.Errorf("parseName(%q) = %q, %q; want %q, %q", test.input, short, long, test.short, test.long)
			}
		})
	}
}

func TestFlagSet__setup(t *testing.T) {
	type testcase struct {
		name  string
		setup func(fs *FlagSet) map[string]any
		cases [][]string
		wants []map[string]any
	}

	tests := []testcase{
		{
			name: "Bool,BoolVar",
			setup: func(fs *FlagSet) map[string]any {
				t := fs.Bool("t,test", false, "")
				v := new(bool)
				fs.BoolVar(v, "v,value", true, "")
				return map[string]any{
					"test":  t,
					"value": v,
				}
			},
			cases: [][]string{
				{"-t"},
				{"-v=0"},
				{"-t=false", "-v=false"},
				{"--test", "--value=false"},
			},
			wants: []map[string]any{
				{"test": true, "value": true},
				{"test": false, "value": false},
				{"test": false, "value": false},
				{"test": true, "value": false},
			},
		},
		{
			name: "Duration,DurationVar",
			setup: func(fs *FlagSet) map[string]any {
				d := fs.Duration("d,duration", 0, "")
				d2 := new(time.Duration)
				fs.DurationVar(d2, "d2", 0, "")
				return map[string]any{
					"duration":  d,
					"duration2": d2,
				}
			},
			cases: [][]string{
				{"-d=0s"},
				{"-d2=0s"},
				{"-d2=0"},
				{"-d=1s", "-d2=1s"},
				{"--duration=1s", "--d2=1s"},
			},
			wants: []map[string]any{
				{"duration": 0 * time.Second, "duration2": 0 * time.Second},
				{"duration": 0 * time.Second, "duration2": 0 * time.Second},
				{"duration": 0 * time.Second, "duration2": 0 * time.Second},
				{"duration": 1 * time.Second, "duration2": 1 * time.Second},
				{"duration": 1 * time.Second, "duration2": 1 * time.Second},
			},
		},
		{
			name: "Float64,Float64Var",
			setup: func(fs *FlagSet) map[string]any {
				f := fs.Float64("f,float", 0, "")
				f2 := new(float64)
				fs.Float64Var(f2, "f2", 0, "")
				return map[string]any{
					"float":  f,
					"float2": f2,
				}
			},
			cases: [][]string{
				{"-f=0"},
				{"-f2=0"},
				{"-f=1.0", "-f2=1.0"},
				{"--float=1.0", "--f2=1.0"},
			},
			wants: []map[string]any{
				{"float": float64(0), "float2": float64(0)},
				{"float": float64(0), "float2": float64(0)},
				{"float": float64(1.0), "float2": float64(1.0)},
				{"float": float64(1.0), "float2": float64(1.0)},
			},
		},
		{
			name: "Int,IntVar",
			setup: func(fs *FlagSet) map[string]any {
				i := fs.Int("i,int", 0, "")
				i2 := new(int)
				fs.IntVar(i2, "i2", 0, "")
				return map[string]any{
					"int":  i,
					"int2": i2,
				}
			},
			cases: [][]string{
				{"-i=0"},
				{"-i2=0"},
				{"-i=1", "-i2=1"},
				{"--int=1", "--i2=2"},
			},
			wants: []map[string]any{
				{"int": int(0), "int2": int(0)},
				{"int": int(0), "int2": int(0)},
				{"int": int(1), "int2": int(1)},
				{"int": int(1), "int2": int(2)},
			},
		},
		{
			name: "Int64,Int64Var",
			setup: func(fs *FlagSet) map[string]any {
				i64 := fs.Int64("i,int64", 0, "")
				i642 := new(int64)
				fs.Int64Var(i642, "i642", 0, "")
				return map[string]any{
					"int64":  i64,
					"int642": i642,
				}
			},
			cases: [][]string{
				{"-i=0"},
				{"-i642=0"},
				{"-i=1", "-i642=1"},
				{"--int64=1", "--i642=1"},
			},
			wants: []map[string]any{
				{"int64": int64(0), "int642": int64(0)},
				{"int64": int64(0), "int642": int64(0)},
				{"int64": int64(1), "int642": int64(1)},
				{"int64": int64(1), "int642": int64(1)},
			},
		},
		{
			name: "String,StringVar",
			setup: func(fs *FlagSet) map[string]any {
				s := fs.String("s,string", "", "")
				s2 := new(string)
				fs.StringVar(s2, "s2", "", "")
				return map[string]any{
					"string":  s,
					"string2": s2,
				}
			},
			cases: [][]string{
				{"-s="},
				{"-s2="},
				{"-s=test", "-s2=test"},
				{"--string=test", "--s2=test"},
			},
			wants: []map[string]any{
				{"string": "", "string2": ""},
				{"string": "", "string2": ""},
				{"string": "test", "string2": "test"},
				{"string": "test", "string2": "test"},
			},
		},
		{
			name: "Uint,UintVar",
			setup: func(fs *FlagSet) map[string]any {
				u := fs.Uint("u,uint", 0, "")
				u2 := new(uint)
				fs.UintVar(u2, "u2", 0, "")
				return map[string]any{
					"uint":  u,
					"uint2": u2,
				}
			},
			cases: [][]string{
				{"-u=0"},
				{"-u2=0"},
				{"-u=1", "-u2=1"},
				{"--uint=1", "--u2=1"},
			},
			wants: []map[string]any{
				{"uint": uint(0), "uint2": uint(0)},
				{"uint": uint(0), "uint2": uint(0)},
				{"uint": uint(1), "uint2": uint(1)},
				{"uint": uint(1), "uint2": uint(1)},
			},
		},
		{
			name: "Uint64,Uint64Var",
			setup: func(fs *FlagSet) map[string]any {
				u64 := fs.Uint64("u,uint64", 0, "")
				u642 := new(uint64)
				fs.Uint64Var(u642, "u642", 0, "")
				return map[string]any{
					"uint64":  u64,
					"uint642": u642,
				}
			},
			cases: [][]string{
				{"-u=0"},
				{"-u642=0"},
				{"-u=1", "-u642=1"},
				{"--uint64=1", "--u642=22"},
			},
			wants: []map[string]any{
				{"uint64": uint64(0), "uint642": uint64(0)},
				{"uint64": uint64(0), "uint642": uint64(0)},
				{"uint64": uint64(1), "uint642": uint64(1)},
				{"uint64": uint64(1), "uint642": uint64(22)},
			},
		},
		{
			name: "TextVar",
			setup: func(fs *FlagSet) map[string]any {
				t := new(time.Time)
				fs.TextVar(t, "t,time", time.Now(), "")
				return map[string]any{
					"time": t,
				}
			},
			cases: [][]string{
				{"-t=2026-01-21T00:00:00Z"},
				{"--time=2026-01-21T00:00:00Z"},
			},
			wants: []map[string]any{
				{"time": time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC)},
				{"time": time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC)},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for i, c := range test.cases {
				t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
					defer func() {
						if r := recover(); r != nil {
							t.Errorf("Parse(%q) panicked: %v", c, r)
						}
					}()

					want := test.wants[i]

					fs := NewFlagSet(test.name)

					flagValues := test.setup(fs)
					err := fs.Parse(c)
					if err != nil {
						t.Errorf("Parse(%q) error: %v", c, err)
					}

					for name, value := range want {
						got, exists := flagValues[name]
						if !exists {
							t.Errorf("result[%q] does not exist", name)
						}
						if reflect.TypeOf(got).Kind() == reflect.Pointer {
							got = reflect.Indirect(reflect.ValueOf(got)).Interface()
						}
						if got != value {
							t.Errorf("result[%q] = %v; want %v", name, got, value)
						}
					}
				})
			}
		})
	}

}
