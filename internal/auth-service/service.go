package auth_service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"

	"github.com/Miroshinsv/disko_go/internal/auth-service/models"
	roleService "github.com/Miroshinsv/disko_go/internal/role-service"
	userService "github.com/Miroshinsv/disko_go/internal/user-service"
	dbConnector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
)

const (
	jwtConstString     = "lololo"
	JWTAuthAudience    = "auth"
	JWTRefreshAudience = "refresh"
)

var (
	self *Service = nil

	errorInvalidToken = errors.New("invalid token")
)

type Service struct {
	log  loggerService.ILogger
	conn dbConnector.IConnector
}

func (h Service) RegisterUser(u models.User) (*userService.Users, error) {
	if u.Email == nil || !emailRegex.MatchString(*u.Email) {
		return &userService.Users{}, errorInvalidEmail
	}

	if u.Password == nil || *u.Password == "" {
		return &userService.Users{}, errorInvalidPassword
	}

	if u.Role == nil {
		return &userService.Users{}, errorUnknownRole
	}

	var role = &roleService.Roles{}
	h.conn.GetConnection().Where(fmt.Sprintf("id=%d", *u.Role)).Find(&role)
	if role == nil || role.ID == 0 {
		return &userService.Users{}, errorUnknownRole
	}

	var existing = &userService.Users{}
	h.conn.GetConnection().Where(fmt.Sprintf("email LIKE '%s'", *u.Email)).Find(existing)
	if existing.ID != 0 {
		return existing, errorExistingUser
	}

	var dbUser = &userService.Users{
		Model:      gorm.Model{},
		FirstName:  "",
		SureName:   "",
		MiddleName: "",
		Email:      *u.Email,
		Phone:      "",
		Password:   fmt.Sprintf("%x", md5.Sum([]byte(*u.Password))),
		Roles:      []*roleService.Roles{role},
	}

	db := h.conn.GetConnection().Create(dbUser)

	return dbUser, db.Error
}

func (h Service) LoginUser(u models.User) (*userService.Users, error) {
	if u.Email == nil || !emailRegex.MatchString(*u.Email) {
		return &userService.Users{}, errorInvalidEmail
	}

	if u.Password == nil || *u.Password == "" {
		return &userService.Users{}, errorInvalidPassword
	}

	var existing = &userService.Users{}
	h.conn.GetConnection().Where(fmt.Sprintf("email LIKE '%s'", *u.Email)).Find(existing)
	if existing.ID == 0 {
		return existing, errorInvalidEmailOrPassword
	}

	if existing.Password != fmt.Sprintf("%x", md5.Sum([]byte(*u.Password))) {
		return &userService.Users{}, errorInvalidEmailOrPassword
	}

	return existing, nil
}

func (h Service) GenerateAuthJWT(dbUser *userService.Users) string {
	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Issuer:    strconv.Itoa(int(dbUser.ID)),
		Audience:  JWTAuthAudience,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtConstString))
	if err != nil {
		panic(err)
	}

	return ss
}

func (h Service) GenerateRefreshJWT(dbUser *userService.Users) string {
	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(48 * time.Hour).Unix(),
		Issuer:    strconv.Itoa(int(dbUser.ID)),
		Audience:  JWTRefreshAudience,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtConstString))
	if err != nil {
		panic(err)
	}

	return ss
}

func (h Service) GetUserByJWT(jwtToken string, jwtType string) (*userService.Users, error) {
	var (
		dbUser = &userService.Users{}
		claims = &jwt.StandardClaims{}
	)

	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConstString), nil
	})
	if err != nil || !token.Valid {
		return dbUser, errorInvalidToken
	}

	if claims.Audience != jwtType {
		return dbUser, errorInvalidToken
	}

	h.conn.GetConnection().Where(fmt.Sprintf("id = %s", claims.Issuer)).Find(dbUser)
	if dbUser.ID == 0 {
		return dbUser, errorInvalidToken
	}

	return dbUser, nil
}

func (h Service) UpdateTokens(jwtRefreshToken string) (map[string]string, error) {
	dbUser, err := h.GetUserByJWT(jwtRefreshToken, JWTRefreshAudience)
	if err != nil {
		return make(map[string]string, 0), err
	}

	return map[string]string{
		"auth":    h.GenerateAuthJWT(dbUser),
		"refresh": h.GenerateRefreshJWT(dbUser),
	}, nil
}

func MustNewAuthService(log loggerService.ILogger, conn dbConnector.IConnector) *Service {
	if self == nil {
		self = &Service{
			log:  log,
			conn: conn,
		}
	}

	return self
}

func GetAuthService() *Service {
	if self == nil {
		panic("service not init")
	}

	return self
}
