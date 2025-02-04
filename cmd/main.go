/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"poly-bridge/chainfeelisten"
	"poly-bridge/coinpricelisten"
	"poly-bridge/common"
	"poly-bridge/conf"
	"poly-bridge/crosschaineffect"
	"poly-bridge/crosschainlisten"
	"poly-bridge/crosschainstats"
	"runtime"
	"strings"
	"syscall"

	"github.com/astaxie/beego/logs"
	"github.com/urfave/cli"
)

var (
	logLevelFlag = cli.UintFlag{
		Name:  "loglevel",
		Usage: "Set the log level to `<level>` (0~6). 0:Trace 1:Debug 2:Info 3:Warn 4:Error 5:Fatal 6:MaxLevel",
		Value: 1,
	}

	configPathFlag = cli.StringFlag{
		Name:  "cliconfig",
		Usage: "Server config file `<path>`",
		Value: "./conf/config_testnet.json",
	}

	logDirFlag = cli.StringFlag{
		Name:  "logdir",
		Usage: "log directory",
		Value: "./Log/",
	}
)

//getFlagName deal with short flag, and return the flag name whether flag name have short name
func getFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "poly-bridge Service"
	app.Action = StartServer
	app.Version = "1.0.0"
	app.Copyright = "Copyright in 2019 The Ontology Authors"
	app.Flags = []cli.Flag{
		logLevelFlag,
		configPathFlag,
		logDirFlag,
	}
	app.Commands = []cli.Command{}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func StartServer(ctx *cli.Context) {
	for true {
		startServer(ctx)
		sig := waitSignal()
		stopServer()
		if sig != syscall.SIGHUP {
			break
		} else {
			continue
		}
	}
}

func startServer(ctx *cli.Context) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/bridge_server.log"}`)
	configFile := ctx.GlobalString(getFlagName(configPathFlag))
	config := conf.NewConfig(configFile)
	if config == nil {
		logs.Error("startServer - read config failed!")
		return
	}
	{
		conf, _ := json.Marshal(config)
		logs.Info("%s\n", string(conf))
	}
	common.SetupChainsSDK(config)
	crosschainlisten.StartCrossChainListen(config.Server, config.Backup, config.ChainListenConfig, config.DBConfig)
	if config.Backup {
		return
	}
	coinpricelisten.StartCoinPriceListen(config.Server, config.CoinPriceUpdateSlot, config.CoinPriceListenConfig, config.DBConfig)
	chainfeelisten.StartFeeListen(config.Server, config.FeeUpdateSlot, config.FeeListenConfig, config.DBConfig)
	crosschaineffect.StartCrossChainEffect(config.Server, config.EventEffectConfig, config.DBConfig)
	crosschainstats.StartCrossChainStats(config.Server, config.StatsConfig, config.DBConfig)
}

func waitSignal() os.Signal {
	exit := make(chan os.Signal, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(sc)
	go func() {
		for sig := range sc {
			logs.Info("cross chain listen received signal:(%s).", sig.String())
			exit <- sig
			close(exit)
			break
		}
	}()
	sig := <-exit
	return sig
}

func stopServer() {
	crosschainlisten.StopCrossChainListen()
	coinpricelisten.StopCoinPriceListen()
	chainfeelisten.StopFeeListen()
	crosschaineffect.StopCrossChainEffect()
	crosschainstats.StopCrossChainStats()
}

func main() {
	if err := setupApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
