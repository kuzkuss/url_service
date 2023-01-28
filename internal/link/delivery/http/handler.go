package delivery

import (
	"net/http"

	"github.com/pkg/errors"

	linkUsecase "github.com/kuzkuss/url_service/internal/link/usecase"
	"github.com/kuzkuss/url_service/models"
	"github.com/kuzkuss/url_service/pkg"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type Delivery struct {
	LinkUC linkUsecase.UseCaseI
}

// CreateShortLink godoc
// @Summary      CreateShortLink
// @Description  create short link
// @Tags     link
// @Accept	 application/json
// @Produce  application/json
// @Param    original_link body models.Link true "link data"
// @Success 201 {object} pkg.Response{body=models.Link} "short link created"
// @Failure 405 {object} echo.HTTPError "method not allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Router   /create [post]
func (del *Delivery) CreateShortLink(c echo.Context) error {
	var link models.Link
	err := c.Bind(&link)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	if ok, err := isRequestValid(&link); !ok {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	err = del.LinkUC.CreateShortLink(&link)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	link.OriginalLink = ""
	return c.JSON(http.StatusCreated, pkg.Response{Body: link})
}

// GetOriginalLink godoc
// @Summary      GetOriginalLink
// @Description  get original link by short link
// @Tags     link
// @Param short_link path string  true  "Short link"
// @Produce  application/json
// @Success  200 {object} pkg.Response{body=models.Link} "success get link"
// @Failure 405 {object} echo.HTTPError "method not allowed"
// @Failure 404 {object} echo.HTTPError "not found"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Router   /get/{short_link} [get]
func (del *Delivery) GetOriginalLink(c echo.Context) error {
	link, err := del.LinkUC.GetOriginalLink(c.Param("short_link"))
	if err != nil {
		causeErr := errors.Cause(err)
		switch {
		case errors.Is(causeErr, models.ErrNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
		}
	}

	return c.JSON(http.StatusOK, pkg.Response{Body: models.Link{OriginalLink: link}})
}

func isRequestValid(link interface{}) (bool, error) {
	validate := validator.New()
	err := validate.Struct(link)
	if err != nil {
		return false, err
	}
	return true, nil
}

func New(e *echo.Echo, linkUC linkUsecase.UseCaseI) {
	handler := &Delivery{
		LinkUC: linkUC,
	}

	e.POST("/create", handler.CreateShortLink)
	e.GET("/get/:short_link", handler.GetOriginalLink)
}
