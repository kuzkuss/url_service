package usecase_test

import (
	"testing"

	linkUsecase "github.com/kuzkuss/url_service/internal/link/usecase"
	linkMocks "github.com/kuzkuss/url_service/internal/link/repository/mocks"
	"github.com/kuzkuss/url_service/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCaseGet struct {
	ArgData string
	ExpectedRes string
	Error error
}

type TestCaseCreate struct {
	ArgData *models.Link
	Error error
}

func TestUsecaseCreateShortLink(t *testing.T) {
	linkSuccess := models.Link {
		OriginalLink: "original_link_success",
		ShortLink: "short_link_success",
	}

	linkError := models.Link {
		OriginalLink: "original_link_error",
		ShortLink: "",
	}

	linkConflict := models.Link {
		OriginalLink: "original_link_conflict",
		ShortLink: "short_link_conflict",
	}

	createErr := errors.New("error")

	mockLinkRepo := linkMocks.NewRepositoryI(t)

	mockLinkRepo.On("SelectLinkByOriginalLink", linkSuccess.OriginalLink).Return("", models.ErrNotFound)
	mockLinkRepo.On("CreateLink", &linkSuccess).Return(nil)
	mockLinkRepo.On("SelectLinkByOriginalLink", linkConflict.OriginalLink).Return(linkConflict.ShortLink, nil)
	mockLinkRepo.On("SelectLinkByOriginalLink", linkError.OriginalLink).Return("", models.ErrNotFound)
	mockLinkRepo.On("CreateLink", &linkError).Return(createErr)

	usecase := linkUsecase.New(mockLinkRepo)

	cases := map[string]TestCaseCreate {
		"success": {
			ArgData:   &linkSuccess,
			Error: nil,
		},
		"conflict": {
			ArgData:   &linkConflict,
			Error: nil,
		},
		"error": {
			ArgData:   &linkError,
			Error: createErr,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := usecase.CreateShortLink(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockLinkRepo.AssertExpectations(t)
}

func TestUsecaseGetOriginalLink(t *testing.T) {
	linkSuccess := models.Link {
		OriginalLink: "original_link_success",
		ShortLink: "short_link_success",
	}

	linkError := models.Link {
		OriginalLink: "",
		ShortLink: "short_link_error",
	}

	linkNotFound := models.Link {
		OriginalLink: "",
		ShortLink: "short_link_not_found",
	}

	getErr := errors.New("error")

	mockLinkRepo := linkMocks.NewRepositoryI(t)

	mockLinkRepo.On("SelectLinkByShortLink", linkSuccess.ShortLink).Return(linkSuccess.OriginalLink, nil)
	mockLinkRepo.On("SelectLinkByShortLink", linkError.ShortLink).Return(linkError.OriginalLink, getErr)
	mockLinkRepo.On("SelectLinkByShortLink", linkNotFound.ShortLink).Return(linkNotFound.OriginalLink, models.ErrNotFound)

	usecase := linkUsecase.New(mockLinkRepo)

	cases := map[string]TestCaseGet {
		"success": {
			ArgData:   linkSuccess.ShortLink,
			ExpectedRes: linkSuccess.OriginalLink,
			Error: nil,
		},
		"error": {
			ArgData:   linkError.ShortLink,
			ExpectedRes: linkError.OriginalLink,
			Error: getErr,
		},
		"not_found": {
			ArgData:   linkNotFound.ShortLink,
			ExpectedRes: linkNotFound.OriginalLink,
			Error: models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			actualRes, err := usecase.GetOriginalLink(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))

			if err == nil {
				assert.Equal(t, test.ExpectedRes, actualRes)
			}
		})
	}
	mockLinkRepo.AssertExpectations(t)
}

