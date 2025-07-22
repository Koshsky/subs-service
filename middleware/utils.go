package middleware

import (
	"net/http"

	"github.com/fatih/color"
)

func colorForStatus(code int) string {
	switch {
	case isSuccess(code):
		return color.GreenString("%d", code)
	case isRedirect(code):
		return color.WhiteString("%d", code)
	case isClientError(code):
		return color.YellowString("%d", code)
	default: // server errors
		return color.RedString("%d", code)
	}
}

func isSuccess(code int) bool {
	return code >= http.StatusOK && code < http.StatusMultipleChoices
}

func isRedirect(code int) bool {
	return code >= http.StatusMultipleChoices && code < http.StatusBadRequest
}

func isClientError(code int) bool {
	return code >= http.StatusBadRequest && code < http.StatusInternalServerError
}

func isServerError(code int) bool {
	return code >= http.StatusInternalServerError
}
