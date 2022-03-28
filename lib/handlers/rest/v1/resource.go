package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/holmes89/magpie/lib"

	"github.com/segmentio/ksuid"
)

type ResourceService interface {
	Get(context.Context, string, string) (lib.Resource, error)
	GetAll(context.Context, string) ([]lib.Resource, error)
	Create(context.Context, lib.Resource) error
}

type resourceHandler struct {
	service ResourceService
}

func MakeV1ResourceHandler(mr *mux.Router, service ResourceService) {

	r := mr.PathPrefix("/r").Subrouter()

	h := &resourceHandler{
		service: service,
	}

	r.HandleFunc("/{resource}/{id}", h.Find).Methods("GET", "OPTIONS")
	r.HandleFunc("/{resource}", h.FindAll).Methods("GET", "OPTIONS")
	r.HandleFunc("/{resource}", h.Create).Methods("POST")
}

func (h *resourceHandler) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	resource := vars["resource"]
	id := vars["id"]

	res, err := h.service.Get(ctx, resource, id)
	if err != nil {
		makeError(w, http.StatusBadRequest, "unable to find results", "get")
		return
	}

	encodeResponse(r.Context(), w, res)
}

func (h *resourceHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	resource := vars["resource"]

	res, err := h.service.GetAll(ctx, resource)
	if err != nil {
		makeError(w, http.StatusBadRequest, "unable to find results", "get")
		return
	}

	encodeResponse(r.Context(), w, res)
}

func (h *resourceHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	resource := vars["resource"]

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	req := struct {
		Name string                 `json:"name"`
		Meta map[string]interface{} `json:"meta_data"`
	}{}
	if err := json.Unmarshal(b, &req); err != nil {
		fmt.Println("unable to unmarshall resource")
		makeError(w, http.StatusBadRequest, "unable to unmarshall resource", "post")
	}

	res := lib.Resource{
		ResourceID: ksuid.New().String(),
		Name:       req.Name,
		Meta:       req.Meta,
		Type:       strings.ToLower(resource),
		CreatedAt:  time.Now().UTC(),
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
