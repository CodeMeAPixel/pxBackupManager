package backup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"pxBackupManager/types"
)

// DiscordEmbedField represents a field in a Discord embed
type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// DiscordEmbed represents an embed in a Discord message
type DiscordEmbed struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Color       int                 `json:"color"`
	Fields      []DiscordEmbedField `json:"fields"`
	Timestamp   string              `json:"timestamp"`
}

// DiscordMessage represents a Discord webhook message
type DiscordMessage struct {
	Content string         `json:"content"`
	Embeds  []DiscordEmbed `json:"embeds"`
}

// SendDiscordNotification sends a backup result notification to Discord
func SendDiscordNotification(webhookURL string, results []types.BackupResult, summary string) error {
	if webhookURL == "" {
		return fmt.Errorf("discord webhook URL is empty")
	}

	// Determine color and title based on success/failure
	allSuccess := true
	for _, r := range results {
		if !r.Success {
			allSuccess = false
			break
		}
	}

	color := 3066993 // Green for success
	title := "✓ Backup Completed Successfully"
	if !allSuccess {
		color = 15158332 // Red for failure
		title = "✗ Backup Failed"
	}

	// Create embed
	embed := DiscordEmbed{
		Title:       title,
		Description: summary,
		Color:       color,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	// Add fields for each backup result
	for _, result := range results {
		status := "✓ Success"
		if !result.Success {
			status = "✗ Failed"
		}

		fieldValue := fmt.Sprintf("%s - %dms", status, result.Duration)
		if result.Size > 0 {
			fieldValue = fmt.Sprintf("%s\nSize: %.2f MB", fieldValue, float64(result.Size)/1024/1024)
		}
		if result.S3URL != "" {
			fieldValue = fmt.Sprintf("%s\nS3: %s", fieldValue, result.S3URL)
		}
		if result.Message != "" && !result.Success {
			fieldValue = fmt.Sprintf("%s\nError: %s", fieldValue, result.Message)
		}

		embed.Fields = append(embed.Fields, DiscordEmbedField{
			Name:   result.Service,
			Value:  fieldValue,
			Inline: true,
		})
	}

	// Create message
	message := DiscordMessage{
		Embeds: []DiscordEmbed{embed},
	}

	// Marshal to JSON
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal Discord message: %w", err)
	}

	// Send to Discord
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send Discord webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Discord webhook returned status %d", resp.StatusCode)
	}

	fmt.Println("Discord notification sent successfully")
	return nil
}
