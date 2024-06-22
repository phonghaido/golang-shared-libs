package errorwrapper

import (
	"fmt"
	"net/http"
)

func InvalidJSON() error {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid json format"))
}
