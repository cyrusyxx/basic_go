package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Email and password regexp pattern
const (
	emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

// UserHandler Struct
type UserHandler struct {
	emailRegExp   *regexp.Regexp
	passwordRegex *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		emailRegExp:   regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
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
	}
	// Handle the keyword pattern
	isPassword, err := h.passwordRegex.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "Password input failed\n")
	}
	if !isPassword {
		ctx.String(http.StatusOK, "Password pattern is wrong\n")
	}

	// Return the string to browser
	ctx.String(http.StatusOK, "HELLO ITS IN SIGNUP")
}

func (h *UserHandler) Profile(ctx *gin.Context) {

}

func (h *UserHandler) Login(ctx *gin.Context) {

}

func (h *UserHandler) Edit(ctx *gin.Context) {

}
