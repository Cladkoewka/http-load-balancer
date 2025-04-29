package proxy

import (
	"github.com/Cladkoewka/http-load-balancer/internal/balancer"
	"github.com/Cladkoewka/http-load-balancer/internal/logger"
	"net/http"
	"net/http/httputil"
)

func ProxyHandler(lb *balancer.LoadBalancer) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		backendURL, err := lb.GetNextBackendURL()
		if err != nil {
			http.Error(rw, "Error getting backend URL", http.StatusInternalServerError)
			return
		}

		logger.Log.Info("forwarding request", "url", r.URL.String(), "backend", backendURL.String())

		proxy := httputil.NewSingleHostReverseProxy(backendURL)
		proxy.ServeHTTP(rw, r)
	}
}
