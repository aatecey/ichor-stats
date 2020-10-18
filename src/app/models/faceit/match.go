package faceit

type Matches struct {
	MatchItem []MatchItem `json:"items"`
}

type MatchItem struct {
	MatchId	string `json:"match_id"`
}

type Match struct {
	Rounds []Rounds `json:"rounds"`
}

type Rounds struct {
	MatchStats MatchStats `json:"round_stats"`
	Teams []Teams `json:"teams"`
}

type MatchStats struct {
	Map string `json:"Map"`
	Score string `json:"Score"`
}

type Teams struct {
	MatchStats MatchStats `json:"round_stats"`
	Players []Players `json:"players"`
}

type Players struct {
	ID string `json:"player_id"`
	Nickname string `json:"nickname"`
	Stats PlayerStats `json:"player_stats"`
}

type PlayerStats struct {
	Kills string `json:"Kills"`
	Assists string `json:"Assists"`
	Deaths string `json:"Deaths"`
	KD string `json:"K/D Ratio"`
	KR string `json:"K/R Ratio"`
	Result string `json:"Result"`
}
