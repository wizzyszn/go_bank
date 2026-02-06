package db

import (
	"fmt"
	"log"
	"os"
)

func (db *DB) RunMigration(filepath string) error {

	content, err := os.ReadFile(filepath)

	if err != nil {
		return fmt.Errorf("Error reading migration file: %w", err)
	}
	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("Error execting migration: %w ", err)
	}

	log.Printf("Migration %s executed successfully", filepath)
	return nil
}

func (db *DB) ResetDatabase(schemaPath string) error {

	log.Println("WARNING: Resetting database - all data will be lost!")

	content, err := os.ReadFile(schemaPath)

	if err != nil {
		return fmt.Errorf("An Error occurred while read from the file: %w", err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("An Occured while restting the database %w", err)
	}

	log.Println("Database reset successfully")

	return nil
}
