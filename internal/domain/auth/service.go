package auth

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/loongkirin/go-family-finance/pkg/database"
	"github.com/loongkirin/go-family-finance/pkg/database/repository"
	"github.com/loongkirin/go-family-finance/pkg/http/request"
	"github.com/loongkirin/go-family-finance/pkg/http/response"
	"github.com/loongkirin/go-family-finance/pkg/util"
)

type Service interface {
	Create(ctx context.Context, req *request.DataRequest[CreateUserDTO]) (*response.DataResponse[UserDTO], error)
	Update(ctx context.Context, req *request.DataRequest[UpdateUserDTO]) (*response.DataResponse[UserDTO], error)
	Delete(ctx context.Context, id string) (bool, error)
	FindById(ctx context.Context, id string) (*response.DataResponse[UserDTO], error)
	Query(ctx context.Context, req *request.Query) (*response.DataListResponse[UserDTO], error)
}

type service struct {
	repo repository.Repository[User]
}

func NewService(repo repository.Repository[User]) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, req *request.DataRequest[CreateUserDTO]) (*response.DataResponse[UserDTO], error) {
	tenantId := ctx.Value("tenant_id").(string)
	user := &User{
		Name:            req.Data.Name,
		Mobile:          req.Data.Mobile,
		Password:        req.Data.Password,
		Active:          true,
		TenantBaseModel: database.NewTenantBaseModel(tenantId, util.GenerateId()),
	}

	user, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &response.DataResponse[UserDTO]{
		Data: UserDTO{
			Id:     user.Id,
			Name:   user.Name,
			Mobile: user.Mobile,
		},
	}, nil
}

func (s *service) Update(ctx context.Context, req *request.DataRequest[UpdateUserDTO]) (*response.DataResponse[UserDTO], error) {
	user := &User{
		Name:            req.Data.Name,
		Mobile:          req.Data.Mobile,
		Password:        req.Data.Password,
		Active:          req.Data.Active,
		TenantBaseModel: database.NewTenantBaseModel(req.Data.TenantId, req.Data.Id),
	}

	user, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &response.DataResponse[UserDTO]{
		Data: UserDTO{
			Id:     user.Id,
			Name:   user.Name,
			Mobile: user.Mobile,
		},
	}, nil
}

func (s *service) Delete(ctx context.Context, id string) (bool, error) {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindById(ctx context.Context, id string) (*response.DataResponse[UserDTO], error) {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &response.DataResponse[UserDTO]{
		Data: UserDTO{
			Id:     user.Id,
			Name:   user.Name,
			Mobile: user.Mobile,
		},
	}, nil
}

func (s *service) Query(ctx context.Context, req *request.Query) (*response.DataListResponse[UserDTO], error) {
	dbQuery := &database.DbQuery{}
	// err := copier.Copy(dbQuery, req)
	err := copier.CopyWithOption(dbQuery, req, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, err
	}
	users, err := s.repo.Query(ctx, dbQuery)
	if err != nil {
		return nil, err
	}
	userDTOs := make([]UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = UserDTO{
			Id:     user.Id,
			Name:   user.Name,
			Mobile: user.Mobile,
		}
	}
	return &response.DataListResponse[UserDTO]{
		DataList: userDTOs,
	}, nil
}
