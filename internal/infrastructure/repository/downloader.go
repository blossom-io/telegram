package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Downloader interface {
	IsDownloaderEnabled(ctx context.Context, chatID int64) (bool, error)
}

var _ Downloader = (*repository)(nil)

func (r *repository) IsDownloaderEnabled(ctx context.Context, chatID int64) (enabled bool, err error) {
	q, a, err := r.DB.Sq.Select("is_downloader_enabled").From("settings").Where("chat_telegram_id = $1", chatID).ToSql()
	if err != nil {
		return false, fmt.Errorf("IsDownloaderEnabled - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return false, fmt.Errorf("IsDownloaderEnabled - r.prepare: %w", err)
	}

	err = stmt.QueryRowContext(ctx, a...).Scan(&enabled)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("IsDownloaderEnabled - stmt.QueryRowContext: %w", err)
	}

	return enabled, nil
}
