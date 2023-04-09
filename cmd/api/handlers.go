package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gitlab.com/nezaysr/go-saham.git/config"
	data "gitlab.com/nezaysr/go-saham.git/data"
)

func heartbeat(c echo.Context) error {
	name := c.Param("your_name")
	payload := jsonResponse{
		Error:   false,
		Message: "Hello " + name + " Don't worry i'm alive",
		Data:    "http://localhost:3000/" + name,
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func SigninHandler(c echo.Context) error {
	tokenStr, err := userSignin(c)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusInternalServerError)
	}

	if err := setJWTToCookie(c, tokenStr); err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusInternalServerError)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Successfully Login",
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func userSignin(c echo.Context) (string, error) {
	var requestPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := readJSON(c.Response().Writer, c.Request(), &requestPayload)
	if err != nil {
		return "", err
	}

	signinPayload := data.SigninPayload{
		Username: requestPayload.Username,
		Password: requestPayload.Password,
	}

	tokenString, err := data.Signin(signinPayload)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func UserSignout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(c.Response().Writer, cookie)

	payload := jsonResponse{
		Error:   false,
		Message: "Successfully logged out",
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func GetUsers(rdb *config.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		pageSize := 10 // default page size
		page := 1      // default page

		if pageSizeParam := c.QueryParam("pageSize"); pageSizeParam != "" {
			pageSize, _ = strconv.Atoi(pageSizeParam)
		}

		if pageParam := c.QueryParam("page"); pageParam != "" {
			page, _ = strconv.Atoi(pageParam)
		}

		users, err := data.GetUserList(rdb, page, pageSize)
		if err != nil {
			return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
		}

		payload := jsonResponse{
			Error:   false,
			Message: "Users list",
			Data:    users,
		}

		return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
	}
}

func GetAUser(c echo.Context) error {
	userIDRaw := c.Param("user_id")

	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	user, err := data.GetUserByID(userID)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "User with id " + userIDRaw,
		Data:    user,
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func CreateAUser(c echo.Context) error {
	var requestPayload struct {
		Username string `gorm:"size:255;not null;unique" json:"username"`
		Fullname string `gorm:"size:255;not null" json:"fullname"`
	}

	generatedPassword := strings.ReplaceAll(uuid.New().String(), "-", "")

	err := readJSON(c.Response().Writer, c.Request(), &requestPayload)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusInternalServerError)
	}

	user := data.InsertUserPayload{
		Username: requestPayload.Username,
		Fullname: requestPayload.Fullname,
		Password: generatedPassword,
	}

	newID, err := data.PostNewUser(user)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusInternalServerError)
	}

	newIDString := strconv.Itoa(newID)

	payload := jsonResponse{
		Error:   false,
		Message: "User with id " + newIDString + " has been created, keep your password below",
		Data:    generatedPassword,
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func UpdateAUser(c echo.Context) error {
	userIDRaw := c.Param("user_id")

	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	var requestPayload struct {
		Fullname     *string `json:"fullname,omitempty"`
		FirstOrderId *int    `json:"first_order_id,omitempty"`
	}

	err = readJSON(c.Response().Writer, c.Request(), &requestPayload)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusInternalServerError)
	}

	user := data.UpdateUserPayload{
		ID:           userID,
		Fullname:     requestPayload.Fullname,
		FirstOrderId: requestPayload.FirstOrderId,
	}

	err = data.UpdateAUserByID(user)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "User with id " + userIDRaw + " has been updated",
		Data:    "user updated",
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func DeleteAUser(c echo.Context) error {
	userIDRaw := c.Param("user_id")

	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	err = data.DeleteAUserByID(userID)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "User with id " + userIDRaw + " has been deleted",
		Data:    "user deleted",
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func GetOrderItemList(rdb *config.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		pageSize := 10 // default page size
		page := 1      // default page

		if pageSizeParam := c.QueryParam("pageSize"); pageSizeParam != "" {
			pageSize, _ = strconv.Atoi(pageSizeParam)
		}

		if pageParam := c.QueryParam("page"); pageParam != "" {
			page, _ = strconv.Atoi(pageParam)
		}

		order_item, err := data.GetOrderItemList(rdb, page, pageSize)
		if err != nil {
			return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
		}

		payload := jsonResponse{
			Error:   false,
			Message: "Order Item list",
			Data:    order_item,
		}

		return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
	}
}

func GetAnOrderItem(c echo.Context) error {
	orderItemIDRaw := c.Param("order_item_id")

	orderItemID, err := strconv.Atoi(orderItemIDRaw)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	orderItem, err := data.GetOrderItemByID(orderItemID)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Order Item with id " + orderItemIDRaw,
		Data:    orderItem,
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func CreateAnOrderItem(c echo.Context) error {
	var requestPayload struct {
		Name      string    `gorm:"size:255;not null;unique" json:"name"`
		Price     int       `json:"price"`
		ExpiredAt time.Time `json:"expired_at,omitempty"`
	}

	err := readJSON(c.Response().Writer, c.Request(), &requestPayload)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusInternalServerError)
	}

	order_item := data.InsertOrderItemPayload{
		Name:      requestPayload.Name,
		Price:     requestPayload.Price,
		ExpiredAt: requestPayload.ExpiredAt,
	}

	newID, err := data.PostNewOrderItem(order_item)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusInternalServerError)
	}

	newIDString := strconv.Itoa(newID)

	payload := jsonResponse{
		Error:   false,
		Message: "Order Item with id " + newIDString + " has been created",
		Data:    "order item created",
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func UpdateAnOrderItem(c echo.Context) error {
	orderItemIDRaw := c.Param("order_item_id")

	orderItemID, err := strconv.Atoi(orderItemIDRaw)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	var requestPayload struct {
		Name      string    `json:"name,omitempty"`
		Price     int       `json:"price,omitempty"`
		ExpiredAt time.Time `json:"expired_at,omitempty"`
	}

	err = readJSON(c.Response().Writer, c.Request(), &requestPayload)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusInternalServerError)
	}

	order_item := data.UpdateOrderItemPayload{
		ID:        orderItemID,
		Name:      requestPayload.Name,
		Price:     requestPayload.Price,
		ExpiredAt: requestPayload.ExpiredAt,
	}

	err = data.UpdateOrderItemByID(order_item)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Order Item with id " + orderItemIDRaw + " has been updated",
		Data:    "order item updated",
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func DeleteAnOrderItem(c echo.Context) error {
	orderItemIDRaw := c.Param("order_item_id")

	orderItemID, err := strconv.Atoi(orderItemIDRaw)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	err = data.DeleteOrderItemByID(orderItemID)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Order Item with id " + orderItemIDRaw + " has been deleted",
		Data:    "order item deleted",
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func UserGetOrderItem(c echo.Context) error {
	id := c.Get("id").(string)
	var requestPayload struct {
		Descriptions *string `json:"descriptions,omitempty"`
	}

	orderItemIDRaw := c.Param("order_item_id")

	orderItemID, err := strconv.Atoi(orderItemIDRaw)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	user, err := data.GetUserByID(userID)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	orderItem, err := data.GetOrderItemByID(orderItemID)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	if user.FirstOrderId == nil {
		user := data.UpdateUserPayload{
			ID:           userID,
			FirstOrderId: &orderItem.ID,
		}

		err = data.UpdateAUserByID(user)
		if err != nil {
			return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
		}
	}

	order_history := data.InsertOrderHistoryPayload{
		UserId:       userID,
		OrderItemId:  orderItemID,
		Descriptions: requestPayload.Descriptions,
	}

	_, err = data.PostAnOrderHistory(order_history)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: user.Fullname + " just bought a " + orderItem.Name,
		Data:    "Order History created",
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func UserRemoveOrderItem(c echo.Context) error {
	id := c.Get("id").(string)
	orderHistoryIDRaw := c.Param("order_history_id")

	orderHistoryID, err := strconv.Atoi(orderHistoryIDRaw)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	err = data.RemoveAnOrderHistory(orderHistoryID)
	if err != nil {
		return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Order History with id " + orderHistoryIDRaw + " owned by user id " + id + " has been deleted",
		Data:    "order history deleted",
	}

	return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
}

func GetOrderHistories(rdb *config.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		pageSize := 10 // default page size
		page := 1      // default page

		if pageSizeParam := c.QueryParam("pageSize"); pageSizeParam != "" {
			pageSize, _ = strconv.Atoi(pageSizeParam)
		}

		if pageParam := c.QueryParam("page"); pageParam != "" {
			page, _ = strconv.Atoi(pageParam)
		}

		order_histories, err := data.GetAllOrderHistories(rdb, page, pageSize)
		if err != nil {
			return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
		}

		payload := jsonResponse{
			Error:   false,
			Message: "Order Histories",
			Data:    order_histories,
		}

		return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
	}
}

func GetAnUsersOrderHistories(rdb *config.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Get("id").(string)

		pageSize := 10 // default page size
		page := 1      // default page

		if pageSizeParam := c.QueryParam("pageSize"); pageSizeParam != "" {
			pageSize, _ = strconv.Atoi(pageSizeParam)
		}

		if pageParam := c.QueryParam("page"); pageParam != "" {
			page, _ = strconv.Atoi(pageParam)
		}

		order_histories, err := data.GetOrderHistoriesByUserID(rdb, page, pageSize, id)
		if err != nil {
			return errorJSON(c.Response().Writer, err, http.StatusBadRequest)
		}

		payload := jsonResponse{
			Error:   false,
			Message: "Order Histories",
			Data:    order_histories,
		}

		return writeJSON(c.Response().Writer, http.StatusAccepted, payload)
	}
}
