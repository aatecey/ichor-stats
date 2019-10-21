package helpers

import "github.com/bwmarrin/discordgo"

type Embed struct {
	*discordgo.MessageEmbed
}

const (
	EmbedLimitDescription = 2048
	EmbedLimitFieldValue  = 1024
)

func NewEmbed() *Embed {
	return &Embed{&discordgo.MessageEmbed{}}
}

func (e *Embed) SetTitle(name string) *Embed {
	e.Title = name
	return e
}

func (e *Embed) SetDescription(description string) *Embed {
	if len(description) > EmbedLimitDescription {
		description = description[:EmbedLimitDescription]
	}
	e.Description = description
	return e
}

func (e *Embed) AddField(name, value string, inline bool) *Embed {
	if len(value) > EmbedLimitFieldValue {
		value = value[:EmbedLimitFieldValue]
	}

	if len(name) > EmbedLimitFieldValue {
		name = name[:EmbedLimitFieldValue]
	}

	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:  name,
		Value: value,
		Inline: inline,
	})

	return e
}

func (e *Embed) SetFooter(args ...string) *Embed {
	iconURL := ""
	text := ""
	proxyURL := ""

	switch {
	case len(args) > 2:
		proxyURL = args[2]
		fallthrough
	case len(args) > 1:
		iconURL = args[1]
		fallthrough
	case len(args) > 0:
		text = args[0]
	case len(args) == 0:
		return e
	}

	e.Footer = &discordgo.MessageEmbedFooter{
		IconURL:      iconURL,
		Text:         text,
		ProxyIconURL: proxyURL,
	}

	return e
}

func (e *Embed) SetImage(args ...string) *Embed {
	var URL string
	var proxyURL string

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		URL = args[0]
	}
	if len(args) > 1 {
		proxyURL = args[1]
	}
	e.Image = &discordgo.MessageEmbedImage{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return e
}

func (e *Embed) SetThumbnail(args ...string) *Embed {
	var URL string
	var proxyURL string

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		URL = args[0]
	}
	if len(args) > 1 {
		proxyURL = args[1]
	}
	e.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return e
}

func (e *Embed) SetAuthor(args ...string) *Embed {
	var (
		name     string
		iconURL  string
		URL      string
		proxyURL string
	)

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		name = args[0]
	}
	if len(args) > 1 {
		iconURL = args[1]
	}
	if len(args) > 2 {
		URL = args[2]
	}
	if len(args) > 3 {
		proxyURL = args[3]
	}

	e.Author = &discordgo.MessageEmbedAuthor{
		Name:         name,
		IconURL:      iconURL,
		URL:          URL,
		ProxyIconURL: proxyURL,
	}

	return e
}

func (e *Embed) SetURL(URL string) *Embed {
	e.URL = URL
	return e
}

func (e *Embed) SetColor(clr int) *Embed {
	e.Color = clr
	return e
}

func (e *Embed) InlineAllFields() *Embed {
	for _, v := range e.Fields {
		v.Inline = true
	}
	return e
}
