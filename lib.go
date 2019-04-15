package drift

import "log"

// Initialize will check for drift migrations table and create one if it doesn't exist
func Initialize(repo Repository) error {
	exists, err := repo.CheckForDriftMigrationsTable()
	if err != nil {
		return err
	}

	if !exists {
		log.Println("Migrations table does not exist... creating...")

		err := repo.CreateDriftMigrationsTable()
		if err != nil {
			return err
		}
	}
	return nil
}
