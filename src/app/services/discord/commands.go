package discord

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"ichor-stats/src/app/services/api"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!stats") {
		requesterID := GetRequesterID(m.Author.ID)
		url := "https://open.faceit.com/data/v4/players/" + requesterID + "/stats/csgo"

		// Create a Bearer string by appending string access token
		var bearer = "Bearer " + api.FACEIT_API_KEY

		// Create a new request using http
		req, err := http.NewRequest("GET", url, nil)

		// add authorization header to the req
		req.Header.Add("Authorization", bearer)

		// Send req using http Client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERRO] -", err)
		}

		var stats Stats
		err = json.NewDecoder(resp.Body).Decode(&stats)
		if err != nil {
			log.Println(err)
			return
		}

		url = "https://open.faceit.com/data/v4/players/" + requesterID
		req, err = http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", bearer)

		// Send req using http Client
		client = &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERRO] -", err)
		}

		var user User
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			log.Println(err)
			return
		}

		embed := discordgo.MessageEmbed{
			Title: user.Games.CSGO.Name,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "ELO",
					Value:  strconv.Itoa(user.Games.CSGO.ELO),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Skill Level",
					Value:  strconv.Itoa(user.Games.CSGO.SkillLevel),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Average K/D Ratio",
					Value:  stats.Lifetime.AverageKD,
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   "Average Headshots %",
					Value:  stats.Lifetime.AverageHeadshots,
					Inline: true,
				},
			},
		}
		_, err = s.ChannelMessageSendEmbed(api.CHANNEL_ID, &embed)
		if err != nil {
			log.Println(err)
		}
	}
}

func GetRequesterID(discordID string) string {
	if discordID == "210457267710066689" {
		return "0d94613d-b736-46ba-b8cd-d2159ddad705"
	} else if discordID == "210449893892947969" {
		return "b26df7d4-8517-4ec6-ab58-708487e5fe60"
	} else if discordID == "210438278623526913" {
		return "b0a57a5a-2f7a-481c-aaa8-8013a83378e3"
	}

	return ""
}

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

func SendMessage(message string) {
	discord, err := discordgo.New("Bot " + api.DISCORD_BOT_ID)
	if err != nil {
		fmt.Println(err)
	}

	_, err = discord.ChannelMessageSend(api.CHANNEL_ID, message)
	if err != nil {
		fmt.Println(err)
	}
}
