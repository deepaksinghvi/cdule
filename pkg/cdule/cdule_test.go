package cdule

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_GetWorkerID(t *testing.T) {
	workerID := getWorkerID()
	require.NotEmpty(t, workerID)
}

func createScheduler() (error, Cdule) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	os.Remove(dirname + "/sqlite.db")

	cdule := Cdule{}
	cdule.NewCdule("./resources", "config_in_memory")
	return err, cdule
}
