package db

import (
	"fmt"
	"strings"
	"time"
)

func (db *DB) AddReference(projectID int, bibtexEntry, apaFormat, sourceType string) (*Reference, error) {
	// Check if reference already exists
	exists, err := db.ReferenceExists(projectID, apaFormat)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing reference: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("reference already exists in this project")
	}

	// Get next reference number for this project
	var nextNum int
	err = db.conn.QueryRow(`SELECT COALESCE(MAX(reference_num), 0) + 1 FROM citations WHERE project_id = ?`, projectID).Scan(&nextNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get next reference number: %w", err)
	}

	query := `INSERT INTO citations (project_id, reference_num, bibtex_entry, apa_format, source_type) VALUES (?, ?, ?, ?, ?)`
	result, err := db.conn.Exec(query, projectID, nextNum, bibtexEntry, apaFormat, sourceType)
	if err != nil {
		return nil, fmt.Errorf("failed to add reference: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Reference{
		ID:           int(id),
		ProjectID:    projectID,
		ReferenceNum: nextNum,
		BibtexEntry:  bibtexEntry,
		APAFormat:    apaFormat,
		SourceType:   sourceType,
		CreatedAt:    time.Now(),
	}, nil
}

func (db *DB) ListReferences(projectID int) ([]*Reference, error) {
	query := `SELECT id, project_id, reference_num, bibtex_entry, apa_format, source_type, created_at 
	          FROM citations 
	          WHERE project_id = ? 
	          ORDER BY reference_num`
	
	rows, err := db.conn.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var references []*Reference
	for rows.Next() {
		var r Reference
		err := rows.Scan(&r.ID, &r.ProjectID, &r.ReferenceNum, &r.BibtexEntry, &r.APAFormat, &r.SourceType, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		references = append(references, &r)
	}

	return references, rows.Err()
}

func (db *DB) DeleteReference(id int) error {
	query := `DELETE FROM citations WHERE id = ?`
	_, err := db.conn.Exec(query, id)
	return err
}

func (db *DB) GetReference(id int) (*Reference, error) {
	query := `SELECT id, project_id, reference_num, bibtex_entry, apa_format, source_type, created_at FROM citations WHERE id = ?`
	
	var r Reference
	err := db.conn.QueryRow(query, id).Scan(&r.ID, &r.ProjectID, &r.ReferenceNum, &r.BibtexEntry, &r.APAFormat, &r.SourceType, &r.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("reference not found")
	}

	return &r, nil
}

func (db *DB) ReferenceExists(projectID int, apaFormat string) (bool, error) {
	query := `SELECT COUNT(*) FROM citations WHERE project_id = ? AND apa_format = ?`
	var count int
	err := db.conn.QueryRow(query, projectID, apaFormat).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (db *DB) GetReferencesByNumbers(projectID int, numbers []int) ([]*Reference, error) {
	if len(numbers) == 0 {
		return []*Reference{}, nil
	}
	
	// Build query with placeholders
	placeholders := make([]string, len(numbers))
	args := make([]interface{}, len(numbers)+1)
	args[0] = projectID
	
	for i, num := range numbers {
		placeholders[i] = "?"
		args[i+1] = num
	}
	
	query := fmt.Sprintf(`
		SELECT id, project_id, reference_num, bibtex_entry, apa_format, source_type, created_at 
		FROM citations 
		WHERE project_id = ? AND reference_num IN (%s)
		ORDER BY reference_num
	`, strings.Join(placeholders, ","))
	
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var references []*Reference
	for rows.Next() {
		var r Reference
		err := rows.Scan(&r.ID, &r.ProjectID, &r.ReferenceNum, &r.BibtexEntry, &r.APAFormat, &r.SourceType, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		references = append(references, &r)
	}
	
	return references, rows.Err()
}

func (db *DB) UpdateReferenceNumbers(projectID int) error {
	// Get all references for the project ordered by current reference_num
	query := `SELECT id, reference_num FROM citations WHERE project_id = ? ORDER BY reference_num`
	
	rows, err := db.conn.Query(query, projectID)
	if err != nil {
		return err
	}
	defer rows.Close()
	
	// Collect all references
	type refUpdate struct {
		id int
		newNum int
	}
	var updates []refUpdate
	newNum := 1
	
	for rows.Next() {
		var id, oldNum int
		if err := rows.Scan(&id, &oldNum); err != nil {
			return err
		}
		updates = append(updates, refUpdate{id: id, newNum: newNum})
		newNum++
	}
	
	if err := rows.Err(); err != nil {
		return err
	}
	
	// Update each reference with its new number
	updateQuery := `UPDATE citations SET reference_num = ? WHERE id = ?`
	for _, u := range updates {
		if _, err := db.conn.Exec(updateQuery, u.newNum, u.id); err != nil {
			return fmt.Errorf("failed to update reference %d: %w", u.id, err)
		}
	}
	
	return nil
}