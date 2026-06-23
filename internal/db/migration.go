package db

import "database/sql"

func migrate(conn *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS samples (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  node_id TEXT NOT NULL,
  seq INTEGER,
  measured_at TEXT NOT NULL,
  received_at TEXT NOT NULL,

  adxl_x REAL,
  adxl_y REAL,
  adxl_z REAL,
  adxl_mag REAL,
  adxl_rms REAL,
  adxl_peak REAL,

  piezo_raw INTEGER,
  piezo_min INTEGER,
  piezo_max INTEGER,
  piezo_peak INTEGER,

  created_at TEXT NOT NULL
)`,
		`CREATE INDEX IF NOT EXISTS idx_samples_node_measured
ON samples (node_id, measured_at)`,
		`CREATE INDEX IF NOT EXISTS idx_samples_measured
ON samples (measured_at)`,
	}

	for _, statement := range statements {
		if _, err := conn.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}
