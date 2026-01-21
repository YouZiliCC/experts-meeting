package meeting

import (
	"experts-meeting/llm"
	"fmt"
	"strings"
)

// Expert represents an AI expert in the meeting
type Expert struct {
	Name   string
	Role   string
	Client *llm.Client
}

// NewExpert creates a new expert
func NewExpert(name, role string, client *llm.Client) *Expert {
	return &Expert{
		Name:   name,
		Role:   role,
		Client: client,
	}
}

// Speak generates expert's opinion
func (e *Expert) Speak(topic string, history []string, round int) (string, error) {
	prompt := e.buildPrompt(topic, history, round)
	return e.Client.Chat(prompt)
}

// Summarize generates final summary (for host)
func (e *Expert) Summarize(topic string, history []string) (string, error) {
	prompt := e.buildSummaryPrompt(topic, history)
	return e.Client.Chat(prompt)
}

func (e *Expert) buildPrompt(topic string, history []string, round int) string {
	context := strings.Join(history, "\n")
	return fmt.Sprintf(`You are a %s participating in an expert meeting.
This is round %d of the discussion.

Previous discussion:
%s

Please provide your brief professional opinion on this topic (within 100 words).`, 
		e.Role, round, context)
}

func (e *Expert) buildSummaryPrompt(topic string, history []string) string {
	context := strings.Join(history, "\n")
	return fmt.Sprintf(`You are a %s and the host of this meeting.

Discussion history:
%s

Please provide a concise summary of the discussion (within 200 words), highlighting key points and conclusions.`,
		e.Role, context)
}
