package repositories

import "github.com/Koshsky/subs-service/auth-service/internal/models"

//go:generate mockery --name=UserRepositoryInterface --output=./mocks --outpkg=mocks --filename=UserRepositoryInterface.go
type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	ValidateUser(email, password string) (*models.User, error)
}

//go:generate mockery --name=DatabaseInterface --output=./mocks --outpkg=mocks --filename=DatabaseInterface.go
type DatabaseInterface interface {
	Create(value interface{}) DatabaseInterface
	Where(query interface{}, args ...interface{}) DatabaseInterface
	First(dest interface{}, conds ...interface{}) DatabaseInterface
	GetError() error
}

// Interface compliance checks - will fail at compile time if interfaces are not implemented
var _ UserRepositoryInterface = (*UserRepository)(nil)
var _ DatabaseInterface = (*GormAdapter)(nil)
