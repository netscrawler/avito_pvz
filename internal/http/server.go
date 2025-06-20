//nolint:staticcheck
package httpserver

import (
	"avito_pvz/internal/http/gen"
	"avito_pvz/internal/models/domain"
	"context"

	"github.com/oapi-codegen/runtime/types"
)

type JWTGenerator interface {
	GenerateToken(email, role string) (string, error)
}

type UserProvider interface {
	Create(
		ctx context.Context,
		email domain.Email,
		password string,
		role domain.Role,
	) (*domain.User, error)
	Auth(ctx context.Context, email string, password string) (*string, error)
}

type PVZProvider interface {
	List(ctx context.Context, params domain.Params) (*[]domain.PVZAgregate, error)
	Create(ctx context.Context, city domain.PvzCity) (*domain.PVZ, error)
}

type ReceptionProvider interface {
	CloseLastReception(ctx context.Context, pvzID domain.PVZID) (*domain.Reception, error)
	Create(ctx context.Context, pvzID domain.PVZID) (*domain.Reception, error)
}

type ProductProvider interface {
	Create(ctx context.Context, protduct domain.ProductToAdd) (*domain.Product, error)
	DeleteLast(ctx context.Context, pvzID domain.PVZID) error
}

type Server struct {
	jwt       JWTGenerator
	user      UserProvider
	pvz       PVZProvider
	reception ReceptionProvider
	product   ProductProvider
}

// (POST /dummyLogin).
func (s *Server) PostDummyLogin(
	ctx context.Context,
	request gen.PostDummyLoginRequestObject,
) (gen.PostDummyLoginResponseObject, error) {
	role := domain.Role(request.Body.Role)
	if !role.IsValid() {
		return gen.PostDummyLogin400JSONResponse{
			Message: ErrInvalidRole.Error(),
		}, ErrInvalidRole
	}

	token, err := s.jwt.GenerateToken("dummy", string(request.Body.Role))
	if err != nil {
		return gen.PostDummyLogin400JSONResponse{
			Message: err.Error(),
		}, err
	}

	response := gen.PostDummyLogin200JSONResponse(token)

	return response, nil
}

// (POST /login).
func (s *Server) PostLogin(
	ctx context.Context,
	request gen.PostLoginRequestObject,
) (gen.PostLoginResponseObject, error) {
	email, password := request.Body.Email, request.Body.Password

	token, err := s.user.Auth(ctx, string(email), password)
	if err != nil {
		return gen.PostLogin401JSONResponse{
			Message: err.Error(),
		}, err
	}

	resp := gen.PostLogin200JSONResponse(*token)

	return resp, nil
}

// (POST /products).
func (s *Server) PostProducts(
	ctx context.Context,
	request gen.PostProductsRequestObject,
) (gen.PostProductsResponseObject, error) {
	pvzId, typeName := request.Body.PvzId, request.Body.Type

	toAdd := domain.ProductToAdd{
		UUID: domain.PVZID(pvzId),
		Type: domain.ProductType(typeName),
	}

	product, err := s.product.Create(ctx, toAdd)
	if err != nil {
		return gen.PostProducts400JSONResponse{
			Message: err.Error(),
		}, err
	}

	return gen.PostProducts201JSONResponse{
		DateTime:    &product.CreatedAt,
		Id:          (*types.UUID)(&product.ID),
		ReceptionId: product.ReceptionID,
		Type:        gen.ProductType(product.Type),
	}, nil
}

// (GET /pvz).
func (s *Server) GetPvz(
	ctx context.Context,
	request gen.GetPvzRequestObject,
) (gen.GetPvzResponseObject, error) {
	params := domain.NewParamsFromDTO(request.Params)

	ag, err := s.pvz.List(ctx, *params)
	if err != nil {
		return gen.GetPvz200JSONResponse{}, err
	}

	resp := domain.AggregateToPvzResponse(*ag)

	return gen.GetPvz200JSONResponse(resp), nil
}

// (POST /pvz).
func (s *Server) PostPvz(
	ctx context.Context,
	request gen.PostPvzRequestObject,
) (gen.PostPvzResponseObject, error) {
	city := request.Body.City

	pvz, err := s.pvz.Create(ctx, domain.PvzCity(city))
	if err != nil {
		return gen.PostPvz400JSONResponse{
			Message: err.Error(),
		}, err
	}

	return gen.PostPvz201JSONResponse{
		City:             gen.PVZCity(pvz.City),
		Id:               (*types.UUID)(pvz.ID),
		RegistrationDate: &pvz.RegistrationDate,
	}, nil
}

// (POST /pvz/{pvzId}/close_last_reception).
func (s *Server) PostPvzPvzIdCloseLastReception(
	ctx context.Context,
	request gen.PostPvzPvzIdCloseLastReceptionRequestObject,
) (gen.PostPvzPvzIdCloseLastReceptionResponseObject, error) {
	recId := request.PvzId

	rec, err := s.reception.CloseLastReception(ctx, domain.PVZID(recId))
	if err != nil {
		return gen.PostPvzPvzIdCloseLastReception400JSONResponse{
			Message: err.Error(),
		}, err
	}

	r := rec.ToDTO()

	return gen.PostPvzPvzIdCloseLastReception200JSONResponse(r), nil
}

// (POST /pvz/{pvzId}/delete_last_product).
func (s *Server) PostPvzPvzIdDeleteLastProduct(
	ctx context.Context,
	request gen.PostPvzPvzIdDeleteLastProductRequestObject,
) (gen.PostPvzPvzIdDeleteLastProductResponseObject, error) {
	pvzId := request.PvzId

	err := s.product.DeleteLast(ctx, domain.PVZID(pvzId))
	if err != nil {
		return gen.PostPvzPvzIdDeleteLastProduct400JSONResponse{
			Message: err.Error(),
		}, err
	}

	return gen.PostPvzPvzIdDeleteLastProduct200Response{}, nil
}

// (POST /receptions).
func (s *Server) PostReceptions(
	ctx context.Context,
	request gen.PostReceptionsRequestObject,
) (gen.PostReceptionsResponseObject, error) {
	pvzId := request.Body.PvzId

	rec, err := s.reception.Create(ctx, domain.PVZID(pvzId))
	if err != nil {
		return gen.PostReceptions400JSONResponse{
			Message: err.Error(),
		}, err
	}

	return gen.PostReceptions201JSONResponse(rec.ToDTO()), nil
}

// (POST /register).
func (s *Server) PostRegister(
	ctx context.Context,
	request gen.PostRegisterRequestObject,
) (gen.PostRegisterResponseObject, error) {
	email := request.Body.Email
	password := request.Body.Password
	role := request.Body.Role

	user, err := s.user.Create(ctx, domain.Email(email), password, domain.Role(role))
	if err != nil {
		return gen.PostRegister400JSONResponse{
			Message: err.Error(),
		}, err
	}

	return gen.PostRegister201JSONResponse(*user.ToDto()), nil
}

func NewServer(
	jwt JWTGenerator,
	user UserProvider,
	pvz PVZProvider,
	reception ReceptionProvider,
	product ProductProvider,
) *Server {
	return &Server{
		jwt:       jwt,
		user:      user,
		pvz:       pvz,
		reception: reception,
		product:   product,
	}
}
