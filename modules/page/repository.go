package page

import (
	"base-site-api/models"
	"github.com/jinzhu/gorm"
)

// Repository interface for Page model
type Repository interface {
	FindCategories() ([]*models.PageCategory, error)
	FindCategoryBySlug(slug string) (*models.PageCategory, error)
	FindBySlug(slug string) (*models.Page, []*models.Page, error)
	FindAllByCategorySlug(categorySlug string) ([]*models.Page, error)
	Update(page *models.Page, id uint) error
	Store(page *models.Page, userID uint) (uint, error)
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository return instance of repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db,
	}
}

// FindCategories return all page categories
func (r *repository) FindCategories() ([]*models.PageCategory, error) {
	var c []*models.PageCategory

	if err := r.db.Find(&c).Error; err != nil {
		return nil, err
	}

	return c, nil
}

// FindCategoryBySlug return category based on slug
func (r *repository) FindCategoryBySlug(slug string) (*models.PageCategory, error) {
	var c models.PageCategory

	if err := r.db.Model(&models.PageCategory{}).Where("slug = ?", slug).First(&c).Error; err != nil {
		return nil, err
	}

	return &c, nil
}

// FindAllByCategorySlug return pages for specific page category
func (r *repository) FindAllByCategorySlug(categorySlug string) ([]*models.Page, error) {
	var pages []*models.Page
	category, err := r.FindCategoryBySlug(categorySlug)
	if err != nil {
		return nil, err
	}

	if err := r.db.Set("gorm:auto_preload", true).Model(&models.Page{}).Where("category_id = ?", category.ID).Where("parent_page_id is null").Find(&pages).Error; err != nil {
		return nil, err
	}

	return pages, nil
}

// FindBySlug return page by slug
func (r *repository) FindBySlug(slug string) (*models.Page, []*models.Page, error) {
	var c models.Page
	var child []*models.Page

	if err := r.db.Model(c).Where("slug = ?", slug).Find(&c).Error; err != nil {
		return nil, nil, err
	}

	if err := r.db.Debug().Model(&models.Page{}).Where("parent_page_id = ?", c.ID).Find(&child).Error; err != nil {
		return &c, nil, err
	}

	return &c, child, nil
}

// Find return page by id
func (r *repository) Find(id uint) (*models.Page, error) {
	page := models.Page{}

	if err := r.db.First(&page, id).Error; err != nil {
		return nil, err
	}

	return &page, nil
}

// Update the page and validate input
func (r *repository) Update(page *models.Page, id uint) error {
	p, err := r.Find(id)
	if err != nil {
		return err
	}

	if page.Title != "" {
		p.Title = page.Title
	}

	if page.Body != "" {
		p.Body = page.Body
	}

	if page.Slug != "" {
		p.Slug = page.Slug
	}

	return r.db.Save(p).Error
}

// Store the page and validate input
func (r *repository) Store(page *models.Page, userID uint) (uint, error) {
	page.UserID = userID
	if err := r.db.Create(page).Error; err != nil {
		return 0, err
	}

	return page.ID, nil
}

// Delete the page by id
func (r *repository) Delete(id uint) error {
	if err := r.db.Where("id = ?", id).Delete(&models.Page{}).Error; err != nil {
		return err
	}

	return nil
}
