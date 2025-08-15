package authgo

import (
	"context"

	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
)

func (s *AuthData) GetServices(ctx context.Context, data *sso.ServicesRequest) ([]*sso.StructureServices, error) {

	switch {
	case data.GetName() != "":
		services, err := s.DB.Services.GetServicesByName(ctx, data.GetName())
		if err != nil {
			return nil, err
		}
		return services, nil
	case data.GetId() != 0:
		service, err := s.DB.Services.GetServiceById(ctx, data.GetId())
		if err != nil {
			return nil, err
		}
		return []*sso.StructureServices{service}, nil
	default:
		services, err := s.DB.Services.GetServicesList(ctx)
		if err != nil {
			return nil, err
		}
		return services, nil
	}
}
