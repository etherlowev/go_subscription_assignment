package api

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"onlineSubscriptions/src/main/impl/service"
	"strconv"
)

type Router struct {
	service service.SubscriptionService
	router  *mux.Router
}

func handleBadRequest(e error, w http.ResponseWriter) bool {
	if e != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return true
	}
	return false
}

func handleErr(e error, w http.ResponseWriter) bool {
	if e != nil {
		log.Printf("err: %v", e)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return true
	}
	return false
}

func NewRouter(service service.SubscriptionService) (*Router, error) {
	router := mux.NewRouter()
	router.StrictSlash(true)
	subrouter := router.PathPrefix("/api/subscriptions").Subrouter()

	r := &Router{
		router:  router,
		service: service,
	}

	subrouter.HandleFunc("/list", r.ListSubs).Methods(http.MethodGet)
	subrouter.HandleFunc("/{id}", r.FindSub).Methods(http.MethodGet)

	return r, nil
}

// FindSub
// @Summary Get subscription by id
// @Description Get a single subscription by ID
// @Tags subscription
// @Produce json
// @Param id path uuid true "subscription ID"
// @Success 200 {object} model.Subscription
// @Router /{id} [get]
func (r *Router) FindSub(w http.ResponseWriter, req *http.Request) {
	pathVariables := mux.Vars(req)
	id := pathVariables["id"]

	subId, err := uuid.Parse(id)
	if handleBadRequest(err, w) {
		return
	}

	subscription, err := r.service.Find(context.Background(), subId)
	if handleErr(err, w) {
		return
	}

	jsonResponse, err := json.Marshal(subscription)
	if handleErr(err, w) {
		return
	}

	_, err = w.Write(jsonResponse)
	if handleErr(err, w) {
		return
	}
	w.WriteHeader(200)
}

func (r *Router) ListSubs(w http.ResponseWriter, req *http.Request) {
	urlQuery := req.URL.Query()
	page, err := strconv.Atoi(urlQuery.Get("page"))
	if handleBadRequest(err, w) {
		return
	}

	perPage, err := strconv.Atoi(urlQuery.Get("perPage"))
	if handleBadRequest(err, w) {
		return
	}

	list, err := r.service.List(context.Background(), page, perPage)
	if handleErr(err, w) {
		return
	}

	jsonResponse, err := json.Marshal(list)
	if handleErr(err, w) {
		return
	}

	_, err = w.Write(jsonResponse)
	if handleErr(err, w) {
		return
	}
	w.WriteHeader(200)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
