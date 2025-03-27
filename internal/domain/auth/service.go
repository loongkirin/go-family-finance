package auth

import (
	"context"
	"errors"

	"github.com/loongkirin/go-family-finance/pkg/database"
	"github.com/loongkirin/go-family-finance/pkg/database/repository"
	"github.com/loongkirin/go-family-finance/pkg/http/request"
	"github.com/loongkirin/go-family-finance/pkg/http/response"
	"github.com/loongkirin/go-family-finance/pkg/oauth"
	"github.com/loongkirin/go-family-finance/pkg/util"
)

type AuthService interface {
	Register(ctx context.Context, req *request.DataRequest[RegisterDTO]) (*response.DataResponse[UserDTO], error)
	Login(ctx context.Context, req *request.DataRequest[LoginDTO]) (*response.DataResponse[UserDTO], error)
}

type service struct {
	userRepo         repository.Repository[User]
	oauthSessionRepo repository.Repository[OAuthSession]
	tenantRepo       repository.Repository[Tenant]
	oauthMaker       oauth.OAuthMaker
}

func NewAuthService(
	userRepo repository.Repository[User],
	oauthSessionRepo repository.Repository[OAuthSession],
	tenantRepo repository.Repository[Tenant],
	oauthMaker oauth.OAuthMaker,
) AuthService {
	return &service{
		userRepo:         userRepo,
		oauthSessionRepo: oauthSessionRepo,
		tenantRepo:       tenantRepo,
		oauthMaker:       oauthMaker,
	}
}

func (s *service) Login(ctx context.Context, req *request.DataRequest[LoginDTO]) (*response.DataResponse[UserDTO], error) {
	user, err := s.findUserByPhone(ctx, req.Data.Phone)
	if err != nil {
		return nil, err
	}
	if !util.BcryptVerify(req.Data.Password, user.Password) {
		return nil, errors.New("密码错误")
	}
	tenantId := ""
	accessToken, _, authErr := s.oauthMaker.GenerateAccessToken(user.Id, user.Email, user.Phone, user.Name)
	if authErr != nil {
		return nil, authErr
	}

	refreshToken, claims, authErr := s.oauthMaker.GenerateRefreshToken(user.Id, user.Email, user.Phone, user.Name)
	if authErr != nil {
		return nil, authErr
	}

	session := &OAuthSession{
		UserId:          user.Id,
		Phone:           user.Phone,
		UserName:        user.Name,
		UserAgent:       ctx.Value("user_agent").(string),
		ClientIp:        ctx.Value("client_ip").(string),
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		ExpiredAt:       claims.ExpiredAt.UnixMilli(),
		TenantBaseModel: database.NewTenantBaseModel(tenantId, claims.Id),
	}
	session, err = s.oauthSessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}
	return &response.DataResponse[UserDTO]{
		Data: UserDTO{
			UserId:   user.Id,
			UserName: user.Name,
			Phone:    user.Phone,
			Email:    user.Email,
			TenantDTO: TenantDTO{
				TenantId: tenantId,
			},
			OAuthDTO: OAuthDTO{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				ExpiredAt:    claims.ExpiredAt.UnixMilli(),
				SessionId:    session.Id,
			},
		},
	}, nil
}

func (s *service) Register(ctx context.Context, req *request.DataRequest[RegisterDTO]) (*response.DataResponse[UserDTO], error) {
	tenant, err := s.findTenantByName(ctx, req.Data.TenantName)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		tenant = &Tenant{
			Name:        req.Data.TenantName,
			DbBaseModel: database.NewDbBaseModel(util.GenerateId()),
		}
		tenant, err = s.tenantRepo.Create(ctx, tenant)
		if err != nil {
			return nil, err
		}
	}
	password, err := util.BcryptHash(req.Data.Password)
	if err != nil {
		return nil, err
	}

	dbUser, err := s.findUserByPhone(ctx, req.Data.Phone)
	if err != nil {
		return nil, err
	}
	if dbUser != nil {
		return nil, errors.New("用户已存在")
	}

	user := &User{
		Name:            req.Data.UserName,
		Phone:           req.Data.Phone,
		Email:           req.Data.Email,
		Password:        password,
		TenantBaseModel: database.NewTenantBaseModel(tenant.Id, util.GenerateId()),
	}
	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &response.DataResponse[UserDTO]{
		Data: UserDTO{
			UserId:   user.Id,
			UserName: user.Name,
			Phone:    user.Phone,
			Email:    user.Email,
			TenantDTO: TenantDTO{
				TenantId:   tenant.Id,
				TenantName: tenant.Name,
			},
		},
	}, nil
}

func (s *service) findUserByPhone(ctx context.Context, phone string) (*User, error) {
	if len(phone) == 0 {
		return nil, errors.New("手机号不能为空")
	}
	wheres := []database.DbQueryWhere{}
	filters := []database.DbQueryFilter{database.NewDbQueryFilter("phone", []interface{}{phone}, database.EQ, "String")}
	wheres = append(wheres, database.NewDbQueryWhere(filters, database.AND))
	query := &database.DbQuery{
		QueryWheres: wheres,
		PageSize:    10,
		PageNumber:  1,
	}
	users, err := s.userRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("用户不存在")
	}
	return &users[0], nil
}

func (s *service) findTenantByName(ctx context.Context, name string) (*Tenant, error) {
	if len(name) == 0 {
		return nil, errors.New("租户名称不能为空")
	}
	wheres := []database.DbQueryWhere{}
	filters := []database.DbQueryFilter{database.NewDbQueryFilter("name", []interface{}{name}, database.EQ, "String")}
	wheres = append(wheres, database.NewDbQueryWhere(filters, database.AND))
	query := &database.DbQuery{
		QueryWheres: wheres,
		PageSize:    10,
		PageNumber:  1,
	}
	tenants, err := s.tenantRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(tenants) == 0 {
		return nil, nil
	}
	return &tenants[0], nil
}
