<h1 align="center">Drxutils</h1>

<div align="center">

_A collection of personal CLI utilities and scripts for my Arch Linux setup_

[![stars](https://img.shields.io/github/stars/druxorey/drxutils?color=8BE9FD&labelColor=191A21&style=for-the-badge)](https://github.com/druxorey/drxutils/stargazers)
[![size](https://img.shields.io/github/repo-size/druxorey/drxutils?label=Size&color=50FA7B&labelColor=191A21&style=for-the-badge)](https://github.com/druxorey/drxutils)
[![Visitors](https://api.visitorbadge.io/api/visitors?path=https%3A%2F%2Fgithub.com%2Fdruxorey%2Fdrxutils&label=Views&labelColor=%23191A21&countColor=%23FFB86C)](https://visitorbadge.io/status?path=https%3A%2F%2Fgithub.com%2Fdruxorey%2Fdrxutils)
[![license](https://img.shields.io/github/license/druxorey/drxutils?color=FF5555&labelColor=191A21&style=for-the-badge)](https://github.com/druxorey/drxutils/blob/main/LICENSE)

</div>

## About

This repository serves as a companion to my [dotfiles](https://github.com/druxorey/dotfiles). It contains all the custom Bash, Python, and Go scripts I have written to automate tasks, manage system updates, handle media, and improve my overall workflow in Arch Linux. These tools are designed to be minimal, fast, and primarily terminal-driven. They integrate deeply with my custom keybindings and workflow.

## Installation

The primary and intended way to deploy these tools is through my automated [`sysupdate`](bash/sysupdate) script. It automatically handles cloning this repository, updating paths, copying Bash scripts, and compiling Go projects into your `~/.local/bin` directory.

However, if you wish to use or install these utilities manually, you can do so by following the instructions below.

### Installing Bash & Python Scripts

The Bash and Python utilities are standalone scripts. You only need to copy them to a directory included in your system's `$PATH` and ensure they have execution permissions.

```bash
# 1. Clone the repository and navigate into it
git clone https://github.com/druxorey/drxutils.git
cd drxutils

# 2. Ensure your local binary directory exists
mkdir -p ~/.local/bin

# 3. Copy the desired script (e.g., 'compress') to the bin directory
cp bash/compress ~/.local/bin/

# 4. Grant execution permissions to the script so it can run globally
chmod +x ~/.local/bin/compress
```

### Installing Go Scripts

Go projects are a bit more complex as they require the Go compiler to build the source code into a standalone executable binary before you can use them.

```bash
# 1. Navigate to the specific Go project's directory (e.g., 'packages')
cd drxutils/packages

# 2. Download and verify all required Go module dependencies for the project
go mod tidy

# 3. Compile the source code and output the binary directly to your local bin path
go build -o ~/.local/bin/packages
```
## Overview for Go Scripts

Each Go project resides in its own subdirectory with a dedicated README detailing its installation and usage.

- **`packages`**: A package manager helper for Arch Linux that categorizes and searches for packages across Pacman and the AUR, maintaining a local registry for system bootstrapping. See the [packages README](packages/README.md) for more details.

## Overview for Bash & Python Scripts

### [`compress`](bash/compress)

A wrapper around `ffmpeg` designed to easily compress video files using the `libx265` video codec. It features 4 predefined compression levels (0 to 3) that automatically adjust the Constant Rate Factor (CRF) and the audio bitrate. Level 0 copies the original audio, while level 3 applies maximum compression with a slower preset and a 96k AAC audio bitrate.

**Usage:** `compress [-o OUTPUT_FILE] [-l COMPRESSION_LEVEL] INPUT_FILE`

**Examples:**

```bash
# Compresses 'video.mp4' using compression level 2 (High compression, CRF 28)
# and saves the resulting file as 'video_small.mp4'
compress -o video_small.mp4 -l 2 video.mp4

# Uses the default compression level (0) and automatically names 
# the output file based on the original name (e.g., video_0.mp4)
compress video.mp4
```

### [`cortana`](bash/cortana)

A background service manager for running local `ollama` . It checks the system processes to ensure no duplicate instances are running. After a brief wait, it boots up `open-webui` via `uv run` and confirms when the server is accessible on localhost. It also handles gracefully killing both processes when requested.

**Usage:** `cortana {start|stop|status}`

**Examples:**

```bash
# Starts both Ollama and Open WebUI services safely in the background
cortana start

# Checks the current process tree and reports if the service is RUNNING or STOPPED
cortana status

# Kills both the open-webui server and the ollama daemon
cortana stop
```

### [`crater`](bash/crater)

A boilerplate generator that speeds up the creation of new scripts and source code files. It reads from predefined template files stored in `~/.local/share/crater`. It supports multiple languages (Bash, C, C++, Go, HTML, Java, LaTeX, Lua, Python, Rust) and automatically assigns the proper file extension. If multiple templates exist for a single language, it provides an interactive selection menu.

**Usage:** `crater [-n FILE_NAME] [-i INDEX] FILE_TYPE`

**Examples:**

```bash
# Interactively prompts you to select an HTML template and creates 'default.html'
crater html

# Creates a new C++ file named 'main.cpp' using the template located at index 0,
# bypassing the interactive selection menu
crater -n main -i 0 cpp
```

### [`dirbyte`](bash/dirbyte)

A recursive directory analysis tool written in Python. It traverses a specified path, calculating the exact size of all contents, and formats the output into human-readable units (Kilobytes, Megabytes, or Gigabytes). It features flags to include hidden files (`-d`), perform recursive inner-folder scans (`-r`), and a script mode (`-s`) that outputs only the raw numerical total for use in other scripts.

**Usage:** `dirbyte [-r] [-d] [-s] [-k|-m|-g] PATH`

**Examples:**

```bash
# Recursively (-r) calculates the total size of the Downloads directory
# and displays the final output in Megabytes (-m)
dirbyte -r -m /home/druxorey/Downloads

# Calculates the size of the current directory, including hidden files (-d),
# but outputs only the raw number (-s) in Gigabytes (-g)
dirbyte -d -s -g .
```

### [`dotbak`](bash/dotbak)

My personal dotfiles backup utility. It acts as a wrapper for `rsync`, configured with specific exclude rules (like ignoring `.git/` or `*.gitignore` files) to cleanly mirror active configurations from `~/.config` and `/etc` into my local dotfiles repository workspace. Additionally, it reads Brave Browser's internal JSON bookmarks file, parses it using `awk`, and converts it into a clean YAML format for backup.

**Usage:** `dotbak [BACKUP_DIRECTORY]`

**Examples:**

```bash
# Synchronizes all predefined system configurations into the provided path
dotbak ~/Workspace/dotfiles/

# Uses the default hardcoded path ($HOME/Workspace/dotfiles/)
dotbak
```

### [`ex`](bash/ex)

An alias for automatic archive extractor. Instead of forcing you to remember the specific unarchiving commands and flags for every single format (`tar`, `unzip`, `unrar`, `gunzip`, `7z`, `unzstd`, etc.), this script evaluates the file extension of the provided archive and silently invokes the correct tool to extract its contents in the current directory.

**Usage:** `ex INPUT_FILE`

**Examples:**

```bash
# Automatically detects the .tar.gz extension and runs 'tar xzf'
ex archive.tar.gz

# Automatically detects the .zip extension and runs '7z x'
ex compressed_folder.zip
```

### [`gfix`](bash/gfix)

A quick Git version control helper. It verifies if the current directory is a valid git repository and if there are enough commits. If the checks pass, it executes `git reset --soft HEAD~1`. This action undoes the last commit but keeps all the file modifications staged in your working directory, allowing you to easily fix a typo or add a forgotten file to your previous commit.

**Usage:** `gfix`

**Examples:**

```bash
# Reverts the very last commit and leaves its changes ready in the staging area
gfix
```

### [`grun`](bash/grun)

A C and C++ compiler wrapper to compile and runner source code. Useful to quickly run small programs. It detects the file extension (`.c` or `.cpp`), invokes `gcc` or `g++` with the `-Wall` flag (and any extra provided flags), executes the resulting binary immediately, and finally deletes the binary to keep your workspace clean.

**Usage:** `grun [-f "FLAGS"] INPUT_FILE`

**Examples:**

```bash
# Compiles 'main.cpp', executes it, and deletes the binary right after it finishes
grun main.cpp

# Compiles 'math.c' injecting the '-O3' optimization flag, then runs and cleans it
grun -f "-O3" math.c
```

### [`hdmi`](bash/hdmi)

A dual-monitor management script tailored for my laptop setup using `xrandr`. When prompted to connect, it displays an interactive menu to select a desired resolution (1330x768, 1600x900, or 1920x1080). It then automatically configures the HDMI output to sit to the right of the primary laptop display (`eDP-1`). Finally, it reloads the BSPWM window manager to properly recognize the new screen space.

**Usage:** `hdmi {connect|disconnect}`

**Examples:**

```bash
# Opens the interactive menu, applies the chosen resolution, and sets up the second monitor
hdmi connect

# Turns off the HDMI-1 output and reloads the window manager
hdmi disconnect
```

### [`lofi`](bash/lofi)

A background music streamer. It pulls audio streams from a hardcoded array of Lofi YouTube URLs. It uses `yt-dlp` to fetch the stream data and pipes it directly into `ffplay` running in headless mode. While the music plays, it launches the `cava` audio visualizer in your terminal. Exiting the script (closing `cava`) automatically intercepts the signal and kills the background audio processes.

**Usage:** `lofi [-d] [-s SOURCE_URL] [VOLUME]`

**Examples:**

```bash
# Plays a randomly selected Lofi stream from the default array at 60% volume
lofi 60

# Plays a specific YouTube stream and runs in debug mode (-d), 
# which shows the standard ffplay output instead of the cava visualizer
lofi -d -s "http://youtu.be/some_id" 80
```

### [`mdtopdf`](bash/mdtopdf)

A Markdown-to-PDF `pandoc` wrapper powered by the `xelatex` PDF engine. It takes Obsidian-flavored markdown files and converts them into beautifully styled documents. It includes a custom Lua filter that intercepts Obsidian callouts (e.g., `[!INFO]`, `[!WARNING]`) and transforms them into colored LaTeX `tcolorbox` elements. It also adjusts table widths and strips internal wiki-links.

**Usage:** `mdtopdf FILE_NAME`

**Examples:**

```bash
# Compiles 'notes.md' into 'notes.pdf', applying the custom Lua filters, 
# syntax highlighting, and LaTeX styling instructions defined in the script
mdtopdf notes.md
```

### [`obsidian-clean`](bash/obsidian-clean)

A Git maintenance utility specifically built for my Obsidian vault. Over time, taking notes generates thousands of automated commits. This script navigates to the vault, creates an orphan branch (a branch with no history), stages all current files, commits them as a fresh «Clean branch», deletes the old `main` branch, and force-pushes the new history. This wipes the commit history completely to save storage space while preserving all current files.

**Usage:** `obsidian-clean`

**Examples:**

```bash
# Wipes the git commit history of the Obsidian vault and creates a fresh starting point
obsidian-clean
```

### [`open`](bash/open)

A terminal-based file opener. Instead of typing the name of specific applications, you pass a file to this script. It checks the file extension and automatically launches the correct GUI or CLI program associated with it in the background (`mpv` for video and audio formats, `nsxiv` for images and SVGs, and `zathura` for PDF documents).

**Usage:** `open INPUT_FILE`

**Examples:**

```bash
# Automatically launches the image file using the 'nsxiv' image viewer
open picture.png

# Automatically launches the video file using the 'mpv' media player
open video.mkv
```

### [`sysupdate`](bash/sysupdate)

My Arch Linux system maintenance and cleanup script. Running this executes a sequence of automated tasks:

1. Prompts for `sudo` and keeps it alive.
2. Updates `pacman` and AUR (`yay`) packages.
3. Upgrades Discord by pulling the latest Vencord installer (if Discord was updated).
4. Updates Yazi file manager plugins via `ya`.
5. Pulls the latest version of this `drxutils` repository, automatically recompiling any modified Go projects and installing updated bash scripts into `~/.local/bin`.
6. Cleans system cache (preserving a specific whitelist like `huggingface` models) and empties the trash bin if it exceeds a certain size threshold.
    

**Usage:** `sysupdate`

**Examples:**

```bash
# Runs the full system update, drxutils synchronization, and cache cleanup process
sysupdate
```    

## License

This project is licensed under the GPL-3.0 License. See the [LICENSE](LICENSE "null") file for more details.