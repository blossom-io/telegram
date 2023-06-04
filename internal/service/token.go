package service

import "context"

type Tokener interface {
	GetOwnerIDsByInviteKey(ctx context.Context, inviteKey string) (ownerTwitchID int64, ownerTelegramID int64, err error)
}

func (svc *service) GetOwnerIDsByInviteKey(ctx context.Context, inviteKey string) (ownerTwitchID int64, ownerTelegramID int64, err error) {
	return svc.repo.GetOwnerIDsByInviteKey(ctx, inviteKey)
}
