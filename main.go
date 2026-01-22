package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

type MapList []map[string]string

func (m *MapList) Set(value string) error {
	return json.Unmarshal([]byte(value), m)
}
func (m *MapList) String() string {
	data, _ := json.Marshal(m)
	return string(data)
}

func main() {

	var configPath, inputText, inputPath, outputPath string
	var discussionRounds, comperePosition int
	var experts MapList
	globalFlags := []cli.Flag{
		&cli.PathFlag{
			Name:        "config",
			Aliases:     []string{"c", "C"},
			Usage:       "Path to config `FILE`",
			Value:       "./EMConfig.json",
			Destination: &configPath,
		},
		&cli.StringFlag{
			Name:        "input-text",
			Aliases:     []string{"t", "T"},
			Usage:       "Input `TEXT` for the meeting (e.g. topic or question)",
			Destination: &inputText,
		},
		altsrc.NewPathFlag(&cli.PathFlag{
			Name:        "input-path",
			Aliases:     []string{"i", "I"},
			Usage:       "Path to input `FILE`",
			Value:       "./input.txt",
			Destination: &inputPath,
		}),
		altsrc.NewPathFlag(&cli.PathFlag{
			Name:        "output-path",
			Aliases:     []string{"o", "O"},
			Usage:       "Path to output `FILE`",
			Value:       "./output.txt",
			Destination: &outputPath,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "rounds",
			Aliases:     []string{"r", "R"},
			Usage:       "Number of discussion rounds (If only 1 round, experts will just give their opinions without reading others' opinions)",
			Value:       2,
			Destination: &discussionRounds,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "compere-position",
			Aliases:     []string{"p", "P"},
			Usage:       "Position of the compere (Based on 1)",
			Value:       1,
			Destination: &comperePosition,
		}),
		altsrc.NewGenericFlag(&cli.GenericFlag{
			Name:  "experts",
			Usage: "Specifies the LLM expert teams (Only recommended in the config file)",
			Value: &experts,
		}),
	}
	commandRun := &cli.Command{
		Name:    "run",
		Aliases: []string{"start"},
		Usage:   "Run an experts meeting session",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v", "V"},
				Usage:   "Enable verbose output",
				Value:   false,
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Print("run command", configPath, inputPath, outputPath)
			return nil
		},
	}
	commandCheckConfig := &cli.Command{
		Name:  "check-config, check",
		Usage: "Check all the configuration file before running",
		Action: func(c *cli.Context) error {
			fmt.Print("status command", configPath, inputPath, outputPath)
			return nil
		},
	}
	commandVersion := &cli.Command{
		Name:  "version(developing)",
		Usage: "Print the version number of experts-meeting",
	}
	commandAdd := &cli.Command{
		Name:     "add(developing)",
		Category: "config-edit",
		Usage:    "Add a new expert to the configuration",
	}
	commandRemove := &cli.Command{
		Name:     "remove(developing)",
		Category: "config-edit",
		Usage:    "Remove an expert from the configuration",
	}
	app := &cli.App{
		Name:    "experts-meeting",
		Version: "v26.1.0",
		Authors: []*cli.Author{
			{
				Name:  "YouZiliCC",
				Email: "y2609984873@gmail.com",
			},
		},
		Usage:    "Simulate a meeting among AI experts using LLMs",
		Flags:    globalFlags,
		Commands: []*cli.Command{commandRun, commandCheckConfig, commandVersion, commandAdd, commandRemove},
		Action: func(*cli.Context) error {
			fmt.Print("here", configPath, inputPath, outputPath)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	// parse flags
	// altsrc.InitInputSourceWithContext(flags, altsrc.NewFileSource(configPath))

	// load config

	// organize prompt

	// create llm client

	// meeting logic

}
