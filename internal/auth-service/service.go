package auth_service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"

	"github.com/Miroshinsv/disko_go/internal/auth-service/models"
	roleService "github.com/Miroshinsv/disko_go/internal/role-service"
	userService "github.com/Miroshinsv/disko_go/internal/user-service"
	config_service "github.com/Miroshinsv/disko_go/pkg/config-service"
	dbConnector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
)

const (
	jwtConstString     = "lololo"
	JWTAuthAudience    = "auth"
	JWTRefreshAudience = "refresh"
	vkApiURL           = "https://api.vk.com/method/%s"
	vkApiVersion       = "5.124"
	vkApiSelfMethod    = "account.getProfileInfo"
	vkDefPassword      = "!123321123321@"
)

var (
	self *Service = nil

	errorInvalidToken = errors.New("invalid token")
)

type Service struct {
	log    loggerService.ILogger
	conn   dbConnector.IConnector
	vkConf *oauth2.Config
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

func (h Service) LoginSocial(token string) (*userService.Users, error) {
	ctx := context.Background()
	h.vkConf.RedirectURL = "http://mighty-beach-02870.herokuapp.com/auth/vk/"
	vkToken, err := h.vkConf.Exchange(ctx, token)
	if err != nil {
		fmt.Println("Error" + err.Error())
		return (*userService.Users)(nil), err
	}
	var params = make(url.Values)
	params.Add("access_token", vkToken.AccessToken)
	params.Add("v", vkApiVersion)

	client := h.vkConf.Client(ctx, vkToken)
	resp, err := client.PostForm(fmt.Sprintf(vkApiURL, vkApiSelfMethod), params)
	if err != nil {
		return (*userService.Users)(nil), err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return (*userService.Users)(nil), err
	}

	var vkResp models.VKResponse
	err = json.Unmarshal(body, &vkResp)
	if err != nil {
		return (*userService.Users)(nil), err
	}

	var existing = &userService.Users{}
	fmt.Printf("VK RESPONSE ID: %d", vkResp.Response.ID)
	h.conn.GetConnection().Where(fmt.Sprintf("email LIKE '%d@vk.com'", vkResp.Response.ID)).Find(existing)

	if existing.ID != 0 {
		return existing, nil
	}

	// @todo: roles?
	var dbUser = &userService.Users{
		Model:     gorm.Model{},
		FirstName: vkResp.Response.FirstName,
		SureName:  vkResp.Response.LastName,
		Email:     fmt.Sprintf("email LIKE '%d@vk.com'", vkResp.Response.ID),
		Password:  fmt.Sprintf("%x", md5.Sum([]byte(vkDefPassword))),
	}

	db := h.conn.GetConnection().Create(dbUser)
	return dbUser, db.Error
}

func MustNewAuthService(log loggerService.ILogger, confService config_service.IConfig, conn dbConnector.IConnector) *Service {
	if self == nil {
		var confAuth Config
		err := confService.Convert(&confAuth)
		if err != nil {
			log.Fatal("error on initializing auth service", err, nil)
		}

		conf := &oauth2.Config{
			ClientID:     confAuth.VKClientID,
			ClientSecret: confAuth.VKClientSecret,
			Scopes:       []string{},
			Endpoint:     vk.Endpoint,
		}

		self = &Service{
			log:    log,
			conn:   conn,
			vkConf: conf,
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
