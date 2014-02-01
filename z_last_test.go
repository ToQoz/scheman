package scheman

import (
	"log"
	"testing"
)

func TestStopMysqld(t *testing.T) {
	log.Print("Stopping Mysqld")
	Mysqld.Stop()
}
