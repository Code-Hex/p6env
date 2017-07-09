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
	perl5     string
	rakudoURL string
	// Command line options
	List    bool
	Verbose bool
	As      string
}

func initInstall() *cobra.Command {
	install := new(Install)
	cmd := exec.Command("perl", "-e", `print $^X`)
	cmd.Env = os.Environ()
	perl5, err := cmd.CombinedOutput()
	if err != nil {
		panic(errors.Wrap(err, "Could not found perl5"))
	}
	install.perl5 = string(perl5)
	install.rakudoURL = git.GetRemote("github.com", "rakudo", "rakudo")

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
	b, err := git.LsRemote([]string{"--tags"}, i.rakudoURL)
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
