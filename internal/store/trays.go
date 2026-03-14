package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexr/garden-app/internal/models"
)

func (s *SQLiteStore) AddTray(ctx context.Context, t *models.Tray) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`INSERT INTO trays (name, rows, cols) VALUES (?, ?, ?)`,
		t.Name, t.Rows, t.Cols,
	)
	if err != nil {
		return 0, fmt.Errorf("insert tray: %w", err)
	}
	id, _ := res.LastInsertId()

	for r := 0; r < t.Rows; r++ {
		for c := 0; c < t.Cols; c++ {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO tray_cells (tray_id, row, col) VALUES (?, ?, ?)`,
				id, r, c,
			); err != nil {
				return 0, fmt.Errorf("insert tray cell (%d,%d): %w", r, c, err)
			}
		}
	}
	return id, tx.Commit()
}

func (s *SQLiteStore) ListTrays(ctx context.Context) ([]models.Tray, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, rows, cols, created_at FROM trays ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("query trays: %w", err)
	}
	defer rows.Close()

	var trays []models.Tray
	for rows.Next() {
		var t models.Tray
		if err := rows.Scan(&t.ID, &t.Name, &t.Rows, &t.Cols, &t.CreatedAt); err != nil {
			return nil, err
		}
		trays = append(trays, t)
	}
	return trays, rows.Err()
}

func (s *SQLiteStore) GetTray(ctx context.Context, id int64) (*models.Tray, error) {
	var t models.Tray
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, rows, cols, created_at FROM trays WHERE id = ?`, id,
	).Scan(&t.ID, &t.Name, &t.Rows, &t.Cols, &t.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get tray: %w", err)
	}

	cellRows, err := s.db.QueryContext(ctx,
		`SELECT id, tray_id, row, col, seed_id, label, status,
		        sown_at, germinated_at, failed_at, notes
		 FROM tray_cells WHERE tray_id = ? ORDER BY row, col`, id)
	if err != nil {
		return nil, fmt.Errorf("query tray cells: %w", err)
	}
	defer cellRows.Close()

	t.Cells = make([][]models.TrayCell, t.Rows)
	for i := range t.Cells {
		t.Cells[i] = make([]models.TrayCell, t.Cols)
	}
	for cellRows.Next() {
		var c models.TrayCell
		if err := cellRows.Scan(
			&c.ID, &c.TrayID, &c.Row, &c.Col,
			&c.SeedID, &c.Label, &c.Status,
			&c.SownAt, &c.GerminatedAt, &c.FailedAt,
			&c.Notes,
		); err != nil {
			return nil, err
		}
		if c.Row >= 0 && c.Row < t.Rows && c.Col >= 0 && c.Col < t.Cols {
			t.Cells[c.Row][c.Col] = c
		}
	}
	return &t, cellRows.Err()
}

func (s *SQLiteStore) RemoveTray(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM trays WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("tray %d not found", id)
	}
	return nil
}

func (s *SQLiteStore) GetTrayCell(ctx context.Context, id int64) (*models.TrayCell, error) {
	var c models.TrayCell
	err := s.db.QueryRowContext(ctx,
		`SELECT id, tray_id, row, col, seed_id, label, status,
		        sown_at, germinated_at, failed_at, notes
		 FROM tray_cells WHERE id = ?`, id,
	).Scan(
		&c.ID, &c.TrayID, &c.Row, &c.Col,
		&c.SeedID, &c.Label, &c.Status,
		&c.SownAt, &c.GerminatedAt, &c.FailedAt,
		&c.Notes,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get tray cell %d: %w", id, err)
	}
	return &c, nil
}

func (s *SQLiteStore) SetTrayCell(ctx context.Context, c *models.TrayCell) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE tray_cells
		SET seed_id=?, label=?, status=?, sown_at=?, germinated_at=?, failed_at=?, notes=?
		WHERE id=?`,
		c.SeedID, c.Label, c.Status,
		c.SownAt, c.GerminatedAt, c.FailedAt,
		c.Notes, c.ID,
	)
	if err != nil {
		return fmt.Errorf("set tray cell: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("tray cell %d not found", c.ID)
	}
	return nil
}

func (s *SQLiteStore) ClearTrayCell(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE tray_cells
		SET seed_id=NULL, label='', status='empty',
		    sown_at=NULL, germinated_at=NULL, failed_at=NULL, notes=''
		WHERE id=?`, id)
	return err
}
