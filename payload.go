package driplimit

type Payload interface {
	ServiceToken() string
}

type payload struct {
	serviceToken string
}

func (a *payload) ServiceToken() string {
	if a == nil {
		return ""
	}
	return a.serviceToken
}
