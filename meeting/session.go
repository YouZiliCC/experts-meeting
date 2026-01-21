package meeting

import (
	"fmt"
)

// Session represents a meeting session
type Session struct {
	Topic   string
	Experts []*Expert
	Host    *Expert
	Rounds  int
	History []string
}

// NewSession creates a new meeting session
func NewSession(topic string, experts []*Expert, hostName string, rounds int) (*Session, error) {
	// Find host
	var host *Expert
	for _, expert := range experts {
		if expert.Name == hostName {
			host = expert
			break
		}
	}
	
	if host == nil {
		return nil, fmt.Errorf("host '%s' not found among experts", hostName)
	}

	return &Session{
		Topic:   topic,
		Experts: experts,
		Host:    host,
		Rounds:  rounds,
		History: []string{fmt.Sprintf("Topic: %s", topic)},
	}, nil
}

// Run executes the meeting
func (s *Session) Run() error {
	fmt.Printf("=== Meeting Started ===\n")
	fmt.Printf("Topic: %s\n", s.Topic)
	fmt.Printf("Experts: ")
	for i, expert := range s.Experts {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s (%s)", expert.Name, expert.Role)
	}
	fmt.Printf("\nHost: %s\n", s.Host.Name)
	fmt.Printf("Rounds: %d\n\n", s.Rounds)

	// Discussion rounds
	for round := 1; round <= s.Rounds; round++ {
		fmt.Printf("\n--- Round %d ---\n\n", round)

		for _, expert := range s.Experts {
			fmt.Printf("[%s] Speaking...\n", expert.Name)

			response, err := expert.Speak(s.Topic, s.History, round)
			if err != nil {
				fmt.Printf("[%s] Error: %v\n", expert.Name, err)
				continue
			}

			fmt.Printf("[%s]: %s\n\n", expert.Name, response)
			s.History = append(s.History, fmt.Sprintf("[Round %d][%s]: %s", round, expert.Name, response))
		}
	}

	// Host summary
	fmt.Printf("\n--- Host Summary ---\n\n")
	fmt.Printf("[%s] Summarizing...\n", s.Host.Name)

	summary, err := s.Host.Summarize(s.Topic, s.History)
	if err != nil {
		return fmt.Errorf("summary failed: %v", err)
	}

	fmt.Printf("[%s Summary]: %s\n", s.Host.Name, summary)
	fmt.Printf("\n=== Meeting Ended ===\n")

	return nil
}
