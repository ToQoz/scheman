package scheman

import (
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestStopMysqld(t *testing.T) {
	log.Print("Stopping Mysqld")
	Mysqld.Stop()
}

func TestThisFileIsLastTestFile(t *testing.T) {
	cmd := exec.Command("go", "list", "-f", "{{.TestGoFiles}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err.Error())
	}

	tests := strings.Split(strings.Trim(string(output), " \n[]"), " ")
	lastTest := tests[len(tests)-1]

	if lastTest != "z_last_test.go" {
		t.Fatalf("expected last_test is z_last_test.go, but got %v", lastTest)
	}
}
