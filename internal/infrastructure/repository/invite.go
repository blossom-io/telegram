package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Inviter interface {
	GetSubchatInviteLinkByTwitchID(ctx context.Context, chatTelegramID int64, ownerTwitchID int64) (string, error)
	SetSubchatInviteLinkByTwitchID(ctx context.Context, chatTelegramID int64, ownerTwitchID int64, inviteLink string) error
}

var _ Inviter = (*repository)(nil)

func (r *repository) GetSubchatInviteLinkByTwitchID(ctx context.Context, chatTelegramID int64, ownerTwitchID int64) (inviteLink string, err error) {
	q, a, err := r.DB.Sq.Select("subchat_telegram_invite_link").From("invite").
		Where("chat_telegram_id = $1", chatTelegramID).
		Where("twitch_id = $2", ownerTwitchID).ToSql()
	if err != nil {
		return "", fmt.Errorf("GetSubchatInviteLinkByTwitchID - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return "", fmt.Errorf("GetSubchatInviteLinkByTwitchID - r.prepare: %w", err)
	}

	err = stmt.QueryRowContext(ctx, a...).Scan(&inviteLink)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("GetSubchatInviteLinkByTwitchID - Exec: %w", err)
	}

	return inviteLink, nil
}

func (r *repository) SetSubchatInviteLinkByTwitchID(ctx context.Context, chatTelegramID int64, twitchID int64, inviteLink string) error {
	q, a, err := r.DB.Sq.Insert("invite").
		Columns("twitch_id", "chat_telegram_id", "subchat_telegram_invite_link").
		Values(twitchID, chatTelegramID, inviteLink).
		Suffix("ON CONFLICT (twitch_id, chat_telegram_id) DO UPDATE SET subchat_telegram_invite_link = EXCLUDED.subchat_telegram_invite_link").
		ToSql()
	if err != nil {
		return fmt.Errorf("SetSubchatInviteLinkByTwitchID - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return fmt.Errorf("SetSubchatInviteLinkByTwitchID - r.prepare: %w", err)
	}

	_, err = stmt.ExecContext(ctx, a...)
	if err != nil {
		return fmt.Errorf("SetSubchatInviteLinkByTwitchID - Exec: %w", err)
	}

	return nil
}
