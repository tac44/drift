package drift

// Repository represents a data provider
type Repository interface {
	CreateDriftMigrationsTable() error
	CheckForDriftMigrationsTable() (bool, error)
}
