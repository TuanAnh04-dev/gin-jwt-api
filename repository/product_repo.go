package repository

import models "go-jwt-api/model"

type ProductRepo interface {
	AddNewProduct(product models.Product) error
	FindProductByName(name string) (models.Product, error)
	GetListProduct() ([]models.Product, error)
}
