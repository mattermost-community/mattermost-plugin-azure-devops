package services

type IServices interface {
	Auth() IAuthService
}

type Services struct {
	authService AuthService
}

func (service *Services) Auth() IAuthService {
	return &service.authService
}

func InitServices() IServices {
	services := Services{
		authService: AuthService{},
	}

	return &services
}
