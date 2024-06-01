package rest

import (
	"net/http"

	service "github.com/Sahil-4555/mvc/services"
	"github.com/Sahil-4555/mvc/shared/common"
	"github.com/Sahil-4555/mvc/shared/log"
	"github.com/Sahil-4555/mvc/shared/message"
	validators "github.com/Sahil-4555/mvc/validators"
	"github.com/gin-gonic/gin"
)

// swagger:route POST /v1/api/signup
// Sign up for a new account.
//
// Responses:
//
//	200: SignupResponse
//	400: BadRequestError
func SignUp(c *gin.Context) {
	log.GetLog().Info("INFO : ", "User Controller Called(SignUp).")

	var req common.SignUpReq

	// Decode the request body into struct and failed if any error occur
	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	// struct field validation
	if resp, ok := validators.ValidateStruct(req, "SignUpRequest"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	// call service
	resp := service.SignUp(req)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Response_SignIn(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "SignUp successfully...")
}

func SignIn(c *gin.Context) {
	log.GetLog().Info("INFO : ", "User Controller Called(SignIn).")

	var req common.SignInReq

	// Decode the request body into struct and failed if any error occur
	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	// struct field validation
	if resp, ok := validators.ValidateStruct(req, "SignInRequest"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	resp := service.SignIn(req)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Response_SignIn(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "SignIn successfully...")
}
