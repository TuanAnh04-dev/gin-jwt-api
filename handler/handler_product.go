package handler

import (
	"encoding/json"
	"fmt"
	"go-jwt-api/config"
	driver "go-jwt-api/driver"
	models "go-jwt-api/model"
	repoImpl "go-jwt-api/repository/repoimpl"
	"go-jwt-api/validation"
	"net/http"
	"strings"

	rsp "go-jwt-api/response"

	"github.com/dgrijalva/jwt-go"
)

type AddProductData struct {
	Name     string  `json:"productName" validate:"required,min=2,max=100"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	Quantity int     `json:"quantity" validate:"required,gt=0"`
	// jwt.StandardClaims `json:"_"`

}

type FindOneProductRequest struct {
	Name string `json:"productName"`
}

func GetOneProduct(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")

	if tokenHeader == "" {
		rsp.ResponseErr(w, http.StatusForbidden)
		return
	}

	splitted := strings.Split(tokenHeader, " ") // Bearer jwt_token
	if len(splitted) != 2 {
		rsp.ResponseErr(w, http.StatusForbidden)
		return
	}

	tokenPart := splitted[1]
	fmt.Println("=================Token part===================" + tokenPart)
	tk := &Claims{}

	token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		// fmt.Println(err)
		rsp.ResponseErr(w, http.StatusInternalServerError)
		return
	}

	if token.Valid {
		product := models.Product{}
		var req FindOneProductRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			rsp.ResponseErr(w, http.StatusBadRequest)
			return
		}
		fmt.Println("Param From request(handler): ", req.Name)
		product, err = repoImpl.NewProductRepo(driver.Mongo.Client.
			Database(config.DB_NAME)).
			FindProductByName(req.Name)
		fmt.Println("Found product: ", product)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Product not found!", http.StatusNotFound)
			return
		}
		rsp.ResponseOk(w, product)

	}
}
func GetListProduct(w http.ResponseWriter, r *http.Request) {
	// tokenHeader := r.Header.Get("Authorization")

	// if tokenHeader == "" {
	// 	rsp.ResponseErr(w, http.StatusForbidden)
	// 	return
	// }

	// splitted := strings.Split(tokenHeader, " ") // Bearer jwt_token
	// if len(splitted) != 2 {
	// 	rsp.ResponseErr(w, http.StatusForbidden)
	// 	return
	// }

	// tokenPart := splitted[1]
	// fmt.Println("=================Token part===================" + tokenPart)
	// tk := &Claims{}

	// token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
	// 	return jwtKey, nil
	// })

	// if err != nil {
	// 	// fmt.Println(err)
	// 	rsp.ResponseErr(w, http.StatusInternalServerError)
	// 	return
	// }

	// if token.Valid {
	var rs = []models.Product{}

	rs, err := repoImpl.NewProductRepo(driver.Mongo.Client.
		Database(config.DB_NAME)).GetListProduct()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Product not found!", http.StatusNotFound)
		return
	}
	rsp.ResponseOk(w, rs)

	// }
}
func AddNewProduct(w http.ResponseWriter, r *http.Request) {
	var addData AddProductData
	err := json.NewDecoder(r.Body).Decode(&addData)
	if err != nil {
		rsp.ResponseErr(w, http.StatusBadRequest)
		return
	}
	if errors := validation.ValidateStruct(addData); errors != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errors})
		return
	}
	_, err = repoImpl.NewProductRepo(driver.Mongo.Client.
		Database(config.DB_NAME)).
		FindProductByName(addData.Name)

	if err == nil {
		rsp.ResponseErr(w, http.StatusConflict)
		return
	}
	newProduct := models.Product{
		Name:     addData.Name,
		Price:    addData.Price,
		Quantity: addData.Quantity,
	}
	// fmt.Println("add_Data: ", addData)
	// fmt.Println(newProduct)
	err = repoImpl.NewProductRepo(driver.Mongo.Client.Database(config.DB_NAME)).AddNewProduct(newProduct)

	if err != nil {
		fmt.Println("Khong thanh cong", err)
		rsp.ResponseErr(w, http.StatusInternalServerError)
	}
	rsp.ResponseOk(w, newProduct)
}
