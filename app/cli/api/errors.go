package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"plandex/auth"

	"github.com/plandex/plandex/shared"
)

func handleApiError(r *http.Response, errBody []byte) *shared.ApiError {
	// Check if the response is JSON
	if r.Header.Get("Content-Type") != "application/json" {
		return &shared.ApiError{
			Type:   shared.ApiErrorTypeOther,
			Status: r.StatusCode,
			Msg:    string(errBody),
		}
	}

	var apiError shared.ApiError
	if err := json.Unmarshal(errBody, &apiError); err != nil {
		log.Printf("Error unmarshalling JSON: %v\n", err)
		return &shared.ApiError{
			Type:   shared.ApiErrorTypeOther,
			Status: r.StatusCode,
			Msg:    string(errBody),
		}
	}

	return &apiError
}

func refreshTokenIfNeeded(apiErr *shared.ApiError) (bool, *shared.ApiError) {
	if apiErr.Type == shared.ApiErrorTypeInvalidToken {
		err := auth.RefreshInvalidToken()
		if err != nil {
			return false, &shared.ApiError{Type: shared.ApiErrorTypeOther, Msg: fmt.Sprintf("error refreshing invalid token")}
		}
		return true, nil
	}
	return false, apiErr
}