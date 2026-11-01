package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"

	"github.com/krishhna24/procscope/internal/version"
)

const (
	exitOK    = 0
	exitError = 1 //runtime failure
	exitUsage = 2 //The user typed something we dont understand
)

const usage = `Usage: procscope <command> [flags]

  Commands:
    version    Print version and build information
    help       Show this help

  Run 'procscope <command> --help' for a command's flags.
  `

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	stop()
	os.Exit(code)
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	_ = ctx //used by commands that do real work (ps,top)
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return exitUsage
	}

	switch args[0] {
	case "version":
		return runVersion(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return exitOK
	default:
		fmt.Fprintf(stderr, "procscope: unknown command %q\n\n%s", args[0], usage)
		return exitUsage
	}
}

func runVersion(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.SetOutput(stderr)
	short := fs.Bool("short", false, "print only the version number")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return exitOK //pflag already printed the flag list
		}
		fmt.Fprintf(stderr, "procscope version: %v\n", err)
		fs.PrintDefaults()
		return exitUsage
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "procscope version: unexpected argument %q\n", fs.Arg(0))
		return exitUsage
	}
	info := version.Get()

	if *short {
		fmt.Fprintln(stdout, info.Version)
		return exitOK
	}
	fmt.Fprintln(stdout, info.String())
	return exitOK
}
