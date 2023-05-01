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
func (r *repository) IsSubchatExistsAndActive(ctx context.Context, ownerTwitchID int64) (subchatTelegramID int64, err error) {
	q, a, err := r.DB.Sq.Select("subchat_telegram_id").From("subchat").
		Where("twitch_id = $1", ownerTwitchID).
		Suffix("AND NOT disabled").
		Suffix("AND subchat_telegram_id IS NOT NULL").ToSql()
	if err != nil {
		return 0, fmt.Errorf("IsSubchatExistsAndActive - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return 0, fmt.Errorf("IsSubchatExistsAndActive - r.prepare: %w", err)
	}

	err = stmt.QueryRowContext(ctx, a...).Scan(&subchatTelegramID)
	if err != nil {
		return 0, fmt.Errorf("IsSubchatExistsAndActive - Exec: %w", err)
	}

	return subchatTelegramID, nil
}
