package in_memory

import (
	"sync"

	"github.com/kuzkuss/url_service/internal/link/repository"
	"github.com/kuzkuss/url_service/models"
)

type linkRepository struct {
    mx sync.RWMutex
    store  map[string]string
}

func New() repository.RepositoryI {
	return &linkRepository {
		store: make(map[string]string),
	}
}

func (dbLink *linkRepository) CreateLink(link *models.Link) error {
	dbLink.mx.Lock()
    dbLink.store[link.ShortLink] = link.OriginalLink
	dbLink.mx.Unlock()
	return nil
}

func (dbLink *linkRepository) SelectLinkByOriginalLink(originalLink string) (string, error) {
	for key, val := range dbLink.store {
		if val == originalLink {
			return key, nil
		}
	}

	return "", models.ErrNotFound
}

func (dbLink *linkRepository) SelectLinkByShortLink(shortLink string) (string, error) {
	dbLink.mx.RLock()
    val, ok := dbLink.store[shortLink]
	dbLink.mx.RUnlock()
	if !ok {
		return "", models.ErrNotFound
	}
    return val, nil
}


