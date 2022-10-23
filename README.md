# HoloCure Save File Transfer Tool

![Tasukeru GUI demo](https://i.imgur.com/HDohzzB.png)

This is a small tool to import [HoloCure](https://kay-yu.itch.io/holocure) save files from one PC to another.

## Download

Grab a release from [Releases](https://github.com/DaniruKun/tasukeru/releases) and save it anywhere on your PC.
Pick the executable matching your architecture (note: HoloCure currently only runs on Windows).

## Usage

1. Get the save file from the source PC at `Users\[your username]\AppData\Local\HoloCure\save.dat` and move it to the target PC
2. Play HoloCure **at least once** one the target PC
3. Launch `Tasukeru.exe`
4. Open the save file you want to import
5. Press `Import`
6. The save should now be imported

## Build

There are 2 build options: as a GUI app, and as a CLI:

### GUI

Install [fyne-cross](https://github.com/fyne-io/fyne-cross).
Then run

```shell
make compile-windows
```

### CLI

```shell
make compile-cli
```

This will produce binaries for each platform in the `bin` directory.

## Advanced

You can manually call the executable and pass arguments directly.

If you pass a **single argument**, then `saveA.dat` will be merged into the `save.dat` found in the system's HoloCure cache directory (e.g. `Local\HoloCure\save.dat`).

This is equivalent to drag n dropping the `saveA.dat` on top of the executable.

E.g. `tasukeru-windows-amd64.exe saveA.dat`

If you pass **2 arguments**, then `saveA.dat` will be merged into `save.dat` or whatever path is given and replace it.

E.g. `tasukeru-windows-amd64.exe saveA.dat save.dat` will produce the patched `save.dat` in the current directory.

On Unix systems you can quickly inspect a save file with

```sh
base64 --decode -i save.dat
```
