package errorwrapper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d - %s", e.StatusCode, e.Msg)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

// Echo Error Wrapper
type EchoAPIFunc func(c echo.Context) error

func EchoErrorWrapper(h EchoAPIFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := h(c); err != nil {
			if apiErr, ok := err.(APIError); ok {
				return EchoWriteJSON(c, apiErr.StatusCode, apiErr)
			} else {
				errResp := map[string]any{
					"statusCode": http.StatusInternalServerError,
					"msg":        "internal server error",
				}
				return EchoWriteJSON(c, http.StatusInternalServerError, errResp)
			}
		}
		return nil
	}
}

func EchoWriteJSON(c echo.Context, statusCode int, v any) error {
	return c.JSON(statusCode, v)
}

// HTTP Error Wrapper
type HTTPAPIFunc func(w http.ResponseWriter, r *http.Request) error

func HTTPErrorWrapper(h HTTPAPIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				HTTPWriteJSON(w, apiErr.StatusCode, apiErr)
			} else {
				errResp := map[string]any{
					"statusCode": http.StatusInternalServerError,
					"msg":        "internal server error",
				}
				HTTPWriteJSON(w, http.StatusInternalServerError, errResp)
			}
		}
	}
}

func HTTPWriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
