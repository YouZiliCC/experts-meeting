package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ExpertConfig represents configuration for a single expert
type ExpertConfig struct {
	Name    string
	Role    string
	APIKey  string
	BaseURL string
	Model   string
}

// MeetingConfig represents the overall meeting configuration
type MeetingConfig struct {
	Rounds  int
	Host    string
	Experts []ExpertConfig
}

// Load reads configuration from file
func Load(path string) (*MeetingConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &MeetingConfig{
		Rounds:  2,
		Experts: make([]ExpertConfig, 0),
	}

	var currentExpert *ExpertConfig
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse section header [expert.Name]
		if strings.HasPrefix(line, "[expert.") && strings.HasSuffix(line, "]") {
			if currentExpert != nil {
				cfg.Experts = append(cfg.Experts, *currentExpert)
			}
			name := strings.TrimSuffix(strings.TrimPrefix(line, "[expert."), "]")
			currentExpert = &ExpertConfig{Name: name}
			continue
		}

		// Parse key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Global settings
		if currentExpert == nil {
			switch key {
			case "rounds":
				if r, err := strconv.Atoi(value); err == nil {
					cfg.Rounds = r
				}
			case "host":
				cfg.Host = value
			}
			continue
		}

		// Expert-specific settings
		switch key {
		case "role":
			currentExpert.Role = value
		case "api_key":
			currentExpert.APIKey = value
		case "base_url":
			currentExpert.BaseURL = value
		case "model":
			currentExpert.Model = value
		}
	}

	// Don't forget the last expert
	if currentExpert != nil {
		cfg.Experts = append(cfg.Experts, *currentExpert)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Validate
	if len(cfg.Experts) == 0 {
		return nil, fmt.Errorf("no experts configured")
	}

	return cfg, nil
}

// GetExpert returns expert config by name
func (c *MeetingConfig) GetExpert(name string) *ExpertConfig {
	for i := range c.Experts {
		if c.Experts[i].Name == name {
			return &c.Experts[i]
		}
	}
	return nil
}
