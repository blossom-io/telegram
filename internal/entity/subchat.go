package entity

type Chat struct {
	Dates
	OwnerTwitchID  int64 `json:"ownerTwitchId"`
	ChatTelegramID int64 `json:"chatTelegramId"`
	Disabled       bool  `json:"disabled"`
}
