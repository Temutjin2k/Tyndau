package middleware

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Temutjin2k/Tyndau/api-gateway/config"
	"github.com/Temutjin2k/Tyndau/api-gateway/internal/utils"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func ReverseProxyHandler(svc *config.Service, targetURL string) http.HandlerFunc {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	return func(w http.ResponseWriter, r *http.Request) {
		if svc.Status != "up" {
			utils.WriteJSON(w, http.StatusServiceUnavailable, ErrorResponse{Message: "Service " + svc.Name + " is currently down."})
			return
		}

		proxy.ServeHTTP(w, r)
	}
}
