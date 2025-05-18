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

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"

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

func GetOneProduct(c *gin.Context) {
	tokenHeader := c.GetHeader("Authorization")

	if tokenHeader == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Authorization header is missing"})
		return
	}

	splitted := strings.Split(tokenHeader, " ") // Bearer jwt_token
	if len(splitted) != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Authorization header is missing"})
		return
	}

	tokenPart := splitted[1]
	// fmt.Println("=================Token part===================" + tokenPart)
	tk := &Claims{}

	token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		// fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token không hợp lệ hoặc đã hết hạn"})
		return
	}

	if token.Valid {
		// product := models.Product{}
		var req FindOneProductRequest
		err := json.NewDecoder(c.Request.Body).Decode(&req)
		if err != nil {
			rsp.ResponseErr(c.Writer, http.StatusBadRequest)
			return
		}
		fmt.Println("Param From request(handler): ", req.Name)
		product, err := repoImpl.NewProductRepo(driver.Mongo.Client.
			Database(config.DB_NAME)).
			FindProductByName(req.Name)
		fmt.Println("Found product: ", product)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"product": product})

	}
}
func GetListProduct(c *gin.Context) {
	var rs = []models.Product{}

	rs, err := repoImpl.NewProductRepo(driver.Mongo.Client.
		Database(config.DB_NAME)).GetListProduct()
	if err != nil {
		// fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product list"})
		return
	}
	// rsp.ResponseOk(w, rs)
	c.JSON(http.StatusOK, rs)

	// }
}
func AddNewProduct(c *gin.Context) {
	var addData AddProductData
	err := json.NewDecoder(c.Request.Body).Decode(&addData)
	if err != nil {
		rsp.ResponseErr(c.Writer, http.StatusBadRequest)
		return
	}
	if errors := validation.ValidateStruct(addData); errors != nil {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(map[string]interface{}{"errors": errors})
		return
	}
	_, err = repoImpl.NewProductRepo(driver.Mongo.Client.
		Database(config.DB_NAME)).
		FindProductByName(addData.Name)

	if err == nil {
		// rsp.ResponseErr(c.Writer, http.StatusConflict)
		c.JSON(http.StatusConflict, gin.H{"error": "San pham da ton tai"})
		return
	}
	newProduct := models.Product{
		Name:     addData.Name,
		Price:    addData.Price,
		Quantity: addData.Quantity,
	}
	err = repoImpl.NewProductRepo(driver.Mongo.Client.Database(config.DB_NAME)).AddNewProduct(newProduct)

	if err != nil {
		fmt.Println("Khong thanh cong", err)
		// rsp.ResponseErr(w, http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add new product"})
	}
	rsp.ResponseOk(c.Writer, newProduct)
}
