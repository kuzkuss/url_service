package delivery_test

import (
	"context"
	"testing"

	linkDelivery "github.com/kuzkuss/url_service/internal/link/delivery/grpc"
	linkMocks "github.com/kuzkuss/url_service/internal/link/usecase/mocks"
	"github.com/kuzkuss/url_service/models"
	link "github.com/kuzkuss/url_service/proto/link"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCaseGet struct {
	ArgData *link.ShortLink
	ExpectedRes *link.OriginalLink
	Error error
}

type TestCaseCreate struct {
	ArgData *link.OriginalLink
	Error error
}

func TestGrpcDeliveryCreateShortLink(t *testing.T) {
	linkSuccess := models.Link {
		OriginalLink: "original_link_success",
	}

	linkError := models.Link {
		OriginalLink: "original_link_error",
	}

	mockPbOriginalLinkSuccess := link.OriginalLink {
		OriginalLink: linkSuccess.OriginalLink,
	}
	mockPbOriginalLinkError := link.OriginalLink {
		OriginalLink: linkError.OriginalLink,
	}

	createErr := errors.New("error")
	ctx := context.Background()

	mockLinkUsecase := linkMocks.NewUseCaseI(t)

	mockLinkUsecase.On("CreateShortLink", &linkSuccess).Return(nil)
	mockLinkUsecase.On("CreateShortLink", &linkError).Return(createErr)

	delivery := linkDelivery.New(mockLinkUsecase)

	cases := map[string]TestCaseCreate {
		"success": {
			ArgData:   &mockPbOriginalLinkSuccess,
			Error: nil,
		},
		"error": {
			ArgData:   &mockPbOriginalLinkError,
			Error: createErr,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := delivery.CreateShortLink(ctx, test.ArgData)
			require.Equal(t, test.Error, err)
		})
	}
	mockLinkUsecase.AssertExpectations(t)
}

func TestGrpcDeliveryGetOriginalLink(t *testing.T) {
	mockPbShortLinkSuccess := link.ShortLink {
		ShortLink: "short_link_success",
	}
	mockPbShortLinkError := link.ShortLink {
		ShortLink: "short_link_error",
	}

	mockPbOriginalLinkSuccess := link.OriginalLink {
		OriginalLink: "original_link_success",
	}

	mockPbOriginalLinkError := link.OriginalLink {
		OriginalLink: "",
	}

	getErr := errors.New("error")

	ctx := context.Background()

	mockLinkUsecase := linkMocks.NewUseCaseI(t)

	mockLinkUsecase.On("GetOriginalLink", mockPbShortLinkSuccess.ShortLink).
										Return(mockPbOriginalLinkSuccess.OriginalLink, nil)
	mockLinkUsecase.On("GetOriginalLink", mockPbShortLinkError.ShortLink).
										Return(mockPbOriginalLinkError.OriginalLink, getErr)

	delivery := linkDelivery.New(mockLinkUsecase)

	cases := map[string]TestCaseGet {
		"success": {
			ArgData:   &mockPbShortLinkSuccess,
			ExpectedRes: &mockPbOriginalLinkSuccess,
			Error: nil,
		},
		"error": {
			ArgData:   &mockPbShortLinkError,
			ExpectedRes: &mockPbOriginalLinkError,
			Error: getErr,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			actualRes, err := delivery.GetOriginalLink(ctx, test.ArgData)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.ExpectedRes, actualRes)
			}
		})
	}
	mockLinkUsecase.AssertExpectations(t)
}

