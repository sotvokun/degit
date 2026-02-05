package alflag

import (
	"encoding"
	"flag"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

var (
	FlagSetShortNameRegexPattern = regexp.MustCompile(`^[a-zA-Z]$`)
	FlagSetLongNameRegexpPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-_]*[a-zA-Z0-9]+$`)
)

// FlagSet is a wrapper around [flag.FlagSet] with some restrictions to make it focused on command line parsing:
//   - [flag.FlagSet.Output] is set to [io.Discard].
//   - [flag.FlagSet.Usage] is set to an empty function.
//   - [flag.FlagSet.ErrorHandling] is set to [flag.ContinueOnError].
type FlagSet struct {
	flagset *flag.FlagSet
}

func NewFlagSet(name string) *FlagSet {
	flagset := flag.NewFlagSet(name, flag.ContinueOnError)
	flagset.SetOutput(io.Discard)
	flagset.Usage = func() {}

	return &FlagSet{flagset}
}

func (f *FlagSet) Arg(i int) string {
	return f.flagset.Arg(i)
}

func (f *FlagSet) Args() []string {
	return f.flagset.Args()
}

func (f *FlagSet) Bool(name string, value bool, usage ...string) *bool {
	ptr := new(bool)
	*ptr = value

	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.BoolVar(ptr, short, value, usageValue)
	}
	if long != "" {
		f.flagset.BoolVar(ptr, long, value, usageValue)
	}
	return ptr
}

func (f *FlagSet) BoolFunc(name string, usage string, fn func(string) error) {
	short, long := f.parseName(name)
	if short != "" {
		f.flagset.BoolFunc(short, usage, fn)
	}
	if long != "" {
		f.flagset.BoolFunc(long, usage, fn)
	}
}

func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage ...string) {
	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.BoolVar(p, short, value, usageValue)
	}
	if long != "" {
		f.flagset.BoolVar(p, long, value, usageValue)
	}
}

func (f *FlagSet) Duration(name string, value time.Duration, usage ...string) *time.Duration {
	ptr := new(time.Duration)
	*ptr = value

	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.DurationVar(ptr, short, value, usageValue)
	}
	if long != "" {
		f.flagset.DurationVar(ptr, long, value, usageValue)
	}
	return ptr
}

func (f *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage ...string) {
	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.DurationVar(p, short, value, usageValue)
	}
	if long != "" {
		f.flagset.DurationVar(p, long, value, usageValue)
	}
}

func (f *FlagSet) Float64(name string, value float64, usage ...string) *float64 {
	ptr := new(float64)
	*ptr = value

	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.Float64Var(ptr, short, value, usageValue)
	}
	if long != "" {
		f.flagset.Float64Var(ptr, long, value, usageValue)
	}
	return ptr
}

func (f *FlagSet) Float64Var(p *float64, name string, value float64, usage ...string) {
	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}
	short, long := f.parseName(name)
	if short != "" {
		f.flagset.Float64Var(p, short, value, usageValue)
	}
	if long != "" {
		f.flagset.Float64Var(p, long, value, usageValue)
	}
}

func (f *FlagSet) Func(name string, usage string, fn func(string) error) {
	short, long := f.parseName(name)
	if short != "" {
		f.flagset.Func(short, usage, fn)
	}
	if long != "" {
		f.flagset.Func(long, usage, fn)
	}
}

func (f *FlagSet) Init(name string) {
	f.flagset.Init(name, flag.ContinueOnError)
}

func (f *FlagSet) Int(name string, value int, usage ...string) *int {
	ptr := new(int)
	*ptr = value

	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.IntVar(ptr, short, value, usageValue)
	}
	if long != "" {
		f.flagset.IntVar(ptr, long, value, usageValue)
	}
	return ptr
}

func (f *FlagSet) IntVar(p *int, name string, value int, usage ...string) {
	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.IntVar(p, short, value, usageValue)
	}
	if long != "" {
		f.flagset.IntVar(p, long, value, usageValue)
	}
}

func (f *FlagSet) Int64(name string, value int64, usage ...string) *int64 {
	ptr := new(int64)
	*ptr = value

	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.Int64Var(ptr, short, value, usageValue)
	}
	if long != "" {
		f.flagset.Int64Var(ptr, long, value, usageValue)
	}
	return ptr
}

func (f *FlagSet) Int64Var(p *int64, name string, value int64, usage ...string) {
	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.Int64Var(p, short, value, usageValue)
	}
	if long != "" {
		f.flagset.Int64Var(p, long, value, usageValue)
	}
}

func (f *FlagSet) NArg() int {
	return f.flagset.NArg()
}

func (f *FlagSet) NFlag() int {
	return f.flagset.NFlag()
}

func (f *FlagSet) Name() string {
	return f.flagset.Name()
}

func (f *FlagSet) Parse(arguments []string) error {
	return f.flagset.Parse(arguments)
}

func (f *FlagSet) Parsed() bool {
	return f.flagset.Parsed()
}

func (f *FlagSet) Set(name string, value string) error {
	return f.flagset.Set(name, value)
}

func (f *FlagSet) String(name string, value string, usage ...string) *string {
	ptr := new(string)
	*ptr = value

	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.StringVar(ptr, short, value, usageValue)
	}
	if long != "" {
		f.flagset.StringVar(ptr, long, value, usageValue)
	}
	return ptr
}

func (f *FlagSet) StringVar(p *string, name string, value string, usage ...string) {
	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.StringVar(p, short, value, usageValue)
	}
	if long != "" {
		f.flagset.StringVar(p, long, value, usageValue)
	}
}

func (f *FlagSet) TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage ...string) {
	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.TextVar(p, short, value, usageValue)
	}
	if long != "" {
		f.flagset.TextVar(p, long, value, usageValue)
	}
}

func (f *FlagSet) Uint(name string, value uint, usage ...string) *uint {
	ptr := new(uint)
	*ptr = value

	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.UintVar(ptr, short, value, usageValue)
	}
	if long != "" {
		f.flagset.UintVar(ptr, long, value, usageValue)
	}
	return ptr
}

func (f *FlagSet) UintVar(p *uint, name string, value uint, usage ...string) {
	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.UintVar(p, short, value, usageValue)
	}
	if long != "" {
		f.flagset.UintVar(p, long, value, usageValue)
	}
}

func (f *FlagSet) Uint64(name string, value uint64, usage ...string) *uint64 {
	ptr := new(uint64)
	*ptr = value

	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.Uint64Var(ptr, short, value, usageValue)
	}
	if long != "" {
		f.flagset.Uint64Var(ptr, long, value, usageValue)
	}
	return ptr
}

func (f *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage ...string) {
	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}

	short, long := f.parseName(name)
	if short != "" {
		f.flagset.Uint64Var(p, short, value, usageValue)
	}
	if long != "" {
		f.flagset.Uint64Var(p, long, value, usageValue)
	}
}

func (f *FlagSet) Var(value Value, name string, usage ...string) {
	usageValue := ""
	if len(usage) > 0 {
		usageValue = usage[0]
	}
	short, long := f.parseName(name)
	if short != "" {
		f.flagset.Var(value, short, usageValue)
	}
	if long != "" {
		f.flagset.Var(value, long, usageValue)
	}
}

func (f *FlagSet) parseName(name string) (string, string) {
	parts := strings.SplitN(name, ",", 2)
	short := ""
	long := ""

	if len(parts) == 1 {
		name := strings.TrimSpace(parts[0])
		if FlagSetShortNameRegexPattern.MatchString(name) {
			short = name
		} else if FlagSetLongNameRegexpPattern.MatchString(name) {
			long = name
		} else {
			panic(fmt.Sprintf("invalid flag name: %s", name))
		}
		return short, long
	}

	short = strings.TrimSpace(parts[0])
	long = strings.TrimSpace(parts[1])
	if !FlagSetShortNameRegexPattern.MatchString(short) {
		panic(fmt.Sprintf("invalid short flag name: %s", short))
	}
	if !FlagSetLongNameRegexpPattern.MatchString(long) {
		panic(fmt.Sprintf("invalid long flag name: %s", long))
	}
	return short, long
}
