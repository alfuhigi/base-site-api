package article

import (
	"base-site-api/models"
	// need to by for database
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/jinzhu/gorm"
)

// GormRepository implementation of repository with gorm.DB
type GormRepository struct {
	db *gorm.DB
}

// NewRepository return instance of GormRepository
func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db,
	}
}

// Find published article by id
func (r *GormRepository) Find(id uint) (*models.Article, error) {
	article := models.Article{}

	if err := r.db.Set("gorm:auto_preload", true).Where("published = 1").First(&article, id).Error; err != nil {
		return nil, err
	}

	return &article, nil
}

// FindBySlug published article by slug
func (r *GormRepository) FindBySlug(slug string) (*models.Article, error) {
	article := models.Article{}

	if err := r.db.Set("gorm:auto_preload", true).Where("published = 1").Where("slug = ?", slug).First(&article).Error; err != nil {
		return nil, err
	}

	return &article, nil
}

// FindAll list all published articles also with order
func (r *GormRepository) FindAll(order string, offset int, limit int) ([]*models.Article, int, error) {
	var articles []*models.Article
	var count int

	if err := r.db.Set("gorm:auto_preload", true).Model(&models.Article{}).Offset(offset).Limit(limit).Order(order).Where("published = 1").Find(&articles).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return articles, count, nil
}

// Update the article
func (r *GormRepository) Update(article *models.Article, id uint) error {
	a, err := r.Find(id)
	if err != nil {
		return err
	}
	if article.Title != "" {
		a.Title = article.Title
	}
	if article.Body != "" {
		a.Body = article.Body
	}
	if article.Short != "" {
		a.Short = article.Short
	}
	if article.Slug != "" {
		a.Slug = article.Slug
	}
	if article.Viewed != 0 {
		a.Viewed = article.Viewed
	}
	a.Published = article.Published

	return r.db.Save(a).Error
}

// Store new article in db and return ID
func (r *GormRepository) Store(article *models.Article, userID uint) (uint, error) {
	article.UserID = userID
	if err := r.db.Create(article).Error; err != nil {
		return 0, err
	}

	return article.ID, nil
}

// Delete article by ID
func (r *GormRepository) Delete(id uint) error {
	if err := r.db.Where("id = ?", id).Delete(&models.Article{}).Error; err != nil {
		return err
	}

	return nil
}
