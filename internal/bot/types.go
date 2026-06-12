package bot

// Event represents a minimal OneBot v11 event payload.
type Event struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	SubType     string `json:"sub_type"`
	UserID      int64  `json:"user_id"`
	GroupID     int64  `json:"group_id"`
	RawMessage  string `json:"raw_message"`
	MessageID   int32  `json:"message_id"`
	SelfID      int64  `json:"self_id"`
	Sender      struct {
		Nickname string `json:"nickname"`
	} `json:"sender"`
}

// OneBot message type constants.
const (
	EventTypeMessage     = "message"
	MessageTypePrivate   = "private"
	MessageTypeGroup     = "group"
)

// SendMessageRequest is the payload for OneBot send_msg API.
type SendMessageRequest struct {
	MessageType string `json:"message_type"`
	UserID      int64  `json:"user_id,omitempty"`
	GroupID     int64  `json:"group_id,omitempty"`
	Message     string `json:"message"`
}
