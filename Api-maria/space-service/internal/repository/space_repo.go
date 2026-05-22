package repository

import (
	"errors"

	"coworking/space-service/internal/models"

	"gorm.io/gorm"
)

// SpaceRepository defines the data-access contract for Space records.
type SpaceRepository interface {
	Create(space *models.Space) error
	GetAll(activeOnly bool) ([]models.Space, error)
	GetByID(id string) (*models.Space, error)
	Update(space *models.Space) error
	Delete(id string) error
	GetAvailable(spaceType string, minCapacity int) ([]models.Space, error)
}

type spaceRepo struct {
	db *gorm.DB
}

// NewSpaceRepository returns a concrete SpaceRepository backed by the given
// GORM database connection.
func NewSpaceRepository(db *gorm.DB) SpaceRepository {
	return &spaceRepo{db: db}
}

// Create inserts a new Space record.
func (r *spaceRepo) Create(space *models.Space) error {
	return r.db.Create(space).Error
}

// GetAll returns all spaces. When activeOnly is true, only active spaces are
// returned.
func (r *spaceRepo) GetAll(activeOnly bool) ([]models.Space, error) {
	var spaces []models.Space
	query := r.db.Model(&models.Space{})
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	if err := query.Order("created_at DESC").Find(&spaces).Error; err != nil {
		return nil, err
	}
	return spaces, nil
}

// GetByID fetches a single Space by its UUID primary key.
func (r *spaceRepo) GetByID(id string) (*models.Space, error) {
	var space models.Space
	if err := r.db.Where("id = ?", id).First(&space).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &space, nil
}

// Update persists changes to an existing Space record.
func (r *spaceRepo) Update(space *models.Space) error {
	return r.db.Save(space).Error
}

// Delete soft-deactivates a space by setting is_active = false instead of
// removing the row so that historical invoices retain referential integrity.
func (r *spaceRepo) Delete(id string) error {
	result := r.db.Model(&models.Space{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetAvailable returns active spaces filtered by optional spaceType and a
// minimum capacity threshold.
func (r *spaceRepo) GetAvailable(spaceType string, minCapacity int) ([]models.Space, error) {
	var spaces []models.Space
	query := r.db.Model(&models.Space{}).Where("is_active = ?", true)

	if spaceType != "" {
		query = query.Where("type = ?", spaceType)
	}
	if minCapacity > 0 {
		query = query.Where("capacity >= ?", minCapacity)
	}

	if err := query.Order("price_per_hour ASC").Find(&spaces).Error; err != nil {
		return nil, err
	}
	return spaces, nil
}
