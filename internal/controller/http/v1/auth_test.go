package v1

import (
	"account-management/internal/controller/http/v1/mocks"
	"account-management/internal/entity"
	"account-management/internal/service"
	"bytes"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()

	v := validator.New()

	// Добавьте вашу кастомную функцию валидации для 'password'
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		lower := regexp.MustCompile(`[a-z]`)
		upper := regexp.MustCompile(`[A-Z]`)
		digit := regexp.MustCompile(`\d`)
		length := regexp.MustCompile(`.{8,}`)

		return lower.MatchString(password) && upper.MatchString(password) && digit.MatchString(password) && length.MatchString(password)
	})

	e.Validator = &CustomValidator{validator: v}

	mockAuth := mocks.NewMockAuth(ctrl)

	r := &authRoutes{
		authService: mockAuth,
	}

	tests := []struct {
		name           string
		input          signUpInput
		body           []byte
		mockReturnID   int64
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name:           "Successful user creation",
			input:          signUpInput{Username: "test", Password: "Password1"},
			mockReturnID:   1,
			mockReturnErr:  nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid request body",
			body:           []byte("invalid json"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "User already exists",
			input:          signUpInput{Username: "test", Password: "Password1"},
			mockReturnID:   0,
			mockReturnErr:  service.ErrUserAlreadyExists,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := tt.body
			if body == nil {
				body, _ = json.Marshal(tt.input)
			}

			httpReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			rec := httptest.NewRecorder()
			httpReq.Header.Set("Content-Type", "application/json")

			context := e.NewContext(httpReq, rec)

			if tt.mockReturnErr != nil || tt.mockReturnID != 0 {
				mockAuth.EXPECT().CreateUserWithAccount(gomock.Any(), gomock.Eq(entity.AuthCreateUserInput{
					Username: tt.input.Username,
					Password: tt.input.Password,
				})).Return(tt.mockReturnID, tt.mockReturnErr)
			}

			err := r.signUp(context)
			if tt.expectedStatus >= 400 {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
