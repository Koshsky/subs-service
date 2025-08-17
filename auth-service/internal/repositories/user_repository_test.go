package repositories_test

import (
	"errors"
	"testing"

	"github.com/Koshsky/subs-service/auth-service/internal/models"
	"github.com/Koshsky/subs-service/auth-service/internal/repositories"
	"github.com/Koshsky/subs-service/auth-service/internal/repositories/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDB *mocks.DatabaseInterface
	repo   *repositories.UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	suite.mockDB = new(mocks.DatabaseInterface)
	suite.repo = &repositories.UserRepository{DB: suite.mockDB}
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.mockDB.AssertExpectations(suite.T())
}

func (suite *UserRepositoryTestSuite) TestNewUserRepository() {
	// Arrange & Act
	repo := repositories.NewUserRepository(nil)

	// Assert
	suite.NotNil(repo)
	suite.NotNil(repo.DB)
}

func (suite *UserRepositoryTestSuite) TestCreateUser_NilUser() {
	// Arrange
	var user *models.User = nil

	// Act & Assert - should panic
	suite.Panics(func() {
		suite.repo.CreateUser(user)
	})
}

// CreateUser Tests
func (suite *UserRepositoryTestSuite) TestCreateUser_Success() {
	// Arrange
	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	suite.mockDB.On("Create", mock.AnythingOfType("*models.User")).Return(suite.mockDB)
	suite.mockDB.On("GetError").Return(nil)

	// Act
	err := suite.repo.CreateUser(user)

	// Assert
	suite.Require().NoError(err)
	suite.Require().NotNil(user)
	suite.NotEqual("password123", user.Password) // Password should be hashed
	suite.NotEmpty(user.Password)

	// Verify that bcrypt hash is valid
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	suite.Require().NoError(bcryptErr)
}

func (suite *UserRepositoryTestSuite) TestCreateUser_DatabaseError() {
	// Arrange
	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}
	dbError := errors.New("database error")

	suite.mockDB.On("Create", mock.AnythingOfType("*models.User")).Return(suite.mockDB)
	suite.mockDB.On("GetError").Return(dbError)

	// Act
	err := suite.repo.CreateUser(user)

	// Assert
	suite.Require().Error(err)
	suite.Require().NotNil(user)
	suite.Contains(err.Error(), "cannot create user")
	suite.ErrorAs(err, &dbError)

	// Verify that password was hashed before DB call
	suite.NotEqual("password123", user.Password)
	suite.NotEmpty(user.Password)

	// Verify that bcrypt hash is valid
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	suite.Require().NoError(bcryptErr)
}

func (suite *UserRepositoryTestSuite) TestCreateUser_BcryptError() {
	// Arrange
	user := &models.User{
		Email:    "test@example.com",
		Password: string(make([]byte, 73)), // 73 bytes > bcrypt limit
	}

	// Act
	err := suite.repo.CreateUser(user)

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "bcrypt")
}

// GetUserByEmail Tests
func (suite *UserRepositoryTestSuite) TestGetUserByEmail_Success() {
	// Arrange
	expectedEmail := "test@example.com"

	suite.mockDB.On("Where", "email = ?", expectedEmail).Return(suite.mockDB)
	suite.mockDB.On("First", mock.AnythingOfType("*models.User")).Return(suite.mockDB)
	suite.mockDB.On("GetError").Return(nil)

	// Act
	user, err := suite.repo.GetUserByEmail(expectedEmail)

	// Assert
	suite.Require().NoError(err)
	suite.Require().NotNil(user)
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail_NotFound() {
	// Arrange
	dbError := errors.New("record not found")

	suite.mockDB.On("Where", "email = ?", "nonexistent@example.com").Return(suite.mockDB)
	suite.mockDB.On("First", mock.AnythingOfType("*models.User")).Return(suite.mockDB)
	suite.mockDB.On("GetError").Return(dbError)

	// Act
	user, err := suite.repo.GetUserByEmail("nonexistent@example.com")

	// Assert
	suite.Error(err)
	suite.Nil(user)
	suite.ErrorAs(err, &dbError)
}

// ValidateUser Tests
func (suite *UserRepositoryTestSuite) TestValidateUser_Success() {
	// Arrange
	email := "test@example.com"
	password := "correctpassword"

	suite.mockDB.On("Where", "email = ?", email).Return(suite.mockDB)
	suite.mockDB.On("First", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		user.Email = email
		// Generate correct hash for the password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}).Return(suite.mockDB)
	suite.mockDB.On("GetError").Return(nil)

	// Act
	user, err := suite.repo.ValidateUser(email, password)

	// Assert
	suite.Require().NoError(err)
	suite.Require().NotNil(user)
	suite.Equal(email, user.Email)
}

func (suite *UserRepositoryTestSuite) TestValidateUser_UserNotFound() {
	// Arrange
	suite.mockDB.On("Where", "email = ?", "nonexistent@example.com").Return(suite.mockDB)
	suite.mockDB.On("First", mock.AnythingOfType("*models.User")).Return(suite.mockDB)
	suite.mockDB.On("GetError").Return(errors.New("record not found"))

	// Act
	user, err := suite.repo.ValidateUser("nonexistent@example.com", "password123")

	// Assert
	suite.Error(err)
	suite.Nil(user)
	suite.Contains(err.Error(), "cannot get user by email")
	suite.Contains(err.Error(), "record not found")
}

func (suite *UserRepositoryTestSuite) TestValidateUser_InvalidPassword() {
	// Arrange
	suite.mockDB.On("Where", "email = ?", "test@example.com").Return(suite.mockDB)
	suite.mockDB.On("First", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		user.Email = "test@example.com"
		user.Password = "correctpassword"
	}).Return(suite.mockDB)
	suite.mockDB.On("GetError").Return(nil)

	// Act
	user, err := suite.repo.ValidateUser("test@example.com", "wrongpassword")

	// Assert
	suite.Error(err)
	suite.Nil(user)
	suite.Contains(err.Error(), "authentication failed")
	suite.Contains(err.Error(), "bcrypt")
}

// Run tests
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
