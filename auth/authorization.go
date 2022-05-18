package auth

import "kit/event_bus"

type Auth struct {
	eb *event_bus.EventBus
}

func NewAuth(eb *event_bus.EventBus) *Auth {
	return &Auth{eb}
}

func (a *Auth) IsAuthorized(id string, perm string) bool {
	return true
}

func (a *Auth) RegisterPermissions(perms []string) {
	a.eb.Publish("permissions.register", perms)
}
