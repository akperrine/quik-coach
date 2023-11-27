package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akperrine/quik-coach/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FindAll() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockUserService) FindOne(email, password string) map[string]interface{} {
	args := m.Called(email, password)
	
	return args.Get(0).(map[string]interface{})
}

func (m *MockUserService) CreateUser(user domain.User) (*mongo.InsertOneResult, error) {
	args := m.Called(user)

	// Need to check nil because this method returns nil for an object currently
	if args.Get(0) != nil {
		return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
	}

	return nil, args.Error(1)
}

func TestUserController_GetAllUsers(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserService.On("FindAll").Return([]byte(`[{"id": "1", "name": "John"}]`), nil)

	userController := &UserController{
		UserService: mockUserService,
	}

	req, err := http.NewRequest("GET", "/users", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(userController.GetAllUsers)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var users []map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "John", users[0]["name"])
}

func TestUserController_RegisterUser(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserService.On("CreateUser", mock.Anything).Return(nil, nil)

	userController := &UserController{
		UserService: mockUserService,
	}

	user := &domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}

	userJSON, err := json.Marshal(user)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/users/register", bytes.NewReader(userJSON))
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(userController.registerUser)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response string
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user successfully created", response)
}

func TestUserController_LoginUser(t *testing.T) {
	mockUserService := new(MockUserService)
	userController := &UserController{
		UserService: mockUserService,
	}

	mockUser := &domain.User{
		ID:        "1",
		FirstName: "A",
		LastName:  "Wade",
		Email:     "john@example.com",
	}

	mockUserService.On("FindOne", "john@example.com", "password123").Return(map[string]interface{}{"message":"logged in",
	"status": true, 
	"user": mockUser,
	}, nil)

	
	reqBody := []byte(`{
		"email": "john@example.com",
		"password": "password123"
	}`)

	req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	res := httptest.NewRecorder()
	
	userController.loginUser(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var response map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["status"])
}