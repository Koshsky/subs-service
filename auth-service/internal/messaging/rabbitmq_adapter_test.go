package messaging

import (
	"fmt"
	"testing"

	"github.com/Koshsky/subs-service/auth-service/internal/config"
	"github.com/Koshsky/subs-service/auth-service/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	messagingMocks "github.com/Koshsky/subs-service/auth-service/internal/messaging/mocks"
)

type RabbitMQAdapterTestSuite struct {
	suite.Suite
	mockPublisher *messagingMocks.IRabbitMQPublisher
	mockConn      *messagingMocks.IRabbitMQConn
	adapter       IMessageBroker
	config        *config.Config
}

func (suite *RabbitMQAdapterTestSuite) SetupTest() {
	suite.config = &config.Config{
		RabbitMQ: config.RabbitMQConfig{
			Exchange: "test_exchange",
		},
	}
	suite.mockPublisher = messagingMocks.NewIRabbitMQPublisher(suite.T())
	suite.mockConn = messagingMocks.NewIRabbitMQConn(suite.T())
	suite.adapter = &RabbitMQAdapter{
		publisher: suite.mockPublisher,
		conn:      suite.mockConn,
		config:    suite.config,
	}
}

func (suite *RabbitMQAdapterTestSuite) TearDownTest() {
	suite.mockPublisher.AssertExpectations(suite.T())
	suite.mockConn.AssertExpectations(suite.T())
}

// ===== CONSTRUCTOR TESTS =====

func (suite *RabbitMQAdapterTestSuite) TestNewRabbitMQAdapter_InvalidConfig() {
	// Arrange
	cfg := &config.Config{
		RabbitMQ: config.RabbitMQConfig{
			URL:      "invalid://url",
			Exchange: "test_exchange",
		},
	}

	// Act
	adapter, err := NewRabbitMQAdapter(cfg)

	// Assert
	suite.Require().Error(err)
	suite.Require().Nil(adapter)
	suite.Contains(err.Error(), "failed to connect to RabbitMQ")
}

// ===== PUBLISH USER CREATED TESTS =====

func (suite *RabbitMQAdapterTestSuite) TestPublishUserCreated_Success() {
	// Arrange
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}
	suite.mockPublisher.On("Publish",
		[]byte(`{"user_id":"`+user.ID.String()+`","email":"test@example.com"}`),
		[]string{"user.created"},
		mock.AnythingOfType("func(*rabbitmq.PublishOptions)"),
		mock.AnythingOfType("func(*rabbitmq.PublishOptions)"),
	).Return(nil)

	// Act
	err := suite.adapter.PublishUserCreated(user)

	// Assert
	suite.Require().NoError(err)
	suite.mockPublisher.AssertExpectations(suite.T())
}

func (suite *RabbitMQAdapterTestSuite) TestPublishUserCreated_NilPublisher() {
	// Arrange
	adapter := &RabbitMQAdapter{
		publisher: nil,
		conn:      suite.mockConn,
		config:    suite.config,
	}
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	// Act
	err := adapter.PublishUserCreated(user)

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "publisher is not initialized")
}

func (suite *RabbitMQAdapterTestSuite) TestPublishUserCreated_NilConfig() {
	// Arrange
	adapter := &RabbitMQAdapter{
		publisher: suite.mockPublisher,
		conn:      suite.mockConn,
		config:    nil,
	}
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	// Act
	err := adapter.PublishUserCreated(user)

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "config is not initialized")
}

func (suite *RabbitMQAdapterTestSuite) TestPublishUserCreated_PublisherError() {
	// Arrange
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}
	expectedError := fmt.Errorf("publisher error")
	suite.mockPublisher.On("Publish",
		[]byte(`{"user_id":"`+user.ID.String()+`","email":"test@example.com"}`),
		[]string{"user.created"},
		mock.AnythingOfType("func(*rabbitmq.PublishOptions)"),
		mock.AnythingOfType("func(*rabbitmq.PublishOptions)"),
	).Return(expectedError)

	// Act
	err := suite.adapter.PublishUserCreated(user)

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "publisher error")
	suite.mockPublisher.AssertExpectations(suite.T())
}

func (suite *RabbitMQAdapterTestSuite) TestPublishUserCreated_NilUser() {
	// Arrange
	var user *models.User = nil

	// Act
	err := suite.adapter.PublishUserCreated(user)

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "user cannot be nil")
}

// ===== PUBLISH USER DELETED TESTS =====

func (suite *RabbitMQAdapterTestSuite) TestPublishUserDeleted_Success() {
	// Arrange
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}
	suite.mockPublisher.On("Publish",
		[]byte(`{"user_id":"`+user.ID.String()+`"}`),
		[]string{"user.deleted"},
		mock.AnythingOfType("func(*rabbitmq.PublishOptions)"),
		mock.AnythingOfType("func(*rabbitmq.PublishOptions)"),
	).Return(nil)

	// Act
	err := suite.adapter.PublishUserDeleted(user)

	// Assert
	suite.Require().NoError(err)
	suite.mockPublisher.AssertExpectations(suite.T())
}

func (suite *RabbitMQAdapterTestSuite) TestPublishUserDeleted_NilPublisher() {
	// Arrange
	adapter := &RabbitMQAdapter{
		publisher: nil,
		conn:      suite.mockConn,
		config:    suite.config,
	}
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	// Act
	err := adapter.PublishUserDeleted(user)

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "publisher is not initialized")
}

func (suite *RabbitMQAdapterTestSuite) TestPublishUserDeleted_NilConfig() {
	// Arrange
	adapter := &RabbitMQAdapter{
		publisher: suite.mockPublisher,
		conn:      suite.mockConn,
		config:    nil,
	}
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	// Act
	err := adapter.PublishUserDeleted(user)

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "config is not initialized")
}

func (suite *RabbitMQAdapterTestSuite) TestPublishUserDeleted_PublisherError() {
	// Arrange
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}
	expectedError := fmt.Errorf("publisher error")
	suite.mockPublisher.On("Publish",
		[]byte(`{"user_id":"`+user.ID.String()+`"}`),
		[]string{"user.deleted"},
		mock.AnythingOfType("func(*rabbitmq.PublishOptions)"),
		mock.AnythingOfType("func(*rabbitmq.PublishOptions)"),
	).Return(expectedError)

	// Act
	err := suite.adapter.PublishUserDeleted(user)

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "publisher error")
	suite.mockPublisher.AssertExpectations(suite.T())
}

func (suite *RabbitMQAdapterTestSuite) TestPublishUserDeleted_NilUser() {
	// Arrange
	var user *models.User = nil

	// Act
	err := suite.adapter.PublishUserDeleted(user)

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "user cannot be nil")
}

// ===== CLOSE TESTS =====

func (suite *RabbitMQAdapterTestSuite) TestClose_Success() {
	// Arrange
	suite.mockPublisher.On("Close").Return()
	suite.mockConn.On("Close").Return(nil)

	// Act & Assert
	suite.NotPanics(func() {
		suite.adapter.Close()
	})
	suite.mockPublisher.AssertExpectations(suite.T())
	suite.mockConn.AssertExpectations(suite.T())
}

func (suite *RabbitMQAdapterTestSuite) TestClose_MultipleCalls() {
	// Arrange
	suite.mockPublisher.On("Close").Return().Times(2)
	suite.mockConn.On("Close").Return(nil).Times(2)

	// Act & Assert
	suite.NotPanics(func() {
		suite.adapter.Close()
		suite.adapter.Close()
	})
	suite.mockPublisher.AssertExpectations(suite.T())
	suite.mockConn.AssertExpectations(suite.T())
}

func TestRabbitMQAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(RabbitMQAdapterTestSuite))
}
