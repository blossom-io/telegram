package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Personer interface {
	GetTelegramByTwitchID(ctx context.Context, twitchID int64) (telegramID int64, telegramUsername string, err error)
	GetTelegramByTelegramUsernameOrID(ctx context.Context, telegramID int64, telegramUsername string) (int64, string, error)
	GetTwitchByTelegramUsernameOrID(ctx context.Context, telegramID int64, telegramUsername string) (int64, string, error)
	LinkTelegramToTwitchID(ctx context.Context, twitchID int64, telegramID int64, telergamUsername string) error
}

var _ Personer = (*repository)(nil)

// LinkTelegramToTwitchID links telegram id and user to twitch user if not linked yet.
func (r *repository) LinkTelegramToTwitchID(ctx context.Context, twitchID int64, telegramID int64, telergamUsername string) error {
	q, a, err := r.DB.Sq.Update("person").
		Set("telegram_id", telegramID).
		Set("telegram_username", telergamUsername).
		Set("updated_at", time.Now()).
		Where("twitch_id = ?", twitchID).
		Where("telegram_id IS NULL").ToSql()
	if err != nil {
		return fmt.Errorf("LinkTelegramToTwitchID - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return fmt.Errorf("LinkTelegramToTwitchID - r.prepare: %w", err)
	}

	_, err = stmt.ExecContext(ctx, a...)
	if err != nil {
		return fmt.Errorf("LinkTelegramToTwitchID - Exec: %w", err)
	}

	return nil
}

// IsTwitchIDLinkedToTelegram checks if twitch id is linked to telegram id.
func (r *repository) GetTelegramByTwitchID(ctx context.Context, twitchID int64) (telegramID int64, telegramUsername string, err error) {
	q, a, err := r.DB.Sq.Select("telegram_id, telegram_username").
		From("person").Where("twitch_id = ?", twitchID).Limit(1).ToSql()
	if err != nil {
		return 0, "", fmt.Errorf("GetTelegramByTwitchID - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return 0, "", fmt.Errorf("GetTelegramByTwitchID - r.prepare: %w", err)
	}

	var (
		id   sql.NullInt64
		user sql.NullString
	)

	row := stmt.QueryRowContext(ctx, a...)
	err = row.Scan(&id, &user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", nil
		}
		return 0, "", fmt.Errorf("GetTelegramByTwitchID - row.Scan: %w", err)
	}

	return id.Int64, user.String, nil
}

// GetTelegramByTelegramUsernameOrID ...
func (r *repository) GetTelegramByTelegramUsernameOrID(ctx context.Context, telegramID int64, telegramUsername string) (int64, string, error) {
	q, a, err := r.DB.Sq.Select("telegram_id, telegram_username").
		From("person").Where("telegram_id = ? OR telegram_username = ?", telegramID, telegramUsername).Limit(1).ToSql()
	if err != nil {
		return 0, "", fmt.Errorf("GetTelegramByTelegramUsernameOrID - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return 0, "", fmt.Errorf("GetTelegramByTelegramUsernameOrID - r.prepare: %w", err)
	}

	var (
		id   sql.NullInt64
		user sql.NullString
	)

	row := stmt.QueryRowContext(ctx, a...)
	err = row.Scan(&id, &user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", nil
		}
		return 0, "", fmt.Errorf("GetTelegramByTelegramUsernameOrID - row.Scan: %w", err)
	}

	return id.Int64, user.String, nil
}

func (r *repository) GetTwitchByTelegramUsernameOrID(ctx context.Context, telegramID int64, telegramUsername string) (int64, string, error) {
	q, a, err := r.DB.Sq.Select("twitch_id, twitch_username").
		From("person").Where("telegram_id = ? OR telegram_username = ?", telegramID, telegramUsername).Limit(1).ToSql()
	if err != nil {
		return 0, "", fmt.Errorf("GetTwitchByTelegramUsernameOrID - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return 0, "", fmt.Errorf("GetTwitchByTelegramUsernameOrID - r.prepare: %w", err)
	}

	var (
		id   sql.NullInt64
		user sql.NullString
	)

	row := stmt.QueryRowContext(ctx, a...)
	err = row.Scan(&id, &user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", nil
		}
		return 0, "", fmt.Errorf("GetTwitchByTelegramUsernameOrID - row.Scan: %w", err)
	}

	return id.Int64, user.String, nil
}
