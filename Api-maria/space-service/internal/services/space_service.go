package services

import (
	"errors"

	"coworking/space-service/internal/models"
	"coworking/space-service/internal/repository"
)

// SpaceService encapsulates business logic for managing co-working spaces.
type SpaceService struct {
	repo repository.SpaceRepository
}

// NewSpaceService constructs a SpaceService with the provided repository.
func NewSpaceService(repo repository.SpaceRepository) *SpaceService {
	return &SpaceService{repo: repo}
}

// CreateSpaceInput holds the validated fields required to create a space.
type CreateSpaceInput struct {
	Name         string           `json:"name"           binding:"required"`
	Type         models.SpaceType `json:"type"           binding:"required"`
	Capacity     int              `json:"capacity"       binding:"required,min=1"`
	PricePerHour float64          `json:"price_per_hour" binding:"required,gt=0"`
	Amenities    string           `json:"amenities"`
}

// UpdateSpaceInput holds the fields that may be updated on an existing space.
type UpdateSpaceInput struct {
	Name         *string           `json:"name"`
	Type         *models.SpaceType `json:"type"`
	Capacity     *int              `json:"capacity"`
	PricePerHour *float64          `json:"price_per_hour"`
	Amenities    *string           `json:"amenities"`
	IsActive     *bool             `json:"is_active"`
}

// CreateSpace validates and persists a new Space record.
func (s *SpaceService) CreateSpace(input CreateSpaceInput) (*models.Space, error) {
	if err := validateSpaceType(input.Type); err != nil {
		return nil, err
	}

	space := &models.Space{
		Name:         input.Name,
		Type:         input.Type,
		Capacity:     input.Capacity,
		PricePerHour: input.PricePerHour,
		Amenities:    input.Amenities,
		IsActive:     true,
	}

	if err := s.repo.Create(space); err != nil {
		return nil, err
	}
	return space, nil
}

// GetSpaces returns all spaces. Pass activeOnly=true to exclude deactivated ones.
func (s *SpaceService) GetSpaces(activeOnly bool) ([]models.Space, error) {
	return s.repo.GetAll(activeOnly)
}

// GetSpace retrieves a single space by ID, returning an error if not found.
func (s *SpaceService) GetSpace(id string) (*models.Space, error) {
	space, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if space == nil {
		return nil, errors.New("space not found")
	}
	return space, nil
}

// UpdateSpace applies partial updates to an existing space.
func (s *SpaceService) UpdateSpace(id string, input UpdateSpaceInput) (*models.Space, error) {
	space, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if space == nil {
		return nil, errors.New("space not found")
	}

	if input.Name != nil {
		space.Name = *input.Name
	}
	if input.Type != nil {
		if err := validateSpaceType(*input.Type); err != nil {
			return nil, err
		}
		space.Type = *input.Type
	}
	if input.Capacity != nil {
		if *input.Capacity < 1 {
			return nil, errors.New("capacity must be at least 1")
		}
		space.Capacity = *input.Capacity
	}
	if input.PricePerHour != nil {
		if *input.PricePerHour <= 0 {
			return nil, errors.New("price_per_hour must be greater than 0")
		}
		space.PricePerHour = *input.PricePerHour
	}
	if input.Amenities != nil {
		space.Amenities = *input.Amenities
	}
	if input.IsActive != nil {
		space.IsActive = *input.IsActive
	}

	if err := s.repo.Update(space); err != nil {
		return nil, err
	}
	return space, nil
}

// DeleteSpace soft-deletes a space by ID.
func (s *SpaceService) DeleteSpace(id string) error {
	space, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if space == nil {
		return errors.New("space not found")
	}
	return s.repo.Delete(id)
}

// GetAvailableSpaces filters active spaces by optional type and minimum capacity.
func (s *SpaceService) GetAvailableSpaces(spaceType string, minCapacity int) ([]models.Space, error) {
	if spaceType != "" {
		if err := validateSpaceType(models.SpaceType(spaceType)); err != nil {
			return nil, err
		}
	}
	return s.repo.GetAvailable(spaceType, minCapacity)
}

// validateSpaceType returns an error when the provided type is not one of the
// defined SpaceType constants.
func validateSpaceType(t models.SpaceType) error {
	switch t {
	case models.SpaceTypeDesk, models.SpaceTypeOffice,
		models.SpaceTypeMeetingRoom, models.SpaceTypeLounge:
		return nil
	}
	return errors.New("invalid space type: must be DESK, OFFICE, MEETING_ROOM, or LOUNGE")
}
