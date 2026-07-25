package bismarck

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type (
	Command struct {
		ctx         *Context
		name        string
		opts        []*option
		subcmdRule  *regexp.Regexp
		middlewares []MiddlewareFunc
		debug       bool `bismarck:"long=debug,desc=run as debug mode"`
		report      bool `bismarck:"long=report,desc=report when command is finished without error"`
		help        bool `bismarck:"short=h,long=help,desc=this help"`
	}
)

func New(name string) *Command {
	// Rules (defined in bismarck.go)
	ruleShortOpt = regexp.MustCompile(`^[\da-zA-Z]$`)
	ruleLongOpt = regexp.MustCompile(`^[\da-zA-Z][\w\-]{0,13}[\da-zA-Z]$`)

	cmd := &Command{
		ctx:        newContext(name, os.Args[1:]),
		name:       name,
		opts:       []*option{},
		subcmdRule: regexp.MustCompile(`^[\da-z][\da-z_\-]{0,13}[\da-z]$`),
		debug:      false,
		report:     false,
		help:       false,
	}
	cmd.ctx.cmd = cmd
	return cmd
}

func (cmd *Command) Use(middleware MiddlewareFunc) {
	cmd.middlewares = append(cmd.middlewares, middleware)
}

func (cmd *Command) parse() {
	v := reflect.ValueOf(cmd).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		if f.Kind() == reflect.Pointer {
			continue
		}
		tag := ft.Tag.Get(tagName)
		if tag == "" {
			tag = ft.Tag.Get(tagAlias)
			if tag == "" {
				continue
			}
		}
		opt := newOption(f, ft, cmd.ctx)
		cmd.opts = append(cmd.opts, opt)
	}
}

func (cmd *Command) Add(subCmd SubCommand, name, desc string) error {
	if !cmd.subcmdRule.MatchString(name) {
		return fmt.Errorf("sub command name is wrong (%s)", cmd.subcmdRule.String())
	}
	t := reflect.TypeOf(subCmd)
	if t.Kind() != reflect.Pointer {
		return fmt.Errorf(`pass a pointer type that implements the SubCommand interface`)
	}
	p := t.Elem()
	if p.Kind() == reflect.Pointer {
		return fmt.Errorf("double pointers are not allowed")
	}
	return cmd.addSubCmd(subCmd, name, desc, false)
}

func (cmd *Command) addSubCmd(subCmd SubCommand, name, desc string, force bool) error {
	bucket := newSubCommand(subCmd, name, desc)
	if bucket.Name == helpCmdName && !force {
		return fmt.Errorf(`"help" is reserved`)
	}
	bucket.ctx = cmd.ctx
	if _, ok := cmd.ctx.subCmds[bucket.Name]; !ok {
		cmd.ctx.scOrder = append(cmd.ctx.scOrder, bucket.Name)
	}
	cmd.ctx.subCmds[bucket.Name] = bucket
	return nil
}

func (cmd *Command) route(args []string) (int, error) {
	idx := -1
	arg := ""
	for idx, arg = range args {
		flg := false
		for _, o := range cmd.opts {
			if o.isShort(arg) || o.isLong(arg) {
				switch o.Name {
				case "debug":
					flg = true
					cmd.debug = true
				case "report":
					flg = true
					cmd.report = true
				case "help":
					flg = true
					cmd.help = true
				}
			}
		}
		if !flg && strings.HasPrefix(arg, "-") {
			return -1, fmt.Errorf(`unknown option "%s"`, arg)
		} else if !flg {
			return idx, nil
		}
	}
	return -1, nil
}

func (cmd *Command) showReport(ctx *Context) {
	subCmdTime := ctx.subCmdFinish.Sub(ctx.subCmdStart)
	cmdTime := time.Since(ctx.cmdStart)
	idx := -1
	for k, v := range ctx.rawArgs {
		if v == ctx.subCmd.Name {
			idx = k
			break
		}
	}
	gOpts := ctx.rawArgs[:idx]
	sOpts := ctx.rawArgs[idx+1:]
	buf := strings.Builder{}
	buf.WriteString("\n")
	buf.WriteString("------------------------------------------------------------\n")
	buf.WriteString(fmt.Sprintf(" %v Command\n", ctx.cmd.name))
	buf.WriteString("------------------------------------------------------------\n")
	buf.WriteString(fmt.Sprintf(" Options:    %v\n", strings.Join(gOpts, " ")))
	buf.WriteString(fmt.Sprintf(" SubCommand: %v\n", ctx.subCmd.Name))
	buf.WriteString(fmt.Sprintf(" SubOptions: %v\n", strings.Join(sOpts, " ")))
	buf.WriteString(fmt.Sprintf(" DateTime:   %v\n", ctx.cmdStart.In(ctx.loc)))
	buf.WriteString(fmt.Sprintf(" ExecTime:   %v\n", cmdTime))
	buf.WriteString(fmt.Sprintf(" SubTime:    %v\n", subCmdTime))
	buf.WriteString("------------------------------------------------------------\n")
	println(buf.String())
}

func (cmd *Command) applyMiddleware(h HandlerFunc) HandlerFunc {
	for i := len(cmd.middlewares) - 1; i >= 0; i-- {
		h = cmd.middlewares[i](h)
	}
	return h
}

func (cmd *Command) Run() error {
	var err error

	// Add Help Command
	err = cmd.addSubCmd(new(help), "help", "this help", true)
	if err != nil {
		return err
	}

	// Parse Options
	cmd.parse()

	// Routing
	rowArgs := cmd.ctx.RawArgs()
	subCmdIdx, err := cmd.route(rowArgs)
	if err != nil {
		return err
	}
	subCmdName := ""
	if cmd.help { // passed -h or --help option.
		subCmdName = helpCmdName
	} else if subCmdIdx < 0 { // passed no sub command.
		subCmdName = helpCmdName
	} else {
		subCmdName = cmd.ctx.RawArg(subCmdIdx)
	}

	// Debug Mode
	debugMode = cmd.debug

	// Copy Sub Command
	if _, ok := cmd.ctx.subCmds[subCmdName]; !ok {
		return fmt.Errorf(`unknown command "%s" is specified`, subCmdName)
	}
	cmd.ctx.subCmd = cmd.ctx.subCmds[subCmdName]

	// Init
	cmd.ctx.subCmd.cmd.Init()

	// parse arguments for global options
	cmd.ctx.subCmd.parse()

	// Parse Argument for Sub Command options
	cmd.ctx.args, err = cmd.ctx.subCmd.parseArgs(rowArgs[subCmdIdx+1:])
	if err != nil {
		return err
	}

	// Run
	cmd.ctx.subCmdStart = time.Now().UTC()
	h := cmd.applyMiddleware(cmd.ctx.subCmd.cmd.Run)
	err = h(cmd.ctx)
	cmd.ctx.subCmdFinish = time.Now().UTC()

	// Report
	if cmd.report && subCmdName != helpCmdName {
		cmd.showReport(cmd.ctx)
	}

	return err
}
