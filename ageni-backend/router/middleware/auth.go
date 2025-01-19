package middleware

import (
	"github.com/readonme/open-studio/common"
	xjwt "github.com/readonme/open-studio/common/auth/jwt"
	"github.com/readonme/open-studio/conf"
	"github.com/readonme/open-studio/lib"
	"github.com/readonme/open-studio/user"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/readonme/open-studio/common/response"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if lib.GetTokenFromGinContext(c) == "" {
			if strings.Contains(c.FullPath(), "/studio/bot_list") ||
				strings.Contains(c.FullPath(), "/studio/bot_tabs") {
				c.Next()
				return
			}
		}

		token := lib.GetTokenFromGinContext(c)

		if token == "" {
			c.Abort()
			response.JSONFail(c, common.ErrTokenEmpty, nil)
			return
		}

		tokenConf := conf.Conf.JWTToken
		userClaim := &xjwt.AccountClaims{}

		if err := xjwt.Jwt.Decode(token, tokenConf, userClaim); err != nil {
			c.Abort()
			response.JSONFail(c, err, nil)
			return
		}

		lib.SetContextUid(c, userClaim.Account.Id)
		model, err := user.GetUser(userClaim.Account.Id)
		if err != nil {
			c.Abort()
			response.JSONFail(c, err, nil)
			return
		}

		lib.SetUserProfile(c, model)
		lib.SetContextUid(c, model.ID)
		c.Next()
	}
}
