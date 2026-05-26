package telegram

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Sender is the interface used by the handler to send messages.
// Allows easy mocking in tests.
type Sender interface {
	Send(chatID int64, text string) error
}

// Client wraps the Telegram Bot API.
type Client struct {
	bot *botapi.BotAPI
}

// NewClient creates a new Client and validates the bot token.
func NewClient(token string) (*Client, error) {
	bot, err := botapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Client{bot: bot}, nil
}

// Send delivers a text message to the given chat ID.
func (c *Client) Send(chatID int64, text string) error {
	msg := botapi.NewMessage(chatID, text)
	_, err := c.bot.Send(msg)
	return err
}
