package auth

import (
	"context"
	"fmt"

	"github.com/loongkirin/gdk/database/model"
	"github.com/loongkirin/gdk/database/query"
	"github.com/loongkirin/gdk/database/repository"
	"github.com/loongkirin/gdk/net/http/request"
	"github.com/loongkirin/gdk/net/http/response"
	"github.com/loongkirin/gdk/oauth"
	"github.com/loongkirin/gdk/util"
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

	// fmt.Println("login user", user)
	if !util.BcryptVerify(req.Data.Password, user.Password) {
		return nil, ErrPasswordInvalid
	}
	tenantId := user.TenantId
	// fmt.Println("login tenantId", tenantId)
	accessToken, _, authErr := s.oauthMaker.GenerateAccessToken(user.Id, user.Email, user.Phone, user.Name)
	// fmt.Println("login accessToken", accessToken)
	if authErr != nil {
		return nil, authErr
	}
	refreshToken, claims, authErr := s.oauthMaker.GenerateRefreshToken(user.Id, user.Email, user.Phone, user.Name)
	// fmt.Println("login refreshToken", refreshToken)
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
		TenantBaseModel: model.NewTenantBaseModel(tenantId, claims.Id),
	}

	session, err = s.oauthSessionRepo.Add(ctx, session)
	if err != nil {
		fmt.Println("login session error", err)
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
	if err != nil && err != ErrTenantNotFound {
		return nil, err
	}
	if tenant == nil {
		tenant = &Tenant{
			Name:        req.Data.TenantName,
			DbBaseModel: model.NewDbBaseModel(util.GenerateId()),
		}
		tenant, err = s.tenantRepo.Add(ctx, tenant)
		if err != nil {
			return nil, err
		}
	}
	password, err := util.BcryptHash(req.Data.Password)
	if err != nil {
		return nil, err
	}

	dbUser, err := s.findUserByPhone(ctx, req.Data.Phone)
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}
	if dbUser != nil {
		return nil, ErrUserExists
	}

	user := &User{
		Name:            req.Data.UserName,
		Phone:           req.Data.Phone,
		Email:           req.Data.Email,
		Password:        password,
		TenantBaseModel: model.NewTenantBaseModel(tenant.Id, util.GenerateId()),
	}
	user, err = s.userRepo.Add(ctx, user)
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
		return nil, ErrPhoneRequired
	}
	wheres := []query.DbQueryWhere{}
	filters := []query.DbQueryFilter{query.NewDbQueryFilter("phone", []interface{}{phone}, query.EQ, "String")}
	wheres = append(wheres, query.NewDbQueryWhere(filters, query.AND))
	query := &query.DbQuery{
		QueryWheres: wheres,
		PageSize:    10,
		PageNumber:  1,
	}
	users, err := s.userRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return &users[0], nil
}

func (s *service) findTenantByName(ctx context.Context, name string) (*Tenant, error) {
	if len(name) == 0 {
		return nil, ErrTenantNameRequired
	}
	wheres := []query.DbQueryWhere{}
	filters := []query.DbQueryFilter{query.NewDbQueryFilter("name", []interface{}{name}, query.EQ, "String")}
	wheres = append(wheres, query.NewDbQueryWhere(filters, query.AND))
	query := &query.DbQuery{
		QueryWheres: wheres,
		PageSize:    10,
		PageNumber:  1,
	}
	tenants, err := s.tenantRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(tenants) == 0 {
		return nil, ErrTenantNotFound
	}
	return &tenants[0], nil
}
