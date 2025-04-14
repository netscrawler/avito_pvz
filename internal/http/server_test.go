package httpserver_test

import (
	"net/http"
	"testing"
	"time"

	httpserver "avito_pvz/internal/http"
	"avito_pvz/internal/http/gen"
	"avito_pvz/internal/models"
	"avito_pvz/internal/models/domain"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Импортируем типы моков
type (
	MockJWTGenerator      = httpserver.MockJWTGenerator
	MockUserProvider      = httpserver.MockUserProvider
	MockPVZProvider       = httpserver.MockPVZProvider
	MockReceptionProvider = httpserver.MockReceptionProvider
	MockProductProvider   = httpserver.MockProductProvider
)

func TestServer_PostDummyLogin(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    gen.PostDummyLoginJSONRequestBody
		mockJWTToken   string
		mockJWTError   error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful login with admin role",
			requestBody: gen.PostDummyLoginJSONRequestBody{
				Role: gen.PostDummyLoginJSONBodyRoleModerator,
			},
			mockJWTToken:   "test-token",
			expectedStatus: http.StatusOK,
			expectedBody:   "test-token",
		},
		{
			name: "successful login with pvz role",
			requestBody: gen.PostDummyLoginJSONRequestBody{
				Role: gen.PostDummyLoginJSONBodyRoleEmployee,
			},
			mockJWTToken:   "test-token",
			expectedStatus: http.StatusOK,
			expectedBody:   "test-token",
		},
		{
			name: "invalid role",
			requestBody: gen.PostDummyLoginJSONRequestBody{
				Role: "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gen.Error{
				Message: "InvalidRole",
			},
		},
		{
			name: "jwt generation error",
			requestBody: gen.PostDummyLoginJSONRequestBody{
				Role: gen.PostDummyLoginJSONBodyRoleModerator,
			},
			mockJWTError:   models.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: gen.Error{
				Message: "InternalError",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// // Создаем моки
			// mockJWT := &MockJWTGenerator{}
			// mockUser := &MockUserProvider{}
			// mockPVZ := &MockPVZProvider{}
			// mockReception := &MockReceptionProvider{}
			// mockProduct := &MockProductProvider{}
			//
			// // Настраиваем ожидания для моков только для валидных ролей
			// if tt.requestBody.Role == gen.PostDummyLoginJSONBodyRoleModerator ||
			// 	tt.requestBody.Role == gen.PostDummyLoginJSONBodyRoleEmployee {
			// 	if tt.mockJWTError == nil {
			// 		mockJWT.On("GenerateToken", "dummy", string(tt.requestBody.Role)).
			// 			Return(tt.mockJWTToken, nil)
			// 	} else {
			// 		mockJWT.On("GenerateToken", "dummy", string(tt.requestBody.Role)).
			// 			Return("", tt.mockJWTError)
			// 	}
			// }

			// Создаем сервер с моками
			// server := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			//
			// // Создаем тестовый запрос
			// body, err := json.Marshal(tt.requestBody)
			// require.NoError(t, err)
			//
			// req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
			// rec := httptest.NewRecorder()

			// Вызываем тестируемый метод
			// server.PostDummyLogin(rec, req)

			// Проверяем статус код
			// assert.Equal(t, tt.expectedStatus, rec.Code)
			//
			// // Проверяем тело ответа
			// if tt.expectedStatus == http.StatusOK {
			// 	var response string
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// } else {
			// 	var response gen.Error
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// }
			//
			// // Проверяем, что все ожидания были выполнены
			// mockJWT.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func TestServer_PostLogin(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    gen.PostLoginJSONRequestBody
		mockUserToken  *string
		mockUserError  error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful login",
			requestBody: gen.PostLoginJSONRequestBody{
				Email:    openapi_types.Email("test@example.com"),
				Password: "password123",
			},
			mockUserToken:  stringPtr("test-token"),
			expectedStatus: http.StatusOK,
			expectedBody:   "test-token",
		},
		{
			name: "invalid credentials",
			requestBody: gen.PostLoginJSONRequestBody{
				Email:    openapi_types.Email("test@example.com"),
				Password: "wrong-password",
			},
			mockUserError:  models.ErrInvalidPassword,
			expectedStatus: http.StatusBadRequest,
			expectedBody: gen.Error{
				Message: models.ErrInvalidPassword.Error(),
			},
		},
		{
			name: "internal error",
			requestBody: gen.PostLoginJSONRequestBody{
				Email:    openapi_types.Email("test@example.com"),
				Password: "password123",
			},
			mockUserError:  models.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: gen.Error{
				Message: models.ErrInternal.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем моки
			// mockJWT := &MockJWTGenerator{}
			// mockUser := &MockUserProvider{}
			// mockPVZ := &MockPVZProvider{}
			// mockReception := &MockReceptionProvider{}
			// mockProduct := &MockProductProvider{}
			//
			// // Настраиваем ожидания для моков
			// mockUser.On("Auth", mock.Anything, string(tt.requestBody.Email), tt.requestBody.Password).
			// 	Return(tt.mockUserToken, tt.mockUserError)
			//
			// // Создаем сервер с моками
			// // server := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			// //
			// // // Создаем тестовый запрос
			// // body, err := json.Marshal(tt.requestBody)
			// // require.NoError(t, err)
			// //
			// // req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			// // rec := httptest.NewRecorder()
			// //
			// // // Вызываем тестируемый метод
			// // server.PostLogin(rec, req)
			//
			// // Проверяем статус код
			// assert.Equal(t, tt.expectedStatus, rec.Code)
			//
			// // Проверяем тело ответа
			// if tt.expectedStatus == http.StatusOK {
			// 	var response string
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// } else {
			// 	var response gen.Error
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// }

			// Проверяем, что все ожидания были выполнены
			// mockUser.AssertExpectations(t)
		})
	}
}

func TestServer_PostProducts(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    gen.PostProductsJSONRequestBody
		mockProduct    *domain.Product
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful product creation",
			requestBody: gen.PostProductsJSONRequestBody{
				PvzId: openapi_types.UUID(uuid.Max),
				Type:  gen.PostProductsJSONBodyType(gen.ProductTypeЭлектроника),
			},
			mockProduct: &domain.Product{
				ID:          uuid.UUID{},
				CreatedAt:   time.Time{},
				Type:        domain.ProductType(gen.ProductTypeЭлектроника),
				ReceptionID: uuid.Max,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid reception ID",
			requestBody: gen.PostProductsJSONRequestBody{
				PvzId: openapi_types.UUID(uuid.Max),
				Type:  gen.PostProductsJSONBodyType(gen.ProductTypeЭлектроника),
			},
			mockProduct: &domain.Product{
				ID:          uuid.UUID{},
				CreatedAt:   time.Time{},
				Type:        domain.ProductType(gen.ProductTypeЭлектроника),
				ReceptionID: uuid.Max,
			},

			mockError: models.ErrReceptionDontExist,
			expectedBody: gen.Error{
				Message: models.ErrReceptionDontExist.Error(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "reception already closed",
			requestBody: gen.PostProductsJSONRequestBody{
				PvzId: openapi_types.UUID(uuid.Max),
				Type:  gen.PostProductsJSONBodyType(gen.ProductTypeЭлектроника),
			},
			mockProduct: &domain.Product{
				ID:          uuid.UUID{},
				CreatedAt:   time.Time{},
				Type:        domain.ProductType(gen.ProductTypeЭлектроника),
				ReceptionID: uuid.Max,
			},

			mockError:      models.ErrReceptionAlreadyClosed,
			expectedStatus: http.StatusBadRequest,
			expectedBody: gen.Error{
				Message: models.ErrReceptionAlreadyClosed.Error(),
			},
		},
		{
			name: "internal error",
			requestBody: gen.PostProductsJSONRequestBody{
				PvzId: openapi_types.UUID(uuid.Max),
				Type:  gen.PostProductsJSONBodyType(gen.ProductTypeЭлектроника),
			},

			mockError:      models.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: gen.Error{
				Message: models.ErrInternal.Error(),
			},
		},
		{
			name: "invalid product type",
			requestBody: gen.PostProductsJSONRequestBody{
				PvzId: openapi_types.UUID(uuid.Max),
				Type:  gen.PostProductsJSONBodyType("any"),
			},

			mockError:      models.ErrInvalidProductType,
			expectedStatus: http.StatusBadRequest,
			expectedBody: gen.Error{
				Message: models.ErrInvalidProductType.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// // Create mocks
			// mockJWT := &MockJWTGenerator{}
			// mockUser := &MockUserProvider{}
			// mockPVZ := &MockPVZProvider{}
			// mockReception := &MockReceptionProvider{}
			// mockProduct := &MockProductProvider{}
			//
			// // Set up expectations for product mock
			// mockProduct.On("Create", mock.Anything, mock.MatchedBy(func(p domain.ProductToAdd) bool {
			// 	return p.UUID == domain.PVZID(tt.requestBody.PvzId) &&
			// 		p.Type == domain.ProductType(tt.requestBody.Type)
			// })).
			// 	Return(tt.mockProduct, tt.mockError)
			//
			// // Create server with mocks
			// server := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			//
			// // Create test request
			// body, err := json.Marshal(tt.requestBody)
			// require.NoError(t, err)
			//
			// req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
			// rec := httptest.NewRecorder()
			//
			// // Call the tested method
			// server.PostProducts(rec, req)
			//
			// // Check status code
			// assert.Equal(t, tt.expectedStatus, rec.Code)
			//
			// // Check response body
			// if tt.expectedStatus == http.StatusCreated {
			// 	var response gen.Product
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.mockProduct.ToDto(), response)
			// } else {
			// 	var response gen.Error
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// }
			//
			// // Verify that all expectations were met
			// mockProduct.AssertExpectations(t)
		})
	}
}

// Add JSON error test cases
func TestServer_PostProducts_JSONErrors(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "invalid json",
			requestBody:    `{"invalid": json}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"ErrorWithJsonDecode"}`,
		},
		{
			name:           "empty body",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"ErrorWithJsonDecode"}`,
		},
		{
			name:           "invalid reception_id format",
			requestBody:    `{"name": "Test", "reception_id": "invalid-uuid", "count": 1}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"InvalidPVZId"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mockJWT := httpserver.NewMockJWTGenerator(t)
			// mockUser := httpserver.NewMockUserProvider(t)
			// mockPVZ := httpserver.NewMockPVZProvider(t)
			// mockReception := httpserver.NewMockReceptionProvider(t)
			// mockProduct := httpserver.NewMockProductProvider(t)
			//
			// s := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			//
			// req := httptest.NewRequest(
			// 	http.MethodPost,
			// 	"/products",
			// 	strings.NewReader(tt.requestBody),
			// )
			// rec := httptest.NewRecorder()
			//
			// s.PostProducts(rec, req)
			//
			// assert.Equal(t, tt.expectedStatus, rec.Code)
			// assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestServer_GetPvz(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	page := 1
	limit := 10

	tests := []struct {
		name           string
		params         gen.GetPvzParams
		mockPVZList    *domain.PVZList
		mockPVZError   error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful get pvz list",
			params: gen.GetPvzParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Page:      &page,
				Limit:     &limit,
			},
			mockPVZList: &domain.PVZList{
				{
					City: "Moscow",
				},
				{
					City: "Saint Petersburg",
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: []gen.PVZ{
				{
					City: "Moscow",
				},
				{
					City: "Saint Petersburg",
				},
			},
		},
		{
			name: "internal error",
			params: gen.GetPvzParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Page:      &page,
				Limit:     &limit,
			},
			mockPVZError:   models.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: gen.Error{
				Message: models.ErrInternal.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем моки
			// mockJWT := &MockJWTGenerator{}
			// mockUser := &MockUserProvider{}
			// mockPVZ := &MockPVZProvider{}
			// mockReception := &MockReceptionProvider{}
			// mockProduct := &MockProductProvider{}
			//
			// // Настраиваем ожидания для моков
			// mockPVZ.On("List", mock.Anything, domain.Params{
			// 	StartDate: tt.params.StartDate,
			// 	EndDate:   tt.params.StartDate, // В реализации сервера EndDate = StartDate
			// 	Page:      tt.params.Page,
			// 	Limit:     tt.params.Limit,
			// }).Return(tt.mockPVZList, tt.mockPVZError)
			//
			// // Создаем сервер с моками
			// server := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			//
			// // Создаем тестовый запрос
			// req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
			// rec := httptest.NewRecorder()
			//
			// // Вызываем тестируемый метод
			// server.GetPvz(rec, req, tt.params)
			//
			// // Проверяем статус код
			// assert.Equal(t, tt.expectedStatus, rec.Code)
			//
			// // Проверяем тело ответа
			// if tt.expectedStatus == http.StatusOK {
			// 	var response []gen.PVZ
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// } else {
			// 	var response gen.Error
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// }
			//
			// // Проверяем, что все ожидания были выполнены
			// mockPVZ.AssertExpectations(t)
		})
	}
}

func TestServer_PostPvz(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    gen.PostPvzJSONRequestBody
		mockPVZ        *domain.PVZ
		mockPVZError   error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful pvz creation",
			requestBody: gen.PostPvzJSONRequestBody{
				City: "Moscow",
			},
			mockPVZ: &domain.PVZ{
				City: "Moscow",
			},
			expectedStatus: http.StatusCreated,
			expectedBody: domain.PVZ{
				City: "Moscow",
			},
		},
		{
			name: "invalid city",
			requestBody: gen.PostPvzJSONRequestBody{
				City: "InvalidCity",
			},
			mockPVZError:   models.ErrInvalidCity,
			expectedStatus: http.StatusBadRequest,
			expectedBody: gen.Error{
				Message: models.ErrInvalidCity.Error(),
			},
		},
		{
			name: "internal error",
			requestBody: gen.PostPvzJSONRequestBody{
				City: "Moscow",
			},
			mockPVZError:   models.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: gen.Error{
				Message: models.ErrInternal.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем моки
			// mockJWT := &MockJWTGenerator{}
			// mockUser := &MockUserProvider{}
			// mockPVZ := &MockPVZProvider{}
			// mockReception := &MockReceptionProvider{}
			// mockProduct := &MockProductProvider{}
			//
			// // Настраиваем ожидания для моков
			// mockPVZ.On("Create", mock.Anything, domain.PVZCity(tt.requestBody.City)).
			// 	Return(tt.mockPVZ, tt.mockPVZError)
			//
			// // Создаем сервер с моками
			// server := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			//
			// // Создаем тестовый запрос
			// body, err := json.Marshal(tt.requestBody)
			// require.NoError(t, err)
			//
			// req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
			// rec := httptest.NewRecorder()
			//
			// // Вызываем тестируемый метод
			// server.PostPvz(rec, req)
			//
			// // Проверяем статус код
			// assert.Equal(t, tt.expectedStatus, rec.Code)
			//
			// // Проверяем тело ответа
			// if tt.expectedStatus == http.StatusCreated {
			// 	var response domain.PVZ
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// } else {
			// 	var response gen.Error
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// }
			//
			// // Проверяем, что все ожидания были выполнены
			// mockPVZ.AssertExpectations(t)
		})
	}
}

func TestServer_PostPvz_JSONErrors(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "invalid json",
			requestBody:    `{"invalid": json}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"ErrorWithJsonDecode"}`,
		},
		{
			name:           "empty body",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"ErrorWithJsonDecode"}`,
		},
		{
			name:           "invalid city format",
			requestBody:    `{"city": 123}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"ErrorWithJsonDecode"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mockJWT := httpserver.NewMockJWTGenerator(t)
			// mockUser := httpserver.NewMockUserProvider(t)
			// mockPVZ := httpserver.NewMockPVZProvider(t)
			// mockReception := httpserver.NewMockReceptionProvider(t)
			// mockProduct := httpserver.NewMockProductProvider(t)
			//
			// s := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			//
			// req := httptest.NewRequest(http.MethodPost, "/pvz", strings.NewReader(tt.requestBody))
			// rec := httptest.NewRecorder()
			//
			// s.PostPvz(rec, req)
			//
			// assert.Equal(t, tt.expectedStatus, rec.Code)
			// assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestServer_PostPvzPvzIdCloseLastReception(t *testing.T) {
	tests := []struct {
		name           string
		pvzID          openapi_types.UUID
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "successful close last reception",
			pvzID:          openapi_types.UUID(uuid.New()),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "reception already closed",
			pvzID:          openapi_types.UUID(uuid.New()),
			mockError:      models.ErrReceptionAlreadyClosed,
			expectedStatus: http.StatusBadRequest,
			expectedBody: gen.Error{
				Message: models.ErrReceptionAlreadyClosed.Error(),
			},
		},
		{
			name:           "reception not found",
			pvzID:          openapi_types.UUID(uuid.New()),
			mockError:      models.ErrReceptionDontExist,
			expectedStatus: http.StatusBadRequest,
			expectedBody: gen.Error{
				Message: models.ErrReceptionDontExist.Error(),
			},
		},
		{
			name:           "internal error",
			pvzID:          openapi_types.UUID(uuid.New()),
			mockError:      models.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: gen.Error{
				Message: models.ErrInternal.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем моки
			// mockJWT := &MockJWTGenerator{}
			// mockUser := &MockUserProvider{}
			// mockPVZ := &MockPVZProvider{}
			// mockReception := &MockReceptionProvider{}
			// mockProduct := &MockProductProvider{}
			//
			// // Настраиваем ожидания для моков
			// mockReception.On("CloseLastReception", mock.Anything, domain.PVZID(tt.pvzID)).
			// 	Return(tt.mockError)
			//
			// // Создаем сервер с моками
			// server := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			//
			// // Создаем тестовый запрос
			// req := httptest.NewRequest(
			// 	http.MethodPost,
			// 	"/pvz/"+tt.pvzID.String()+"/close_last_reception",
			// 	nil,
			// )
			// rec := httptest.NewRecorder()
			//
			// // Вызываем тестируемый метод
			// server.PostPvzPvzIdCloseLastReception(rec, req, tt.pvzID)
			//
			// // Проверяем статус код
			// assert.Equal(t, tt.expectedStatus, rec.Code)
			//
			// // Проверяем тело ответа только в случае ошибки
			// if tt.expectedStatus != http.StatusOK {
			// 	var response gen.Error
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// }
			//
			// // Проверяем, что все ожидания были выполнены
			// mockReception.AssertExpectations(t)
		})
	}
}

func TestServer_PostPvzPvzIdDeleteLastProduct(t *testing.T) {
	tests := []struct {
		name           string
		pvzID          openapi_types.UUID
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "successful delete last product",
			pvzID:          openapi_types.UUID(uuid.New()),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "product not found",
			pvzID:          openapi_types.UUID(uuid.New()),
			mockError:      models.ErrProductNotFound,
			expectedStatus: http.StatusBadRequest,
			expectedBody: gen.Error{
				Message: models.ErrProductNotFound.Error(),
			},
		},
		{
			name:           "internal error",
			pvzID:          openapi_types.UUID(uuid.New()),
			mockError:      models.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: gen.Error{
				Message: models.ErrInternal.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем моки
			// mockJWT := &MockJWTGenerator{}
			// mockUser := &MockUserProvider{}
			// mockPVZ := &MockPVZProvider{}
			// mockReception := &MockReceptionProvider{}
			// mockProduct := &MockProductProvider{}
			//
			// // Настраиваем ожидания для моков
			// mockProduct.On("DeleteLast", mock.Anything, domain.PVZID(tt.pvzID)).
			// 	Return(tt.mockError)
			//
			// // Создаем сервер с моками
			// server := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			//
			// // Создаем тестовый запрос
			// req := httptest.NewRequest(
			// 	http.MethodPost,
			// 	"/pvz/"+tt.pvzID.String()+"/delete_last_product",
			// 	nil,
			// )
			// rec := httptest.NewRecorder()
			//
			// // Вызываем тестируемый метод
			// server.PostPvzPvzIdDeleteLastProduct(rec, req, tt.pvzID)
			//
			// // Проверяем статус код
			// assert.Equal(t, tt.expectedStatus, rec.Code)
			//
			// // Проверяем тело ответа только в случае ошибки
			// if tt.expectedStatus != http.StatusOK {
			// 	var response gen.Error
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// }
			//
			// // Проверяем, что все ожидания были выполнены
			// mockProduct.AssertExpectations(t)
		})
	}
}

func TestServer_PostReceptions(t *testing.T) {
	type testCase struct {
		name           string
		pvzID          domain.PVZID
		mockReception  *domain.Reception
		mockError      error
		expectedStatus int
		expectedBody   string
	}

	testCases := []testCase{
		{
			name:           "success",
			pvzID:          domain.PVZID(uuid.New()),
			mockReception:  &domain.Reception{},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "reception already exists",
			pvzID:          domain.PVZID(uuid.New()),
			mockReception:  nil,
			mockError:      models.ErrReceptionAlreadyExists,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"ReceptionAlreadyClosed"}`,
		},
		{
			name:           "pvz not found",
			pvzID:          domain.PVZID(uuid.New()),
			mockReception:  nil,
			mockError:      models.ErrPVZNotFound,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"PVZNotFound"}`,
		},
		{
			name:           "internal error",
			pvzID:          domain.PVZID(uuid.New()),
			mockReception:  nil,
			mockError:      models.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"InternalError"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// mockJWT := httpserver.NewMockJWTGenerator(t)
			// mockUser := httpserver.NewMockUserProvider(t)
			// mockPVZ := httpserver.NewMockPVZProvider(t)
			// mockReception := httpserver.NewMockReceptionProvider(t)
			// mockProduct := httpserver.NewMockProductProvider(t)
			//
			// s := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			//
			// mockReception.EXPECT().
			// 	Create(mock.Anything, tc.pvzID).
			// 	Return(tc.mockReception, tc.mockError)
			//
			// reqBody := gen.PostReceptionsJSONRequestBody{
			// 	PvzId: openapi_types.UUID(tc.pvzID),
			// }
			// body, err := json.Marshal(reqBody)
			// require.NoError(t, err)
			//
			// rec := httptest.NewRecorder()
			//
			// ctx:= context.Background()
			// s.PostReceptions(ctx, req)
			//
			// assert.Equal(t, tc.expectedStatus, rec.Code)
			// if tc.expectedBody != "" {
			// 	assert.JSONEq(t, tc.expectedBody, rec.Body.String())
			// }
		})
	}
}

func TestServer_PostReceptions_JSONErrors(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "invalid json",
			requestBody:    `{"invalid": json`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"ErrorWithJsonDecode"}`,
		},
		{
			name:           "empty body",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"ErrorWithJsonDecode"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mockJWT := httpserver.NewMockJWTGenerator(t)
			// mockUser := httpserver.NewMockUserProvider(t)
			// mockPVZ := httpserver.NewMockPVZProvider(t)
			// mockReception := httpserver.NewMockReceptionProvider(t)
			// mockReception.AssertNotCalled(t, "Create")
			//
			// mockProduct := httpserver.NewMockProductProvider(t)
			//
			// s := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			//
			// req := httptest.NewRequest(
			// 	http.MethodPost,
			// 	"/receptions",
			// 	strings.NewReader(tt.requestBody),
			// )
			// rec := httptest.NewRecorder()
			//
			// s.PostReceptions(rec, req)
			//
			// assert.Equal(t, tt.expectedStatus, rec.Code)
			// assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestServer_PostRegister(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    gen.PostRegisterJSONRequestBody
		mockUser       *domain.User
		mockUserError  error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful registration",
			requestBody: gen.PostRegisterJSONRequestBody{
				Email:    "test@example.com",
				Password: "password123",
				Role:     "admin",
			},
			mockUser: &domain.User{
				Email: "test@example.com",
				Role:  "admin",
			},
			expectedStatus: http.StatusCreated,
			expectedBody: domain.User{
				Email: "test@example.com",
				Role:  "admin",
			},
		},
		{
			name: "user already exists",
			requestBody: gen.PostRegisterJSONRequestBody{
				Email:    "test@example.com",
				Password: "password123",
				Role:     "admin",
			},
			mockUserError:  models.ErrUserAlreadyExist,
			expectedStatus: http.StatusBadRequest,
			expectedBody: gen.Error{
				Message: models.ErrUserAlreadyExist.Error(),
			},
		},
		{
			name: "internal error",
			requestBody: gen.PostRegisterJSONRequestBody{
				Email:    "test@example.com",
				Password: "password123",
				Role:     "admin",
			},
			mockUserError:  models.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: gen.Error{
				Message: models.ErrInternal.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// // Создаем моки
			// mockJWT := &MockJWTGenerator{}
			// mockUser := &MockUserProvider{}
			// mockPVZ := &MockPVZProvider{}
			// mockReception := &MockReceptionProvider{}
			// mockProduct := &MockProductProvider{}
			//
			// // Настраиваем ожидания для моков
			// mockUser.On("Create", mock.Anything, domain.Email(tt.requestBody.Email), tt.requestBody.Password, domain.UserRole(tt.requestBody.Role)).
			// 	Return(tt.mockUser, tt.mockUserError)
			//
			// // Создаем сервер с моками
			// server := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			//
			// // Создаем тестовый запрос
			// body, err := json.Marshal(tt.requestBody)
			// require.NoError(t, err)
			//
			// req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			// rec := httptest.NewRecorder()
			//
			// // Вызываем тестируемый метод
			// server.PostRegister(rec, req)
			//
			// // Проверяем статус код
			// assert.Equal(t, tt.expectedStatus, rec.Code)
			//
			// // Проверяем тело ответа
			// if tt.expectedStatus == http.StatusCreated {
			// 	var response domain.User
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// } else {
			// 	var response gen.Error
			// 	err := json.NewDecoder(rec.Body).Decode(&response)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, tt.expectedBody, response)
			// }
			//
			// // Проверяем, что все ожидания были выполнены
			// mockUser.AssertExpectations(t)
		})
	}
}

func TestServer_PostRegister_JSONErrors(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "invalid json",
			requestBody:    `{"invalid": json}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"ErrorWithJsonDecode"}`,
		},
		{
			name:           "empty body",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"ErrorWithJsonDecode"}`,
		},
		{
			name:           "invalid email format",
			requestBody:    `{"email": 123, "password": "pass", "role": "admin"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"ErrorWithJsonDecode"}`,
		},
		{
			name:           "invalid role format",
			requestBody:    `{"email": "test@example.com", "password": "pass", "role": 123}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"ErrorWithJsonDecode"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mockJWT := httpserver.NewMockJWTGenerator(t)
			// mockUser := httpserver.NewMockUserProvider(t)
			// mockPVZ := httpserver.NewMockPVZProvider(t)
			// mockReception := httpserver.NewMockReceptionProvider(t)
			// mockProduct := httpserver.NewMockProductProvider(t)
			//
			// s := httpserver.NewServer(mockJWT, mockUser, mockPVZ, mockReception, mockProduct)
			//
			// req := httptest.NewRequest(
			// 	http.MethodPost,
			// 	"/register",
			// 	strings.NewReader(tt.requestBody),
			// )
			// rec := httptest.NewRecorder()
			//
			// s.PostRegister(rec, req)
			//
			// assert.Equal(t, tt.expectedStatus, rec.Code)
			// assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
