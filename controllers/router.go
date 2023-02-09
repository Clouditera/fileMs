package controllers

import (
	"fileMS/controllers/api"
	"fileMS/pkg/config"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"github.com/mlogclub/simple/web"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func Router() {
	app := iris.New()
	app.Logger().SetLevel("warn")
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
		MaxAge:           600,
		AllowedMethods:   []string{iris.MethodGet, iris.MethodPost, iris.MethodOptions, iris.MethodHead, iris.MethodDelete, iris.MethodPut},
		AllowedHeaders:   []string{"*"},
	}))
	app.AllowMethods(iris.MethodOptions)

	app.OnAnyErrorCode(func(ctx iris.Context) {
		path := ctx.Path()
		var err error
		if strings.Contains(path, "/api/admin/") {
			err = ctx.JSON(web.JsonErrorCode(ctx.GetStatusCode(), "Http error"))
		}
		if err != nil {
			logrus.Error(err)
		}
	})

	app.Any("/", func(i iris.Context) {
		_, _ = i.HTML("<h1>Powered by fileMS</h1>")
	})

	// api
	mvc.Configure(app.Party("/api"), func(m *mvc.Application) {
		m.Party("/file").Handle(new(api.FileController)) // 文件管理
	})

	server := &http.Server{Addr: ":" + config.Instance.Port}
	handleSignal(server)
	err := app.Run(iris.Server(server), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 true,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		EnableOptimizations:               true,
		TimeFormat:                        "2006-01-02 15:04:05",
		Charset:                           "UTF-8",
	}))
	if err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
}

func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		logrus.Infof("got signal [%s], exiting now", s)
		if err := server.Close(); nil != err {
			logrus.Errorf("server close failed: " + err.Error())
		}

		// TODO close minio

		logrus.Infof("Exited")
		os.Exit(0)
	}()
}
