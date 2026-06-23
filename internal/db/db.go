package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/youhey/kendo-server/internal/model"
)

type Store struct {
	conn *sql.DB
}

func Open(path string) (*Store, error) {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create db directory: %w", err)
		}
	}

	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	conn.SetMaxOpenConns(1)

	store := &Store{conn: conn}
	if err := store.configure(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	if err := migrate(conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("migrate sqlite: %w", err)
	}

	return store, nil
}

func (s *Store) Close() error {
	return s.conn.Close()
}

func (s *Store) configure() error {
	statements := []string{
		`PRAGMA journal_mode = WAL`,
		`PRAGMA synchronous = NORMAL`,
		`PRAGMA busy_timeout = 5000`,
	}

	for _, statement := range statements {
		if _, err := s.conn.Exec(statement); err != nil {
			return fmt.Errorf("configure sqlite %q: %w", statement, err)
		}
	}

	return nil
}

func (s *Store) SaveSample(ctx context.Context, sample model.Sample) error {
	now := time.Now().Format(time.RFC3339Nano)
	_, err := s.conn.ExecContext(ctx, `INSERT INTO samples (
  node_id,
  seq,
  measured_at,
  received_at,
  adxl_x,
  adxl_y,
  adxl_z,
  adxl_mag,
  adxl_rms,
  adxl_peak,
  piezo_raw,
  piezo_min,
  piezo_max,
  piezo_peak,
  created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sample.NodeID,
		nullableInt(sample.Seq),
		sample.MeasuredAt,
		now,
		sample.ADXL.X,
		sample.ADXL.Y,
		sample.ADXL.Z,
		sample.ADXL.Mag,
		sample.ADXL.RMS,
		sample.ADXL.Peak,
		sample.Piezo.Raw,
		sample.Piezo.Min,
		sample.Piezo.Max,
		sample.Piezo.Peak,
		now,
	)
	if err != nil {
		return fmt.Errorf("insert sample: %w", err)
	}

	return nil
}

func (s *Store) RecentSamples(ctx context.Context, nodeID string, limit int) ([]model.SampleRecord, error) {
	const baseQuery = `SELECT
  id,
  node_id,
  seq,
  measured_at,
  received_at,
  adxl_x,
  adxl_y,
  adxl_z,
  adxl_mag,
  adxl_rms,
  adxl_peak,
  piezo_raw,
  piezo_min,
  piezo_max,
  piezo_peak,
  created_at
FROM samples`

	var (
		rows *sql.Rows
		err  error
	)
	if nodeID == "" {
		rows, err = s.conn.QueryContext(ctx, baseQuery+` ORDER BY measured_at DESC, id DESC LIMIT ?`, limit)
	} else {
		rows, err = s.conn.QueryContext(ctx, baseQuery+` WHERE node_id = ? ORDER BY measured_at DESC, id DESC LIMIT ?`, nodeID, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("query recent samples: %w", err)
	}
	defer rows.Close()

	records := make([]model.SampleRecord, 0)
	for rows.Next() {
		var (
			record model.SampleRecord
			seq    sql.NullInt64
		)
		if err := rows.Scan(
			&record.ID,
			&record.NodeID,
			&seq,
			&record.MeasuredAt,
			&record.ReceivedAt,
			&record.ADXL.X,
			&record.ADXL.Y,
			&record.ADXL.Z,
			&record.ADXL.Mag,
			&record.ADXL.RMS,
			&record.ADXL.Peak,
			&record.Piezo.Raw,
			&record.Piezo.Min,
			&record.Piezo.Max,
			&record.Piezo.Peak,
			&record.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan sample: %w", err)
		}
		if seq.Valid {
			value := seq.Int64
			record.Seq = &value
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate samples: %w", err)
	}

	return records, nil
}

func nullableInt(value *int64) any {
	if value == nil {
		return nil
	}
	return *value
}
