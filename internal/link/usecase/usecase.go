package usecase

import (
	"github.com/pkg/errors"
	"crypto/sha256"
	"math/big"
	"math/rand"

	linkRep "github.com/kuzkuss/url_service/internal/link/repository"
	"github.com/kuzkuss/url_service/models"
)

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")

type UseCaseI interface {
	GetOriginalLink(link string) (string, error)
	CreateShortLink(link *models.Link) (error)
}

type useCase struct {
	linkRepository linkRep.RepositoryI
}

func New(linkRepository linkRep.RepositoryI) UseCaseI {
	return &useCase{
		linkRepository: linkRepository,
	}
}

func (uc *useCase) CreateShortLink(link *models.Link) (error) {
	shortLink, err := uc.linkRepository.SelectLinkByOriginalLink(link.OriginalLink)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return errors.Wrap(err, "link repository error")
	} else if err == nil {
		link.ShortLink = shortLink
		return nil
	}

	link.ShortLink, err = generateShortLink(link.OriginalLink)
	if err != nil {
		return errors.Wrap(err, "generation short link error")
	}

	err = uc.linkRepository.CreateLink(link)
	if err != nil {
		return errors.Wrap(err, "link repository error")
	}

	return nil
}

func (uc *useCase) GetOriginalLink(link string) (string, error) {
	gotLink, err := uc.linkRepository.SelectLinkByShortLink(link)
	if err != nil {
		return "", errors.Wrap(err, "link repository error")
	}

	return gotLink, nil
}

func generateShortLink(originalLink string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(originalLink))
	if err != nil {
		return "", err
	}
	num := new(big.Int).SetBytes(h.Sum(nil)).Uint64()

	res := encode(num)

	return res, nil
}

func encode(num uint64) string {
	res := make([]rune, 10)
	for idx := range res {
		if num > 0 {
			res[idx] = alphabet[num % uint64(len(alphabet))]
			num /= uint64(len(alphabet))
		} else {
			res[idx] = alphabet[rand.Intn(len(alphabet))]
		}
	}
	return string(res)
}

