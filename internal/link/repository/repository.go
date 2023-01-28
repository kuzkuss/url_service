package repository

import (
	"github.com/kuzkuss/url_service/models"
)

type RepositoryI interface {
	SelectLinkByOriginalLink(originalLink string) (string, error)
	SelectLinkByShortLink(shortLink string) (string, error)
	CreateLink(link *models.Link) (error)
}
