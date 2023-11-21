package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"nunu_ginblog/cmd/server/wire"
	"nunu_ginblog/pkg/config"
	"nunu_ginblog/pkg/log"
)

// @title                      Nunu Example API
// @version                    1.0.0
// @description                This is a sample server celler server.
// @termsOfService             http://swagger.io/terms/
// @contact.name               API Support
// @contact.url                http://www.swagger.io/support
// @contact.email              support@swagger.io
// @license.name               Apache 2.0
// @license.url                http://www.apache.org/licenses/LICENSE-2.0.html
// @host                       localhost:8000
// @securityDefinitions.apiKey Bearer
// @in                         header
// @name                       Authorization
// @externalDocs.description   OpenAPI
// @externalDocs.url           https://swagger.io/resources/open-api/
func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger := log.NewLog(conf)

	app, cleanup, err := wire.NewWire(conf, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	logger.Info("server start", zap.String("host", "http://127.0.0.1:"+conf.GetString("http.port")))
	logger.Info("docs addr", zap.String("addr", fmt.Sprintf("http://127.0.0.1:%d/swagger/index.html", conf.GetInt("http.port"))))
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
