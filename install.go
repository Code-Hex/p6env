package p6env

import (
	"os/exec"
	"strings"

	"os"

	"bytes"
	"path"

	"bufio"

	"github.com/Code-Hex/p6env/internal/git"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Install struct for "install" sub command
type Install struct {
	perl5    string
	rakudo   string
	backends []backend
	// Command line options
	List    bool
	Verbose bool
	As      string
}

// backend struct represent perl6 backend
type backend struct {
	name      string
	configure []string
	requires  []string
}

func initInstall() *cobra.Command {
	cmd := exec.Command("perl", "-e", `print $^X`)
	cmd.Env = os.Environ()
	perl5, err := cmd.CombinedOutput()
	if err != nil {
		panic(errors.Wrap(err, "Could not found perl5"))
	}

	perl := string(perl5)
	rakudo := git.GetRemote("github.com", "rakudo", "rakudo")
	moar := git.GetRemote("github.com", "MoarVM", "MoarVM")
	nqp := git.GetRemote("github.com", "perl6", "nqp")
	install := &Install{
		perl5:  perl,
		rakudo: rakudo,
		backends: []backend{
			backend{
				name: "jvm",
				configure: []string{
					perl, "Configure.pl", "--backends=jvm", "--gen-nqp", `--git-reference=\"$git_reference\"`, "--make-install",
				},
				requires: []string{rakudo, nqp},
			},
			backend{
				name: "moar",
				configure: []string{
					perl, "Configure.pl", "--backends=moar", "--gen-moar", `--git-reference=\"$git_reference\"`, "--make-install",
				},
				requires: []string{rakudo, nqp, moar},
			},
			backend{
				name: "moar-blead",
				configure: []string{
					perl, "Configure.pl", "--backends=moar", "--gen-moar=master", "--gen-nqp=master", `--git-reference=\"$git_reference\"`, "--make-install",
				},
				requires: []string{rakudo, nqp, moar},
			},
		},
	}

	installCmd := &cobra.Command{
		Use:           "install",
		Short:         "Install a Perl6 version",
		Long:          "Install a Perl6 version",
		RunE:          install.run,
		SilenceErrors: true,
	}

	installCmd.Flags().BoolVarP(&install.List, "list", "l", false, "List all available versions")
	installCmd.Flags().BoolVarP(&install.Verbose, "verbose", "v", false, "Verbose mode: print compilation status to stdout")
	installCmd.Flags().StringVar(&install.As, "as", "", "Install the definition as <NAME>")

	return installCmd
}

func (i *Install) run(cmd *cobra.Command, args []string) error {
	if i.List {
		list, err := i.available()
		if err != nil {
			return errors.Wrap(err, "Failed to get rakudo available list")
		}
		os.Stdout.Write(list)
		return nil
	}
	return cmd.Usage()
}

func (i *Install) available() ([]byte, error) {
	b, err := git.LsRemote([]string{"--tags"}, i.rakudo)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.WriteString("Available Rakudo versions:\n")

	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		sep := strings.Split(line, "\t")
		v := sep[1]
		if strings.HasSuffix(v, `^{}`) {
			base := path.Base(strings.TrimSuffix(v, `^{}`))
			if ('0' <= base[0] && base[0] <= '9') || base[0] == 'v' {
				buf.WriteString("  " + base + "\n")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
