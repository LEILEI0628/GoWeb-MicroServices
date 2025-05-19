package handler_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	jwtx "github.com/LEILEI0628/GinPro/middleware/jwt"
	"github.com/LEILEI0628/GoWeb/internal/domain"
	"github.com/LEILEI0628/GoWeb/internal/service"
	svcmocks "github.com/LEILEI0628/GoWeb/internal/service/mocks"
	"github.com/LEILEI0628/GoWeb/internal/web/handler"
	"github.com/LEILEI0628/GoWeb/internal/web/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncrypt(t *testing.T) {
	password := "hello@123"
	value, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(value))
	err = bcrypt.CompareHashAndPassword(value, []byte(password))
	assert.NoError(t, err)
}

func TestNil(t *testing.T) {
	testTypeAssert(nil)
}

func testTypeAssert(c any) {
	_, ok := c.(jwtx.UserClaims)
	println(ok)
}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.UserServiceInterface

		reqBody string

		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				usersvc := svcmocks.NewMockUserServiceInterface(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(nil)
				// 注册成功是 return nil
				return usersvc
			},

			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功！",
		},
		{
			name: "参数不对，bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				usersvc := svcmocks.NewMockUserServiceInterface(ctrl)
				// 注册成功是 return nil
				return usersvc
			},

			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world123"
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				usersvc := svcmocks.NewMockUserServiceInterface(ctrl)
				return usersvc
			},

			reqBody: `
{
	"email": "123@q",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱格式错误",
		},
		{
			name: "两次输入密码不匹配",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				usersvc := svcmocks.NewMockUserServiceInterface(ctrl)
				return usersvc
			},

			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world1234",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				usersvc := svcmocks.NewMockUserServiceInterface(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello123",
	"confirmPassword": "hello123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位且包含数字、特殊字符",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				usersvc := svcmocks.NewMockUserServiceInterface(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(service.ErrUserEmailDuplicated)
				// 注册成功是 return nil
				return usersvc
			},

			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				usersvc := svcmocks.NewMockUserServiceInterface(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(errors.New("随便一个 error"))
				// 注册成功是 return nil
				return usersvc
			},

			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			// 用不上 codeSvc
			h := handler.NewUserHandler(tc.mock(ctrl))
			r := router.NewUserRouters(h)
			r.RegisterRouters(server)

			req, err := http.NewRequest(http.MethodPost,
				"/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			// 数据是 JSON 格式
			req.Header.Set("Content-Type", "application/json")
			// 这里你就可以继续使用 req

			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口。
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写回到 resp 里
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())

		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersvc := svcmocks.NewMockUserServiceInterface(ctrl)

	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
		Return(errors.New("mock error"))

	//usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
	//	Email: "124@qq.com",
	//}).Return(errors.New("mock error"))

	err := usersvc.SignUp(context.Background(), domain.User{
		Email: "123@qq.com",
	})
	t.Log(err)
}
