package faceit

type User struct {
	Games       Games `json:"games"`
}

type Games struct {
	CSGO       CSGO `json:"csgo"`
}

type CSGO struct {
	SkillLevel int `json:"skill_level"`
	ELO        int `json:"faceit_elo"`
	Name       string `json:"game_player_name"`
}

type Stats struct {
	ID       string `json:"player_id"`
	Lifetime Lifetime `json:"lifetime"`
	Segment []SegmentStats `json:"segments"`
}

type Lifetime struct {
	AverageHeadshots    string `json:"Average Headshots %"`
	AverageKD           string `json:"Average K/D Ratio"`
	CurrentWinStreak    string `json:"Current Win Streak"`
	RecentResults		[]string `json:"Recent Results"`
}

type SegmentStats struct {
	CsMap       string `json:"label"`
	LifetimeMapStats    LifetimeMapStats `json:"stats"`
}

type LifetimeMapStats struct {
	Assists       string `json:"Assists"`
	Kills       string `json:"Kills"`
	Deaths       string `json:"Deaths"`
}