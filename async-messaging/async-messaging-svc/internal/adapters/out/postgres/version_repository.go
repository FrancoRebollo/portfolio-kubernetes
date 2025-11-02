package repository

type VersionRepository struct {
	dbPost PostgresDB
}

func NewVersionRepository(dbPost PostgresDB) *VersionRepository {
	return &VersionRepository{
		dbPost: dbPost,
	}
}
