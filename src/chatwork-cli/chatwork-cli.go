package main

import (
	"fmt"
	"os"
	"os/user"
	"io/ioutil"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/hashicorp/hcl"
	chatwork "github.com/yoppi/go-chatwork"
)

func main() {
	app := cli.NewApp()
	app.Name = "chatwork-cli"
	app.Version = Version
	app.Usage = "chatwork-cli [message]"
	app.Author = "Kentaro Kawano"
	app.Email = "kawano.kentaro@synergy101.jp"
	app.Flags = []cli.Flag {
		cli.StringFlag {
			Name: "roomid, r",
			Usage: "Room ID",
		},
	}
	app.Action = doMain
	app.Run(os.Args)
}

type Config struct {
	ApiKey string
}

// LoadConfig loads the CLI configuration from ".terraformrc" files.
func LoadConfig(path string) (*Config, error) {
	// Read the HCL file and prepare for parsing
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"Error reading %s: %s", path, err)
	}

	// Parse it
	obj, err := hcl.Parse(string(d))
	if err != nil {
		return nil, fmt.Errorf(
			"Error parsing %s: %s", path, err)
	}

	// Build up the result
	var result Config
	if err := hcl.DecodeObject(&result, obj); err != nil {
		return nil, err
	}

	return &result, nil
}

func doMain(c *cli.Context) {
	// get HomeDir
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Load Config
	config, err := LoadConfig(user.HomeDir + "/.chatworkrc")
	if err != nil {
		fmt.Println(err)
		return
	}

	rid := c.String("roomid")
	if rid == "" {
		fmt.Println("roomid is required")
		return
	}

	msg := ""
	if len(c.Args()) > 0 {
		msg = strings.Join(c.Args(), " ")
	} else {
		buf, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
			return
		}
		msg = string(buf)
	}

	if msg == "" {
		fmt.Println("no message")
		return
	}

	cw := chatwork.NewClient(config.ApiKey)
	cw.PostRoomMessage(rid, msg)
}
