package crypto

import (
	"time"

	"github.com/Sahil-4555/mvc/configs/middleware"
	"github.com/Sahil-4555/mvc/shared/log"
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
