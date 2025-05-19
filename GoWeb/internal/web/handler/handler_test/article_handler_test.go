package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/LEILEI0628/GinPro/GinX"
	jwtx "github.com/LEILEI0628/GinPro/middleware/jwt"
	loggerx "github.com/LEILEI0628/GinPro/middleware/logger"
	"github.com/LEILEI0628/GoWeb/internal/domain"
	"github.com/LEILEI0628/GoWeb/internal/service"
	svcmocks "github.com/LEILEI0628/GoWeb/internal/service/mocks"
	"github.com/LEILEI0628/GoWeb/internal/web/handler"
	"github.com/LEILEI0628/GoWeb/internal/web/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.ArticleServiceInterface

		reqBody string

		wantCode int
		wantRes  ginx.Result
	}{
		{
			name: "新建并发表",
			mock: func(ctrl *gomock.Controller) service.ArticleServiceInterface {
				svc := svcmocks.NewMockArticleServiceInterface(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `
{
	"title":"我的标题",
	"content": "我的内容"
}
`,
			wantCode: 200,
			wantRes: ginx.Result{
				Data: float64(1),
				Msg:  "OK",
			},
		},
		{
			name: "publish失败",
			mock: func(ctrl *gomock.Controller) service.ArticleServiceInterface {
				svc := svcmocks.NewMockArticleServiceInterface(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("publish error"))
				return svc
			},
			reqBody: `
{
	"title":"我的标题",
	"content": "我的内容"
}
`,
			wantCode: 200,
			wantRes: ginx.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("claims", &jwtx.UserClaims{
					UID: 123,
				})
			})
			// 用不上 codeSvc
			h := handler.NewArticleHandler(tc.mock(ctrl), &loggerx.NoneLogger{})
			r := router.NewArticleRouters(h)
			r.RegisterRouters(server)

			req, err := http.NewRequest(http.MethodPost,
				"/articles/publish", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes ginx.Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)
		})
	}
}
