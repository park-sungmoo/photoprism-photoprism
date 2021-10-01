package search

import (
	"os"
	"testing"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	log = logrus.StandardLogger()
	log.SetLevel(logrus.DebugLevel)

	if err := os.Remove(".test.db"); err == nil {
		log.Debugln("removed .test.db")
	}

	db := entity.InitTestDb(os.Getenv("PHOTOPRISM_TEST_DRIVER"), os.Getenv("PHOTOPRISM_TEST_DSN"))
	defer db.Close()

	code := m.Run()

	os.Exit(code)
}
