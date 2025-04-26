package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/domain/entities"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/domain/repositories"
)

// UserService error definitions
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserTaskType definitions
const (
	UserTaskTypeWelcomeEmail = "welcome_email"
	UserTaskTypeProfileUpdate = "profile_update"
)

// UserService implements user-related use cases
type UserService struct {
	userRepo repositories.UserRepository
	taskRepo repositories.TaskRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repositories.UserRepository, taskRepo repositories.TaskRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		taskRepo: taskRepo,
	}
}

// CreateUser creates a new user and queues a welcome email task
func (s *UserService) CreateUser(
	ctx context.Context, 
	username, 
	email, 
	password, 
	firstName, 
	lastName string,
) (*entities.User, error) {
	// Check if user with email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Check if user with username already exists
	existingUser, err = s.userRepo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Create new user
	user, err := entities.NewUser(username, email, password, firstName, lastName)
	if err != nil {
		return nil, err
	}

	// Save user to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Queue welcome email task
	if err := s.queueWelcomeEmailTask(ctx, user); err != nil {
		// Log error but continue - don't fail user creation if task queuing fails
		fmt.Printf("Failed to queue welcome email task: %v\n", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetUserByUsername retrieves a user by their username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(
	ctx context.Context,
	id string,
	username string,
	email string,
	firstName string,
	lastName string,
) (*entities.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// If email is changing, check if it's already taken
	if email != "" && email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, ErrUserAlreadyExists
		}
	}

	// If username is changing, check if it's already taken
	if username != "" && username != user.Username {
		existingUser, err := s.userRepo.GetByUsername(ctx, username)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, ErrUserAlreadyExists
		}
	}

	if err := user.Update(username, email, firstName, lastName); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Queue profile update notification task
	if err := s.queueProfileUpdateTask(ctx, user); err != nil {
		// Log error but continue - don't fail update if task queuing fails
		fmt.Printf("Failed to queue profile update task: %v\n", err)
	}

	return user, nil
}

// ListUsers retrieves a list of users with pagination
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*entities.User, int, error) {
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.userRepo.Delete(ctx, id)
}

// UpdatePassword updates a user's password
func (s *UserService) UpdatePassword(ctx context.Context, id string, password string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := user.UpdatePassword(password); err != nil {
		return err
	}

	return s.userRepo.Update(ctx, user)
}

// queueWelcomeEmailTask creates and queues a task to send welcome email
func (s *UserService) queueWelcomeEmailTask(ctx context.Context, user *entities.User) error {
	// Create task data with user information needed for email
	type welcomeEmailData struct {
		UserID    string `json:"user_id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	data := welcomeEmailData{
		UserID:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create task
	task, err := entities.NewTask(UserTaskTypeWelcomeEmail, jsonData)
	if err != nil {
		return err
	}

	// Save task to repository
	return s.taskRepo.Create(ctx, task)
}

// queueProfileUpdateTask creates and queues a task to process profile update
func (s *UserService) queueProfileUpdateTask(ctx context.Context, user *entities.User) error {
	// Create task data with user information
	type profileUpdateData struct {
		UserID    string `json:"user_id"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	data := profileUpdateData{
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create task
	task, err := entities.NewTask(UserTaskTypeProfileUpdate, jsonData)
	if err != nil {
		return err
	}

	// Save task to repository
	return s.taskRepo.Create(ctx, task)
}
