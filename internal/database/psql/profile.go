package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Profile struct {
	id        int            `db:"u.id"`
	email     string         `db:"u.email"`
	login     string         `db:"u.login"`
	createdAt time.Time      `db:"u.created_at"`
	confirmed bool           `db:"u.confirmed"`
	name      string         `db:"up.name"`
	lastName  string         `db:"up.last_name"`
	phone     sql.NullString `db:"up.phone"`
	services  []byte         `db:"services"`
	sessions  []byte         `db:"sessions"`
	photo     string         `db:"up.photo_url"`
}

type Services struct {
	ServiceID   int       `json:"service_id"`
	ServiceName string    `json:"service_name"`
	Active      bool      `json:"active"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type Session struct {
	IP        string    `json:"ip"`
	Device    string    `json:"device"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

func (s *DatabaseProvider) GetProfile(ctx context.Context, userId int) (*sso.ProfileResponse, error) {
	var profile Profile

	query := `
	SELECT 
		u.id, 
		u.email, 
		u.login, 
		u.created_at, 
		u.confirmed,
		up.name, 
		up.last_name, 
		up.phone,
		json_agg(
			DISTINCT jsonb_build_object(
				'service_id', p.service_id,
				'service_name', s.name,
				'active', p.active,
				'expires_at', (p.expires_at AT TIME ZONE 'UTC')::timestamptz
			)
		) AS permissions,
		json_agg(
			DISTINCT jsonb_build_object(
				'ip', at.ip,
				'device', at.device,
				'created_at', (at.created_at AT TIME ZONE 'UTC')::timestamptz,
				'is_active', at.is_active
			)
		) AS sessions,
		up.photo_url
	FROM users u
	LEFT JOIN users_profile up ON u.id = up.user_id
	LEFT JOIN permission p ON u.id = p.user_id
	LEFT JOIN services s ON p.service_id = s.id
	LEFT JOIN active_tokens at ON u.id = at.user_id
	WHERE u.id = $1
	GROUP BY u.id, u.email, u.login, u.created_at, u.confirmed, 
			up.name, up.last_name, up.phone, up.photo_url;
	`
	row := s.db.QueryRowContext(ctx, query, userId)
	if err := row.Scan(
		&profile.id,
		&profile.email,
		&profile.login,
		&profile.createdAt,
		&profile.confirmed,
		&profile.name,
		&profile.lastName,
		&profile.phone,
		&profile.services,
		&profile.sessions,
		&profile.photo,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "profile not found")
		}
		message := fmt.Sprintf("failed to get profile: %v", err)
		s.log.ErrorContext(
			ctx,
			message,
			"method", "GetProfile",
			"point", "get.profile",
			"message", message,
		)
		return nil, status.Error(codes.Internal, "database error")
	}

	var services []*Services
	if err := json.Unmarshal(profile.services, &services); err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Internal, "server error")
	}

	var servicesResponse []*sso.Services
	for _, service := range services {
		servicesResponse = append(servicesResponse, &sso.Services{
			Id:        int32(service.ServiceID),
			Name:      service.ServiceName,
			Active:    service.Active,
			ExpiresAt: timestamppb.New(service.ExpiresAt),
		})
	}

	var sessions []*Session
	if err := json.Unmarshal(profile.sessions, &sessions); err != nil {
		return nil, status.Error(codes.Internal, "server error")
	}

	var sessionsResponse []*sso.Session
	for _, session := range sessions {
		sessionsResponse = append(sessionsResponse, &sso.Session{
			Ip:       session.IP,
			Device:   session.Device,
			Created:  timestamppb.New(session.CreatedAt),
			IsActive: session.IsActive,
		})
	}

	return &sso.ProfileResponse{
		Id:        int64(profile.id),
		Email:     profile.email,
		Login:     profile.login,
		CreatedAt: timestamppb.New(profile.createdAt),
		Confirmed: profile.confirmed,
		Name:      profile.name,
		LastName:  profile.lastName,
		Services:  servicesResponse,
		Sessions:  sessionsResponse,
		Photo:     profile.photo,
	}, nil
}

func (s *DatabaseProvider) GetMiniProfile(ctx context.Context, userId int) (*sso.MiniProfileResponse, error) {
	var profile Profile

	query := `
		SELECT 
			u.confirmed,
			up.name,
			up.photo_url,
			u.id
		FROM users u
		LEFT JOIN users_profile up ON u.id = up.user_id
		WHERE u.id = $1
	`
	row := s.db.QueryRowContext(ctx, query, userId)
	if err := row.Scan(
		&profile.confirmed,
		&profile.name,
		&profile.photo,
		&profile.id,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "profile not found")
		}
		message := fmt.Sprintf("failed to get mini profile: %v", err)
		s.log.ErrorContext(
			ctx,
			message,
			"method", "GetMiniProfile",
			"point", "get.mini.profile",
			"message", message,
		)
		return nil, status.Error(codes.Internal, "database error")
	}

	return &sso.MiniProfileResponse{
		Confirmed: profile.confirmed,
		Name:      profile.name,
		Photo:     profile.photo,
		Id:        int64(profile.id),
	}, nil
}
