package faceit

type Webhook struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	MatchID string `json:"id"`
}