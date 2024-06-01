package crypto

import (
	"time"

	"chat-demo-golang/configs/middleware"
	"chat-demo-golang/shared/log"
)

type UserTokenData struct {
	ID        string
	CreatedAt time.Time
}

func (u *UserTokenData) TimeStamp() {
	u.CreatedAt = time.Now()
}

func GenerateAuthToken(tokenData UserTokenData) string {
	tokenData.TimeStamp()
	token, err := middleware.GenerateToken(&tokenData)
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
	}
	return token
}
