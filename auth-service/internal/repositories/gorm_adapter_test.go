package repositories_test

import (
	"testing"

	"github.com/Koshsky/subs-service/auth-service/internal/repositories"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormAdapterTestSuite struct {
	suite.Suite
}

// TestUser - helper struct for tests
type TestUser struct {
	ID    uint   `gorm:"primarykey"`
	Email string `gorm:"uniqueIndex"`
}

// setupTestDB creates in-memory SQLite database for tests
func (suite *GormAdapterTestSuite) setupTestDB() (*gorm.DB, repositories.IDatabase) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	adapter := repositories.NewGormAdapter(db)

	err = db.AutoMigrate(&TestUser{})
	suite.Require().NoError(err)

	return db, adapter
}

// ===== CONSTRUCTOR TESTS =====

func (suite *GormAdapterTestSuite) TestNewGormAdapter_NilDB() {
	// Arrange & Act
	adapter := repositories.NewGormAdapter(nil)

	// Assert
	suite.Require().NotNil(adapter)
	suite.Require().IsType(&repositories.GormAdapter{}, adapter)
}

func (suite *GormAdapterTestSuite) TestNewGormAdapter_Success() {
	// Arrange
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// Act
	adapter := repositories.NewGormAdapter(db)

	// Assert
	suite.Require().NotNil(adapter)
	suite.Require().IsType(&repositories.GormAdapter{}, adapter)
}

// ===== METHOD TESTS =====

func (suite *GormAdapterTestSuite) TestCreateWithRealDB() {
	// Arrange
	_, adapter := suite.setupTestDB()
	user := &TestUser{Email: "test@example.com"}

	// Act
	createResult := adapter.Create(user)

	// Assert
	suite.Require().NotNil(createResult)
	suite.Require().IsType(&repositories.GormAdapter{}, createResult)
	suite.Require().NoError(createResult.GetError())
}

func (suite *GormAdapterTestSuite) TestCreateWithNilDB() {
	// Arrange
	adapter := repositories.NewGormAdapter(nil)
	user := &TestUser{Email: "test@example.com"}

	// Act
	createResult := adapter.Create(user)

	// Assert
	suite.Require().NotNil(createResult)
	suite.Require().IsType(&repositories.GormAdapter{}, createResult)
	suite.Require().Error(createResult.GetError())
	suite.Contains(createResult.GetError().Error(), "database is nil")
}

func (suite *GormAdapterTestSuite) TestWhereWithRealDB() {
	// Arrange
	_, adapter := suite.setupTestDB()

	// Act
	whereResult := adapter.Where("email = ?", "test@example.com")

	// Assert
	suite.Require().NotNil(whereResult)
	suite.Require().IsType(&repositories.GormAdapter{}, whereResult)
	suite.Require().NoError(whereResult.GetError())
}

func (suite *GormAdapterTestSuite) TestWhereWithNilDB() {
	// Arrange
	adapter := repositories.NewGormAdapter(nil)

	// Act
	whereResult := adapter.Where("email = ?", "test@example.com")

	// Assert
	suite.Require().NotNil(whereResult)
	suite.Require().IsType(&repositories.GormAdapter{}, whereResult)
	suite.Require().Error(whereResult.GetError())
	suite.Contains(whereResult.GetError().Error(), "database is nil")
}

func (suite *GormAdapterTestSuite) TestFirstWithRealDB() {
	// Arrange
	_, adapter := suite.setupTestDB()

	// Create test data
	user := &TestUser{Email: "test@example.com"}
	adapter.Create(user)

	// Act
	var foundUser TestUser
	firstResult := adapter.First(&foundUser, "email = ?", "test@example.com")

	// Assert
	suite.Require().NotNil(firstResult)
	suite.Require().IsType(&repositories.GormAdapter{}, firstResult)
	suite.Require().NoError(firstResult.GetError())
	suite.Equal("test@example.com", foundUser.Email)
}

func (suite *GormAdapterTestSuite) TestFirstWithNilDB() {
	// Arrange
	adapter := repositories.NewGormAdapter(nil)

	// Act
	var foundUser TestUser
	firstResult := adapter.First(&foundUser, "email = ?", "test@example.com")

	// Assert
	suite.Require().NotNil(firstResult)
	suite.Require().IsType(&repositories.GormAdapter{}, firstResult)
	suite.Require().Error(firstResult.GetError())
	suite.Contains(firstResult.GetError().Error(), "database is nil")
}

func (suite *GormAdapterTestSuite) TestGetErrorWithNilDB() {
	// Arrange
	adapter := repositories.NewGormAdapter(nil)

	// Act
	err := adapter.GetError()

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "database is nil")
}

func (suite *GormAdapterTestSuite) TestMethodChaining() {
	// Arrange
	_, adapter := suite.setupTestDB()

	// Create test data
	users := []TestUser{
		{Email: "user1@test.com"},
		{Email: "user2@test.com"},
		{Email: "user3@test.com"},
	}

	for _, user := range users {
		adapter.Create(&user)
	}

	// Act - check method chaining
	var foundUser TestUser
	result := adapter.Where("email = ?", "user2@test.com").First(&foundUser)

	// Assert
	suite.Require().NotNil(result)
	suite.Require().IsType(&repositories.GormAdapter{}, result)
	suite.Require().NoError(result.GetError())
	suite.Equal("user2@test.com", foundUser.Email)
}

// Run tests
func TestGormAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(GormAdapterTestSuite))
}
