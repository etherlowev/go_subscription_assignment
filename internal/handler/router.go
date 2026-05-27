package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	_ "onlineSubscriptions/docs"
	"onlineSubscriptions/internal/middleware"
	"onlineSubscriptions/internal/services"
	"os"
)

var logger = log.New(os.Stdout, "router: ", log.LstdFlags|log.Lshortfile)

type errorResponse struct {
	Msg string `json:"msg,omitempty"`
}

type Router struct {
	service services.SubscriptionService
	router  *mux.Router
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.Printf("json encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := errorResponse{msg}

	if err != nil {
		logger.Printf("%s: %v", msg, err)
	}

	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		logger.Printf("Json error has occurred when encoding response: %v", encodeErr)
	}
}

func NewRouter(service services.SubscriptionService) (*Router, error) {
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
	subrouter.HandleFunc("/sum", r.CalculateSum).Methods(http.MethodGet)
	subrouter.HandleFunc("/create", r.CreateSubscription).Methods(http.MethodPost)

	subrouter.HandleFunc("/{id}", r.FindSub).Methods(http.MethodGet)
	subrouter.HandleFunc("/{id}", r.UpdateSubscription).Methods(http.MethodPut)
	subrouter.HandleFunc("/{id}", r.DeleteSubscription).Methods(http.MethodDelete)

	return r, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
