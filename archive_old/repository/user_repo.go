package repository

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"crowdfund/model"
)

type UserRepository interface {
	GetUserByID(id int) (*model.User, error)
	CreateUser(user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	LoginUser(user *model.User) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(id int) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(user *model.User) (*model.User, error) {
	//validate user data
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return nil, errors.New("name, email, and password are required")
	}
	if len(user.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	userPass := user.Password
	userPassHash, err := bcrypt.GenerateFromPassword([]byte(userPass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(userPassHash)

	if err := r.db.Omit("id").Create(user).Error; err != nil {
		return nil, err
	}

	if err := r.db.Last(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(user *model.User) (*model.User, error) {
	//validate user data
	if user.ID == 0 || user.Name == "" || user.Email == "" {
		return nil, errors.New("id, name, and email are required")
	}

	userPass := user.Password
	userPassHash, err := bcrypt.GenerateFromPassword([]byte(userPass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(userPassHash)

	if err := r.db.Model(user).Updates(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) LoginUser(user *model.User) (*model.User, error) {
	var userDb model.User

	//validate user data
	if user.Email == "" || user.Password == "" {
		return nil, errors.New("email and password are required")
	}

	if err := r.db.Where("email = ?", user.Email).First(&userDb).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(user.Password)); err != nil {
		return nil, err
	}

	return &userDb, nil
}
