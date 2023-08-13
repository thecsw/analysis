package discord

import "time"

// DiscordExport is the struct that represents the JSON file exported from Discord.
type DiscordExport struct {
	Guild      Guild     `json:"guild"`
	Channel    Channel   `json:"channel"`
	DateRange  DateRange `json:"dateRange"`
	ExportedAt time.Time `json:"exportedAt"`
	Messages   []Message `json:"messages"`
}

// DateRange is the struct that represents the date range object in the JSON file exported from Discord.
type DateRange struct {
	After  time.Time `json:"after"`
	Before time.Time `json:"before"`
}

// Guild is the struct that represents the guild object in the JSON file exported from Discord.
type Guild struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IconURL string `json:"iconUrl"`
}

// Channel is the struct that represents the channel object in the JSON file exported from Discord.
type Channel struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	CategoryID string `json:"categoryId"`
	Category   string `json:"category"`
	Name       string `json:"name"`
	Topic      string `json:"topic"`
}

// Message is the struct that represents the message object in the JSON file exported from Discord.
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

// Author is the struct that represents the author object in the JSON file exported from Discord.
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

// Role is the struct that represents the role object in the JSON file exported from Discord.
type Role struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Position int    `json:"position"`
}

// Embed is the struct that represents the embed object in the JSON file exported from Discord.
type Embed struct {
	Title       string       `json:"title"`
	URL         string       `json:"url"`
	Timestamp   time.Time    `json:"timestamp"`
	Description string       `json:"description"`
	Color       string       `json:"color"`
	Author      EmbedAuthor  `json:"author"`
	Image       EmbedImage   `json:"image"`
	Images      []EmbedImage `json:"images"`
	Fields      []Field      `json:"fields"`
}

// EmbedAuthor is the struct that represents the embed author object in the JSON file exported from Discord.
type EmbedAuthor struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// EmbedImage is the struct that represents the embed image object in the JSON file exported from Discord.
type EmbedImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Field is the struct that represents the field object in the JSON file exported from Discord.
type Field struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	IsInline bool   `json:"isInline"`
}

// Interaction is the struct that represents the interaction object in the JSON file exported from Discord.
type Interaction struct {
	ID   string          `json:"id"`
	Name string          `json:"name"`
	User InteractionUser `json:"user"`
}

// InteractionUser is the struct that represents the interaction user object in the JSON file exported from Discord.
type InteractionUser Mention

// Mention is the struct that represents the mention object in the JSON file exported from Discord.
type Mention struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Discriminator string `json:"discriminator"`
	Nickname      string `json:"nickname"`
	Color         string `json:"color"`
	IsBot         bool   `json:"isBot"`
	Roles         []Role `json:"roles"`
	AvatarURL     string `json:"avatarUrl"`
}

// Reaction is the struct that represents the reaction object in the JSON file exported from Discord.
type Reaction struct {
	Emoji Emoji `json:"emoji"`
	Count int   `json:"count"`
}

// Emoji is the struct that represents the emoji object in the JSON file exported from Discord.
type Emoji struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	IsAnimated bool   `json:"isAnimated"`
	ImageURL   string `json:"imageUrl"`
}

// Attachment is the struct that represents the attachment object in the JSON file exported from Discord.
type Attachment struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	FileName      string `json:"fileName"`
	FileSizeBytes int    `json:"fileSizeBytes"`
}

// Sticker is the struct that represents the sticker object in the JSON file exported from Discord.
type Sticker struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Format    string `json:"format"`
	SourceURL string `json:"sourceUrl"`
}
