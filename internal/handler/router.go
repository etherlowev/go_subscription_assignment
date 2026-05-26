package handler

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	_ "onlineSubscriptions/docs"
	"onlineSubscriptions/internal/models"
	"onlineSubscriptions/internal/service"
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

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	subrouter := router.PathPrefix("/api/subscriptions").Subrouter()

	r := &Router{
		router:  router,
		service: service,
	}

	subrouter.HandleFunc("/list", r.ListSubs).Methods(http.MethodGet)
	subrouter.HandleFunc("/{id}", r.FindSub).Methods(http.MethodGet)

	subrouter.HandleFunc("/{id}", r.UpdateSubscription).Methods(http.MethodPut)
	subrouter.HandleFunc("/", r.CreateSubscription).Methods(http.MethodPost)
	subrouter.HandleFunc("/{id}", r.DeleteSubscription).Methods(http.MethodDelete)

	return r, nil
}

// FindSub godoc
// @Summary Get subscription by id
// @Description Get a single subscription by ID
// @Tags subscription
// @Produce json
// @Param id path uuid true "subscription ID"
// @Success 200 {object} models.Subscription
// @Router /{id} [get]
func (r *Router) FindSub(w http.ResponseWriter, req *http.Request) {
	pathVariables := mux.Vars(req)
	id := pathVariables["id"]
	subId, err := uuid.Parse(id)
	if handleBadRequest(err, w) {
		return
	}

	subscription, err := r.service.GetSubscriptionById(context.Background(), subId)
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
}

// ListSubs godoc
// @Summary List subscription page
// @Description Get a page of subscriptions
// @Tags subscription
// @Produce json
// @Param page path int true "page number"
// @Param perPage path int true "amount of entries in page"
// @Success 200 {array} models.Subscription
// @Router /{id} [get]
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

	list, err := r.service.ListSubscriptions(context.Background(), page, perPage)
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
}

func (r *Router) CreateSubscription(w http.ResponseWriter, req *http.Request) {
	var subscriptionRequest models.SubscriptionRequest
	err := json.NewDecoder(req.Body).Decode(&subscriptionRequest)
	if handleErr(err, w) {
		return
	}

	saved, err := r.service.InsertSubscription(context.Background(), &subscriptionRequest)
	if handleErr(err, w) {
		return
	}

	if saved {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotModified)
	}
}

func (r *Router) UpdateSubscription(w http.ResponseWriter, req *http.Request) {
	pathVariables := mux.Vars(req)
	id := pathVariables["id"]

	subId, err := uuid.Parse(id)

	handleBadRequest(err, w)

	var request models.SubscriptionRequest
	err = json.NewDecoder(req.Body).Decode(&request)
	if handleErr(err, w) {
		return
	}

	saved, err := r.service.UpdateSubscription(context.Background(), subId, &request)
	if handleErr(err, w) {
		return
	}

	if saved {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotModified)
	}
}

func (r *Router) DeleteSubscription(w http.ResponseWriter, req *http.Request) {
	var subscription models.Subscription
	if err := json.NewDecoder(req.Body).Decode(&subscription); err != nil {
		handleErr(err, w)
	}

	removed, err := r.service.RemoveSubscription(context.Background(), subscription.Id)
	if handleErr(err, w) {
		return
	}

	if removed {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotModified)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
