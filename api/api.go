package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type API struct {
	Logger *log.Logger
	router mux.Router
}

var information map[string]string = map[string]string{
	"router":     "A software router is just a router that is implemented completely in software: the CPU does all of the work",
	"routing":    "works: if an incoming request URL matches one of the paths, the corresponding handler is called passing",
	"subrouting": "Subrouters can be used to create domain or path \"namespaces\": you define subrouters in a central place and then parts of the app can register its paths relatively to a given subrouter.",
}

func (api *API) Initialize(host string) *mux.Router {
	api.router = *mux.NewRouter()

	api_subrouter := api.router.
		Host(host).
		Subrouter()

	api_subrouter.
		HandleFunc("/health", healthHandler).
		Methods(http.MethodGet).
		Name("health")

	info_subrouter := api_subrouter.
		PathPrefix("/information").
		Methods(http.MethodGet).
		Subrouter()

	info_subrouter.
		HandleFunc("", api.informationGeneralHandler)

	info_subrouter.
		Headers("request-type", "body").
		HandlerFunc(api.informationBodyHandler)

	info_subrouter.
		HandleFunc("/{category:[a-z A-Z]+}", api.informationQueryHandler).
		Name("category")

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

func (api *API) informationBodyHandler(w http.ResponseWriter, r *http.Request) {
	var req InformationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info, ok := information[req.Category]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(InformationResponse{Category: req.Category, Info: info})
	}

}

func (api *API) informationQueryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	category, ok := vars["category"]

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info, ok := information[category]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(InformationResponse{Category: category, Info: info})
	}
}

func (api *API) informationGeneralHandler(w http.ResponseWriter, r *http.Request) {
	var res []InformationCategoryResponse = make([]InformationCategoryResponse, 0)

	for key := range information {
		if url, err := api.router.Get("category").URL("category", key); err != nil {
			continue
		} else {
			res = append(res, InformationCategoryResponse{Name: strings.Title(key), URL: url.String()})
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
