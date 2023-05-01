package repository

import (
	"context"
	"fmt"
	"time"
)

type Personer interface {
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

// // GetShow gets list of all added shows from db.
// func (r *repository) GetShows(ctx context.Context) ([]entity.Show, error) {
// 	q, _, err := r.DB.Sq.Select("title").From("show").PlaceholderFormat(sq.Question).ToSql()
// 	if err != nil {
// 		return nil, fmt.Errorf("showRepo - GetShows - r.Sq: %w", err)
// 	}

// 	rows, err := r.DB.DB.QueryContext(ctx, q)
// 	if err != nil {
// 		return nil, fmt.Errorf("showRepo - GetShows - Query: %w", err)
// 	}
// 	defer rows.Close()

// 	var shows []entity.Show

// 	for rows.Next() {
// 		e := entity.Show{}

// 		err = rows.Scan(&e.Title)
// 		if err != nil {
// 			return nil, fmt.Errorf("showRepo - GetShows - rows.Scan: %w", err)
// 		}

// 		shows = append(shows, e)
// 	}

// 	return shows, nil
// }

// // GetShowByID gets show by id from db.
// func (r *repository) GetShowByID(ctx context.Context, id int64) (show entity.Show, err error) {
// 	q, a, err := r.DB.Sq.Select("*").From("show").Where("id = ?", id).Limit(1).PlaceholderFormat(sq.Dollar).ToSql()
// 	if err != nil {
// 		return show, fmt.Errorf("showRepo - GetShowByID - r.Sq: %w", err)
// 	}

// 	rows, err := r.DB.DB.QueryContext(ctx, q, a...)
// 	if err != nil {
// 		return show, fmt.Errorf("showRepo - GetShowByID - Query: %w", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		err = rows.Scan(&show.Title)
// 		if err != nil {
// 			return show, fmt.Errorf("showRepo - GetShowByID - rows.Scan: %w", err)
// 		}
// 	}

// 	return show, nil
// }

// // AddShow adds show to db.
// func (r *repository) AddShow(ctx context.Context, show entity.Show) error {
// 	q, a, err := r.DB.Sq.Insert("show").
// 		Columns("tvmaze_id, title, location, network, genre, plot, quality, country_code, updated_at").
// 		Values(show.TVMazeID, show.Title, show.Location, show.Network, show.Genre, show.Plot, show.Quality, show.CountryCode, show.UpdatedAt).
// 		ToSql()
// 	if err != nil {
// 		return fmt.Errorf("ShowRepo - AddShow - r.Builder: %w", err)
// 	}

// 	_, err = r.DB.DB.ExecContext(ctx, q, a...)
// 	if err != nil {
// 		return fmt.Errorf("ShowRepo - AddShow - r.Sqlite.DB.ExecContext: %w", err)
// 	}

// 	return nil
// }
