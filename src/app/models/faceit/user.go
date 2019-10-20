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
}

type Lifetime struct {
	AverageHeadshots    string `json:"Average Headshots %"`
	AverageKD           string `json:"Average K/D Ratio"`
}
