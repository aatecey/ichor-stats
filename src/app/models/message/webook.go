package message

type Webhook struct {
	Payload   Payload `json:"payload"`
	Requester string  `json:"third_party_id"`
}

type Payload struct {
	MatchID    string       `json:"id"`
	MatchTeams []MatchTeams `json:"teams"`
}

type MatchTeams struct {
	Name   string   `json:"name"`
	Roster []Roster `json:"roster"`
}

type Roster struct {
	Nickname   string `json:"nickname"`
	SkillLevel int    `json:"game_skill_level"`
	ID         string `json:"id"`
}
