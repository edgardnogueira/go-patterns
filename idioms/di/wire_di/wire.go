//+build wireinject

package wire_di

import (
	"github.com/google/wire"
)

// ProviderSet defines the provider set for the application
var ProviderSet = wire.NewSet(
	// Basic dependencies
	NewConfig,
	NewLogger,
	
	// Database connection
	NewDatabaseConnection,
	
	// API client
	NewAPIClient,
	
	// Repositories
	NewUserRepository,
	NewMessageRepository,
	
	// Define interfaces and their implementations
	wire.Bind(new(UserRepository), new(*DBUserRepository)),
	wire.Bind(new(MessageRepository), new(*DBMessageRepository)),
	
	// Services
	NewNotificationService,
	wire.Bind(new(NotificationService), new(*EmailNotificationService)),
	NewUserService,
	NewMessageService,
	
	// Application
	NewApplication,
)

// InitializeApplication initializes the entire application with its dependencies
// This function signature is used by Wire to generate the dependency injection code
func InitializeApplication() (*Application, error) {
	wire.Build(ProviderSet)
	return &Application{}, nil
}
