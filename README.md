# HoloCure Save File Transfer Tool

This is a small tool to import [HoloCure](https://kay-yu.itch.io/holocure) save files from one PC to another.

## Download

Grab a release from [Releases](https://github.com/DaniruKun/tasukeru/releases) and save it anywhere on your PC.
Pick the executable matching your architecture (note: HoloCure currently only runs on Windows).

## Usage

1. Build a release by running `make` or download a release for your platform from [Releases](https://github.com/DaniruKun/tasukeru/releases)
2. Get the save file from the source PC at `Users\[your username]\AppData\Local\HoloCure\save.dat` and move it to the target PC
3. Play HoloCure **at least once** one the target PC
4. On the target PC, drag and drop the `save.dat` onto `tasukeru-*.exe`
5. When prompted, press `Enter`
6. The save should now be imported
