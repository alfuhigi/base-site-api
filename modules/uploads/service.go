package uploads

import (
	"base-site-api/models"
	"base-site-api/modules"
	"base-site-api/storage"
	"github.com/gosimple/slug"
	"mime/multipart"
)

// Service interface for uploads
type Service interface {
	UploadCategories(typeSlug string) ([]*models.UploadCategory, error)
	UploadsByCategory(categorySlug string, page int, size int) ([]*models.Upload, int, error)
	Store(file *multipart.FileHeader, filename string, categorySlug string) (*storage.UploadFile, error)
	StoreCategory(categoryName string, subPath string, typeSlug string) (uint, error)
	UpdateCategory(categoryName string, subPath string, id uint) error
	Update(description string, id uint) error
	Delete(id uint) error
	DeleteCategory(id uint) error
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
		s3:         storage.NewS3(),
	}
}

type service struct {
	modules.Service
	repository Repository
	s3         *storage.S3Storage
}

// UploadCategories by type slug
func (s *service) UploadCategories(typeSlug string) ([]*models.UploadCategory, error) {
	return s.repository.FindCategoriesByType(typeSlug)
}

// Uploads by category slug
func (s *service) UploadsByCategory(categorySlug string, page int, size int) ([]*models.Upload, int, error) {
	l, o := s.CalculateLimitAndOffset(page, size)
	return s.repository.FindUploadsByCategory(categorySlug, l, o)
}

// Store upload the file and save the row to db with all information about the file itself
func (s *service) Store(file *multipart.FileHeader, filename string, categorySlug string) (*storage.UploadFile, error) {
	f, err := s.s3.Store(file, filename)

	if err != nil {
		return nil, err
	}

	c, err := s.repository.FindCategoryBySlug(categorySlug)

	if err != nil {
		return nil, err
	}

	u := models.Upload{
		File:       f.URL,
		Thumbnail:  f.URLSmall,
		CategoryID: c.ID,
	}

	_, err = s.repository.Store(&u)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (s *service) StoreCategory(categoryName string, subPath string, typeSlug string) (uint, error) {
	t, err := s.repository.FindTypeBySlug(typeSlug)
	if err != nil {
		return 0, err
	}

	c := models.UploadCategory{
		Name:    categoryName,
		SubPath: subPath,
		TypeID:  t.ID,
		Slug:    slug.Make(categoryName),
	}

	return s.repository.StoreCategory(&c)
}

// TODO: rename also in s3 if category name change
// UpdateCategory update the category it self and later also the s3 bucket
func (s *service) UpdateCategory(categoryName string, subPath string, id uint) error {
	c, err := s.repository.FindCategory(id)
	if err != nil {
		return err
	}

	if categoryName == "" {
		categoryName = c.Name
	}

	if subPath == "" {
		subPath = c.SubPath
	}

	return s.repository.UpdateCategory(categoryName, subPath, id)
}

func (s *service) Update(desc string, id uint) error {
	return s.repository.Update(desc, id)
}

func (s *service) Delete(id uint) error {
	return s.repository.Delete(id)
}

func (s *service) DeleteCategory(id uint) error {
	return s.repository.DeleteCategory(id)
}
