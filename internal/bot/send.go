package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const sendMsgPath = "/send_msg"
const maxMessageLength = 4000

// ClientConfig is set by the bot Client before sending.
var ClientConfig struct {
	HTTPURL     string
	AccessToken string
}

// sendMessage sends a text message via the OneBot HTTP API.
func sendMessage(msgType string, targetID int64, groupID int64, text string) error {
	segments := splitMessage(text, maxMessageLength)
	for _, seg := range segments {
		if err := sendSingleMessage(msgType, targetID, groupID, seg); err != nil {
			return err
		}
	}
	return nil
}

func sendSingleMessage(msgType string, targetID int64, groupID int64, text string) error {
	payload := SendMessageRequest{
		MessageType: msgType,
		Message:     text,
	}
	if msgType == MessageTypePrivate {
		payload.UserID = targetID
	} else {
		payload.GroupID = groupID
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal send_msg payload: %w", err)
	}

	apiURL := ClientConfig.HTTPURL
	if apiURL == "" {
		apiURL = "http://127.0.0.1:3000"
	}
	apiURL = strings.TrimSuffix(apiURL, "/") + sendMsgPath

	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("build send_msg request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		if ClientConfig.AccessToken != "" {
			req.Header.Set("Authorization", "Bearer "+ClientConfig.AccessToken)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("send_msg request failed: %w", err)
			continue
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("send_msg returned status %d", resp.StatusCode)
			continue
		}
		return nil
	}
	return lastErr
}

func splitMessage(text string, limit int) []string {
	if len(text) <= limit {
		return []string{text}
	}
	var parts []string
	for len(text) > limit {
		parts = append(parts, text[:limit])
		text = text[limit:]
	}
	if text != "" {
		parts = append(parts, text)
	}
	return parts
}
