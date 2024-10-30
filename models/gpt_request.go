package models

type GptRequest struct {
	Model            string         `json:"model"`
	Messages         []Message      `json:"messages"`
	Temperature      int            `json:"temperature"`
	MaxTokens        int            `json:"max_tokens"`
	TopP             int            `json:"top_p"`
	FrequencyPenalty int            `json:"frequency_penalty"`
	PresencePenalty  int            `json:"presence_penalty"`
	ResponseFormat   ResponseFormat `json:"response_format"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
