package handler

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"net/http"
	"onlineSubscriptions/internal/models"
	"onlineSubscriptions/internal/validators"
)

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
	logger.Printf("Request to find subscription with id: %v", id)

	subId, err := ParseUuidQuery(id)
	if id != "" && err != nil {
		writeError(w, http.StatusBadRequest, "Id format is incorrect (not uuid)", err)
		return
	}

	subscription, err := r.service.GetSubscriptionById(req.Context(), subId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "Subscription not found", nil)
		} else {
			writeError(w, http.StatusBadRequest, "Failed to get subscription by id", err)
		}
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
	rawPage, rawPerPage := urlQuery.Get("page"), urlQuery.Get("perPage")

	logger.Printf("Listing the page of subscriptions. Page number: %v; Page size: %v",
		rawPage, rawPerPage)

	page, err := ParseIntQuery(rawPage)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Page format is incorrect (not integer)", err)
		return
	}

	perPage, err := ParseIntQuery(rawPerPage)
	if err != nil {
		writeError(w, http.StatusBadRequest, "perPage format is incorrect (not integer)", err)
		return
	}

	list, err := r.service.ListSubscriptions(req.Context(), page, perPage)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to get a list of subscriptions", err)
		return
	}

	if list == nil {
		writeJSON(w, http.StatusOK, []models.Subscription{})
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

	if err := json.NewDecoder(req.Body).Decode(&subscriptionRequest); err != nil {
		writeError(w, http.StatusBadRequest, "Failed to decode request body", err)
		return
	}

	if err := subscriptionRequest.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	logger.Printf("Creating subscription with body: %+v", subscriptionRequest)

	saved, err := r.service.InsertSubscription(req.Context(), &subscriptionRequest)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to insert a subscription", err)
		return
	}

	if saved {
		writeJSON(w, http.StatusCreated, nil)
	} else {
		writeJSON(w, http.StatusConflict, nil)
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

	subId, err := ParseUuidQuery(id)

	if err != nil {
		writeError(w, http.StatusBadRequest, "Subscription id format is incorrect (not uuid)", err)
		return
	}

	var subscriptionRequest models.SubscriptionRequest
	err = json.NewDecoder(req.Body).Decode(&subscriptionRequest)

	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to decode request body", err)
		return
	}

	if err := subscriptionRequest.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	logger.Printf("Update request for subscription: %v, with body %+v", id, subscriptionRequest)

	saved, err := r.service.UpdateSubscription(req.Context(), subId, &subscriptionRequest)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "Subscription not found", nil)
		} else {
			writeError(w, http.StatusBadRequest, "Failed to update subscription", err)
		}
		return
	}

	if saved {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotModified)
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
	logger.Printf("Delete request for subscription: %v", id)

	subUid, err := ParseUuidQuery(id)

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
		w.WriteHeader(http.StatusNoContent)
	} else {
		writeError(w, http.StatusNotFound, "Subscription not found", nil)
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

	logger.Printf("Price sum calculation with parameters: "+
		"userId: %v, subName: %v, dateFrom: %v, dateTo: %v",
		userId, subName, dateFrom, dateTo)

	err := validators.ValidateDates(dateFrom, dateTo)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if userId != "" {
		_, err := ParseUuidQuery(userId)
		if err != nil {
			writeError(w, http.StatusBadRequest, "user_id format is incorrect (not uuid)", nil)
			return
		}
	}

	res, err := r.service.SubscriptionSum(req.Context(), userId, subName, dateFrom, dateTo)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get a subscription sum", err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}
