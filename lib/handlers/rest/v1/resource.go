package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/holmes89/magpie/lib"

	"github.com/segmentio/ksuid"
)

type SiteFactory interface {
	Create(context.Context, lib.Site) error
}

type SiteRepository interface {
	Get(ctx context.Context, id string) (lib.Site, error)
	GetAll(ctx context.Context) ([]lib.Site, error)
}

type SiteService interface {
	SiteFactory
	SiteRepository
}

type siteHandler struct {
	service SiteService
}

func MakeV1SiteHandler(mr *mux.Router, service SiteService) {

	r := mr.PathPrefix("/r").Subrouter()

	h := &siteHandler{
		service: service,
	}

	r.HandleFunc("/{id}", h.Find).Methods("GET", "OPTIONS")
	r.HandleFunc("", h.FindAll).Methods("GET", "OPTIONS")
	r.HandleFunc("", h.Create).Methods("POST")
}

func (h *siteHandler) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := h.service.Get(ctx, id)
	if err != nil {
		makeError(w, http.StatusBadRequest, "unable to find results", "get")
		return
	}

	encodeResponse(r.Context(), w, res)
}

func (h *siteHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.service.GetAll(ctx)
	if err != nil {
		makeError(w, http.StatusBadRequest, "unable to find results", "get")
		return
	}

	encodeResponse(r.Context(), w, res)
}

func (h *siteHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	req := struct {
		URL string `json:"url"`
	}{}
	if err := json.Unmarshal(b, &req); err != nil {
		fmt.Println("unable to unmarshall resource")
		makeError(w, http.StatusBadRequest, "unable to unmarshall resource", "post")
	}

	res := lib.Site{
		ID:        ksuid.New().String(),
		Name:      req.URL,
		URL:       req.URL,
		CreatedAt: time.Now().UTC(),
	}

	if err := h.service.Create(ctx, res); err != nil {
		fmt.Println("unable to create resource")
		makeError(w, http.StatusBadRequest, "unable to create resource", "post")
	}
	w.WriteHeader(http.StatusCreated)
	encodeResponse(ctx, w, res)
}

func makeError(w http.ResponseWriter, code int, message string, method string) {
	log.Printf("HTTP Error: %d - %s - %s", code, method, message)
	http.Error(w, message, code)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response) //TODO check error and handle?
}
