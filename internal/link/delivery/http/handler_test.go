package delivery_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kuzkuss/url_service/models"
	"github.com/kuzkuss/url_service/pkg"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	linkDelivery "github.com/kuzkuss/url_service/internal/link/delivery/http"
	linkMocks "github.com/kuzkuss/url_service/internal/link/usecase/mocks"
)

type TestCaseCreate struct {
	ArgData string
	ExpectedResponse string
	Error error
	StatusCode int
}

type TestCaseGet struct {
	ArgData string
	ExpectedResponse string
	Error error
	StatusCode int
}

func TestHttpDeliveryCreateShortLink(t *testing.T) {
	linkSuccess := models.Link {
		OriginalLink: "original_link_success",
		ShortLink: "short_link_success",
	}

	linkInternalError := models.Link {
		OriginalLink: "original_link_internal_error",
	}

	linkInvalid := models.Link{}

	jsonLinkSuccess, err := json.Marshal(linkSuccess)
	assert.NoError(t, err)

	jsonLinkInternalErr, err := json.Marshal(linkInternalError)
	assert.NoError(t, err)

	jsonLinkInvalid, err := json.Marshal(linkInvalid)
	assert.NoError(t, err)

	createErr := errors.New("error")

	mockLinkUsecase := linkMocks.NewUseCaseI(t)

	mockLinkUsecase.On("CreateShortLink", &linkSuccess).Return(nil)
	mockLinkUsecase.On("CreateShortLink", &linkInternalError).Return(createErr)

	response := pkg.Response {
		Body: models.Link { ShortLink: linkSuccess.ShortLink},
	}
	jsonResponse, err := json.Marshal(response)
	assert.NoError(t, err)

	e := echo.New()
	linkDelivery.New(e, mockLinkUsecase)

	delivery := linkDelivery.Delivery {
		LinkUC: mockLinkUsecase,
	}

	cases := map[string]TestCaseCreate {
		"success": {
			ArgData:   string(jsonLinkSuccess),
			ExpectedResponse: string(jsonResponse) + "\n",
			Error: nil,
			StatusCode: http.StatusCreated,
		},
		"bad_request": {
			ArgData:   "aaa",
			Error: &echo.HTTPError{
				Code: http.StatusBadRequest,
				Message: models.ErrBadRequest.Error(),
			},
		},
		"invalid_request": {
			ArgData:   string(jsonLinkInvalid),
			Error: &echo.HTTPError{
				Code: http.StatusBadRequest,
				Message: models.ErrBadRequest.Error(),
			},
		},
		"internal_error": {
			ArgData:   string(jsonLinkInternalErr),
			Error: &echo.HTTPError{
				Code: http.StatusInternalServerError,
				Message: models.ErrInternalServerError.Error(),
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(echo.POST, "/create", strings.NewReader(test.ArgData))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/create")

			err = delivery.CreateShortLink(c)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.StatusCode, rec.Code)
				assert.Equal(t, test.ExpectedResponse, rec.Body.String())
			}
		})
	}

	mockLinkUsecase.AssertExpectations(t)
}

func TestHttpDeliveryGetOriginalLink(t *testing.T) {
	linkSuccess := models.Link {
		OriginalLink: "original_link_success",
		ShortLink: "short_link_success",
	}

	linkInternalError := models.Link {
		OriginalLink: "",
		ShortLink: "short_link_internal_error",
	}

	linkNotFound := models.Link {
		OriginalLink: "",
		ShortLink: "short_link_not_found",
	}

	mockLinkUsecase := linkMocks.NewUseCaseI(t)

	mockLinkUsecase.On("GetOriginalLink", linkSuccess.ShortLink).
										Return(linkSuccess.OriginalLink, nil)
	mockLinkUsecase.On("GetOriginalLink", linkInternalError.ShortLink).
										Return(linkInternalError.OriginalLink, models.ErrInternalServerError)
	mockLinkUsecase.On("GetOriginalLink", linkNotFound.ShortLink).
										Return(linkNotFound.OriginalLink, models.ErrNotFound)

	response := pkg.Response {
		Body: models.Link {OriginalLink: linkSuccess.OriginalLink},
	}
	jsonResponse, err := json.Marshal(response)
	assert.NoError(t, err)

	e := echo.New()
	linkDelivery.New(e, mockLinkUsecase)

	delivery := linkDelivery.Delivery {
		LinkUC: mockLinkUsecase,
	}

	cases := map[string]TestCaseGet {
		"success": {
			ArgData:   linkSuccess.ShortLink,
			ExpectedResponse: string(jsonResponse) + "\n",
			Error: nil,
			StatusCode: http.StatusOK,
		},
		"not_found": {
			ArgData:   linkNotFound.ShortLink,
			Error: &echo.HTTPError{
				Code: http.StatusNotFound,
				Message: models.ErrNotFound.Error(),
			},
		},
		"internal_error": {
			ArgData:   linkInternalError.ShortLink,
			Error: &echo.HTTPError{
				Code: http.StatusInternalServerError,
				Message: models.ErrInternalServerError.Error(),
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(echo.GET, "/get/:short_link", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/get/:short_link")
			c.SetParamNames("short_link")
			c.SetParamValues(test.ArgData)

			err = delivery.GetOriginalLink(c)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.StatusCode, rec.Code)
				assert.Equal(t, test.ExpectedResponse, rec.Body.String())
			}
		})
	}

	mockLinkUsecase.AssertExpectations(t)
}

