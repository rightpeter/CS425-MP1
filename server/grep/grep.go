package grep

import (
	"log"
	"os/exec"
)

// Grep funcion
func Grep(grepRules []string, path string) (string, error) {
	grepRules = append(grepRules, path)
	cmd := exec.Command("grep", grepRules...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
		return "", err
	}
	return string(out), nil
}
