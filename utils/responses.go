package utils

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type PaginatedResponse struct {
	Total int64                    `json:"total"`
	Page  int64                      `json:"page"`
	Pages int                      `json:"pages"`
	Limit int64                      `json:"limit"`
	//Data  []map[string]interface{} `json:"data"`
	Data  []interface{} `json:"data"`
}

func NewPaginatedResponse(page int64, limit int64) PaginatedResponse {
	return PaginatedResponse{
		Total: 0,
		Limit: limit,
		Page:  page,
		Pages: 1,
		Data:  []interface{}{},
	}
}

func SuccessResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

func CreatedResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": data,
	})
}

func AcceptedResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"data": data,
	})
}


func ValidationResponse(c echo.Context, error string) error {
	return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
		"error": error,
	})
}

func ErrorResponse(c echo.Context, error string) error {
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"error": error,
	})
}

func CustomErrorResponse(c echo.Context, error string, code int) error {
	return c.JSON(code, map[string]interface{}{
		"error": error,
	})
}