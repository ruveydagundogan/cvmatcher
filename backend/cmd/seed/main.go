package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/auth"
	iammodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/iam/model"
	pgresiam "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/postgres/iam"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/config"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/database"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/memory"
	iamrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/iam/repository"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := database.NewPostgresPool(ctx, cfg.Database, slog.Default())
	if err != nil {
		fmt.Fprintf(os.Stderr, "database connection failed: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	bcrypt := auth.NewBcryptAuthService(10)
	var userRepo iamrepo.UserRepository
	var roleRepo iamrepo.RoleRepository
	userRepo = pgresiam.NewUserRepository(pool)
	roleRepo = pgresiam.NewRoleRepository(pool)

	_ = memory.NewInMemoryUserRepo()

	users := []struct {
		email     string
		firstName string
		lastName  string
		password  string
		role      string
	}{
		{"admin@cvmatcher.com", "Admin", "User", "admin123", "admin"},
		{"user@cvmatcher.com", "Regular", "User", "user123", "user"},
	}

	for _, u := range users {
		existing, err := userRepo.FindByEmail(ctx, u.email)
		if err == nil && existing != nil {
			fmt.Printf("user %s already exists (id=%s), skipping\n", u.email, existing.ID)
			continue
		}

		hash, err := bcrypt.HashPassword(u.password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "hash password for %s: %v\n", u.email, err)
			continue
		}

		user := iammodel.NewUser(u.email, hash, u.firstName, u.lastName)
		if err := userRepo.Create(ctx, user); err != nil {
			fmt.Fprintf(os.Stderr, "create user %s: %v\n", u.email, err)
			continue
		}

		role, err := roleRepo.FindByName(ctx, u.role)
		if err != nil {
			fmt.Fprintf(os.Stderr, "find role %s: %v\n", u.role, err)
		} else {
			_ = roleRepo.AssignToUser(ctx, user.ID, role.ID)
		}

		fmt.Printf("created %s user: %s / %s\n", u.role, u.email, u.password)
	}

	fmt.Println("seed complete")
}
