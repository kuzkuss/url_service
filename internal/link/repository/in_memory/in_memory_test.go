package in_memory_test

import (
	"testing"

	"github.com/kuzkuss/url_service/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	linkRep "github.com/kuzkuss/url_service/internal/link/repository/in_memory"
)

type TestCaseSelect struct {
	ArgData string
	ExpectedRes string
	Error error
}

type TestCaseCreate struct {
	ArgData *models.Link
	Error error
}

func TestUsecaseCreateLink(t *testing.T) {
	linkSuccess := models.Link {
		OriginalLink: "original_link_success",
		ShortLink: "short_link_success",
	}

	repository := linkRep.New()

	cases := map[string]TestCaseCreate {
		"success": {
			ArgData:   &linkSuccess,
			Error: nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := repository.CreateLink(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
}

func TestUsecaseSelectLinkByShortLink(t *testing.T) {
	linkSuccess := models.Link {
		OriginalLink: "original_link_success",
		ShortLink: "short_link_success",
	}

	linkNotFound := models.Link {
		OriginalLink: "",
		ShortLink: "short_link_not_found",
	}

	repository := linkRep.New()

	cases := map[string]TestCaseSelect {
		"success": {
			ArgData:   linkSuccess.ShortLink,
			ExpectedRes: linkSuccess.OriginalLink,
			Error: nil,
		},
		"not_found": {
			ArgData:   linkNotFound.ShortLink,
			ExpectedRes: linkNotFound.OriginalLink,
			Error: models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			if name == "success" {
				repository.CreateLink(&linkSuccess)
			}
			actualRes, err := repository.SelectLinkByShortLink(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
			assert.Equal(t, test.ExpectedRes, actualRes)
		})
	}
}

func TestUsecaseSelectLinkByOriginalLink(t *testing.T) {
	linkSuccess := models.Link {
		OriginalLink: "original_link_success",
		ShortLink: "short_link_success",
	}

	linkNotFound := models.Link {
		OriginalLink: "original_link_not_found",
		ShortLink: "",
	}

	repository := linkRep.New()

	cases := map[string]TestCaseSelect {
		"success": {
			ArgData:   linkSuccess.OriginalLink,
			ExpectedRes: linkSuccess.ShortLink,
			Error: nil,
		},
		"not_found": {
			ArgData:   linkNotFound.OriginalLink,
			ExpectedRes: linkNotFound.ShortLink,
			Error: models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			if name == "success" {
				repository.CreateLink(&linkSuccess)
			}
			actualRes, err := repository.SelectLinkByOriginalLink(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
			assert.Equal(t, test.ExpectedRes, actualRes)
		})
	}
}


