package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/nezaysr/go-saham.git/config"
)

func Routes(e *echo.Echo, rdb *config.Database) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://*", "http://*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/ping/:your_name", heartbeat)

	// Auth Routes
	authRoutes := e.Group("/auth")
	authRoutes.POST("/si", SigninHandler) //SIGNIN user
	authRoutes.POST("/so", UserSignout)   //SIGNOUT user

	// User Routes
	userRoutes := e.Group("/users")
	userRoutes.Use(AuthenticationMiddleware)
	userRoutes.GET("/gl", RoleRequiredMiddleware(GetUsers(rdb), "admin"))          //GET user list
	userRoutes.GET("/g/:user_id", GetAUser)                                        //GET a user by ID
	userRoutes.POST("/c", RoleRequiredMiddleware(CreateAUser, "admin"))            //CREATE a new user
	userRoutes.PUT("/u/:user_id", UpdateAUser)                                     //UPDATE a user
	userRoutes.DELETE("/d/:user_id", RoleRequiredMiddleware(DeleteAUser, "admin")) //DELETE a user
	userRoutes.GET("/goi/:order_item_id", UserGetOrderItem)                        //GET a user gets an order item
	userRoutes.DELETE("/roi/:order_history_id", UserRemoveOrderItem)               //GET a user removes an order item

	// Order Item Routes
	orderItemRoutes := e.Group("/order_item")
	orderItemRoutes.Use(AuthenticationMiddleware)
	orderItemRoutes.GET("/gl", GetOrderItemList(rdb))                                               //GET order item list
	orderItemRoutes.GET("/g/:order_item_id", GetAnOrderItem)                                        //GET an order item by ID
	orderItemRoutes.POST("/c", RoleRequiredMiddleware(CreateAnOrderItem, "admin"))                  //CREATE a new order item
	orderItemRoutes.PUT("/u/:order_item_id", RoleRequiredMiddleware(UpdateAnOrderItem, "admin"))    //UPDATE an order item
	orderItemRoutes.DELETE("/d/:order_item_id", RoleRequiredMiddleware(DeleteAnOrderItem, "admin")) //DELETE an order item

	// Order Item Routes
	orderHistoriesRoutes := e.Group("/order_histories")
	orderHistoriesRoutes.Use(AuthenticationMiddleware)
	orderHistoriesRoutes.GET("/g", GetAnUsersOrderHistories(rdb))                            //GET an order histories by ID
	orderHistoriesRoutes.GET("/gl", RoleRequiredMiddleware(GetOrderHistories(rdb), "admin")) //GET order histories list
}
