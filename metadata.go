package platform

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

type UserContext struct {
	UserID    string
	Email     string
	Roles     []string
	RequestID string
	TraceID   string
}

func UserFromIncoming(ctx context.Context) UserContext {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return UserContext{}
	}
	get := func(key string) string {
		vals := md.Get(key)
		if len(vals) == 0 {
			return ""
		}
		return vals[0]
	}
	rolesRaw := get(MetadataUserRoles)
	var roles []string
	if rolesRaw != "" {
		for _, r := range strings.Split(rolesRaw, ",") {
			if v := strings.TrimSpace(r); v != "" {
				roles = append(roles, v)
			}
		}
	}
	return UserContext{
		UserID:    get(MetadataUserID),
		Email:     get("x-user-email"),
		Roles:     roles,
		RequestID: get(MetadataRequestID),
		TraceID:   get(MetadataTraceID),
	}
}

func OutgoingMetadata(uc UserContext) metadata.MD {
	md := metadata.Pairs(
		MetadataUserID, uc.UserID,
		MetadataUserRoles, strings.Join(uc.Roles, ","),
		MetadataRequestID, uc.RequestID,
		MetadataTraceID, uc.TraceID,
		"x-user-email", uc.Email,
	)
	return md
}

func AppendOutgoing(ctx context.Context, uc UserContext) context.Context {
	return metadata.NewOutgoingContext(ctx, OutgoingMetadata(uc))
}

func HasRole(uc UserContext, role string) bool {
	for _, r := range uc.Roles {
		if r == role {
			return true
		}
	}
	return false
}
