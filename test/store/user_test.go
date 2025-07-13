package store_test

import (
	"context"
	"testing"

	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/store"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
}

var (
	ctx   = context.Background()
	users = []*model.User{
		{
			Name:     "test_user",
			Email:    "test_user@example.com",
			Password: "test_password",
		},
		{
			Name:     "test_user2",
			Email:    "test_user2@example.com",
			Password: "test_password2",
		},
		{
			Name:     "test_user3",
			Email:    "test_user3@example.com",
			Password: "test_password3",
		},
	}
)

func TestCreateUser(t *testing.T) {
	if err := userRepo.CreateBatch(ctx, users); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
}

func TestQueryUser(t *testing.T) {
	user, err := userRepo.Query(ctx, store.Where("name", "test_user"))
	if err != nil {
		t.Fatalf("failed to query user: %v", err)
	}
	zap.L().Info("query user", zap.Any("user", user))
}

func TestUpdateUser(t *testing.T) {
	// user, err := userRepo.Query(ctx, store.Where("name", "test_user"))
	// if err != nil {
	// 	t.Fatalf("failed to query user: %v", err)
	// }
	// user.Name = "test_user_update"
	user := &model.User{
		Name: "zero",
	}
	if err := userRepo.Update(ctx, user, store.Where("name", "test_user_update")); err != nil {
		t.Fatalf("failed to update user: %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	user, err := userRepo.Query(ctx, store.Where("name", "zero"))
	if err != nil {
		t.Fatalf("failed to query user: %v", err)
	}
	if err := userRepo.Delete(ctx, user); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}
}

func TestTX(t *testing.T) {
	if err := txManager.Transaction(ctx, func(ctx context.Context) error {
		if err := userRepo.Create(ctx, &model.User{
			Name:     "test_user_tx",
			Email:    "test_user_tx@example.com",
			Password: "test_password_tx",
		}); err != nil {
			return err
		}

		if _, err := userRepo.Query(ctx, store.Where("name", "notFound")); err != nil {
			return err
		}

		return nil
	}); err != nil {
		t.Fatalf("failed to transaction: %v", err)
	}
}

func TestCreateTable(t *testing.T) {
	if err := db.AutoMigrate(&Test{}); err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
}
