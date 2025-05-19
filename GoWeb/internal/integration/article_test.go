package integration

import (
	"bytes"
	"encoding/json"
	jwtx "github.com/LEILEI0628/GinPro/middleware/jwt"
	"github.com/LEILEI0628/GoWeb/internal/integration/startup"
	article2 "github.com/LEILEI0628/GoWeb/internal/repository/article"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao/article"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao/po"
	"github.com/LEILEI0628/GoWeb/internal/service"
	"github.com/LEILEI0628/GoWeb/internal/web/handler"
	"github.com/LEILEI0628/GoWeb/internal/web/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ArticleTestSuite 测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetupSuite() {
	// 在所有测试执行之前初始化
	// 手动注册方式（如需对server进行定制时）：
	s.db = startup.InitTestDB()
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", &jwtx.UserClaims{
			UID: 123,
		})
	})
	artDao := article.NewArticleDAO(s.db)
	repo := article2.NewArticleRepository(artDao)
	svc := service.NewArticleService(repo)
	articleHandler := handler.NewArticleHandler(svc, startup.InitLog())
	router.NewArticleRouters(articleHandler).RegisterRouters(s.server)
}

// TearDownTest 每一个都会执行
func (s *ArticleTestSuite) TearDownTest() {
	// 清空所有数据，并且自增主键恢复到 1
	s.db.Exec("TRUNCATE TABLE articles")
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T) // 集成测试准备数据
		after  func(t *testing.T) // 集成测试验证数据
		art    Article            // 预期输入

		// HTTP响应码
		wantCode int
		// HTTP响应中携带帖子ID
		wantRes Result[int64]
	}{
		{
			name: "新建帖子并保存",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				// 验证数据库
				var art po.Article
				err := s.db.Where("id=?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.CreateTime > 0)
				assert.True(t, art.UpdateTime > 0)
				art.CreateTime = 0
				art.UpdateTime = 0
				assert.Equal(t, po.Article{
					Id:       1,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
				}, art)
			},
			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "OK",
			},
		},
		{
			name: "修改帖子并保存",
			before: func(t *testing.T) {
				// 修改数据库
				err := s.db.Create(po.Article{
					Id:         2,
					Title:      "Title标题",
					Content:    "Content内容",
					AuthorId:   123,
					CreateTime: 111,
					UpdateTime: 222,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据库
				var art po.Article
				err := s.db.Where("id=?", 2).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.UpdateTime > 222)
				art.UpdateTime = 0
				assert.Equal(t, po.Article{
					Id:         2,
					Title:      "新的标题",
					Content:    "新的内容",
					AuthorId:   123,
					CreateTime: 111,
				}, art)
			},
			art: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 2,
				Msg:  "OK",
			},
		},
		{
			name: "修改别的创作者的帖子",
			before: func(t *testing.T) {
				// 提前准备数据
				err := s.db.Create(po.Article{
					Id:      3,
					Title:   "我的标题",
					Content: "我的内容",
					// 测试模拟的用户 ID 是123，这里是 789
					// 意味着是123在修改别人的数据
					AuthorId: 789,
					// 跟时间有关的测试，尽量不要用 time.Now()
					// 因为 time.Now() 每次运行都不同，很难断言
					CreateTime: 123,
					UpdateTime: 234,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据库
				var art po.Article
				err := s.db.Where("id=?", 3).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, po.Article{
					Id:         3,
					Title:      "我的标题",
					Content:    "我的内容",
					CreateTime: 123,
					UpdateTime: 234,
					AuthorId:   789,
				}, art)
			},
			art: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 构造请求 -> 执行 -> 验证结果
			tc.before(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/edit", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			// 数据为 JSON 格式
			req.Header.Set("Content-Type", "application/json")
			// 继续使用 req
			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写回到 resp 里
			s.server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)
			tc.after(t)
		})
	}
}

func (s *ArticleTestSuite) TestAll() {
	s.T().Log("这是测试套件")
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	// 业务状态码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
