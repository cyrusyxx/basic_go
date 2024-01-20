package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"webook/webook/internal/domain"
	"webook/webook/internal/service"
)

// Email and password regexp pattern
const (
	emailRegexPattern = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	//emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,70}$`
)

var SigKey = []byte("ukRIDSD0JpWD5Qv0P46Y8IGLjB2uvShj")

// UserHandler Struct
type UserHandler struct {
	emailRegExp   *regexp.Regexp
	passwordRegex *regexp.Regexp
	svc           *service.UserService
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int64
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegExp:   regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:           svc,
	}
}

// RegisterRoutes Register Router
func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	//ug.POST("/login", h.Login)
	ug.POST("/login", h.LoginJWT)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)
}

// SignUp Sign up
func (h *UserHandler) SignUp(ctx *gin.Context) {
	// verify
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignupReq

	// Get the context content
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// Handle the email pattern
	isEmail, err := h.emailRegExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "Email input failed\n")
	}
	if !isEmail {
		ctx.String(http.StatusOK, "Email pattern is wrong\n")
		return
	}
	// Handle the keyword pattern
	isPassword, err := h.passwordRegex.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "Password input failed\n")
	}
	if !isPassword {
		ctx.String(http.StatusOK, "Password pattern is wrong\n")
		return
	}

	// real signup
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		// Success:Return the string to browser
		ctx.String(http.StatusOK, "HELLO ITS IN SIGNUP")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "Email Duplicate !!")
	default:
		ctx.String(http.StatusOK, "System Error !!")
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type SignupReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req SignupReq

	// Get the context content
	if err := ctx.Bind(&req); err != nil {
		return
	}

	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			MaxAge:   300,
			HttpOnly: true,
		})
		err := sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "System Error!!")
			return
		}
		ctx.String(http.StatusOK, "Login Success!!")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "User not found or password is wrong")
	default:
		ctx.String(http.StatusOK, "System Error!!")
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type SignupReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req SignupReq

	// Get the context content
	if err := ctx.Bind(&req); err != nil {
		return
	}

	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		uc := UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "",
				Subject:   "",
				Audience:  nil,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
				NotBefore: nil,
				IssuedAt:  nil,
				ID:        "",
			},
			Uid: u.Id,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
		tokenStr, err := token.SignedString(SigKey)
		if err != nil {
			ctx.String(http.StatusOK, "System Error!!")
			return
		}
		ctx.Header("x-jwt-token", tokenStr)

		ctx.String(http.StatusOK, "Login Success!!")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "User not found or password is wrong")
	default:
		ctx.String(http.StatusOK, "System Error!!")
	}

}

func (h *UserHandler) Edit(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello form edit")
	type EditReq struct {
		NickName    string `json:"nickname"`
		Birthday    string `json:"birthday"`
		Description string `json:"description"`
	}
	var req EditReq

	sess := sessions.Default(ctx)
	uid := sess.Get("userId").(int64)
	//ctx.String(http.StatusOK, "uid is %d", uid)

	if err := ctx.Bind(&req); err != nil {
		return
	}

	if len(req.NickName) > 8 {
		// 处理校验错误
		ctx.String(http.StatusOK, "Invalid nickname length")
		return
	}

	if len(req.Birthday) != 10 {
		// 处理校验错误
		ctx.String(http.StatusOK, "Invalid birthday length")
		return
	}

	if len(req.Description) > 50 {
		// 处理校验错误
		ctx.String(http.StatusOK, "Invalid description length")
		return
	}

	err := h.svc.Edit(ctx, uid, req.NickName, req.Birthday, req.Description)
	if err != nil {
		return
	}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello form profile\n")

	sess := sessions.Default(ctx)
	uid := sess.Get("userId").(int64)

	u, err := h.svc.Profile(ctx, uid)
	if err != nil {
		return
	}

	response := fmt.Sprintf("id=%d\nEmail=%s\nNickName=%s\nBirthday=%s\nDescription=%s\n",
		u.Id, u.Email, u.NickName, u.Birthday, u.Description)

	// 返回用户信息
	ctx.String(http.StatusOK, response)
}
