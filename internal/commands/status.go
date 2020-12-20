package commands

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"
)

// StatusCommand performs a server health check.
var StatusCommand = cli.Command{
	Name:   "status",
	Usage:  "Performs a server health check",
	Action: statusAction,
}

// statusAction shows the server health status
func statusAction(ctx *cli.Context) error {
	conf := config.NewConfig(ctx)
	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("http://%s:%d/api/v1/status", conf.HttpHost(), conf.HttpPort())

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	var status string

	if resp, err := client.Do(req); err != nil {
		return fmt.Errorf("can't connect to %s:%d", conf.HttpHost(), conf.HttpPort())
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("server running at %s:%d, bad status %d\n", conf.HttpHost(), conf.HttpPort(), resp.StatusCode)
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	} else {
		status = string(body)
	}

	message := gjson.Get(status, "status").String()

	if message != "" {
		fmt.Println(message)
	} else {
		fmt.Println("unknown")
	}

	return nil
}
