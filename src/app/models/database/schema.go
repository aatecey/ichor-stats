package database

type Player struct {
	Matches []Match `json:"matches"`
}

type Match struct {
	ID string `json:"id"`
	Map string `json:"map_name"`
	Result string `json:"result"`
	Score string `json:"score"`
	Kills string`json:"kills"`
	Assists string`json:"assists"`
	Deaths string`json:"deaths"`
	Headshots string `json:"headshots"`
	HeadshotPercentage string `json:"headshot_percentage"`
	Pentas string `json:"pentas"`
	Quads string `json:"quads"`
	Triples string `json:"triples"`
	KillDeathRatio string `json:"kill_death_ratio"`
	KillRoundRatio string `json:"kill_round_ratio"`
	MVPs string `json:"mvps"`
}