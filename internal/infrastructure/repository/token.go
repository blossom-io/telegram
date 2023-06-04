package repository

import (
	"context"
	"fmt"
)

type Tokener interface {
	GetOwnerIDsByInviteKey(ctx context.Context, inviteKey string) (ownerTwitchID int64, ownerTelegramID int64, err error)
}

var _ Tokener = (*repository)(nil)

// GetOwnerIDsByInviteKey checks if invite key exists and gets his twitch ID owner
func (r *repository) GetOwnerIDsByInviteKey(ctx context.Context, inviteKey string) (ownerTwitchID int64, ownerTelegramID int64, err error) {
	q, a, err := r.DB.Sq.Select("COALESCE(token.twitch_id, 0), COALESCE(person.telegram_id, 0)").
		From("token").
		LeftJoin("person ON token.twitch_id = person.twitch_id").
		Where("invite_key = $1", inviteKey).ToSql()
	if err != nil {
		return 0, 0, fmt.Errorf("GetOwnerIDsByInviteKey - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return 0, 0, fmt.Errorf("GetOwnerIDsByInviteKey - r.prepare: %w", err)
	}

	err = stmt.QueryRowContext(ctx, a...).Scan(&ownerTwitchID, &ownerTelegramID)
	if err != nil {
		return 0, 0, fmt.Errorf("GetOwnerIDsByInviteKey - Exec: %w", err)
	}

	return ownerTwitchID, ownerTelegramID, nil
}
