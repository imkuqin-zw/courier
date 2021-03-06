// Copyright 2022 The imkuqin-zw Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	config2 "dubbo.apache.org/dubbo-go/v3/config"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/imkuqin-zw/courier/pkg/config"
	"github.com/imkuqin-zw/courier/pkg/config/source"
	"github.com/imkuqin-zw/courier/pkg/config/source/env"
	"github.com/imkuqin-zw/courier/pkg/config/source/file"
	flagSource "github.com/imkuqin-zw/courier/pkg/config/source/flag"
	dubbo "github.com/imkuqin-zw/courier/pkg/config/source/remote/dubbov3"
)

var o *Options
var rc *config2.RootConfig

func initOpts(ops ...Option) {
	o = &Options{
		appName:           baseAppName,
		EnvKeyReplaceRule: map[string]map[string]string{},
	}
	for _, f := range ops {
		f(o)
	}
}

func loadEnvAndFlag() {
	fs := flag.NewFlagSet("config", flag.ExitOnError)
	fs.String("config_dir", "./conf", "default config root dir path")
	o.EnvKeyReplaceRule["dubbo"] = map[string]string{"_": "-"}
	sources := []source.Source{
		env.NewSource(
			env.WithPrefix("COURIER"),
			env.WithStrippedPrefix("COURIER"),
			env.WithKeyReplaceRule(o.EnvKeyReplaceRule),
		),
		flagSource.NewSource(flagSource.WithFlagSet(fs), flagSource.IncludeUnset(true)),
	}
	if err := config.Load(sources...); err != nil {
		panic(fmt.Sprintf("fault to load env and flag source: %s", err.Error()))
	}
}

func loadAppCfgFile() {
	if o.disableAppCfgFile {
		return
	}
	courierFile := filepath.Join(config.Get("config.dir").String("./conf"), "courier.yaml")
	_, err := os.Stat(courierFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(fmt.Sprintf("fault to read courier file config: %s", err.Error()))
	}
	fileSource := file.NewSource(file.WithPath(courierFile), file.WithWatch(false))
	if err := config.Load(fileSource); err != nil {
		panic(fmt.Sprintf("fault to load courier file source: %s", err.Error()))
	}
}

func loadDubboV3() {
	registerPOJO()
	initDubboV3RootConfig()
	initDubboV3ConfigCenter()
	if err := initDubboV3(); err != nil {
		panic("fault to init dubbo v3")
	}
}

func getAppName() string {
	if rc != nil && rc.Application != nil && rc.Application.Name != "" {
		return rc.Application.Name
	}
	return o.appName
}

func initDubboV3RootConfig() {
	dubboFile := filepath.Join(config.Get("config.dir").String("./conf"), "dubbo.yaml")
	_, err := os.Stat(dubboFile)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Sprintf("fault to read courier file config: %s", err.Error()))
		}
	} else {
		fileSource := file.NewSource(file.WithPath(dubboFile), file.WithWatch(false))
		if err := config.Load(fileSource); err != nil {
			panic(fmt.Sprintf("fault to load courier file source: %s", err.Error()))
		}
	}
	fmt.Println(string(config.Bytes()))
	rc = config2.NewRootConfigBuilder().Build()
	if err := config.Get("dubbo").Scan(rc); err != nil {
		panic(fmt.Sprintf("fault to init dubbo v3 root config: %s", err.Error()))
	}
	config2.SetRootConfig(*rc)
}

func initDubboV3ConfigCenter() {
	if err := rc.Logger.Init(); err != nil { // init default logger
		panic(fmt.Sprintf("fault to init dubbo v3 default logger: %s", err.Error()))
	}
	if rc.ConfigCenter.Protocol != "" {
		rc.ConfigCenter.DataId = fmt.Sprintf("%s.dubbo", getAppName())
		if err := rc.ConfigCenter.Init(rc); err != nil {
			panic(fmt.Sprintf("fault to init dubbo v3 dynamic config center: %s", err.Error()))
		}
		if err := rc.Logger.Init(); err != nil { // init logger using config from config center again
			panic(fmt.Sprintf("fault to init dubbo v3 config logger: %s", err.Error()))
		}
	}

	if o.disableAppCfgDynamic {
		return
	}
	if dubboV3Source := dubbo.NewDubbov3Source(fmt.Sprintf("%s.app", getAppName())); dubboV3Source != nil {
		if err := config.Load(dubboV3Source); err != nil {
			panic(fmt.Sprintf("fault to load dubbo v3 source: %s", err.Error()))
		}
	}
}

func registerPOJO() {
	hessian.RegisterPOJO(&common.MetadataInfo{})
	hessian.RegisterPOJO(&common.ServiceInfo{})
	hessian.RegisterPOJO(&common.URL{})
}

func initRouterConfig() error {
	routers := rc.Router
	if len(routers) > 0 {
		for _, r := range routers {
			if err := r.Init(); err != nil {
				return err
			}
		}
		rc.Router = routers
	}

	//chain.SetVSAndDRConfigByte(vsBytes, drBytes)
	return nil
}

func initDubboV3() error {
	if err := rc.Application.Init(); err != nil {
		return err
	}

	// init user define
	if err := rc.Custom.Init(); err != nil {
		return err
	}

	// init protocol
	protocols := rc.Protocols
	if len(protocols) <= 0 {
		protocol := &config2.ProtocolConfig{}
		protocols = make(map[string]*config2.ProtocolConfig, 1)
		protocols[constant.Dubbo] = protocol
		rc.Protocols = protocols
	}
	for _, protocol := range protocols {
		if err := protocol.Init(); err != nil {
			return err
		}
	}

	// init registry
	registries := rc.Registries
	if registries != nil {
		for _, reg := range registries {
			if err := reg.Init(); err != nil {
				return err
			}
		}
	}

	if err := rc.MetadataReport.Init(rc); err != nil {
		return err
	}
	if err := rc.Metric.Init(); err != nil {
		return err
	}
	for _, t := range rc.Tracing {
		if err := t.Init(); err != nil {
			return err
		}
	}
	if err := initRouterConfig(); err != nil {
		return err
	}

	return nil
}

func initProvider() error {
	if o.providerFactory != nil {
		for _, provider := range o.providerFactory() {
			config2.SetProviderService(provider)
		}
	}

	if err := rc.Provider.Init(rc); err != nil {
		return err
	}
	return nil
}

func initConsumer() error {
	if o.consumerFactory != nil {
		for _, consumer := range o.consumerFactory() {
			config2.SetConsumerService(consumer)
		}
	}

	if err := rc.Consumer.Init(rc); err != nil {
		return err
	}
	return nil
}

func initShutdown() error {
	if err := rc.Shutdown.Init(); err != nil {
		return err
	}
	rc.Shutdown.InternalSignal = false
	return nil
}
