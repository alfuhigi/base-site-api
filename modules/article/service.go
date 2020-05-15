package article

import (
	"base-site-api/models"
)

// Service interface for Article model
type Service interface {
	Find(id uint) (*models.Article, error)
	FindAll(sort string) ([]*models.Article, error)
	Update(article *models.Article, id uint) error
	Store(article *models.Article, userID uint) (*models.Article, error)
	Delete(id uint, userID uint) error
}
