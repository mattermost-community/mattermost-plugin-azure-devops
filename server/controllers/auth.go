package controllers

import "net/http"

type IAuthController interface {
	SignIn(w http.ResponseWriter, req *http.Request)
	SignOut()
}

type AuthController struct {
}

func (auth *AuthController) SignIn(w http.ResponseWriter, req *http.Request) {
	// if we've made it here, we're authorized.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"is_tested_new": true}`))
}

func (auth *AuthController) SignOut() {

}
