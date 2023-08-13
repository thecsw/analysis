package main

import "time"

type DiscordExport struct {
	Guild struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		IconURL string `json:"iconUrl"`
	} `json:"guild"`
	Channel struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		CategoryID string `json:"categoryId"`
		Category   string `json:"category"`
		Name       string `json:"name"`
		Topic      any    `json:"topic"`
	} `json:"channel"`
	DateRange struct {
		After  any `json:"after"`
		Before any `json:"before"`
	} `json:"dateRange"`
	ExportedAt time.Time `json:"exportedAt"`
	Messages   []Message `json:"messages"`
}

type Message struct {
	ID                 string       `json:"id"`
	Type               string       `json:"type"`
	Timestamp          time.Time    `json:"timestamp"`
	TimestampEdited    time.Time    `json:"timestampEdited"`
	CallEndedTimestamp time.Time    `json:"callEndedTimestamp"`
	IsPinned           bool         `json:"isPinned"`
	Content            string       `json:"content"`
	Author             Author       `json:"author"`
	Attachments        []Attachment `json:"attachments"`
	Embeds             []Embed      `json:"embeds"`
	Stickers           []Sticker    `json:"stickers"`
	Reactions          []Reaction   `json:"reactions"`
	Mentions           []Mention    `json:"mentions"`
	Interaction        Interaction  `json:"interaction"`
}

type Author struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Discriminator string `json:"discriminator"`
	Nickname      string `json:"nickname"`
	Color         string `json:"color"`
	IsBot         bool   `json:"isBot"`
	Roles         []Role `json:"roles"`
	AvatarURL     string `json:"avatarUrl"`
}

type Role struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Position int    `json:"position"`
}

type Embed struct {
	Title       string `json:"title"`
	URL         any    `json:"url"`
	Timestamp   any    `json:"timestamp"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Author      struct {
		Name string `json:"name"`
		URL  any    `json:"url"`
	} `json:"author"`
	Image struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"image"`
	Images []struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"images"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	IsInline bool   `json:"isInline"`
}

type Interaction struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	User struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		Discriminator string `json:"discriminator"`
		Nickname      string `json:"nickname"`
		Color         string `json:"color"`
		IsBot         bool   `json:"isBot"`
		Roles         []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Color    any    `json:"color"`
			Position int    `json:"position"`
		} `json:"roles"`
		AvatarURL string `json:"avatarUrl"`
	} `json:"user"`
}

type Mention struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Discriminator string `json:"discriminator"`
	Nickname      string `json:"nickname"`
	Color         any    `json:"color"`
	IsBot         bool   `json:"isBot"`
	Roles         []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Color    any    `json:"color"`
		Position int    `json:"position"`
	} `json:"roles"`
	AvatarURL string `json:"avatarUrl"`
}

type Reaction struct {
	Emoji struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Code       string `json:"code"`
		IsAnimated bool   `json:"isAnimated"`
		ImageURL   string `json:"imageUrl"`
	} `json:"emoji"`
	Count int `json:"count"`
}

type Attachment struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	FileName      string `json:"fileName"`
	FileSizeBytes int    `json:"fileSizeBytes"`
}

type Sticker struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Format    string `json:"format"`
	SourceURL string `json:"sourceUrl"`
}
