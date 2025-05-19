package handler

import (
	"errors"
	jwtx "github.com/LEILEI0628/GinPro/middleware/jwt"
	"github.com/LEILEI0628/GinPro/middleware/session"
	"github.com/LEILEI0628/GoWeb/internal/domain"
	"github.com/LEILEI0628/GoWeb/internal/service"
	regexp "github.com/dlclark/regexp2" // 自带的regexp无法处理复杂正则
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// UserHandler User相关的业务处理
type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	service     service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler { // 使用此方法可以提示忘记传参
	const ( // 就近原则和最小化作用域原则
		emailRegexPattern = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$` // ``比""简洁（无需转义）
		//emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&.])[A-Za-z\d@$!%*#?&.]{8,72}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)

	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
		service:     service,
	}
}

func (handler *UserHandler) SignUp(ctx *gin.Context) {
	// 内部结构体：确保只有本方法能访问
	type SignUpReq struct {
		Email           string `json:"email"` // `json:"email"`：标签
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	// Bind方法会根据Content-Type来解析数据到json中，出错直接返回400状态码
	if err := ctx.Bind(&req); err != nil { // 注意：&req
		return
	}

	// 校验操作
	err := handler.checkMessage(ctx, req.Email, req.Password, req.ConfirmPassword)
	if err != nil {
		return
	}

	err = handler.service.SignUp(ctx.Request.Context(), domain.User{Email: req.Email, Password: req.Password})
	if errors.Is(err, service.ErrUserEmailDuplicated) {
		// 使用service.Err...而不直接使用dao.Err...是为了防止跨层耦合
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		// 不要在日志中打印敏感信息
		zap.L().Error("插入用户失败", zap.Error(err))
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "注册成功！")
}

func (handler *UserHandler) SignInByJWT(ctx *gin.Context) {
	type SignInReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req SignInReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	userFind, err := handler.service.SignIn(ctx.Request.Context(), domain.User{Email: req.Email, Password: req.Password})
	if errors.Is(err, service.ErrInvalidEmailOrPassword) {
		ctx.String(http.StatusOK, "账号或密码错误") // 不要明确告知账号不对或密码不对
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	tokenStr, err := handler.getJWTTokenStr(ctx, userFind)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	ctx.Header("x-jwt-token", tokenStr)
	ctx.String(http.StatusOK, "登录成功")
}

func (handler *UserHandler) getJWTTokenStr(ctx *gin.Context, userFind domain.User) (string, error) {
	userClaims := jwtx.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)), // 12小时后过期
		},
		UID:       userFind.Id,
		UserAgent: ctx.Request.UserAgent(),
	}
	return jwtx.CreateJWT([]byte("7x9FpL2QaZ8rT4wY6vBcN1mK3jH5gD7s"), userClaims)
}

func (handler *UserHandler) SignInBySession(ctx *gin.Context) {
	type SignInReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req SignInReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	userFind, err := handler.service.SignIn(ctx.Request.Context(), domain.User{Email: req.Email, Password: req.Password})
	if errors.Is(err, service.ErrInvalidEmailOrPassword) {
		ctx.String(http.StatusOK, "账号或密码错误") // 不要明确告知账号不对或密码不对
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 创建session
	err = session.CreateSession(ctx, userFind.Id, sessions.Options{
		//Secure: true, // 使用https协议
		MaxAge: 60 * 60 * 60,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
}

func (handler *UserHandler) SignOutByJWT(ctx *gin.Context) {
	// 注销session（设置过期）
	session := sessions.Default(ctx)
	session.Options(sessions.Options{
		MaxAge: -1,
	})
	err := session.Save()
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "退出成功")
}

func (handler *UserHandler) EditByJWT(ctx *gin.Context) {
	var userReq domain.UserProfile
	if err := ctx.Bind(&userReq); err != nil {
		return
	}
	_, err := time.Parse(time.DateOnly, userReq.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "生日格式错误")
		return
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*jwtx.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	UID := claims.UID

	err = handler.service.Edit(ctx.Request.Context(), UID, userReq)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, map[string]int{"code": 0})
}

func (handler *UserHandler) ProfileByJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")          // 发生error则claims为空，在类型断言时也可判断，故可以忽略此处错误
	claims, ok := c.(*jwtx.UserClaims) // 类型断言
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	UID := claims.UID
	userFind, err := handler.service.Profile(ctx.Request.Context(), UID)
	if errors.Is(err, service.ErrUserNotFound) {
		ctx.String(http.StatusOK, "账号信息未找到")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, userFind)

}

func (handler *UserHandler) ProfileBySession(ctx *gin.Context) {
	ctx.String(http.StatusOK, "个人信息")

}

func (handler *UserHandler) checkMessage(ctx *gin.Context, email string, password string, confirmPassword string) error {
	ok, err := handler.emailExp.MatchString(email)
	if err != nil {
		// 此处的err是正则表达式错误（使用regexp2时只有超时才会出现err）
		// 不要把后端错误信息（err.Error()）传到前端！
		zap.L().Error("正则表达式错误", zap.Error(err))
		ctx.String(http.StatusOK, "系统错误") // 可以直观的看到请求是否到达服务器
		//context.String(http.StatusInternalServerError, "系统错误") // RESTful风格更符合http规范
		return err
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		//context.String(http.StatusBadRequest, "邮箱格式错误")
		return errors.New("邮箱格式错误")
	}

	if password != confirmPassword {
		zap.L().Warn("password与confirmPassword不一致") // 可能为攻击行为
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return errors.New("两次输入的密码不一致")
	}

	ok, err = handler.passwordExp.MatchString(password)
	if err != nil {
		zap.L().Error("正则表达式处理错误", zap.Error(err))
		ctx.String(http.StatusOK, "系统错误")
		return err
	}

	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位且包含数字、特殊字符")
		return errors.New("密码格式错误")
	}
	return nil
}
