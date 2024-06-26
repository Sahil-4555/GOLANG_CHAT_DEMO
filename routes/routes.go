package routes

import (
	"context"
	"fmt"
	"net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"chat-demo-golang/configs"
	"chat-demo-golang/configs/middleware"
	rest "chat-demo-golang/controllers/rest"
	"chat-demo-golang/controllers/ws"
	"chat-demo-golang/docs"
	"chat-demo-golang/shared/log"

	"github.com/gin-gonic/gin"
)

var server *http.Server

func Run() {
	port := configs.Port()
	log.GetLog().Info("", "Service listen on "+port)
	router := gin.New()
	server = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}
	SetupRoutes(router)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.GetLog().Fatal("", "listen: %s\n", err)
	}
}

func Close(ctx context.Context) error {
	if server != nil {
		return server.Shutdown(ctx)
	}
	return nil
}

func SetupRoutes(r *gin.Engine) {

	docs.SwaggerInfo.Title = "CHAT APPLICATION API"
	docs.SwaggerInfo.Description = "This is a sample doc of Chat Application API."
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.GinMiddleware())

	SetupDefaultEndpoints(r)

	public := r.Group("/v1/api")
	SetupUserRoutes(public)

	private := r.Group("/")
	private.Use(middleware.AuthHandler())
	SetupChannelRoutes(private)
	SetupChatRoutes(private)
	SetupMediaRoutes(private)

	wsServer := ws.NewWebSocketServer()
	go wsServer.Run()
	wsGroup := r.Group("/ws")
	wsGroup.Use(middleware.AuthHandlerWebsocket())
	wsGroup.GET("", func(c *gin.Context) {
		ws.ServeWs(wsServer, c.Writer, c.Request)
	})
}

func SetupDefaultEndpoints(r *gin.Engine) {

	r.GET("/", func(c *gin.Context) {
		htmlString := "<html><body>Welcome to CHAT APPLICATION!</body></html>"
		c.Writer.WriteHeader(http.StatusOK)
		_, err := c.Writer.Write([]byte(htmlString))
		if err != nil {
			return
		}
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func SetupUserRoutes(rg *gin.RouterGroup) {
	rg.POST("/signup", rest.SignUp)
	rg.POST("/login", rest.SignIn)
}

func SetupChannelRoutes(rg *gin.RouterGroup) {
	ChannelRoute := rg.Group("/v1/channel")
	{
		ChannelRoute.POST("/joinRoom", rest.JoinChannel)
		ChannelRoute.GET("/get-direct-channel", rest.GetOneToOneChannelsConnectedWithUser)
		ChannelRoute.GET("/get-group-channel", rest.GetPrivateChannelsConnectedWithUser)
		ChannelRoute.POST("/createGroup", rest.CreateGroup)
		ChannelRoute.POST("/addfavoritechannel", rest.AddRemoveFavouriteChannel)
		ChannelRoute.GET("/get-favourite-channel", rest.GetFavouritesChannel)
		ChannelRoute.PUT("/addmembers", rest.AddMembersToGroup)
		ChannelRoute.PUT("/leavechannel", rest.LeaveChannel)
		ChannelRoute.PUT("/removeuserbyadmin", rest.RemoveUserFromGroupByGroupAdmin)
		ChannelRoute.PUT("/give-admin-rights", rest.GiveAdminRightsToUser)
		ChannelRoute.PUT("/closeconversation", rest.CloseConversation)
		ChannelRoute.GET("/getrecentchannels", rest.GetRecentChannelsOfUsers)
		ChannelRoute.GET("/get-all-direct-channel", rest.GetAllOneToOneChannelConnectedWithUserInOrder)
		ChannelRoute.GET("/join-group/:channelId", rest.JoinGroupWithChannelId)
	}
}

func SetupChatRoutes(rg *gin.RouterGroup) {
	ChatRoute := rg.Group("/v1/chat")
	{
		ChatRoute.PUT("/updatemessage", rest.UpdateMessage)
		ChatRoute.PUT("/deletemessage", rest.DeleteMessage)
		ChatRoute.GET("/searchhandler", rest.SearchHandler)
		ChatRoute.GET("/messages-by-channelid/:id", rest.GetMessagesByChannelId)
		ChatRoute.GET("/get-all-users", rest.GetAllUsers)
		ChatRoute.GET("/get-channel-members/:channelId", rest.GetChannelMembers)
	}
}

func SetupMediaRoutes(rg *gin.RouterGroup) {
	MediaRoute := rg.Group("/v1/media")
	{
		MediaRoute.POST("/uploadmedia", rest.UploadMedia)
	}
}
