package repoimpl

import (
	"context"
	"fmt"
	models "go-jwt-api/model"
	repo "go-jwt-api/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepoImpl struct {
	Db *mongo.Database
}

func NewProductRepo(db *mongo.Database) repo.ProductRepo {
	return &ProductRepoImpl{
		Db: db,
	}
}
func (pri *ProductRepoImpl) AddNewProduct(product models.Product) error {

	bbyte, _ := bson.Marshal(product)

	_, err := pri.Db.Collection("products").InsertOne(context.Background(), bbyte)

	if err != nil {
		return err
	}
	return nil
}

func (pri *ProductRepoImpl) FindProductByName(name string) (models.Product, error) {
	product := models.Product{}

	rs := pri.Db.Collection("products").
		FindOne(context.Background(),
			bson.M{"name": name})
	err := rs.Decode(&product)
	//
	// fmt.Println("Product in implement: ", rs)

	if err != nil {
		return product, err
	}

	return product, nil
}

func (pri *ProductRepoImpl) GetListProduct() ([]models.Product, error) {
	products := []models.Product{}
	collection := pri.Db.Collection("products")
	rs, err := collection.Find(context.Background(), bson.M{})
	fmt.Println(rs)
	if err != nil {
		return nil, err
	}
	for rs.Next(context.Background()) {
		var p models.Product
		if err := rs.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rs.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
