package coze

type RequestBody struct {
	Method string
	URL    string
	Host   string
	Token  string
	Body   any
}

// Coze

type CozeMessage struct {
	Role        string      `json:"role"`
	Type        string      `json:"type"`
	Content     string      `json:"content"`
	ContentType string      `json:"content_type"`
	ExtraInfo   interface{} `json:"extra_info"`
}

type CozeRequestBody struct {
	ConversationID string `json:"conversation_id"`
	BotID          string `json:"bot_id"`
	User           string `json:"user"`
	Query          string `json:"query"`
	Stream         bool   `json:"stream"`
}

type CozeResponseBody struct {
	Messages       []CozeMessage `json:"messages"`
	ConversationID string        `json:"conversation_id"`
	Code           int           `json:"code"`
	Msg            string        `json:"msg"`
}
