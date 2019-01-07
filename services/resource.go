package services

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/moonkeat/chainstack/models"
)

type ResourceService interface {
	ListResources(userID int) ([]models.Resource, error)
}

type resourceService struct {
	DB *sqlx.DB
}

func (s resourceService) ListResources(userID int) ([]models.Resource, error) {
	resources := []models.Resource{}
	err := s.DB.Select(&resources, "SELECT key, created_at FROM resources WHERE user_id = $1", userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return resources, nil
}

func NewResourceService(db *sqlx.DB) ResourceService {
	return &resourceService{
		DB: db,
	}
}
