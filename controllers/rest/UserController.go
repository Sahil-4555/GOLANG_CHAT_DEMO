package rest

import (
	"net/http"

	service "chat-demo-golang/services"
	"chat-demo-golang/shared/common"
	"chat-demo-golang/shared/log"
	"chat-demo-golang/shared/message"
	validators "chat-demo-golang/validators"

	"github.com/gin-gonic/gin"
)

// SignUp
// @Summary SignUp - This API is used to SignUp the user.
// @tags UserMethod
// @security BearerAuth
// @Param SignupRequest body common.SignUpReq true "signup request"
// @router /v1/api/signup [post]
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

// SignIn
// @Summary SignIn - This API is used to SignIn the user.
// @tags UserMethod
// @security BearerAuth
// @Param SignupRequest body common.SignInReq true "signin request"
// @router /v1/api/login [post]
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
