package postgres

import (
	"github.com/kuzkuss/url_service/internal/link/repository"
	"github.com/kuzkuss/url_service/models"
	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type linkRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) repository.RepositoryI {
	return &linkRepository{
		db: db,
	}
}

func (dbLink *linkRepository) CreateLink(link *models.Link) error {
	tx := dbLink.db.Create(link)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table links)")
	}

	return nil
}

func (dbLink *linkRepository) SelectLinkByOriginalLink(originalLink string) (string, error) {
	link := models.Link{}

	tx := dbLink.db.Where("original_link = ?", originalLink).Take(&link)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return "", models.ErrNotFound
	} else if tx.Error != nil {
		return "", errors.Wrap(tx.Error, "database error (table links)")
	}

	return link.ShortLink, nil
}

func (dbLink *linkRepository) SelectLinkByShortLink(shortLink string) (string, error) {
	link := models.Link{}

	tx := dbLink.db.Where("short_link = ?", shortLink).Take(&link)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return "", models.ErrNotFound
	} else if tx.Error != nil {
		return "", errors.Wrap(tx.Error, "database error (table links)")
	}

	return link.OriginalLink, nil
}

