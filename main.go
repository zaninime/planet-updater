package main

import (
	"fmt"

	"go.zanini.me/planet-updater/firmwares"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	typeFlag       = kingpin.Flag("type", "Lamp type").Required().PlaceHolder("[pro|compact]").Enum("pro", "compact")
	destinationArg = kingpin.Arg("destination", "The WiFish address").Required().String()
	portArg        = kingpin.Arg("port", "The WiFish IP port number").Required().Uint16()
)

const (
	assetPROFirmwareName     = "PRO-V14.bin"
	assetCompactFirmwareName = "Compact-V15.bin"
)

func main() {
	kingpin.Version("1.0.0")
	kingpin.Parse()

	var firmwareAsset string

	if *typeFlag == "pro" {
		firmwareAsset = assetPROFirmwareName
	} else {
		firmwareAsset = assetCompactFirmwareName
	}

	fmt.Println(updatePlanet(*destinationArg, int(*portArg), firmwares.MustAsset(firmwareAsset)))
}
