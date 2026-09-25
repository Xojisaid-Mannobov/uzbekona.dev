package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"uzbekona.dev/backend/internal/apperr"
	"uzbekona.dev/backend/internal/config"
	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/repository"
)

type AuthService struct {
	cfg    *config.Config
	admins *repository.AdminRepo
}

// Email topilmaganda ham bcrypt ishlashi uchun (timing attack'dan himoya)
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("uzbekona-dummy-password"), bcrypt.DefaultCost)

type claims struct {
	// PasswordVersion — parol o'zgarsa eski tokenlar bekor bo'ladi
	PasswordVersion string `json:"pv"`
	jwt.RegisteredClaims
}

// Login email/parolni tekshiradi va imzolangan JWT qaytaradi.
func (s *AuthService) Login(ctx context.Context, in *model.LoginInput) (string, *model.Admin, error) {
	invalid := apperr.New(401, "invalid_credentials", "Email yoki parol noto‘g‘ri")

	admin, err := s.admins.GetByEmail(ctx, strings.TrimSpace(in.Email))
	if errors.Is(err, repository.ErrNotFound) {
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(in.Password))
		return "", nil, invalid
	}
	if err != nil {
		return "", nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(in.Password)) != nil {
		return "", nil, invalid
	}

	token, err := s.issueToken(admin)
	if err != nil {
		return "", nil, err
	}
	logErr("last_login yangilanmadi", s.admins.TouchLogin(ctx, admin.ID))
	return token, admin, nil
}

func (s *AuthService) issueToken(admin *model.Admin) (string, error) {
	now := time.Now()
	c := claims{
		PasswordVersion: passwordVersion(admin.PasswordHash),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(admin.ID, 10),
			Issuer:    "uzbekona.dev",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.JWTTTL)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(s.cfg.JWTSecret))
}

// Authenticate tokenni tekshiradi va admin hali mavjudligini, parol o'zgarmaganini tasdiqlaydi.
func (s *AuthService) Authenticate(ctx context.Context, token string) (*model.Admin, error) {
	var c claims
	_, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (any, error) {
		return []byte(s.cfg.JWTSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer("uzbekona.dev"))
	if err != nil {
		return nil, apperr.Unauthorized("Sessiya muddati tugagan. Qaytadan kiring")
	}

	id, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		return nil, apperr.Unauthorized("")
	}
	admin, err := s.admins.GetByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, apperr.Unauthorized("")
	}
	if err != nil {
		return nil, err
	}
	if passwordVersion(admin.PasswordHash) != c.PasswordVersion {
		return nil, apperr.Unauthorized("Parol o‘zgargan. Qaytadan kiring")
	}
	return admin, nil
}

// ChangePassword joriy parolni tekshirib yangisini o'rnatadi va yangi token qaytaradi.
func (s *AuthService) ChangePassword(ctx context.Context, admin *model.Admin, in *model.PasswordInput) (string, error) {
	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(in.CurrentPassword)) != nil {
		return "", apperr.Validation(map[string]string{"current_password": "Joriy parol noto‘g‘ri"})
	}
	hash, err := hashPassword(in.NewPassword)
	if err != nil {
		return "", err
	}
	if err := s.admins.UpdatePassword(ctx, admin.ID, hash); err != nil {
		return "", err
	}
	admin.PasswordHash = hash
	return s.issueToken(admin)
}

func (s *AuthService) ListAdmins(ctx context.Context) ([]model.Admin, error) {
	return s.admins.List(ctx)
}

func (s *AuthService) CreateAdmin(ctx context.Context, in *model.AdminInput) (*model.Admin, error) {
	hash, err := hashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	id, err := s.admins.Create(ctx, strings.TrimSpace(in.Name), strings.ToLower(strings.TrimSpace(in.Email)), hash)
	if repository.IsUniqueViolation(err) {
		return nil, apperr.FieldConflict("email", "Bu email bilan admin mavjud")
	}
	if err != nil {
		return nil, err
	}
	return s.admins.GetByID(ctx, id)
}

func (s *AuthService) DeleteAdmin(ctx context.Context, current *model.Admin, id int64) error {
	if current.ID == id {
		return apperr.BadRequest("O‘zingizni o‘chira olmaysiz")
	}
	return mapErr(s.admins.Delete(ctx, id), "Admin topilmadi")
}

// EnsureAdmin — adminlar jadvali bo'sh bo'lsa, env'dagi ma'lumotlar bilan birinchi adminni yaratadi.
func (s *AuthService) EnsureAdmin(ctx context.Context) error {
	count, err := s.admins.Count(ctx)
	if err != nil || count > 0 {
		return err
	}
	if s.cfg.AdminEmail == "" || s.cfg.AdminPassword == "" {
		slog.Warn("admin yo‘q: ADMIN_EMAIL va ADMIN_PASSWORD berilsa birinchi admin yaratiladi")
		return nil
	}
	if len(s.cfg.AdminPassword) < 10 {
		return errors.New("ADMIN_PASSWORD kamida 10 belgidan iborat bo‘lsin")
	}
	hash, err := hashPassword(s.cfg.AdminPassword)
	if err != nil {
		return err
	}
	if _, err := s.admins.Create(ctx, s.cfg.AdminName, strings.ToLower(s.cfg.AdminEmail), hash); err != nil {
		return err
	}
	slog.Info("birinchi admin yaratildi", "email", s.cfg.AdminEmail)
	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}

func passwordVersion(hash string) string {
	sum := sha256.Sum256([]byte(hash))
	return hex.EncodeToString(sum[:6])
}
