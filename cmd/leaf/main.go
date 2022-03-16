package main

import (
	"github.com/imkuqin-zw/courier/internal/leaf"
	"github.com/imkuqin-zw/courier/pkg/app"
	"github.com/imkuqin-zw/courier/pkg/config"
)

func main() {
	app.Init(app.WithAppName("leaf"))
	start()
	app.WaitSignals()
}

func start() {
	snowflakeEnable := len(config.Get("dubbo.provider.services.SnowflakeUC").Bytes()) > 0
	segmentEnable := len(config.Get("dubbo.provider.services.SegmentUC").Bytes()) > 0
	if snowflakeEnable && segmentEnable {
		app.Start(app.WithProvidersFactory(leaf.NewAllProviderServices))
	} else if snowflakeEnable {
		app.Start(app.WithProvidersFactory(leaf.NewSnowflakeProviderServices))
	} else if segmentEnable {
		app.Start(app.WithProvidersFactory(leaf.NewSegmentProviderServices))
	}
}
