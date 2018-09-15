package grep

import (
	"log"
	"os/exec"
	"strings"
)

// Grep funcion
func Grep(grepRules string, path string) (string, error) {
	args := strings.Fields(grepRules)
	args = append(args, "/tmp/mp1.log")
	cmd := exec.Command("grep", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
		return "", err
	}
	return string(out), nil
}
