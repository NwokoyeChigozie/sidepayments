package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

const (
	ApiType       AuthorizationType = "api"
	AppType       AuthorizationType = "app"
	AuthType      AuthorizationType = "auth"
	BusinessAdmin AuthorizationType = "business_admin"
	Business      AuthorizationType = "business"
)

type (
	AuthorizationType  string
	AuthorizationTypes []AuthorizationType
)

func Authorize(db postgresql.Databases, extReq request.ExternalRequest, authTypes ...AuthorizationType) gin.HandlerFunc {

	return func(c *gin.Context) {
		if len(authTypes) > 0 {

			msg := ""
			for _, v := range authTypes {
				ms, status := v.ValidateAuthorizationRequest(c, db, extReq)
				if status {
					return
				}
				msg = ms
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.UnauthorisedResponse(http.StatusUnauthorized, fmt.Sprint(http.StatusUnauthorized), "Unauthorized", msg))
		}
	}
}

func (at AuthorizationType) in(authTypes AuthorizationTypes) bool {
	for _, v := range authTypes {
		if v == at {
			return true
		}
	}
	return false
}

func (at AuthorizationType) ValidateAuthorizationRequest(c *gin.Context, db postgresql.Databases, extReq request.ExternalRequest) (string, bool) {
	if at == ApiType {
		return at.ValidateApiType(c, extReq)
	} else if at == AppType {
		return at.ValidateAppType(c, extReq)
	} else if at == AuthType {
		return at.ValidateAuthType(c, extReq)
	} else if at == BusinessAdmin {
		return at.ValidateBusinessAdminType(c, extReq)
	} else if at == Business {
		return at.ValidateBusinessType(c, extReq)
	}

	return "authorized", true
}

func (at AuthorizationType) ValidateAuthType(c *gin.Context, extReq request.ExternalRequest) (string, bool) {

	var invalidToken = "Your request was made with invalid credentials."
	authorizationToken := GetHeader(c, "Authorization")
	if authorizationToken == "" {
		return "token not provided", false
	}

	bearerTokenArr := strings.Split(authorizationToken, " ")
	if len(bearerTokenArr) != 2 {
		return invalidToken, false
	}

	bearerToken := bearerTokenArr[1]

	if bearerToken == "" {
		return invalidToken, false
	}

	reqInf, err := extReq.SendExternalRequest(request.ValidateAuthorization, external_models.ValidateAuthorizationReq{
		Type:               string(AuthType),
		AuthorizationToken: bearerToken,
	})
	if err != nil {
		return err.Error(), false
	}

	dataResponse := reqInf.(external_models.ValidateAuthorizationDataModel)
	if !dataResponse.Status {
		return dataResponse.Message, false
	}

	models.MyIdentity = &dataResponse.Data
	return "authorized", true
}

func (at AuthorizationType) ValidateBusinessType(c *gin.Context, extReq request.ExternalRequest) (string, bool) {
	privateKey, publicKey, msg, status := at.getAccessTokens(c)
	if !status {
		return msg, false
	}
	_, err := extReq.SendExternalRequest(request.ValidateAuthorization, external_models.ValidateAuthorizationReq{
		Type:        string(Business),
		VPrivateKey: privateKey,
		VPublicKey:  publicKey,
	})
	if err != nil {
		return err.Error(), false
	}
	return "authorized", true
}

func (at AuthorizationType) ValidateAppType(c *gin.Context, extReq request.ExternalRequest) (string, bool) {
	config := config.GetConfig().App
	appKey := GetHeader(c, "v-app")
	if appKey == "" {
		return "missing app key", false
	}

	if appKey != config.Key {
		return "invalid app key", false
	}

	return "authorized", true
}

func (at AuthorizationType) ValidateBusinessAdminType(c *gin.Context, extReq request.ExternalRequest) (string, bool) {
	privateKey, publicKey, msg, status := at.getAccessTokens(c)
	if !status {
		return msg, false
	}
	_, err := extReq.SendExternalRequest(request.ValidateAuthorization, external_models.ValidateAuthorizationReq{
		Type:        string(BusinessAdmin),
		VPrivateKey: privateKey,
		VPublicKey:  publicKey,
	})
	if err != nil {
		return err.Error(), false
	}
	return "authorized", true
}

func (at AuthorizationType) ValidateApiType(c *gin.Context, extReq request.ExternalRequest) (string, bool) {
	privateKey, publicKey, msg, status := at.getAccessTokens(c)
	if !status {
		return msg, false
	}
	_, err := extReq.SendExternalRequest(request.ValidateAuthorization, external_models.ValidateAuthorizationReq{
		Type:        string(ApiType),
		VPrivateKey: privateKey,
		VPublicKey:  publicKey,
	})
	if err != nil {
		return err.Error(), false
	}
	return msg, status
}

func (at AuthorizationType) getAccessTokens(c *gin.Context) (string, string, string, bool) {
	privateKey := GetHeader(c, "v-private-key")
	publicKey := GetHeader(c, "v-public-key")

	if privateKey == "" && publicKey == "" {
		return privateKey, publicKey, "missing api keys", false
	}

	if privateKey == "" || publicKey == "" {
		return privateKey, publicKey, "either public or private key is missing", false
	}

	return privateKey, publicKey, "authorized", true
}

func GetHeader(c *gin.Context, key string) string {
	header := ""
	if c.GetHeader(key) != "" {
		header = c.GetHeader(key)
	} else if c.GetHeader(strings.ToLower(key)) != "" {
		header = c.GetHeader(strings.ToLower(key))
	} else if c.GetHeader(strings.ToUpper(key)) != "" {
		header = c.GetHeader(strings.ToUpper(key))
	} else if c.GetHeader(strings.Title(key)) != "" {
		header = c.GetHeader(strings.Title(key))
	}
	return header
}
