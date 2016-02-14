package main

/* Planet Updater - easily update your Elos Planet lamps.
Copyright (C) 2015 Francesco Zanini <francesco@zanini.me>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.*/

import (
	"os"
	"time"

	"github.com/gosuri/uiprogress"
	"github.com/mgutz/ansi"
	"github.com/mgutz/logxi/v1"
	"github.com/zaninime/planet-updater/firmwares"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	typeFlag       = kingpin.Flag("type", "Choose lamp type.").Required().PlaceHolder("[pro|compact]").Enum("pro", "compact")
	debugFlag      = kingpin.Flag("debug", "Enable debugging messages.").Bool()
	destinationArg = kingpin.Arg("destination", "The WiFish address").Required().String()
	portArg        = kingpin.Arg("port", "The WiFish IP port number").Required().Uint16()
)

const (
	assetPROFirmwareName     = "PLANET12Ch-V116.bin"
	assetCompactFirmwareName = "COMPACT12Ch-V117.bin"
)

const version = "1.0.0"

var (
	yellowColor = ansi.ColorCode("yellow")
	resetColor  = ansi.ColorCode("reset")
)

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	if *debugFlag {
		log.DefaultLog.SetLevel(log.LevelDebug)
	} else {
		log.DefaultLog.SetLevel(log.LevelInfo)
	}

	var firmwareAsset string
	if *typeFlag == "pro" {
		firmwareAsset = assetPROFirmwareName
	} else {
		firmwareAsset = assetCompactFirmwareName
	}

	log.Info("Planet Updater starting", "version", version)
	ch := make(chan updateProgress)

	tStart := time.Now()

	go func() {
		if err := updatePlanet(*destinationArg, int(*portArg), firmwares.MustAsset(firmwareAsset), ch); err != nil {
			log.Error("Couldn't complete update", "err", err)
			os.Exit(1)
		}
	}()

	var progressBar *uiprogress.Bar

	for msg := range ch {
		switch msg.state {
		case 0:
			log.Info("Resolving", "dest", *destinationArg)
		case 1:
			log.Info("Connecting to lamp", "dest", msg.message)
		case 2:
			log.Info("Handshaking in progress")
		case 3:
			log.Info("Putting lamp in download mode")
		case 4:
			log.Info("Awaiting sync")
		case 5:
			log.Info("Uploading firmware")
			uiprogress.Start()
			progressBar = uiprogress.AddBar(msg.pkts)
			progressBar.PrependElapsed()
			progressBar.AppendCompleted()
		case 6:
			progressBar.Incr()
			//fmt.Printf("\r  Progress: %s% 5.1f%%%s", yellowColor, msg.progress*100.0, resetColor)
		case 7:
			uiprogress.Stop()
			//fmt.Print("\r                      \r")
			log.Info("Awaiting feedback from lamp")
		case 8:
			deltaStart := time.Now().Sub(tStart)
			log.Info("Update completed", "time", deltaStart.String())
		}
	}
}
