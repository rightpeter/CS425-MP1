package tests

import (
	"log"
	"os/exec"
	"testing"
)

func TestEx(t *testing.T) {
	var command string
	command = "./client.bin \"nodes\""
	cmd := exec.Command(command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
	}

	t.Log("-> ", string(out))
	t.Error("Errro!!!")

}
