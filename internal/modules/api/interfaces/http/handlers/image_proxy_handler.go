package handlers

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// ImageProxyHandler proxies external image requests.
type ImageProxyHandler struct {
	client *http.Client
}

// NewImageProxyHandler creates a new image proxy handler.
func NewImageProxyHandler() *ImageProxyHandler {
	return &ImageProxyHandler{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

// ProxyImage fetches an external image and returns it.
func (h *ImageProxyHandler) ProxyImage(w http.ResponseWriter, r *http.Request) {
	rawURL := r.URL.Query().Get("url")
	if rawURL == "" {
		http.Error(w, "missing url parameter", http.StatusBadRequest)
		return
	}

	// Only allow fhr.ru subdomains
	if !strings.Contains(rawURL, ".fhr.ru/") {
		http.Error(w, "host not allowed", http.StatusForbidden)
		return
	}

	resp, err := h.client.Get(rawURL)
	if err != nil {
		http.Error(w, "failed to fetch image", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "upstream error", resp.StatusCode)
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Cache-Control", "public, max-age=86400")
	io.Copy(w, resp.Body)
}
