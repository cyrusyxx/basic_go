package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/webook/internal/domain"
	"webook/webook/internal/service"
)

// Email and password regexp pattern
const (
	emailRegexPattern = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	//emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,70}$`
)

// UserHandler Struct
type UserHandler struct {
	emailRegExp   *regexp.Regexp
	passwordRegex *regexp.Regexp
	svc           *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegExp:   regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:           svc,
	}
}

// Register Router
func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.Login)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)
}

// Sign up
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
			MaxAge:   900,
			HttpOnly: true,
		})
		err := sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "System Error!!")
			return
		}
		ctx.String(http.StatusOK, "Login Seccess!!")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "User not found or password is wrong")
	default:
		ctx.String(http.StatusOK, "System Error!!")
	}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello form profile")
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello form edit")
}
