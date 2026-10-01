package service

import (
	"context"
	"errors"
	"go-notification/internal/apperror"
	"go-notification/internal/model"
	"go-notification/internal/repository"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeNotificationRepository struct {
	notification *model.Notification
	err          error

	receivedID     string
	receivedUserID string
}

func (f *fakeNotificationRepository) Create(
	ctx context.Context,
	notification *model.Notification,
) error {
	return nil
}

func (f *fakeNotificationRepository) GetAll(
	ctx context.Context,
	params repository.GetNotificationsParams,
) ([]model.Notification, int64, error) {
	return nil, 0, nil
}

func (f *fakeNotificationRepository) Update(
	ctx context.Context,
	id string,
	notification *model.Notification,
	userID string,
) (*model.Notification, error) {
	return nil, nil
}

func (f *fakeNotificationRepository) Delete(
	ctx context.Context,
	id string,
	userID string,
) error {
	return nil
}

func (f *fakeNotificationRepository) GetById(
	ctx context.Context,
	id string,
	userID string,
) (*model.Notification, error) {
	f.receivedID = id
	f.receivedUserID = userID

	return f.notification, f.err
}

func TestNotificationService_Get(t *testing.T) {
	notification := &model.Notification{
		ID:        bson.NewObjectID(),
		UserID:    "user-123",
		Title:     "Test Notification",
		Message:   "Hello from test",
		CreatedAt: time.Now(),
	}

	dbErr := errors.New("database connection failed")

	tests := []struct {
		name string

		// Fake repository result
		repoNotification *model.Notification
		repoErr          error

		// Input to service
		userID string
		id     string

		// Expected result
		expectedTitle string
		expectedErr   error
	}{
		{
			name:             "success",
			repoNotification: notification,
			repoErr:          nil,
			userID:           "user-123",
			id:               notification.ID.Hex(),
			expectedTitle:    "Test Notification",
			expectedErr:      nil,
		},
		{
			name:             "notification not found",
			repoNotification: nil,
			repoErr:          apperror.ErrNotificationNotFound,
			userID:           "user-123",
			id:               "notification-id",
			expectedTitle:    "",
			expectedErr:      apperror.ErrNotificationNotFound,
		},
		{
			name:             "repository error",
			repoNotification: nil,
			repoErr:          dbErr,
			userID:           "user-123",
			id:               "notification-id",
			expectedTitle:    "",
			expectedErr:      dbErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			fakeRepo := &fakeNotificationRepository{
				notification: tt.repoNotification,
				err:          tt.repoErr,
			}

			notificationService := NewNotificationService(fakeRepo)

			// Act
			result, err := notificationService.Get(
				context.Background(),
				tt.id,
				tt.userID,
			)

			// Assert: repository received correct arguments
			if fakeRepo.receivedUserID != tt.userID {
				t.Errorf(
					"expected repository userID %q, got %q",
					tt.userID,
					fakeRepo.receivedUserID,
				)
			}

			if fakeRepo.receivedID != tt.id {
				t.Errorf(
					"expected repository id %q, got %q",
					tt.id,
					fakeRepo.receivedID,
				)
			}

			// Assert: expected error
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf(
					"expected error %v, got %v",
					tt.expectedErr,
					err,
				)
			}

			if tt.expectedErr != nil {
				if result != nil {
					t.Errorf(
						"expected result to be nil, got %+v",
						result,
					)
				}

				return
			}

			// Assert: expected result
			if result == nil {
				t.Fatal("expected result, got nil")
			}

			if result.Title != tt.expectedTitle {
				t.Errorf(
					"expected title %q, got %q",
					tt.expectedTitle,
					result.Title,
				)
			}
		})
	}
}
