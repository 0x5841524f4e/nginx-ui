package router

import (
	"bufio"
	"github.com/0xJacky/Nginx-UI/server/api"
	"github.com/0xJacky/Nginx-UI/server/settings"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	r.Use(recovery())

	r.Use(cacheJs())

	r.Use(static.Serve("/", mustFS("")))

	r.NoRoute(func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept")
		if strings.Contains(accept, "text/html") {
			file, _ := mustFS("").Open("index.html")
			defer file.Close()
			stat, _ := file.Stat()
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html",
				bufio.NewReader(file), nil)
			return
		}
	})

	root := r.Group("/api")
	{

		root.GET("settings", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"demo": settings.ServerSettings.Demo,
			})
		})

		root.GET("install", api.InstallLockCheck)
		root.POST("install", api.InstallNginxUI)

		root.POST("/login", api.Login)
		root.DELETE("/logout", api.Logout)

		g := root.Group("/", authRequired())
		{
			g.GET("analytic", api.Analytic)
			g.GET("analytic/init", api.GetAnalyticInit)

			g.GET("users", api.GetUsers)
			g.GET("user/:id", api.GetUser)
			g.POST("user", api.AddUser)
			g.POST("user/:id", api.EditUser)
			g.DELETE("user/:id", api.DeleteUser)

			g.GET("domains", api.GetDomains)
			g.GET("domain/:name", api.GetDomain)

			// Modify site configuration directly
			g.POST("domain/:name", api.EditDomain)

			// Transform NgxConf to nginx configuration
			g.POST("ngx/build_config", api.BuildNginxConfig)
			// Tokenized nginx configuration to NgxConf
			g.POST("ngx/tokenize_config", api.TokenizeNginxConfig)

			g.POST("domain/:name/enable", api.EnableDomain)
			g.POST("domain/:name/disable", api.DisableDomain)
			g.DELETE("domain/:name", api.DeleteDomain)

			g.GET("configs", api.GetConfigs)
			g.GET("config/:name", api.GetConfig)
			g.POST("config", api.AddConfig)
			g.POST("config/:name", api.EditConfig)

			//g.GET("backups", api.GetFileBackupList)
			//g.GET("backup/:id", api.GetFileBackup)

			g.GET("template", api.GetTemplate)

			g.GET("cert/issue/:domain", api.IssueCert)

			// Add domain to auto-renew cert list
			g.POST("cert/:domain", api.AddDomainToAutoCert)
			// Delete domain from auto-renew cert list
			g.DELETE("cert/:domain", api.RemoveDomainFromAutoCert)

			// pty
			g.GET("pty", api.Pty)

			// Nginx log
			g.GET("nginx_log", api.NginxLog)
		}
	}

	return r
}
