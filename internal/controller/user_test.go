package controller

import (
	"crud/internal/mocks"
	"crud/internal/util/response"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const OFFSET = "offset"
const LIMIT = "limit"

func TestUnitBadRequestGetUsers(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		url          string
		errorMessage string
	}{
		{"Non integer offset GetUsers", "/api/v1/user/?offset=one",
			"strconv.Atoi: parsing \"one\": invalid syntax"},
		{"Non integer limit GetUsers", "/api/v1/user/?offset=1&limit=two",
			"strconv.Atoi: parsing \"two\": invalid syntax"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRecorder := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)

			router := gin.Default()
			routerGroup := router.Group("/api/v1")
			controller := NewUserController(nil)
			controller.SetupRoutes(routerGroup)
			router.ServeHTTP(testRecorder, req)

			assert.Equal(t, http.StatusBadRequest, testRecorder.Code)
			responseBody := testRecorder.Body.String()
			statusMessage := response.HTTPStatusMessage{}
			assert.NoError(t, json.Unmarshal([]byte(responseBody), &statusMessage))
			assert.Equal(t, tt.errorMessage, statusMessage.Message)
		})
	}
}

func TestUnitInternalErrorGetUsers(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	expectedOffset := 1
	expectedLimit := 1
	expectedErrorMessage := "some error"

	testRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(testRecorder)
	context.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	context.Params = gin.Params{
		gin.Param{Key: OFFSET, Value: strconv.Itoa(expectedOffset)},
		gin.Param{Key: LIMIT, Value: strconv.Itoa(expectedLimit)}}
	requestContext := context.Request.Context()
	mockService := mocks.NewMockIUserService(t)
	mockService.EXPECT().
		GetUsers(expectedOffset, expectedLimit, &requestContext).
		Return(nil, errors.New(expectedErrorMessage))

	controller := NewUserController(mockService)
	controller.GetUsers(context)

	assert.Equal(t, http.StatusInternalServerError, testRecorder.Code)
	responseBody := testRecorder.Body.String()
	statusMessage := response.HTTPStatusMessage{}
	assert.NoError(t, json.Unmarshal([]byte(responseBody), &statusMessage))
	assert.Equal(t, expectedErrorMessage, statusMessage.Message)
}
