package database

import (
	"context"
	"database/sql"

	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *DatabaseProvider) GetServicesList(ctx context.Context) ([]*sso.StructureServices, error) {
	var (
		services   []*sso.StructureServices
		methodName = "GetServices"
	)

	q := `
		SELECT id, name, created_at
		FROM services
	`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "service not found")
		}
		errMessage := status.Errorf(codes.Internal, "failed to query %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return nil, status.Error(codes.Internal, "database error")
	}
	defer rows.Close()

	for rows.Next() {
		var service ServicesList
		if err := rows.Scan(&service.Id, &service.Name, &service.CreatedAt); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan %v", err)
		}
		services = append(services, &sso.StructureServices{
			Id:        service.Id,
			Name:      service.Name,
			CreatedAt: timestamppb.New(service.CreatedAt),
		})
	}

	return services, nil
}

func (s *DatabaseProvider) GetServicesByName(ctx context.Context, name string) ([]*sso.StructureServices, error) {
	var (
		services   []*sso.StructureServices
		methodName = "GetServicesByName"
	)

	q := `
		SELECT id, name, created_at
		FROM services
		WHERE name LIKE '%' || $1 || '%'
	`

	rows, err := s.db.QueryContext(ctx, q, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "service %v not found", name)
		}
		errMessage := status.Errorf(codes.Internal, "failed to query %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return nil, status.Error(codes.Internal, "database error")
	}
	defer rows.Close()

	for rows.Next() {
		var service ServicesList
		if err := rows.Scan(&service.Id, &service.Name, &service.CreatedAt); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan %v", err)
		}
		services = append(services, &sso.StructureServices{
			Id:        service.Id,
			Name:      service.Name,
			CreatedAt: timestamppb.New(service.CreatedAt),
		})
	}

	return services, nil
}

func (s *DatabaseProvider) GetServiceById(ctx context.Context, id int32) (*sso.StructureServices, error) {
	var (
		service    ServicesList
		methodName = "GetServiceById"
	)

	q := `
		SELECT id, name, created_at
		FROM services
		WHERE id = $1
	`
	rows, err := s.db.QueryContext(ctx, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "service %v not found", id)
		}
		errMessage := status.Errorf(codes.Internal, "failed to query %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return nil, status.Error(codes.Internal, "database error")
	}

	for rows.Next() {
		if err := rows.Scan(&service.Id, &service.Name, &service.CreatedAt); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan %v", err)
		}
	}

	return &sso.StructureServices{
		Id:        service.Id,
		Name:      service.Name,
		CreatedAt: timestamppb.New(service.CreatedAt),
	}, nil
}
