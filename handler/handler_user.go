package handler

import (
	"encoding/json"
	config "go-jwt-api/config"
	driver "go-jwt-api/driver"
	models "go-jwt-api/model"

	repoImpl "go-jwt-api/repository/repoimpl"
	rs "go-jwt-api/response"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("abcdefghijklmnopq")

type Claims struct {
	Email              string `json:"email"`
	DisplayName        string `json:"displayName"`
	jwt.StandardClaims `json:"_"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var regData models.RegistrationData
	err := json.NewDecoder(r.Body).Decode(&regData)
	// fmt.Println("===========", err)

	if err != nil {
		rs.ResponseErr(w, http.StatusBadRequest)
		return
	}

	_, err = repoImpl.NewUserRepo(driver.Mongo.Client.
		Database(config.DB_NAME)).
		FindUserByEmail(regData.Email)

	if err != models.ERR_USER_NOT_FOUND {
		rs.ResponseErr(w, http.StatusConflict)
		return
	}

	user := models.User{
		Email:       regData.Email,
		Password:    regData.Password,
		DisplayName: regData.DisplayName,
	}
	err = repoImpl.NewUserRepo(driver.Mongo.Client.
		Database(config.DB_NAME)).Insert(user)
	if err != nil {
		rs.ResponseErr(w, http.StatusInternalServerError)
		return
	}

	var tokenString string
	tokenString, err = GenToken(user)
	if err != nil {
		rs.ResponseErr(w, http.StatusInternalServerError)
		return
	}

	rs.ResponseOk(w, models.RegisterResponse{
		Token:  tokenString,
		Status: http.StatusOK,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginData models.LoginData
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		rs.ResponseErr(w, http.StatusBadRequest)
		return
	}

	var user models.User
	user, err = repoImpl.NewUserRepo(driver.Mongo.Client.
		Database(config.DB_NAME)).
		CheckLoginInfo(loginData.Email, loginData.Password)
	if err != nil {
		rs.ResponseErr(w, http.StatusUnauthorized)
		return
	}

	var tokenString string
	tokenString, err = GenToken(user)
	if err != nil {
		rs.ResponseErr(w, http.StatusInternalServerError)
		return
	}

	rs.ResponseOk(w, models.RegisterResponse{
		Token:  tokenString,
		Status: http.StatusOK,
	})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")

	if tokenHeader == "" {
		rs.ResponseErr(w, http.StatusForbidden)
		return
	}

	splitted := strings.Split(tokenHeader, " ") // Bearer jwt_token
	if len(splitted) != 2 {
		rs.ResponseErr(w, http.StatusForbidden)
		return
	}

	tokenPart := splitted[1]
	// fmt.Println("=================Token part===================" + tokenPart)
	tk := &Claims{}

	token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		rs.ResponseErr(w, http.StatusInternalServerError)
		return
	}

	if token.Valid {
		rs.ResponseOk(w, token.Claims)
	}
}

func GenToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(3600 * time.Second)
	claims := &Claims{
		Email:       user.Email,
		DisplayName: user.DisplayName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
