package alflag

import (
	"encoding"
	"time"
)

var flagset *FlagSet

func init() {
	flagset = NewFlagSet("")
}

func Arg(i int) string {
	return flagset.Arg(i)
}

func Args() []string {
	return flagset.Args()
}

func Bool(name string, value bool, usage ...string) *bool {
	return flagset.Bool(name, value, usage...)
}

func BoolVar(p *bool, name string, value bool, usage ...string) {
	flagset.BoolVar(p, name, value, usage...)
}

func BoolFunc(name string, usage string, fn func(string) error) {
	flagset.BoolFunc(name, usage, fn)
}

func Duration(name string, value time.Duration, usage ...string) *time.Duration {
	return flagset.Duration(name, value, usage...)
}

func DurationVar(p *time.Duration, name string, value time.Duration, usage ...string) {
	flagset.DurationVar(p, name, value, usage...)
}

func Float64(name string, value float64, usage ...string) *float64 {
	return flagset.Float64(name, value, usage...)
}

func Float64Var(p *float64, name string, value float64, usage ...string) {
	flagset.Float64Var(p, name, value, usage...)
}

func Func(name string, usage string, fn func(string) error) {
	flagset.Func(name, usage, fn)
}

func Init(name string) {
	flagset.Init(name)
}

func Int(name string, value int, usage ...string) *int {
	return flagset.Int(name, value, usage...)
}

func IntVar(p *int, name string, value int, usage ...string) {
	flagset.IntVar(p, name, value, usage...)
}

func Int64(name string, value int64, usage ...string) *int64 {
	return flagset.Int64(name, value, usage...)
}

func Int64Var(p *int64, name string, value int64, usage ...string) {
	flagset.Int64Var(p, name, value, usage...)
}

func NArg() int {
	return flagset.NArg()
}

func NFlag() int {
	return flagset.NFlag()
}

func Name() string {
	return flagset.Name()
}

func Parse(arguments []string) error {
	return flagset.Parse(arguments)
}

func Parsed() bool {
	return flagset.Parsed()
}

func Set(name string, value string) error {
	return flagset.Set(name, value)
}

func String(name string, value string, usage ...string) *string {
	return flagset.String(name, value, usage...)
}

func StringVar(p *string, name string, value string, usage ...string) {
	flagset.StringVar(p, name, value, usage...)
}

func TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage ...string) {
	flagset.TextVar(p, name, value, usage...)
}

func Uint(name string, value uint, usage ...string) *uint {
	return flagset.Uint(name, value, usage...)
}

func UintVar(p *uint, name string, value uint, usage ...string) {
	flagset.UintVar(p, name, value, usage...)
}

func Uint64(name string, value uint64, usage ...string) *uint64 {
	return flagset.Uint64(name, value, usage...)
}

func Uint64Var(p *uint64, name string, value uint64, usage ...string) {
	flagset.Uint64Var(p, name, value, usage...)
}

func Var(value Value, name string, usage ...string) {
	flagset.Var(value, name, usage...)
}
