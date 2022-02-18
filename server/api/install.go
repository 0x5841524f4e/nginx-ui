package api

import (
	model2 "github.com/0xJacky/Nginx-UI/server/model"
	"github.com/0xJacky/Nginx-UI/server/settings"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

func installLockStatus() bool {
	lockPath := settings.ConfPath
	_, err := os.Stat(lockPath)

	return !os.IsNotExist(err)

}

func InstallLockCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"lock": installLockStatus(),
	})
}

type InstallJson struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,max=255"`
	Password string `json:"password" binding:"required,max=255"`
}

func InstallNginxUI(c *gin.Context) {
	// 安装过就别访问了
	if installLockStatus() {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "installed",
		})
		return
	}
	var json InstallJson
	ok := BindAndValid(c, &json)
	if !ok {
		return
	}

	serverSettings := settings.Conf.Section("server")
	serverSettings.Key("JwtSecret").SetValue(uuid.New().String())
	serverSettings.Key("Email").SetValue(json.Email)
	err := settings.Save()
	if err != nil {
		ErrHandler(c, err)
		return
	}

	curd := model2.NewCurd(&model2.Auth{})
	pwd, _ := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	err = curd.Add(&model2.Auth{
		Name:     json.Username,
		Password: string(pwd),
	})
	if err != nil {
		ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})

}
