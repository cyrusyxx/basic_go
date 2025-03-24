package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"webook/webook/internal/domain"
	"webook/webook/internal/service"
	ijwt "webook/webook/internal/web/jwt"
	"webook/webook/pkg/ginx"
)

// Email and password regexp pattern
const (
	emailRegexPattern = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	//emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,70}$`
	bizlogin             = "login"
)

// UserHandler Struct
type UserHandler struct {
	emailRegExp   *regexp.Regexp
	passwordRegex *regexp.Regexp
	usersvc       service.UserService
	codesvc       service.CodeService
	ijwt.Handler
}

func NewUserHandler(svc service.UserService,
	codesvc service.CodeService, jwthdl ijwt.Handler) *UserHandler {
	return &UserHandler{
		emailRegExp:   regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
		usersvc:       svc,
		codesvc:       codesvc,
		Handler:       jwthdl,
	}
}

// RegisterRoutes Register Router
func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/user")
	ug.POST("/signup", h.SignUp)
	//ug.POST("/login", h.Login)
	ug.POST("/login", h.LoginJWT)
	ug.POST("/logout", h.LogoutJWT)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)
	ug.GET("/refresh_token", h.RefreshToken)

	// phone code
	ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
	ug.POST("/login_sms/code/verify", h.VerifySMSLoginCode)
}

// SignUp Sign up
func (h *UserHandler) SignUp(ctx *gin.Context) {
	// Define the request struct
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req SignupReq

	// Get the context content
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// Check the email pattern
	isEmail, err := h.emailRegExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "Email input failed\n")
	}
	if !isEmail {
		ctx.String(http.StatusOK, "Email pattern is wrong\n")
		return
	}

	// Check the keyword pattern
	isPassword, err := h.passwordRegex.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "Password input failed\n")
	}
	if !isPassword {
		ctx.String(http.StatusOK, "Password pattern is wrong\n")
		return
	}

	// Sign up
	err = h.usersvc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	// Check the error
	switch err {
	case nil:
		// Success:Return the string to browser
		ctx.String(http.StatusOK, "HELLO ITS IN SIGNUP")
	case service.ErrDuplicateUser:
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

	u, err := h.usersvc.Login(ctx, req.Email, req.Password)
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
	// Define the request struct
	type SignupReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req SignupReq

	// Get the context content
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// Login
	u, err := h.usersvc.Login(ctx, req.Email, req.Password)

	// Check the error
	switch err {

	// OK
	case nil:
		h.SetJWTToken(ctx, u.Id)
		ctx.String(http.StatusOK, "Login Success!!")

	// Error password or user not found
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "User not found or password is wrong")

	// Other error
	default:
		ctx.String(http.StatusOK, "System Error!!")
	}

}

func (h *UserHandler) Edit(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Hello form edit\n")
	type EditReq struct {
		NickName    string `json:"nickname"`
		Birthday    string `json:"birthday"`
		Description string `json:"description"`
	}
	var req EditReq

	//sess := sessions.Default(ctx)
	//uid := sess.Get("userId").(int64)
	//ctx.String(http.StatusOK, "uid is %d", uid)

	if err := ctx.Bind(&req); err != nil {
		return
	}

	// Verify the length of the input
	if len(req.NickName) > 8 {
		ctx.String(http.StatusOK, "Invalid nickname length")
		return
	}
	if len(req.Birthday) != 10 {
		ctx.String(http.StatusOK, "Invalid birthday length")
		return
	}
	if len(req.Description) > 50 {
		ctx.String(http.StatusOK, "Invalid description length")
		return
	}

	// Get the token from the header
	tokenstr := ctx.GetHeader("Authorization")
	segs := strings.Split(tokenstr, " ")
	if len(segs) != 2 {
		ctx.String(http.StatusOK, "Token is invalid\n")
		return
	}
	tokenstr = segs[1]
	token, err := jwt.ParseWithClaims(tokenstr, &ijwt.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return ijwt.SigKey, nil
		})
	if err != nil {
		ctx.String(http.StatusOK, "%s", err)
		return
	}

	// Get the user id from the token
	uc, ok := token.Claims.(*ijwt.UserClaims)
	if !ok || !token.Valid {
		ctx.String(http.StatusOK, "Token is invalid\n")
		return
	}

	// Edit the user profile
	uid := uc.Uid
	err = h.usersvc.Edit(ctx, uid, req.NickName, req.Birthday, req.Description)
	if err != nil {
		return
	}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	//ctx.String(http.StatusOK, "Hello form profile/n")

	//sess := sessions.Default(ctx)
	//uid := sess.Get("userId").(int64)

	// Get the user id from the token
	tokenstr := ctx.GetHeader("Authorization")
	segs := strings.Split(tokenstr, " ")
	if len(segs) != 2 {
		ctx.String(http.StatusOK, "Token is invalid\n")
		return
	}
	tokenstr = segs[1]
	token, err := jwt.ParseWithClaims(tokenstr, &ijwt.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return ijwt.SigKey, nil
		})
	if err != nil {
		ctx.String(http.StatusOK, "%s", err)
		return
	}

	// Get the user id from the token
	uc, ok := token.Claims.(*ijwt.UserClaims)
	if !ok || !token.Valid {
		ctx.String(http.StatusOK, "Token is invalid\n")
		return
	}

	uid := uc.Uid
	u, err := h.usersvc.Profile(ctx, uid)
	if err != nil {
		return
	}

	//response := fmt.Sprintf("id=%d\nEmail=%s\nNickName=%s\nBirthday=%s\nDescription=%s\n",
	//	u.Id, u.Email, u.NickName, u.Birthday, u.Description)

	// 返回用户信息
	ctx.JSON(http.StatusOK, u)
	//ctx.String(http.StatusOK, response)
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	if req.Phone == "" {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "Please input phone number",
		})
		return
	}

	err := h.codesvc.Send(ctx, bizlogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, ginx.Result{
			Msg: "Send SMS code success",
		})
	case service.ErrCodeSendTooFast:
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "Send SMS code too fast",
		})
	default:
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "Send SMS code failed",
		})
	}
}

func (h *UserHandler) VerifySMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := h.codesvc.Verify(ctx, bizlogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "System error",
		})
		zap.L().Error("Verify SMS code failed", zap.Error(err))
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "Code is wrong, please input again",
		})
		return
	}
	u, err := h.usersvc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "System error",
		})
		return
	}
	h.SetJWTToken(ctx, u.Id)
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "Login success",
	})
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	tokenStr := h.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc,
		func(token *jwt.Token) (interface{}, error) {
			return ijwt.RefreshSigKey, nil
		})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// Check ssid
	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	h.SetJWTToken(ctx, rc.Uid)
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "Refresh token success",
	})
}

func (h *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "System error",
		})
		return
	}
	ctx.String(http.StatusOK, "Logout success")
}
