package core

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var client *redis.Client

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUuid string
	User       User
}

/*
type User struct {
	ID        int64     `bson:"_id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	Roles     []string  `bson:"roles" json:"roles"`
	Active    bool      `bson:"active" json:"active"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
*/

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"email,required"`
	Password  string             `bson:"password" json:"password" validate:"required,min=6"`
	Salt      string             `bson:"salt" json:"salt"`
	Roles     []string           `bson:"roles" json:"roles"`
	Active    bool               `bson:"active" json:"active"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func redisConn() (*redis.Client, error) {
	dsn := os.Getenv("REDIS_DSN")

	fmt.Printf(dsn)

	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}

func CreateToken(userId string) (*TokenDetails, error) {
	var err error

	_id := userId

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = _id

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = _id + "-rf"

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = _id
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)

	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = _id
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)

	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))

	if err != nil {
		return nil, err
	}

	return td, nil
}

func VerifyPassword(password string, passwordHash string, salt string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password+salt))

	if err != nil {
		return err
	}

	return nil
}

func CreateAuth(user []byte, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	//dsn := os.Getenv("REDIS_DSN")
	//
	//fmt.Printf(dsn)
	//
	//client := redis.NewClient(&redis.Options{
	//	Addr:     dsn,
	//	Password: "",
	//	DB:       0,
	//})

	//_, err := client.Ping().Result()

	//useJson, _ := json.Unmarshal(user, *User)

	client, errConn := redisConn()

	if errConn != nil {
		return err
	}

	err := client.Set(td.AccessUuid, user, at.Sub(now)).Err()

	if err != nil {
		return err
	}

	err = client.Set(td.RefreshUuid, user, rt.Sub(now)).Err()

	if err != nil {
		return err
	}

	return nil
}
