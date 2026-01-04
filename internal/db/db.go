package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

type Project struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Reference struct {
	ID            int
	ProjectID     int
	ReferenceNum  int    // Stable number within project
	BibtexEntry   string
	APAFormat     string
	SourceType    string
	CreatedAt     time.Time
}

func NewDB(path string) (*DB, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS citations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			reference_num INTEGER NOT NULL DEFAULT 0,
			bibtex_entry TEXT,
			apa_format TEXT NOT NULL,
			source_type TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_citations_project_id ON citations(project_id)`,
	}

	for _, query := range queries {
		if _, err := db.conn.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}
	
	// Check if reference_num column exists and add it if not
	var count int
	err := db.conn.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('citations') WHERE name='reference_num'`).Scan(&count)
	if err == nil && count == 0 {
		// Column doesn't exist, add it
		if _, err := db.conn.Exec(`ALTER TABLE citations ADD COLUMN reference_num INTEGER DEFAULT 0`); err != nil {
			return fmt.Errorf("failed to add reference_num column: %v", err)
		}
		
		// Update existing references with sequential numbers per project
		rows, err := db.conn.Query(`SELECT DISTINCT project_id FROM citations`)
		if err != nil {
			return err
		}
		defer rows.Close()
		
		var projectIDs []int
		for rows.Next() {
			var pid int
			if err := rows.Scan(&pid); err != nil {
				return err
			}
			projectIDs = append(projectIDs, pid)
		}
		
		for _, pid := range projectIDs {
			if _, err := db.conn.Exec(`
				UPDATE citations 
				SET reference_num = (
					SELECT COUNT(*) 
					FROM citations c2 
					WHERE c2.project_id = citations.project_id 
					AND c2.id <= citations.id
				)
				WHERE project_id = ?
			`, pid); err != nil {
				return err
			}
		}
	}
	
	return nil
}