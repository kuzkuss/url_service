package delivery

import (
	"context"

	linkUsecase "github.com/kuzkuss/url_service/internal/link/usecase"
	"github.com/kuzkuss/url_service/models"
	link "github.com/kuzkuss/url_service/proto/link"
)

type LinkManager struct {
	link.UnimplementedLinksServer
	LinkUC linkUsecase.UseCaseI
}

func New(uc linkUsecase.UseCaseI) link.LinksServer {
	return LinkManager{LinkUC: uc}
}

func (lm LinkManager) CreateShortLink(ctx context.Context, originalLink *link.OriginalLink) (*link.ShortLink, error) {
	modelLink := models.Link {
		OriginalLink: originalLink.OriginalLink,
	}
	err := lm.LinkUC.CreateShortLink(&modelLink)

	resp := &link.ShortLink {
		ShortLink: modelLink.ShortLink,
	}

	return resp, err
}

func (lm LinkManager) GetOriginalLink(ctx context.Context, shortLink *link.ShortLink) (*link.OriginalLink, error) {
	originalLink, err := lm.LinkUC.GetOriginalLink(shortLink.ShortLink)

	resp := &link.OriginalLink {
		OriginalLink: originalLink,
	}

	return resp, err
}
