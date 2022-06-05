package controllers

type IController interface {
	Auth() IAuthController
}

type Controller struct {
	authController AuthController
}

func (controller *Controller) Auth() IAuthController {
	return &controller.authController
}

func InitControllers() IController {
	controller := Controller{
		authController: AuthController{},
	}

	return &controller
}
