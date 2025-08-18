package repositories_test

import (
	"errors"
	"testing"

	"github.com/Koshsky/subs-service/auth-service/internal/models"
	"github.com/Koshsky/subs-service/auth-service/internal/repositories"
	"github.com/Koshsky/subs-service/auth-service/internal/repositories/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDB   *mocks.IDatabase
	userRepo *repositories.UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	suite.mockDB = new(mocks.IDatabase)
	suite.userRepo = &repositories.UserRepository{DB: suite.mockDB}
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.mockDB.AssertExpectations(suite.T())
}

// ===== CONSTRUCTOR TESTS =====

func (suite *UserRepositoryTestSuite) TestNewUserRepository_Success() {
	// Arrange
	mockDB := new(mocks.IDatabase)

	// Act
	repo := repositories.NewUserRepository(mockDB)

	// Assert
	suite.Require().NotNil(repo)
	suite.Equal(mockDB, repo.DB)
}

// ===== CREATE USER TESTS =====

func (suite *UserRepositoryTestSuite) TestCreateUser_Success() {
	// Arrange
	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	suite.mockDB.On("Create", user).Return(suite.mockDB)
	suite.mockDB.On("GetError").Return(nil)

	// Act
	err := suite.userRepo.CreateUser(user)

	// Assert
	suite.Require().NoError(err)
	suite.Require().NotEqual(uuid.Nil, user.ID)
	suite.mockDB.AssertExpectations(suite.T())
}

func (suite *UserRepositoryTestSuite) TestCreateUser_NilDatabase() {
	// Arrange
	repo := &repositories.UserRepository{DB: nil}
	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Act
	err := repo.CreateUser(user)

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "database connection is not initialized")
}

func (suite *UserRepositoryTestSuite) TestCreateUser_WithExistingUUID() {
	// Arrange
	existingUUID := uuid.New()
	user := &models.User{
		ID:       existingUUID,
		Email:    "test@example.com",
		Password: "password123",
	}

	suite.mockDB.On("Create", user).Return(suite.mockDB)
	suite.mockDB.On("GetError").Return(nil)

	// Act
	err := suite.userRepo.CreateUser(user)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(existingUUID, user.ID) // UUID should not change
	suite.mockDB.AssertExpectations(suite.T())
}

func (suite *UserRepositoryTestSuite) TestCreateUser_DatabaseError() {
	// Arrange
	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}
	expectedError := errors.New("database error")

	suite.mockDB.On("Create", user).Return(suite.mockDB)
	suite.mockDB.On("GetError").Return(expectedError)

	// Act
	err := suite.userRepo.CreateUser(user)

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "database error")
	suite.mockDB.AssertExpectations(suite.T())
}

// ===== GET USER BY EMAIL TESTS =====

func (suite *UserRepositoryTestSuite) TestGetUserByEmail_Success() {
	// Arrange
	email := "test@example.com"

	suite.mockDB.On("Where", "email = ?", email).Return(suite.mockDB)
	suite.mockDB.On("First", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		user.ID = uuid.New()
		user.Email = email
		user.Password = "hashed_password"
	}).Return(suite.mockDB)
	suite.mockDB.On("GetError").Return(nil)

	// Act
	result, err := suite.userRepo.GetUserByEmail(email)

	// Assert
	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Equal(email, result.Email)
	suite.mockDB.AssertExpectations(suite.T())
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail_NilDatabase() {
	// Arrange
	repo := &repositories.UserRepository{DB: nil}

	// Act
	user, err := repo.GetUserByEmail("test@example.com")

	// Assert
	suite.Require().Error(err)
	suite.Require().Nil(user)
	suite.Contains(err.Error(), "database connection is not initialized")
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail_UserNotFound() {
	// Arrange
	email := "test@example.com"

	suite.mockDB.On("Where", "email = ?", email).Return(suite.mockDB)
	suite.mockDB.On("First", mock.AnythingOfType("*models.User")).Return(suite.mockDB)
	suite.mockDB.On("GetError").Return(gorm.ErrRecordNotFound)

	// Act
	user, err := suite.userRepo.GetUserByEmail(email)

	// Assert
	suite.Require().Error(err)
	suite.Require().Nil(user)
	suite.Contains(err.Error(), "cannot get user by email")
	suite.mockDB.AssertExpectations(suite.T())
}

// Run tests
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
