package controller

type Controller interface {
	NewController() Controller
	InitRoutes()
}
