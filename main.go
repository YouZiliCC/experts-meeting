package main

import (
	"experts-meeting/config"
	"experts-meeting/llm"
	"experts-meeting/meeting"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Parse flags
	configFile := flag.String("config", "meeting.conf", "Configuration file path")
	inputFile := flag.String("file", "", "Input text file path")
	inputText := flag.String("text", "", "Direct input text")
	hostOverride := flag.String("host", "", "Override host expert name")
	roundsOverride := flag.Int("rounds", 0, "Override discussion rounds")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Apply overrides
	if *hostOverride != "" {
		cfg.Host = *hostOverride
	}
	if *roundsOverride > 0 {
		cfg.Rounds = *roundsOverride
	}

	// Get topic
	topic, err := getTopic(*inputFile, *inputText)
	if err != nil {
		fmt.Printf("Failed to get topic: %v\n", err)
		os.Exit(1)
	}

	// Create experts
	experts := make([]*meeting.Expert, 0, len(cfg.Experts))
	for _, expertCfg := range cfg.Experts {
		client := llm.NewClient(expertCfg.APIKey, expertCfg.BaseURL, expertCfg.Model)
		expert := meeting.NewExpert(expertCfg.Name, expertCfg.Role, client)
		experts = append(experts, expert)
	}

	// Create and run session
	session, err := meeting.NewSession(topic, experts, cfg.Host, cfg.Rounds)
	if err != nil {
		fmt.Printf("Failed to create session: %v\n", err)
		os.Exit(1)
	}

	if err := session.Run(); err != nil {
		fmt.Printf("Meeting failed: %v\n", err)
		os.Exit(1)
	}
}

func getTopic(filePath, text string) (string, error) {
	if text != "" {
		return text, nil
	}

	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	return "", fmt.Errorf("please specify topic using -file or -text")
}
