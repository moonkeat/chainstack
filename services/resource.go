package services

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/moonkeat/chainstack/models"
)

type ResourceService interface {
	CreateResource(userID int) (*models.Resource, error)
	ListResources(userID int) ([]models.Resource, error)
}

type resourceService struct {
	DB *sqlx.DB
}

func (s resourceService) CreateResource(userID int) (*models.Resource, error) {
	key := uuid.NewV4()
	createdAt := time.Now().UTC()
	_, err := s.DB.Query("INSERT INTO resources (key, created_at, user_id) VALUES ($1, $2, $3)", key.String(), createdAt, userID)
	if err != nil {
		return nil, err
	}

	return &models.Resource{Key: key.String(), CreatedAt: createdAt, UserID: userID}, nil
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
