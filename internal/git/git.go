package git

import (
	"fmt"
	"os"
	"os/exec"
)

const git = "git"

var env = os.Environ()

func GetRemote(host, user, project string) string {
	protocol := os.Getenv("GIT_PROTOCOL")
	if protocol == "" {
		protocol = "git"
	}
	if protocol == "ssh" {
		return fmt.Sprintf("git@%s:%s/%s.git", host, user, project)
	}
	return fmt.Sprintf("%s://%s/%s/%s.git", protocol, host, user, project)
}

func LsRemote(opt []string, remote string) ([]byte, error) {
	args := make([]string, 0, len(opt)+2)
	args = append(args, "ls-remote")
	args = append(args, opt...)
	args = append(args, remote)

	cmd := exec.Command(git, args...)
	cmd.Env = env
	o, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return o, nil
}
