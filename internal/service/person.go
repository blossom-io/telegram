package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"blossom/pkg/key"
)

type Personer interface {
	GetTelegramByTwitchID(ctx context.Context, twitchID int64) (telegramID int64, telegramUsername string, err error)
	GetTelegramByTelegramUsernameOrID(ctx context.Context, telegramID int64, telegramUsername string) (int64, string, error)
	GetTwitchByTelegramUsernameOrID(ctx context.Context, telegramID int64, telegramUsername string) (int64, string, error)
	GetSubchatIDByInviteKey(ctx context.Context, inviteKey string) (int64, error)
	LinkTelegramToTwitchID(ctx context.Context, twitchID int64, telegramID int64, telergamUsername string) error
}

var _ Personer = (*service)(nil)

func (svc *service) GetSubchatIDByInviteKey(ctx context.Context, inviteKey string) (int64, error) {
	var ownerTwitchID int64

	k := key.ExtractKey(inviteKey)
	splitKey := strings.Split(k, ":")
	if len(splitKey) == 2 {
		ownerTwitchID, _ = strconv.ParseInt(splitKey[0], 10, 64)
	}

	chatTelegramID, err := svc.repo.IsSubchatExistsAndActive(ctx, ownerTwitchID)
	if err != nil {
		return 0, err
	}

	if chatTelegramID == 0 {
		return 0, fmt.Errorf("subchat is not active or not exists")
	}

	return chatTelegramID, nil
}

func (svc *service) LinkTelegramToTwitchID(ctx context.Context, twitchID int64, telegramID int64, telergamUsername string) error {
	return svc.repo.LinkTelegramToTwitchID(ctx, twitchID, telegramID, telergamUsername)
}

func (svc *service) GetTelegramByTwitchID(ctx context.Context, twitchID int64) (telegramID int64, telegramUsername string, err error) {
	return svc.repo.GetTelegramByTwitchID(ctx, twitchID)
}

func (svc *service) GetTelegramByTelegramUsernameOrID(ctx context.Context, telegramID int64, telegramUsername string) (int64, string, error) {
	return svc.repo.GetTelegramByTelegramUsernameOrID(ctx, telegramID, telegramUsername)
}

func (svc *service) GetTwitchByTelegramUsernameOrID(ctx context.Context, telegramID int64, telegramUsername string) (int64, string, error) {
	return svc.repo.GetTwitchByTelegramUsernameOrID(ctx, telegramID, telegramUsername)
}
