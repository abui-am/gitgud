package ai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/user/gitgud/internal/config"
)

// AutocommitRules represents autocommit configuration rules
type AutocommitRules struct {
	Rules  string
	Source string
	Path   string
}

// Client handles OpenAI API interactions
type Client struct {
	client *openai.Client
	apiKey string
}

// NewClient creates a new AI client using the provided config manager
func NewClient(configManager *config.Manager) (*Client, error) {
	apiKey, err := configManager.GetOpenAIAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAI API key: %v", err)
	}

	client := openai.NewClient(apiKey)

	return &Client{
		client: client,
		apiKey: apiKey,
	}, nil
}

// GenerateCommitMessage generates a commit message based on the git diff
func (c *Client) GenerateCommitMessage(diff, customContext string) (string, error) {
	if diff == "" {
		return "", fmt.Errorf("no changes to commit")
	}

	// Get autocommit rules if they exist
	rules, _ := c.getAutocommitRules()

	// Build the system prompt
	systemPrompt := c.buildSystemPrompt(rules, customContext)

	// Build user prompt with diff
	userPrompt := fmt.Sprintf("Please generate a commit message for the following changes:\n\n```diff\n%s\n```", diff)

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			MaxTokens:   150,
			Temperature: 0.3,
		},
	)

	if err != nil {
		return "", fmt.Errorf("error generating commit message: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	message := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Clean up the message (remove quotes if present)
	message = strings.Trim(message, "\"'")

	return message, nil
}

// GenerateFileCommitMessage generates a commit message for a specific file
func (c *Client) GenerateFileCommitMessage(filename, diff, customContext string) (string, error) {
	if diff == "" {
		return "", fmt.Errorf("no changes found for file: %s", filename)
	}

	// Get autocommit rules if they exist
	rules, _ := c.getAutocommitRules()

	// Build the system prompt
	systemPrompt := c.buildSystemPrompt(rules, customContext)

	// Build user prompt with file-specific context
	userPrompt := fmt.Sprintf("Please generate a commit message for the changes in file '%s':\n\n```diff\n%s\n```", filename, diff)

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			MaxTokens:   150,
			Temperature: 0.3,
		},
	)

	if err != nil {
		return "", fmt.Errorf("error generating commit message for %s: %v", filename, err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI for file: %s", filename)
	}

	message := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Clean up the message (remove quotes if present)
	message = strings.Trim(message, "\"'")

	return message, nil
}

// buildSystemPrompt builds the system prompt for OpenAI
func (c *Client) buildSystemPrompt(rules AutocommitRules, customContext string) string {
	basePrompt := `You are an expert software developer assistant helping to generate concise, informative Git commit messages. 

Generate commit messages that:
1. Start with a clear, imperative verb (e.g., "Add", "Fix", "Update", "Remove", "Refactor")
2. Are concise but descriptive (50 characters or less for the title)
3. Focus on WHAT was changed and WHY when it's not obvious
4. Use conventional commit format when appropriate (feat:, fix:, docs:, etc.)
5. Avoid unnecessary words like "this", "the", "a" when possible

Examples of good commit messages:
- "Add user authentication middleware"
- "Fix memory leak in file processing"
- "Update API documentation for v2.0"
- "Refactor database connection handling"

Respond with ONLY the commit message, no quotes or additional text.`

	// Add custom rules if available
	if rules.Rules != "" {
		basePrompt += fmt.Sprintf("\n\nAdditional project-specific guidelines:\n%s", rules.Rules)
	}

	// Add custom context if provided
	if customContext != "" {
		basePrompt += fmt.Sprintf("\n\nAdditional context: %s", customContext)
	}

	return basePrompt
}

// getAutocommitRules loads autocommit rules from various sources
func (c *Client) getAutocommitRules() (AutocommitRules, error) {
	// Check for .autocommit.md in current directory
	autocommitFile := ".autocommit.md"
	if _, err := os.Stat(autocommitFile); err == nil {
		content, err := os.ReadFile(autocommitFile)
		if err == nil {
			return AutocommitRules{
				Rules:  string(content),
				Source: "current directory",
				Path:   autocommitFile,
			}, nil
		}
	}

	// Check for .autocommit.md in home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		homeAutocommitFile := filepath.Join(homeDir, ".autocommit.md")
		if _, err := os.Stat(homeAutocommitFile); err == nil {
			content, err := os.ReadFile(homeAutocommitFile)
			if err == nil {
				return AutocommitRules{
					Rules:  string(content),
					Source: "home directory",
					Path:   homeAutocommitFile,
				}, nil
			}
		}
	}

	// Check for autocommit rules in .gg directory
	ggDir := ".gg"
	ggAutocommitFile := filepath.Join(ggDir, "autocommit.md")
	if _, err := os.Stat(ggAutocommitFile); err == nil {
		content, err := os.ReadFile(ggAutocommitFile)
		if err == nil {
			return AutocommitRules{
				Rules:  string(content),
				Source: ".gg directory",
				Path:   ggAutocommitFile,
			}, nil
		}
	}

	return AutocommitRules{}, fmt.Errorf("no autocommit rules found")
}
