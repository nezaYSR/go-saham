package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
	"gitlab.com/nezaysr/go-saham.git/config"
	"gitlab.com/nezaysr/go-saham.git/storage"
	"golang.org/x/crypto/bcrypt"
)

func Signin(signinPayload SigninPayload) (string, error) {
	db := storage.GetDBInstance()
	user := &Users{}
	if err := db.Where("username = ?", signinPayload.Username).First(&user).Error; err != nil {
		return "", err
	}

	// Verify the password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signinPayload.Password))
	if err != nil {
		return "", err
	}

	var jwtTokenPayload JWTTokenPayload
	jwtTokenPayload.ID = user.ID
	jwtTokenPayload.Username = user.Username
	jwtTokenPayload.Role = user.Role

	// Signing in jwt
	tokenString := signJWTToken(jwtTokenPayload)
	if tokenString == "" {
		return "", errors.New("failed to sign JWT token")
	}

	return tokenString, nil
}

func signJWTToken(jwtTokenPayload JWTTokenPayload) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       jwtTokenPayload.ID,
		"username": jwtTokenPayload.Username,
		"role":     jwtTokenPayload.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token with a secret key and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return ""
	}

	return tokenString
}

func GetUserList(rdb *config.Database, page int, limit int) ([]Users, error) {
	cacheKey := fmt.Sprintf("user_list:%d:%d", page, limit)
	cachedData, err := rdb.Client.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var users []Users

		err = json.Unmarshal([]byte(cachedData), &users)
		if err != nil {
			log.Printf("Failed to unmarshal cached user data: %v", err)
		}
		return users, nil
	} else if err != redis.Nil {
		log.Printf("Failed to get cached user data: %v", err)
	}

	// User list not found in cache
	db := storage.GetDBInstance()
	offset := (page - 1) * limit

	users := []Users{}
	if err := db.Offset(offset).Limit(limit).Order("id DESC").Find(&users).Error; err != nil {
		return nil, err
	}

	cacheValue, err := json.Marshal(users)
	if err != nil {
		log.Printf("Failed to marshal user data for caching: %v", err)
	}
	cacheTTL := 1 * time.Minute
	err = rdb.Client.Set(context.Background(), cacheKey, cacheValue, cacheTTL).Err()
	if err != nil {
		log.Printf("Failed to store user data in cache: %v", err)
	}

	return users, nil
}

func GetUserByID(user_id int) (*Users, error) {
	db := storage.GetDBInstance()
	user := &Users{}
	if err := db.Select("id, username, fullname,first_order_id,role,created_at, updated_at, deleted_at").Where("id=?", user_id).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func PostNewUser(userPayload InsertUserPayload) (int, error) {
	db := storage.GetDBInstance()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPayload.Password), 12)
	if err != nil {
		return 0, err
	}

	user := &Users{
		Username:  userPayload.Username,
		Fullname:  userPayload.Fullname,
		Password:  string(hashedPassword),
		Role:      "user",
		CreatedAt: time.Now(),
	}

	if err := db.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func UpdateAUserByID(userPayload UpdateUserPayload) error {
	db := storage.GetDBInstance()
	user := &Users{}

	if err := db.Model(user).Where("id = ?", userPayload.ID).UpdateColumns(
		map[string]interface{}{
			"fullname":       userPayload.Fullname,
			"first_order_id": userPayload.FirstOrderId,
			"updated_at": &NullableTime{
				Time:  time.Now(),
				Valid: true,
			},
		},
	).Error; err != nil {
		return err
	}
	return nil
}

func DeleteAUserByID(user_id int) error {
	db := storage.GetDBInstance()
	user := &Users{}

	if err := db.Model(user).Where("id = ?", user_id).Delete(user).Error; err != nil {
		return err
	}

	return nil
}

func GetOrderItemList(rdb *config.Database, page int, limit int) ([]OrdersItem, error) {
	cacheKey := fmt.Sprintf("order_item:%d:%d", page, limit)
	cachedData, err := rdb.Client.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var order_item []OrdersItem

		err = json.Unmarshal([]byte(cachedData), &order_item)
		if err != nil {
			log.Printf("Failed to unmarshal cached order_item data: %v", err)
		}
		return order_item, nil
	} else if err != redis.Nil {
		log.Printf("Failed to get cached order_item data: %v", err)
	}

	db := storage.GetDBInstance()
	offset := (page - 1) * limit

	order_item := []OrdersItem{}
	if err := db.Offset(offset).Limit(limit).Order("id DESC").Find(&order_item).Error; err != nil {
		return nil, err
	}

	cacheValue, err := json.Marshal(order_item)
	if err != nil {
		log.Printf("Failed to marshal order_item data for caching: %v", err)
	}
	cacheTTL := 1 * time.Minute
	err = rdb.Client.Set(context.Background(), cacheKey, cacheValue, cacheTTL).Err()
	if err != nil {
		log.Printf("Failed to store order_item data in cache: %v", err)
	}

	return order_item, nil
}

func GetOrderItemByID(orderItemID int) (*OrdersItem, error) {
	db := storage.GetDBInstance()
	order_item := &OrdersItem{}
	if err := db.Where("id=?", orderItemID).First(order_item).Error; err != nil {
		return nil, err
	}

	return order_item, nil
}

func PostNewOrderItem(orderItemPayload InsertOrderItemPayload) (int, error) {
	db := storage.GetDBInstance()

	order_item := &OrdersItem{
		Name:      orderItemPayload.Name,
		Price:     orderItemPayload.Price,
		ExpiredAt: orderItemPayload.ExpiredAt,
		CreatedAt: time.Now(),
	}

	if err := db.Create(order_item).Error; err != nil {
		return 0, err
	}
	return order_item.ID, nil
}

func UpdateOrderItemByID(orderItemPayload UpdateOrderItemPayload) error {
	db := storage.GetDBInstance()
	order_item := &OrdersItem{}

	if err := db.Model(order_item).Where("id = ?", orderItemPayload.ID).UpdateColumns(
		map[string]interface{}{
			"name":       orderItemPayload.Name,
			"price":      orderItemPayload.Price,
			"expired_at": orderItemPayload.ExpiredAt,
			"updated_at": &NullableTime{
				Time:  time.Now(),
				Valid: true,
			},
		},
	).Error; err != nil {
		return err
	}
	return nil
}

func DeleteOrderItemByID(orderItemID int) error {
	db := storage.GetDBInstance()
	order_item := &OrdersItem{}

	if err := db.Model(order_item).Where("id = ?", orderItemID).Delete(order_item).Error; err != nil {
		return err
	}

	return nil
}

func PostAnOrderHistory(orderHistoryPayload InsertOrderHistoryPayload) (int, error) {
	db := storage.GetDBInstance()

	order_histories := &OrdersHistories{
		UserId:       orderHistoryPayload.UserId,
		OrderItemId:  orderHistoryPayload.OrderItemId,
		Descriptions: orderHistoryPayload.Descriptions,
		CreatedAt:    time.Now(),
	}

	if err := db.Create(order_histories).Error; err != nil {
		return 0, err
	}
	return order_histories.ID, nil
}

func RemoveAnOrderHistory(orderHistoryID int) error {
	db := storage.GetDBInstance()
	order_history := &OrdersHistories{}

	if err := db.Model(order_history).Where("id = ?", orderHistoryID).Delete(order_history).Error; err != nil {
		return err
	}

	return nil
}

func GetOrderHistoriesByUserID(rdb *config.Database, page int, limit int, idRaw string) ([]OrdersHistories, error) {
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("order_histories:%d:%d", page, limit)
	cachedData, err := rdb.Client.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var order_histories []OrdersHistories

		err = json.Unmarshal([]byte(cachedData), &order_histories)
		if err != nil {
			log.Printf("Failed to unmarshal cached order_histories data: %v", err)
		}
		return order_histories, nil

	} else if err != redis.Nil {
		log.Printf("Failed to get cached order_histories data: %v", err)
	}

	db := storage.GetDBInstance()
	offset := (page - 1) * limit

	order_histories := []OrdersHistories{}
	if err := db.Where("user_id = ?", id).Offset(offset).Limit(limit).Order("id DESC").Find(&order_histories).Error; err != nil {
		return nil, err
	}

	cacheValue, err := json.Marshal(order_histories)
	if err != nil {
		log.Printf("Failed to marshal order_histories data for caching: %v", err)
	}
	cacheTTL := 1 * time.Minute
	err = rdb.Client.Set(context.Background(), cacheKey, cacheValue, cacheTTL).Err()
	if err != nil {
		log.Printf("Failed to store order_histories data in cache: %v", err)
	}

	return order_histories, nil
}

func GetAllOrderHistories(rdb *config.Database, page int, limit int) ([]OrdersHistories, error) {
	cacheKey := fmt.Sprintf("order_histories:%d:%d", page, limit)
	cachedData, err := rdb.Client.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var order_histories []OrdersHistories

		err = json.Unmarshal([]byte(cachedData), &order_histories)
		if err != nil {
			log.Printf("Failed to unmarshal cached order_histories data: %v", err)
		}
		return order_histories, nil

	} else if err != redis.Nil {
		log.Printf("Failed to get cached order_histories data: %v", err)
	}

	db := storage.GetDBInstance()
	offset := (page - 1) * limit

	order_histories := []OrdersHistories{}
	if err := db.Offset(offset).Limit(limit).Order("id DESC").Find(&order_histories).Error; err != nil {
		return nil, err
	}

	cacheValue, err := json.Marshal(order_histories)
	if err != nil {
		log.Printf("Failed to marshal order_histories data for caching: %v", err)
	}
	cacheTTL := 1 * time.Minute
	err = rdb.Client.Set(context.Background(), cacheKey, cacheValue, cacheTTL).Err()
	if err != nil {
		log.Printf("Failed to store order_histories data in cache: %v", err)
	}

	return order_histories, nil
}
