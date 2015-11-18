# Planet Updater

[![AGPL v3](https://gnu.org/graphics/agplv3-88x31.png "AGPL v3")](https://gnu.org/licenses/agpl.html)

This simple CLI tool lets you update your [Elos lamps](http://elos.eu/gammaprodotti/illuminazione/planet~64/products.html) directly from the CLI.

## Features

* **Fast**. At least 4x faster that the stock one. In fact, the update requires ~ 4' 30" instead of 18'.
* **Lightweight**. Around 6 MB with firmwares bundled in.
* **Multi-platform**. It's written in Go, so it's supported on every platform supported by Go itself.
* **Ready to use**. You can find the binary releases in the [Releases](https://github.com/zaninime/planet-updater/releases) section.

## Bundled firmwares


| Lamp type | Firmware      | Notes                   |
|-----------|---------------|-------------------------|
| PRO       | v14 (aka 114) | Latest as of 2015-11-18 |
| Compact   | v15 (aka 115) | Latest as of 2015-11-18 |

## Usage

The command usage and switches are visible using the flag ``--help``.

Basically, you just need to launch the executable, specifying the type of the lamp (pro or compact) and the WiFish address to which the lamp is connected to, ie.

```bash
$ # Update PlanetPRO at 192.168.1.50:5000
$ planet-updater --type pro 192.168.1.50 5000

$ # Update PlanetCompact at 192.168.1.100:5000
$ planet-updater --type compact 192.168.1.100 5000
```

If you find the update procedure failing more than once, feel free to run the tool with the ``--debug`` switch and open an issue on GitHub.

## For developers

### How to get it

It's as simple as getting every other Go pkg: ``go get github.com/zaninime/planet-updater``.

### How to hack it

...should I tell you? It's a 2 source files tool. Improvements and ideas are more than welcome.
