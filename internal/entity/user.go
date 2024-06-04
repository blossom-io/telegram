package entity

type User struct {
	Dates
	TwitchID         int64  `json:"twitchId"`
	TwitchUsername   string `json:"twitchUsername"`
	TelegramID       int64  `json:"telegramId"`
	TelegramUsername string `json:"telegramUsername"`
	// ChatTelegramID   int64  `json:"chatTelegramId"`
}
