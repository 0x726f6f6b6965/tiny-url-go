package protos

type ShortenedURL struct {
	Shorten   string `json:"shorten" dynamodbav:"shorten,omitempty"`
	Original  string `json:"original" dynamodbav:"original,omitempty"`
	Owner     string `json:"owner" dynamodbav:"owner,omitempty"`
	CreatedAt int64  `json:"created_at" dynamodbav:"created_at,omitempty"`
	ExpiresAt int64  `json:"expires_at" dynamodbav:"expires_at,omitempty"`
	UpdatedAt int64  `json:"updated_at" dynamodbav:"updated_at,omitempty"`
}
