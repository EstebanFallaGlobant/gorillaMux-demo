package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	Logger *log.Logger
	router mux.Router
}

func (api *API) Initialize() *mux.Router {
	api.router = *mux.NewRouter()

	api.router.
		HandleFunc("/health", healthHandler).
		Name("health")

	return &api.router
}

func (api *API) GetHealthURL() (string, error) {
	url, err := api.router.Get("health").URL()

	if err != nil {
		return "", err
	} else {
		return url.String(), nil
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("ok")
}
