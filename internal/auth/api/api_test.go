package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumane/internal/auth/database"
	mock_database "github.com/and-period/furumane/mock/auth/database"
	mock_cognito "github.com/and-period/furumane/mock/pkg/cognito"
	"github.com/and-period/furumane/pkg/jst"
	"github.com/and-period/furumane/pkg/uuid"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var (
	current   = jst.Now()
	idmock    = uuid.New()
	tokenmock = "access-token"
)

type mocks struct {
	db        *dbmocks
	adminAuth *mock_cognito.MockClient
	userAuth  *mock_cognito.MockClient
}

type dbmocks struct {
	admin *mock_database.MockAdmin
}

type testResponse struct {
	code int
	body interface{}
}

type testOptions struct {
	now  func() time.Time
	uuid func() string
}

type testOption func(opts *testOptions)

func withNow(now time.Time) testOption {
	return func(opts *testOptions) {
		opts.now = func() time.Time {
			return now
		}
	}
}

func withUUID(uuid string) testOption {
	return func(opts *testOptions) {
		opts.uuid = func() string {
			return uuid
		}
	}
}

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		db:        newDBMocks(ctrl),
		adminAuth: mock_cognito.NewMockClient(ctrl),
		userAuth:  mock_cognito.NewMockClient(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbmocks {
	return &dbmocks{
		admin: mock_database.NewMockAdmin(ctrl),
	}
}

func newController(mocks *mocks, opts *testOptions) Controller {
	params := &Params{
		WaitGroup: &sync.WaitGroup{},
		Database: &database.Database{
			Admin: mocks.db.admin,
		},
		AdminAuth: mocks.adminAuth,
		UserAuth:  mocks.userAuth,
	}
	ctrl := NewController(params).(*controller)
	ctrl.now = func() time.Time {
		return opts.now()
	}
	ctrl.uuid = func() string {
		return opts.uuid()
	}
	return ctrl
}

func newRoutes(c Controller, r *gin.Engine) {
	c.Routes(r.Group(""))
}

func testSetup(t *testing.T, ctrl *gomock.Controller, setup func(*mocks), opts ...testOption) (*controller, *testOptions) {
	gin.SetMode(gin.TestMode)
	mocks := newMocks(ctrl)
	dopts := &testOptions{
		now:  jst.Now,
		uuid: uuid.New,
	}
	for i := range opts {
		opts[i](dopts)
	}
	c := newController(mocks, dopts)
	setup(mocks)

	return c.(*controller), dopts
}

func testGet(t *testing.T, setup func(*mocks), expect *testResponse, path string, opts ...testOption) {
	testHTTP(t, setup, expect, newHTTPRequest(t, http.MethodGet, path, nil), opts...)
}

func testPost(t *testing.T, setup func(*mocks), expect *testResponse, path string, body interface{}, opts ...testOption) {
	testHTTP(t, setup, expect, newHTTPRequest(t, http.MethodPost, path, body), opts...)
}

func testPut(t *testing.T, setup func(*mocks), expect *testResponse, path string, body interface{}, opts ...testOption) {
	testHTTP(t, setup, expect, newHTTPRequest(t, http.MethodPut, path, body), opts...)
}

func testDelete(t *testing.T, setup func(*mocks), expect *testResponse, path string, opts ...testOption) {
	testHTTP(t, setup, expect, newHTTPRequest(t, http.MethodDelete, path, nil), opts...)
}

/**
 * testHTTP - HTTPハンドラのテストを実行
 */
func testHTTP(t *testing.T, setup func(*mocks), expect *testResponse, req *http.Request, opts ...testOption) {
	t.Parallel()

	// setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	h, _ := testSetup(t, ctrl, setup, opts...)
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	newRoutes(h, r)

	// test
	r.ServeHTTP(w, req)
	require.Equal(t, expect.code, w.Code)
	if isError(w) || expect.body == nil {
		return
	}

	body, err := json.Marshal(expect.body)
	require.NoError(t, err, err)
	require.JSONEq(t, string(body), w.Body.String())
}

func isError(res *httptest.ResponseRecorder) bool {
	return res.Code < 200 || 300 <= res.Code
}

/**
 * newHTTPRequest - HTTP Request(application/json)を生成
 */
func newHTTPRequest(t *testing.T, method, path string, body interface{}) *http.Request {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var buf []byte
	if body != nil {
		var err error
		buf, err = json.Marshal(body)
		require.NoError(t, err, err)
	}

	req, err := http.NewRequest(method, path, bytes.NewReader(buf))
	require.NoError(t, err, err)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenmock))
	req.Header.Add("adminId", idmock)
	return req
}

func TestController(t *testing.T) {
	t.Parallel()
	h := NewController(&Params{}, WithLogger(zap.NewNop()))
	assert.NotNil(t, h)
}
