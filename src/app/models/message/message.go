package message

type Match struct {
	Player			string      `json:"player_id"`
	Result			string      `json:"result"`
	Map				string		`json:"map"`
	Score			string		`json:"score"`
	Kills			string		`json:"kills"`
	Assists			string		`json:"assists"`
	Deaths			string		`json:"deaths"`
	KillDeathRatio	string		`json:"killDeathRatio"`
	KillRoundRatio	string		`json:"killRoundRatio"`
}