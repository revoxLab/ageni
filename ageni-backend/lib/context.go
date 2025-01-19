package lib

import (
	"context"
	"github.com/readonme/open-studio/dal/model"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/readonme/open-studio/common/log"
	"google.golang.org/grpc/metadata"
)

const (
	TokenKey       = "authorization"
	KeyUid         = "__uid__"
	KeyDevice      = "__device__"
	KeyIp          = "__ip__"
	KeyUserProfile = "__profile__"
)

func GetTokenFromGinContext(c *gin.Context) string {
	return c.Request.Header.Get(strings.ToLower(TokenKey))
}

func RequestContext(c *gin.Context) context.Context {
	m := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
		TokenKey, GetTokenFromGinContext(c),
	))
	return m
}

func SetUserProfile(c *gin.Context, profile *model.User) {
	c.Set(KeyUserProfile, profile)
}

func SetContextUid(c *gin.Context, uid int64) {
	c.Set(KeyUid, uid)
}

func GetContextUid(c *gin.Context) int64 {
	if i, ok := c.Get(KeyUid); ok {
		return i.(int64)
	}
	return 0
}

func GetContextDevice(c *gin.Context) string {
	if i, ok := c.Get(KeyDevice); ok {
		return i.(string)
	}
	return ""
}

func GetContextIp(c *gin.Context) string {
	if i, ok := c.Get(KeyIp); ok {
		return i.(string)
	}
	return ""
}

func GetUserProfile(c *gin.Context) *model.User {
	if i, ok := c.Get(KeyUserProfile); ok {
		return i.(*model.User)
	}
	return &model.User{}
}
func GetTokenFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	token, ok := md[TokenKey]
	if !ok {
		return ""
	}

	return strings.Join(token, "")
}

func GetContextTraceId(c *gin.Context) string {
	if i, ok := c.Get(log.TraceIdKey); ok {
		return i.(string)
	}
	return ""
}
