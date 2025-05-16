package repoimpl

import (
	"context"
	"fmt"
	models "go-jwt-api/model"
	repo "go-jwt-api/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepoImpl struct {
	Db *mongo.Database
}

func NewUserRepo(db *mongo.Database) repo.UserRepo {
	return &UserRepoImpl{
		Db: db,
	}
}

func (mongo *UserRepoImpl) FindUserByEmail(email string) (models.User, error) {
	user := models.User{}
	rs := mongo.Db.Collection("users").FindOne(context.Background(), bson.M{"email": email})

	err := rs.Decode(&user)
	if user == (models.User{}) {
		return user, models.ERR_USER_NOT_FOUND
	}
	if err != nil {
		return user, err
	}

	return user, nil
}

func (mongo *UserRepoImpl) Insert(user models.User) error {
	bbyte, _ := bson.Marshal(user)
	_, err := mongo.Db.Collection("users").InsertOne(context.Background(), bbyte)

	if err != nil {
		return err
	}

	return nil
}

func (mongo *UserRepoImpl) CheckLoginInfo(email, password string) (models.User, error) {
	user := models.User{}
	rs := mongo.Db.Collection("users").
		FindOne(context.Background(), bson.M{"email": email, "password": password})
	err := rs.Decode(&user)
	fmt.Println(err)
	if err != nil {
		return user, err
	}
	return user, nil
}
