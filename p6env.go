package p6env

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	version = "0.0.1"
	msg     = "Next generation Perl6 installer"
	name    = "p6env"
)

type p6env struct {
	*cobra.Command
	// Command line options
	Help       bool
	Version    bool
	StackTrace bool
}

func New() *p6env {
	p6env := new(p6env)
	p6env.rootCmdSetup()
	return p6env
}

func (p *p6env) rootCmdSetup() {
	p.Command = &cobra.Command{
		Use:           name,
		Short:         msg,
		Long:          msg,
		RunE:          p.run,
		SilenceErrors: true,
	}

	p.Command.AddCommand(initInstall())

	p.Flags().BoolVarP(&p.StackTrace, "trace", "t", false, "display detail error messages")
	p.Flags().BoolVarP(&p.Version, "version", "v", false, "print the version")
	p.Flags().BoolVarP(&p.Help, "help", "h", false, "print this message")
}

func (p *p6env) Run() int {
	if e := p.Command.Execute(); e != nil {
		exitCode, err := UnwrapErrors(e)
		if p.StackTrace {
			fmt.Fprintf(os.Stderr, "Error:\n  %+v\n", e)
		} else {
			fmt.Fprintf(os.Stderr, "Error:\n  %v\n", err)
		}
		return exitCode
	}
	return 0
}

func (p *p6env) run(cmd *cobra.Command, args []string) error {
	if p.Version {
		os.Stdout.WriteString(msg + "\n")
		return nil
	}
	return cmd.Usage()
}
