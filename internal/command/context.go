package command

import "log"

type ContextKey string

type Context struct {
	DryRun  bool
	Verbose bool
	Logger  *log.Logger
	dict    map[ContextKey]any
}

func NewContext() *Context {
	return &Context{
		Logger: log.Default(),
	}
}

func (c *Context) Logf(format string, args ...any) {
	if !c.Verbose {
		return
	}
	c.Logger.Printf(format, args...)
}

func (c *Context) Log(args ...any) {
	if !c.Verbose {
		return
	}
	c.Logger.Print(args...)
}

func (c *Context) Logln(args ...any) {
	if !c.Verbose {
		return
	}
	c.Logger.Println(args...)
}

func (c *Context) Set(key ContextKey, value any) {
	if c.dict == nil {
		c.dict = make(map[ContextKey]any)
	}
	c.dict[key] = value
}

func (c *Context) Get(key ContextKey) any {
	value, exists := c.dict[key]
	if !exists {
		return nil
	}
	return value
}
