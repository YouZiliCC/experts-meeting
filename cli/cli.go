package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

type MapList []map[string]string

func (m *MapList) Set(value string) error {
	if err := json.Unmarshal([]byte(value), m); err != nil {
		log.Printf("Failed to unmarshal experts: %v\nValue: %s", err, value)
		return err
	}
	return nil
}

func (m *MapList) String() string {
	data, _ := json.Marshal(m)
	return string(data)
}

func CLIRun(args []string) error {
	var configPath, inputText, inputPath, outputPath string
	var discussionRounds, comperePosition int
	var experts MapList
	globalFlags := []cli.Flag{
		&cli.PathFlag{
			Name:        "config-path",
			Aliases:     []string{"c", "C"},
			Usage:       "Path to config `FILE`",
			Value:       "./EMConfig.json",
			Destination: &configPath,
			Action:      isFileExists,
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
			Name:   "experts",
			Usage:  "Specifies the LLM expert teams (Only recommended in the config file)",
			Value:  &experts,
			Hidden: true,
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
		Action: func(*cli.Context) error {
			// fmt.Print("run command", configPath, inputPath, outputPath)
			return nil
		},
	}
	commandCheckConfig := &cli.Command{
		Name:    "check-config",
		Aliases: []string{"check"},
		Usage:   "Check all the configuration file before running",
		Action: func(*cli.Context) error {
			return printConfig(configPath, inputText, inputPath, outputPath, discussionRounds, comperePosition, experts)
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
			fmt.Println(experts)
			return nil
		},
		Before: altsrc.InitInputSourceWithContext(globalFlags, altsrc.NewJSONSourceFromFlagFunc("config-path")),
	}
	return app.Run(args)
}

func printConfig(configPath, inputText, inputPath, outputPath string, discussionRounds, comperePosition int, experts MapList) error {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("  CONFIGURATION CHECK")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nFile Paths:")
	fmt.Printf("  ├─ Config File:  %s\n", configPath)
	fmt.Printf("  ├─ Input Text:   %s\n", inputText)
	fmt.Printf("  ├─ Input File:   %s\n", inputPath)
	fmt.Printf("  └─ Output File:  %s\n", outputPath)
	fmt.Println("\nMeeting Settings:")
	fmt.Printf("  ├─ Discussion Rounds:  %d\n", discussionRounds)
	fmt.Printf("  └─ Compere Position:   %d\n", comperePosition)
	fmt.Println("\nExperts Configuration:")
	if len(experts) == 0 {
		fmt.Println("  No experts configured!")
	} else {
		for i, expert := range experts {
			isCompere := (i + 1) == comperePosition
			prefix := "├─"
			if i == len(experts)-1 {
				prefix = "└─"
			}
			compereTag := ""
			if isCompere {
				compereTag = " [COMPERE]"
			}

			fmt.Printf("  %s Expert #%d%s\n", prefix, i+1, compereTag)
			fmt.Printf("  %s   Name:      %s\n", getSubPrefix(i, len(experts)), expert["name"])
			fmt.Printf("  %s   Role:      %s\n", getSubPrefix(i, len(experts)), expert["role"])
			fmt.Printf("  %s   Model:     %s\n", getSubPrefix(i, len(experts)), expert["model"])
			fmt.Printf("  %s   Base URL:  %s\n", getSubPrefix(i, len(experts)), expert["base_url"])
			fmt.Printf("  %s   API Key:   %s\n", getSubPrefix(i, len(experts)), maskAPIKey(expert["api_key"]))
			displayPrompt := expert["prompt"]
			if len(displayPrompt) > 60 {
				displayPrompt = displayPrompt[:60] + "..."
			}
			fmt.Printf("  %s   Prompt:    %s\n", getSubPrefix(i, len(experts)), displayPrompt)

			if i < len(experts)-1 {
				fmt.Printf("  %s\n", getSubPrefix(i, len(experts)))
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println()

	return nil
}

func getSubPrefix(index, total int) string {
	if index == total-1 {
		return " "
	}
	return "│"
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

func isFileExists(cCtx *cli.Context, path string) error {
	_, err := os.Stat(path)
	return err
}
