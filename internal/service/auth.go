package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/project/kocokan/internal/model"
	"github.com/project/kocokan/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo, jwtSecret}
}

func (s *AuthService) Register(name, email, password string) (*model.User, error) {
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return nil, errors.New("email sudah terdaftar")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &model.User{Name: name, Email: email, PasswordHash: string(hash)}
	return u, s.userRepo.Create(u)
}

func (s *AuthService) Login(email, password string) (string, *model.User, error) {
	u, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("email atau password salah")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("email atau password salah")
	}
	token, err := s.generateToken(u)
	return token, u, err
}

func (s *AuthService) ValidateToken(tokenStr string) (*model.User, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("signing method invalid")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token tidak valid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("claims tidak valid")
	}
	id := uint(claims["sub"].(float64))
	return s.userRepo.FindByID(id)
}

func (s *AuthService) generateToken(u *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   u.ID,
		"name":  u.Name,
		"email": u.Email,
		"exp":   time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtSecret))
}
