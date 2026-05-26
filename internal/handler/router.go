package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	_ "onlineSubscriptions/docs"
	"onlineSubscriptions/internal/errors"
	"onlineSubscriptions/internal/middleware"
	"onlineSubscriptions/internal/models"
	"onlineSubscriptions/internal/service"
	"regexp"
	"strconv"
)

type Router struct {
	service service.SubscriptionService
	router  *mux.Router
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Printf("%s: %v", msg, err)
	}
	http.Error(w, msg, code)
}

func NewRouter(service service.SubscriptionService) (*Router, error) {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.JsonMiddleware)

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	subrouter := router.PathPrefix("/api/subscriptions").Subrouter()

	r := &Router{
		router:  router,
		service: service,
	}

	subrouter.HandleFunc("/", r.ListSubs).Methods(http.MethodGet)
	subrouter.HandleFunc("/{id}", r.FindSub).Methods(http.MethodGet)
	subrouter.HandleFunc("/sum", r.CalculateSum).Methods(http.MethodGet)

	subrouter.HandleFunc("/create", r.CreateSubscription).Methods(http.MethodPost)
	subrouter.HandleFunc("/{id}", r.UpdateSubscription).Methods(http.MethodPut)
	subrouter.HandleFunc("/{id}", r.DeleteSubscription).Methods(http.MethodDelete)

	return r, nil
}

// FindSub godoc
// @Summary Get subscription by id
// @Description Get a single subscription by ID
// @Tags subscription
// @Produce json
// @Param id path string true "subscription ID"
// @Success 200 {object} models.Subscription
// @Router /api/subscriptions/{id} [get]
func (r *Router) FindSub(w http.ResponseWriter, req *http.Request) {
	pathVariables := mux.Vars(req)
	id := pathVariables["id"]
	subId, err := uuid.Parse(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Id format is incorrect (not uuid)", err)
		return
	}

	subscription, err := r.service.GetSubscriptionById(req.Context(), subId)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to get subscription by id", err)
		return
	}

	writeJSON(w, http.StatusOK, subscription)
}

// ListSubs godoc
// @Summary List subscription page
// @Description Get a page of subscriptions
// @Tags subscription
// @Produce json
// @Param page query int true "page number"
// @Param perPage query int true "amount of entries in page"
// @Success 200 {array} models.Subscription
// @Router /api/subscriptions [get]
func (r *Router) ListSubs(w http.ResponseWriter, req *http.Request) {
	urlQuery := req.URL.Query()
	page, err := strconv.Atoi(urlQuery.Get("page"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Page format is incorrect (not integer)", err)
		return
	}

	perPage, err := strconv.Atoi(urlQuery.Get("perPage"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "perPage format is incorrect (not integer)", err)
		return
	}

	list, err := r.service.ListSubscriptions(req.Context(), page, perPage)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to get a list of subscriptions", err)
		return
	}

	if *list == nil {
		writeJSON(w, http.StatusOK, make([]int, 0))
	} else {
		writeJSON(w, http.StatusOK, list)
	}
}

// CreateSubscription godoc
// @Summary Create a subscription
// @Description Creates a subscription
// @Tags subscription
// @Produce json
// @Param request body models.SubscriptionRequest true "subscription request"
// @Success 201
// @Router /api/subscriptions/create [post]
func (r *Router) CreateSubscription(w http.ResponseWriter, req *http.Request) {
	var subscriptionRequest models.SubscriptionRequest
	err := json.NewDecoder(req.Body).Decode(&subscriptionRequest)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to decode request body", err)
		return
	}

	saved, err := r.service.InsertSubscription(req.Context(), &subscriptionRequest)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to insert a subscription", err)
		return
	}

	if saved {
		writeJSON(w, http.StatusCreated, nil)
	} else {
		writeJSON(w, http.StatusNotModified, nil)
	}
}

// UpdateSubscription godoc
// @Summary Update a subscription
// @Description Updates a subscription
// @Tags subscription
// @Produce json
// @Param id path string true "subscription id"
// @Param request body models.SubscriptionRequest true "subscription request"
// @Success 204
// @Router /api/subscriptions/{id} [put]
func (r *Router) UpdateSubscription(w http.ResponseWriter, req *http.Request) {
	pathVariables := mux.Vars(req)
	id := pathVariables["id"]

	subId, err := uuid.Parse(id)

	if err != nil {
		writeError(w, http.StatusBadRequest, "Subscription id format is incorrect (not uuid)", err)
		return
	}

	var request models.SubscriptionRequest
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to decode request body", err)
		return
	}

	saved, err := r.service.UpdateSubscription(req.Context(), subId, &request)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to update subscription", err)
		return
	}

	if saved {
		writeJSON(w, http.StatusNoContent, nil)
	} else {
		writeJSON(w, http.StatusNotModified, nil)
	}
}

// DeleteSubscription godoc
// @Summary Delete a subscription
// @Description Deletes a subscription
// @Tags subscription
// @Produce json
// @Param id path string true "subscription id"
// @Success 204
// @Router /api/subscriptions/{id} [delete]
func (r *Router) DeleteSubscription(w http.ResponseWriter, req *http.Request) {
	pathVariables := mux.Vars(req)
	id := pathVariables["id"]

	subUid, err := uuid.Parse(id)

	if err != nil {
		writeError(w, http.StatusBadRequest, "Subscription id format is incorrect (not uuid)", err)
		return
	}

	removed, err := r.service.RemoveSubscription(req.Context(), subUid)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to remove a subscription", err)
		return
	}

	if removed {
		writeJSON(w, http.StatusNoContent, nil)
	} else {
		writeJSON(w, http.StatusNotModified, nil)
	}
}

// CalculateSum godoc
// @Summary Calculate a sum price of subscriptions
// @Description Calculates a sum price of subscriptions for user between dates
// @Tags subscription
// @Produce json
// @Param user_id query string false "subscription user id"
// @Param subscription_name query string false "subscription name"
// @Param date_from query string false "subscription start date from"
// @Param date_to query string false "subscription start date to"
// @Success 200
// @Router /api/subscriptions/sum [get]
func (r *Router) CalculateSum(w http.ResponseWriter, req *http.Request) {
	urlQuery := req.URL.Query()
	userId := urlQuery.Get("user_id")
	subName := urlQuery.Get("subscription_name")
	dateFrom := urlQuery.Get("date_from")
	dateTo := urlQuery.Get("date_to")

	regex, err := regexp.Compile("^\\d{2}-\\d{4}$")

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to compile regex", err)
		return
	}

	if !regex.MatchString(dateFrom) {
		dateError := &errors.ParseError{ParsedString: dateFrom}
		writeError(w, http.StatusBadRequest, "Invalid date format", dateError)
		return
	}

	if !regex.MatchString(dateTo) {
		dateError := &errors.ParseError{ParsedString: dateTo}
		writeError(w, http.StatusBadRequest, "Invalid date format", dateError)
		return
	}

	if userId != "" {
		_, err := uuid.Parse(userId)
		if err != nil {
			writeError(w, http.StatusBadRequest, "User id format is incorrect (not uuid)", err)
		}
	}

	res, err := r.service.SubscriptionSum(req.Context(), userId, subName, dateFrom, dateTo)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get a subscription sum", err)
	}

	writeJSON(w, http.StatusOK, res)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
