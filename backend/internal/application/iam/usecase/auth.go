package usecase

import (
	"context"
	"log/slog"
	"time"

	auditmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/audit/model"
	auditrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/audit/repository"
	iammodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/iam/model"
	iamrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/iam/repository"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/iam/service"
	apperrors "github.com/ruveydagundogan/cvmatcher/backend/internal/shared/errors"
)

type RegisterUseCase struct {
	userRepo    iamrepo.UserRepository
	roleRepo    iamrepo.RoleRepository
	authService service.AuthService
	jwtService  service.JWTService
	auditRepo   auditrepo.AuditRepository
	logger      *slog.Logger
}

func NewRegisterUseCase(
	userRepo iamrepo.UserRepository,
	roleRepo iamrepo.RoleRepository,
	authService service.AuthService,
	jwtService service.JWTService,
	auditRepo auditrepo.AuditRepository,
	logger *slog.Logger,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		authService: authService,
		jwtService:  jwtService,
		auditRepo:   auditRepo,
		logger:      logger,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, email, password, firstName, lastName, role string) (string, *iammodel.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Only public roles can self-register
	if role == "" {
		role = "user"
	}
	if role != "user" && role != "hr" {
		return "", nil, apperrors.Validation("invalid role")
	}

	existing, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		uc.logger.Warn("failed to check existing email", "error", err)
	}
	if existing != nil {
		return "", nil, apperrors.AlreadyExists("email already registered")
	}

	hashedPassword, err := uc.authService.HashPassword(password)
	if err != nil {
		return "", nil, apperrors.Internal("failed to hash password", err)
	}

	user := iammodel.NewUser(email, hashedPassword, firstName, lastName)

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return "", nil, apperrors.Internal("failed to create user", err)
	}

	defaultRole, err := uc.roleRepo.FindByName(ctx, role)
	if err == nil && defaultRole != nil {
		if err := uc.roleRepo.AssignToUser(ctx, user.ID, defaultRole.ID); err != nil {
			uc.logger.Warn("failed to assign role", "user_id", user.ID, "error", err)
		}
	}

	token, err := uc.jwtService.GenerateToken(user.ID, user.Email, role)
	if err != nil {
		return "", nil, apperrors.Internal("failed to generate token", err)
	}

	if err := uc.auditRepo.Save(ctx, auditmodel.NewAuditLog(
		user.ID, "register", "user", user.ID,
		"", "", map[string]string{"email": email, "role": role},
	)); err != nil {
		uc.logger.Warn("failed to save audit log", "error", err)
	}

	uc.logger.Info("user registered", "user_id", user.ID, "email", email, "role", role)

	return token, user, nil
}

type LoginUseCase struct {
	userRepo   iamrepo.UserRepository
	roleRepo   iamrepo.RoleRepository
	jwtService service.JWTService
	authService service.AuthService
	auditRepo  auditrepo.AuditRepository
	logger     *slog.Logger
}

func NewLoginUseCase(
	userRepo iamrepo.UserRepository,
	roleRepo iamrepo.RoleRepository,
	jwtService service.JWTService,
	authService service.AuthService,
	auditRepo auditrepo.AuditRepository,
	logger *slog.Logger,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		jwtService:  jwtService,
		authService: authService,
		auditRepo:   auditRepo,
		logger:      logger,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (string, *iammodel.User, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", nil, "", apperrors.Unauthorized("invalid email or password")
	}

	if !uc.authService.ComparePassword(user.PasswordHash, password) {
		return "", nil, "", apperrors.Unauthorized("invalid email or password")
	}

	role := "user"
	userRole, err := uc.roleRepo.GetUserRole(ctx, user.ID)
	if err == nil && userRole != nil {
		role = userRole.RoleName
	}

	token, err := uc.jwtService.GenerateToken(user.ID, user.Email, role)
	if err != nil {
		return "", nil, "", apperrors.Internal("failed to generate token", err)
	}

	if err := uc.auditRepo.Save(ctx, auditmodel.NewAuditLog(
		user.ID, "login", "user", user.ID,
		"", "", nil,
	)); err != nil {
		uc.logger.Warn("failed to save audit log", "error", err)
	}

	uc.logger.Info("user logged in", "user_id", user.ID)

	return token, user, role, nil
}

type GetProfileUseCase struct {
	userRepo iamrepo.UserRepository
	logger   *slog.Logger
}

func NewGetProfileUseCase(userRepo iamrepo.UserRepository, logger *slog.Logger) *GetProfileUseCase {
	return &GetProfileUseCase{userRepo: userRepo, logger: logger}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, userID string) (*iammodel.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, apperrors.NotFound("user not found")
	}
	return user, nil
}

type UpdateProfileUseCase struct {
	userRepo iamrepo.UserRepository
	logger   *slog.Logger
}

func NewUpdateProfileUseCase(userRepo iamrepo.UserRepository, logger *slog.Logger) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{userRepo: userRepo, logger: logger}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID, firstName, lastName string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return apperrors.NotFound("user not found")
	}

	user.FirstName = firstName
	user.LastName = lastName

	return uc.userRepo.Update(ctx, user)
}
