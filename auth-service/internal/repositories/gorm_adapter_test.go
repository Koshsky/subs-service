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
func (suite *GormAdapterTestSuite) setupTestDB() (*gorm.DB, repositories.DatabaseInterface) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.NoError(err)

	adapter := repositories.NewGormAdapter(db)

	err = db.AutoMigrate(&TestUser{})
	suite.NoError(err)

	return db, adapter
}

func (suite *GormAdapterTestSuite) TestNewGormAdapter() {
	// Arrange & Act
	adapter := repositories.NewGormAdapter(nil)

	// Assert
	suite.NotNil(adapter)
	suite.IsType(&repositories.GormAdapter{}, adapter)
}

func (suite *GormAdapterTestSuite) TestAdapterMethodsWithNilDB() {
	// Arrange
	adapter := repositories.NewGormAdapter(nil)
	testCases := []struct {
		name   string
		method func() repositories.DatabaseInterface
	}{
		{
			name: "Create",
			method: func() repositories.DatabaseInterface {
				return adapter.Create(&struct{}{})
			},
		},
		{
			name: "Where",
			method: func() repositories.DatabaseInterface {
				return adapter.Where("email = ?", "test@example.com")
			},
		},
		{
			name: "First",
			method: func() repositories.DatabaseInterface {
				return adapter.First(&struct{}{})
			},
		},
	}

	// Act & Assert
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			result := tc.method()

			suite.NotNil(result)
			suite.IsType(&repositories.GormAdapter{}, result)
			suite.Error(result.GetError())
			suite.Contains(result.GetError().Error(), "database is nil")
		})
	}
}
func (suite *GormAdapterTestSuite) TestCreateWithRealDB() {
	// Arrange
	_, adapter := suite.setupTestDB()
	user := &TestUser{Email: "test@example.com"}

	// Act
	createResult := adapter.Create(user)

	// Assert
	suite.NotNil(createResult)
	suite.IsType(&repositories.GormAdapter{}, createResult)
	suite.NoError(createResult.GetError())
}

func (suite *GormAdapterTestSuite) TestWhereWithRealDB() {
	// Arrange
	_, adapter := suite.setupTestDB()

	// Act
	whereResult := adapter.Where("email = ?", "test@example.com")

	// Assert
	suite.NotNil(whereResult)
	suite.IsType(&repositories.GormAdapter{}, whereResult)
	suite.NoError(whereResult.GetError())
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
	suite.NotNil(firstResult)
	suite.IsType(&repositories.GormAdapter{}, firstResult)
	suite.NoError(firstResult.GetError())
	suite.Equal("test@example.com", foundUser.Email)
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
	suite.NotNil(result)
	suite.IsType(&repositories.GormAdapter{}, result)
	suite.NoError(result.GetError())
	suite.Equal("user2@test.com", foundUser.Email)
}

// Run tests
func TestGormAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(GormAdapterTestSuite))
}
