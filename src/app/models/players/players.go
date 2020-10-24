package players

type PlayerDetails struct {
	DiscordId string
	FaceitId  string
}

var Players = map[string]PlayerDetails {
	"b0a57a5a-2f7a-481c-aaa8-8013a83378e3": PlayerDetails{DiscordId: "210438278623526913", FaceitId: "gingajamie"},
	"0d94613d-b736-46ba-b8cd-d2159ddad705": PlayerDetails{DiscordId: "210457267710066689", FaceitId: "Tecey"},
	"b26df7d4-8517-4ec6-ab58-708487e5fe60": PlayerDetails{DiscordId: "210449893892947969", FaceitId: "gartlady"},
}