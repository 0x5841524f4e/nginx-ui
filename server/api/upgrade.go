package api

import (
	"github.com/0xJacky/Nginx-UI/server/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func GetRelease(c *gin.Context) {
	data, err := service.GetRelease()
	if err != nil {
		ErrHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetCurrentVersion(c *gin.Context) {
	curVer, err := service.GetCurrentVersion()
	if err != nil {
		ErrHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, curVer)
}

func PerformCoreUpgrade(c *gin.Context) {
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// upgrade http to websocket
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("[Error] PerformCoreUpgrade Upgrade", err)
		return
	}
	defer ws.Close()

	_ = ws.WriteJSON(gin.H{
		"status":  "info",
		"message": "Initialing Core Upgrader",
	})

	u, err := service.NewUpgrader()

	if err != nil {
		_ = ws.WriteJSON(gin.H{
			"status":  "error",
			"message": "Initialing core upgrader error",
		})
		log.Println("[Error] PerformCoreUpgrade service.NewUpgrader()", err)
		return
	}
	_ = ws.WriteJSON(gin.H{
		"status":  "info",
		"message": "Downloading latest release",
	})
	tarName, err := u.DownloadLatestRelease()
	if err != nil {
		_ = ws.WriteJSON(gin.H{
			"status":  "error",
			"message": "Download latest release error",
		})
		log.Println("[Error] PerformCoreUpgrade DownloadLatestRelease", err)
		return
	}
	_ = ws.WriteJSON(gin.H{
		"status":  "info",
		"message": "Performing core upgrade",
	})
	_ = os.Remove(u.Release.ExPath)
	// bye, overseer will restart nginx-ui
	err = u.PerformCoreUpgrade(filepath.Dir(u.Release.ExPath), tarName)
	if err != nil {
		_ = ws.WriteJSON(gin.H{
			"status":  "error",
			"message": "Perform core upgrade error",
		})
		log.Println("[Error] PerformCoreUpgrade PerformCoreUpgrade", err)
		return
	}
}
