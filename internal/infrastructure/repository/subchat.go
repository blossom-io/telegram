package repository

import (
	"context"
	"fmt"
)

type Subchater interface {
	IsSubchatExistsAndActive(ctx context.Context, ownerTwitchID int64) (int64, error)
}

var _ Subchater = (*repository)(nil)

// IsSubchatExistsAndActive checks if subchat exists and active.
func (r *repository) IsSubchatExistsAndActive(ctx context.Context, ownerTwitchID int64) (chatTelegramID int64, err error) {
	q, a, err := r.DB.Sq.Select("chat_telegram_id").From("chat").
		Where("twitch_id = $1", ownerTwitchID).
		Suffix("AND NOT disabled").
		Suffix("AND chat_telegram_id IS NOT NULL").ToSql()
	if err != nil {
		return 0, fmt.Errorf("IsSubchatExistsAndActive - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return 0, fmt.Errorf("IsSubchatExistsAndActive - r.prepare: %w", err)
	}

	err = stmt.QueryRowContext(ctx, a...).Scan(&chatTelegramID)
	if err != nil {
		return 0, fmt.Errorf("IsSubchatExistsAndActive - Exec: %w", err)
	}

	return chatTelegramID, nil
}
