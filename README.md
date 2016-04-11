# Deis Client

[![Build Status](https://travis-ci.org/deis/workflow-cli.svg?branch=master)](https://travis-ci.org/deis/workflow-cli)
[![Go Report Card](http://goreportcard.com/badge/deis/workflow-cli)](http://goreportcard.com/report/deis/workflow-cli)
[![Download](https://api.bintray.com/packages/deis/deisci/deis/images/download.svg)](https://bintray.com/deis/deisci/deis/_latestVersion)

`deis` is a command line utility used to interact with the [Deis](http://deis.io) open source PaaS.

Please add any [issues](https://github.com/deis/workflow-cli/issues) you find with this software to the [Deis Workflow CLI Project](https://github.com/deis/workflow-cli).

## Installation

### From Bintray

Install the `deis` client from [bintray](https://bintray.com/) by running

	$ curl -sSL http://deis.io/deis-cli/install-v2.sh | bash

The installer will fetch the latest version of the client into your current directory.

	$ ./deis --version

### From Scratch

To compile the client from scratch, ensure you have Docker installed and run

	$ make bootstrap
	$ make build

`make bootstrap` will fetch all required dependencies, while `make build` will compile and install
the client in the current directory.

	$ ./deis --version

## Usage

Running `deis help` will give you a up to date list of `deis` commands.
To learn more about a command run `deis help <command>`.

## Windows Support

`deis` has experimental support for Windows. To build deis for Windows, you need to install
[go](https://golang.org/) and [glide](https://github.com/Masterminds/glide). Then run the `make.bat` script.

## License

Copyright 2015, Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
