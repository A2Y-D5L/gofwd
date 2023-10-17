# gofwd

![Go Report Card](https://goreportcard.com/badge/github.com/A2Y-D5L/gofwd)
![Build Status](https://github.com/A2Y-D5L/gofwd/workflows/Go/badge.svg)
![License](https://img.shields.io/github/license/A2Y-D5L/gofwd)

Move forward with `gofwd` to effortlessly keep your Go installation in step with the latest release.

## Overview

`gofwd` checks your current Go installation against the latest official release. If a newer version is found, `gofwd` offers a streamlined update process, ensuring you're always working with the most recent stable version of Go.

## Features

- **Automatic Version Checking**: Compares your local Go version with the latest official release.
- **User Confirmation**: Before making any changes, `gofwd` seeks your confirmation.
- **Installation Backup**: Before updating, the current installation is backed up, ensuring peace of mind.
- **Permission Checks**: Ensures the necessary permissions are in place before making changes.

## Installation

To get started with `gofwd`, you can clone the repository and build the tool:

```bash
git clone https://github.com/A2Y-D5L/gofwd.git
cd gofwd
go build
```

This will produce a `gofwd` binary in the current directory.

## Usage

Simply run the `gofwd` binary:

```bash
./gofwd
```

Follow the on-screen prompts to check for updates and proceed with the installation if a newer version is available.

## Contributing

Contributions welcome! If you find a bug or have suggestions, please open an issue. If you'd like to contribute code, please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
