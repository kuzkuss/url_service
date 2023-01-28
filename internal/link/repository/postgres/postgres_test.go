package postgres_test

import (
	"regexp"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kuzkuss/url_service/models"
	linkRep "github.com/kuzkuss/url_service/internal/link/repository/postgres"
)

type TestCaseCreate struct {
	ArgData *models.Link
	Error error
}

type TestCaseSelect struct {
	ArgData string
	ExpectedRes string
	Error error
}

func TestRepositoryCreateLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	gdb.Logger.LogMode(logger.Info)

	linkSuccess := models.Link {
		OriginalLink: "original_link_success",
		ShortLink: "short_link_success",
	}

	linkError := models.Link {
		OriginalLink: "original_link_error",
		ShortLink: "short_link_error",
	}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "links" ("original_link","short_link") VALUES ($1,$2)`)).WithArgs(
			linkSuccess.OriginalLink, linkSuccess.ShortLink).WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectCommit()

	createErr := errors.New("error")

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "links" ("original_link","short_link") VALUES ($1,$2)`)).WithArgs(
			linkError.OriginalLink, linkError.ShortLink).WillReturnError(createErr)

	mock.ExpectRollback()

	repository := linkRep.New(gdb)

	cases := map[string]TestCaseCreate {
		"success": {
			ArgData:   &linkSuccess,
			Error: nil,
		},
		"error": {
			ArgData:   &linkError,
			Error: createErr,
		},
	}

	t.Run("success", func(t *testing.T) {
		err := repository.CreateLink(cases["success"].ArgData)
		require.Equal(t, cases["success"].Error, errors.Cause(err))
	})

	t.Run("error", func(t *testing.T) {
		err := repository.CreateLink(cases["error"].ArgData)
		require.Equal(t, cases["error"].Error, errors.Cause(err))
	})

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRepositorySelectLinkByShortLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	gdb.Logger.LogMode(logger.Info)

	linkSuccess := models.Link {
		OriginalLink: "original_link_success",
		ShortLink: "short_link_success",
	}

	linkError := models.Link {
		OriginalLink: "",
		ShortLink: "short_link_error",
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "links" WHERE short_link = $1 LIMIT 1`)).WithArgs(linkSuccess.ShortLink).
		WillReturnRows(sqlmock.NewRows([]string{"short_link", "original_link"}).
		AddRow(linkSuccess.ShortLink, linkSuccess.OriginalLink))

	getErr := errors.New("error")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "links" WHERE short_link = $1 LIMIT 1`)).WithArgs(linkError.ShortLink).
		WillReturnError(getErr)

	repository := linkRep.New(gdb)

	cases := map[string]TestCaseSelect {
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
	}

	t.Run("success", func(t *testing.T) {
		actualRes, err := repository.SelectLinkByShortLink(cases["success"].ArgData)
		require.Equal(t, cases["success"].Error, errors.Cause(err))
		assert.Equal(t, cases["success"].ExpectedRes, actualRes)
	})

	t.Run("error", func(t *testing.T) {
		actualRes, err := repository.SelectLinkByShortLink(cases["error"].ArgData)
		require.Equal(t, cases["error"].Error, errors.Cause(err))
		assert.Equal(t, cases["error"].ExpectedRes, actualRes)
	})

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func TestRepositorySelectLinkByOriginalLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	gdb.Logger.LogMode(logger.Info)

	linkSuccess := models.Link {
		OriginalLink: "original_link_success",
		ShortLink: "short_link_success",
	}

	linkError := models.Link {
		OriginalLink: "original_link_error",
		ShortLink: "",
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "links" WHERE original_link = $1 LIMIT 1`)).WithArgs(linkSuccess.OriginalLink).
		WillReturnRows(sqlmock.NewRows([]string{"short_link", "original_link"}).
		AddRow(linkSuccess.ShortLink, linkSuccess.OriginalLink))

	getErr := errors.New("error")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "links" WHERE original_link = $1 LIMIT 1`)).WithArgs(linkError.OriginalLink).
		WillReturnError(getErr)

	repository := linkRep.New(gdb)

	cases := map[string]TestCaseSelect {
		"success": {
			ArgData:   linkSuccess.OriginalLink,
			ExpectedRes: linkSuccess.ShortLink,
			Error: nil,
		},
		"error": {
			ArgData:   linkError.OriginalLink,
			ExpectedRes: linkError.ShortLink,
			Error: getErr,
		},
	}

	t.Run("success", func(t *testing.T) {
		actualRes, err := repository.SelectLinkByOriginalLink(cases["success"].ArgData)
		require.Equal(t, cases["success"].Error, errors.Cause(err))
		assert.Equal(t, cases["success"].ExpectedRes, actualRes)
	})

	t.Run("error", func(t *testing.T) {
		actualRes, err := repository.SelectLinkByOriginalLink(cases["error"].ArgData)
		require.Equal(t, cases["error"].Error, errors.Cause(err))
		assert.Equal(t, cases["error"].ExpectedRes, actualRes)
	})

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
