package api

import "onlineSubscriptions/src/main/impl/service"

type subscriptionApi struct {
	subscriptionService service.PostgresSubscriptionService
}
