package services

type IAuthService interface {
	SignIn()
	SignOut()
}

type AuthService struct {
}

func (auth *AuthService) SignIn() {

}

func (auth *AuthService) SignOut() {

}
