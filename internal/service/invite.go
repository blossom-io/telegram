package service

import "context"

type Inviter interface {
	GetSubchatInviteLinkByTwitchID(ctx context.Context, subchatTelegramID int64, ownerTwitchID int64) (string, error)
	SetSubchatInviteLinkByTwitchID(ctx context.Context, subchatTelegramID int64, ownerTwitchID int64, inviteLink string) error
}

var _ Inviter = (*service)(nil)

func (svc *service) GetSubchatInviteLinkByTwitchID(ctx context.Context, subchatTelegramID int64, ownerTwitchID int64) (string, error) {
	return svc.repo.GetSubchatInviteLinkByTwitchID(ctx, subchatTelegramID, ownerTwitchID)
}

func (svc *service) SetSubchatInviteLinkByTwitchID(ctx context.Context, subchatTelegramID int64, ownerTwitchID int64, inviteLink string) error {
	return svc.repo.SetSubchatInviteLinkByTwitchID(ctx, subchatTelegramID, ownerTwitchID, inviteLink)
}
