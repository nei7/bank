package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/nei7/bank/db/mock"
	"github.com/nei7/bank/internal/db"
	"github.com/nei7/bank/util"
	"github.com/stretchr/testify/require"
)

type createUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e createUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.Password)
	if err != nil {
		return false
	}

	e.arg.Password = arg.Password

	return reflect.DeepEqual(e.arg, arg)
}

func (e createUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %s", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return createUserParamsMatcher{arg, password}
}

func TestGetCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"email":     user.Email,
				"full_name": user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,

					Email:    user.Email,
					FullName: user.FullName,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				userBodyMatch(t, recorder.Body, user)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			request.Header.Set("Content-type", "application/json")

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}

}

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(7)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	return db.User{
		FullName: util.RandomOwner(),
		Username: util.RandomOwner(),
		Password: hashedPassword,
		Email:    util.RandomEmail(),
	}, password
}

func userBodyMatch(t *testing.T, body *bytes.Buffer, account db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var expected db.User
	err = json.Unmarshal(data, &expected)
	require.NoError(t, err)

	require.Equal(t, expected.Email, account.Email)
	require.Equal(t, expected.FullName, account.FullName)
	require.Equal(t, expected.Username, account.Username)
}
