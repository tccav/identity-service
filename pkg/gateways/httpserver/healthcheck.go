package httpserver

import "net/http"

// Healthcheck ...
// ShowEntity godoc
// @Summary Check if service is healthy
// @Tags Internal
// @Success 200
// @Router /healthcheck [get]
func Healthcheck(_ http.ResponseWriter, _ *http.Request) {}
