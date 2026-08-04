package bismarck

import (
	"context"
	"database/sql"
	"time"
)

type (
	Context struct {
		cmd          *Command
		subCmd       *SubCommandBucket
		subCmds      map[string]*SubCommandBucket
		scOrder      []string
		loc          *time.Location
		db           *sql.DB
		cmdStart     time.Time
		subCmdStart  time.Time
		subCmdFinish time.Time
		rawArgs      []string
		args         []string
		stores       map[string]any
	}
)

func newContext(cmdName string, rawArgs []string) *Context {
	loc, _ := time.LoadLocation("UTC")
	return &Context{
		cmd:          nil,
		subCmd:       nil,
		subCmds:      map[string]*SubCommandBucket{},
		scOrder:      []string{},
		loc:          loc,
		db:           nil,
		cmdStart:     time.Now().UTC(),
		subCmdStart:  time.Unix(0, 0),
		subCmdFinish: time.Unix(0, 0),
		rawArgs:      rawArgs,
		args:         []string{},
		stores:       map[string]any{},
	}
}

func (ctx *Context) Location() *time.Location {
	return ctx.loc
}

func (ctx *Context) SetLocation(loc *time.Location) {
	ctx.loc = loc
}

func (ctx *Context) Database() (*sql.Conn, error) {
	return ctx.db.Conn(context.Background())
}

func (ctx *Context) SetDatabase(db *sql.DB) {
	ctx.db = db
}

func (ctx *Context) StartTime(command bool) time.Time {
	if command {
		return ctx.cmdStart.In(ctx.loc)
	} else {
		return ctx.subCmdStart.In(ctx.loc)
	}
}

func (ctx *Context) Arg(i int) string {
	if i >= 0 && i < len(ctx.args) {
		return ctx.args[i]
	}
	return ""
}

func (ctx *Context) Args() []string {
	return ctx.args
}

func (ctx *Context) NumArg() int {
	return len(ctx.args)
}

func (ctx *Context) RawArg(i int) string {
	if i >= 0 && i < len(ctx.rawArgs) {
		return ctx.rawArgs[i]
	}
	return ""
}

func (ctx *Context) RawArgs() []string {
	return ctx.rawArgs
}

func (ctx *Context) NumRawArg() int {
	return len(ctx.rawArgs)
}

func (ctx *Context) Get(key string) any {
	if v, ok := ctx.stores[key]; ok {
		return v
	} else {
		return nil
	}
}

func (ctx *Context) Set(key string, value any) {
	ctx.stores[key] = value
}
