package db

import (
	"database/sql"
	"fmt"
	"time"
)

func (db *DB) CreateProject(name string) (*Project, error) {
	query := `INSERT INTO projects (name, updated_at) VALUES (?, ?)`
	result, err := db.conn.Exec(query, name, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return db.GetProject(int(id))
}

func (db *DB) GetProject(id int) (*Project, error) {
	query := `SELECT id, name, created_at, updated_at FROM projects WHERE id = ?`

	var p Project
	err := db.conn.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (db *DB) GetProjectByName(name string) (*Project, error) {
	query := `SELECT id, name, created_at, updated_at FROM projects WHERE name = ?`

	var p Project
	err := db.conn.QueryRow(query, name).Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (db *DB) ListProjects() ([]*Project, error) {
	query := `SELECT id, name, created_at, updated_at FROM projects ORDER BY name`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		var p Project
		err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}

	return projects, rows.Err()
}

func (db *DB) GetOrCreateProject(name string) (*Project, error) {
	project, err := db.GetProjectByName(name)
	if err == nil {
		return project, nil
	}

	return db.CreateProject(name)
}

func (db *DB) DeleteProject(id int) error {
	// Delete all references associated with the project first
	_, err := db.conn.Exec(`DELETE FROM citations WHERE project_id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete project references: %w", err)
	}

	// Delete the project
	result, err := db.conn.Exec(`DELETE FROM projects WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}
