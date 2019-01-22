[![Build Status](https://travis-ci.org/drycc/workflow-cli.svg?branch=master)](https://travis-ci.org/drycc/workflow-cli)

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/4/4c/Anchor_pictogram_yellow.svg/156px-Anchor_pictogram_yellow.svg.png) | Drycc Workflow is the open source fork of Drycc Workflow.<br />Please [go here](https://www.drycc.com/) for more detail. |
|---:|---|
| 08/27/2018 | Team Drycc [blog][] comes online |
| 08/20/2018 | Drycc [#community slack][] goes dark |
| 08/10/2018 | Drycc Workflow [v2.19.4][] fourth patch release |
| 08/08/2018 | [Drycc website][] goes dark, then redirects to Azure Kubernetes Service |
| 08/01/2018 | Drycc Workflow [v2.19.3][] third patch release |
| 07/17/2018 | Drycc Workflow [v2.19.2][] second patch release |
| 07/12/2018 | Drycc Workflow [v2.19.1][] first patch release |
| 06/29/2018 | Drycc Workflow [v2.19.0][] first release in the open source fork of Drycc |
| 06/16/2018 | Drycc Workflow [v2.19][] series is announced |
| 03/01/2018 | End of Drycc Workflow maintenance: critical patches no longer merged |
| 12/11/2017 | Team Drycc [slack community][] invites first volunteers |
| 09/07/2017 | Drycc Workflow [v2.18][] final release before entering maintenance mode |
| 09/06/2017 | Team Drycc [slack community][] comes online |

# Drycc Client

[![Build Status](https://ci.drycc.cc/buildStatus/icon?job=Drycc/workflow-cli/master)](https://ci.drycc.cc/job/Drycc/job/workflow-cli/job/master/)
[![Go Report Card](https://goreportcard.com/badge/github.com/drycc/workflow-cli)](https://goreportcard.com/report/github.com/drycc/workflow-cli)
[![codebeat badge](https://codebeat.co/badges/05d314a8-ca61-4211-b69e-e7a3033662c8)](https://codebeat.co/projects/github-com-drycc-workflow-cli)
[![codecov](https://codecov.io/gh/drycc/workflow-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/drycc/workflow-cli)

Download Links:

- [64 Bit Linux](https://storage.googleapis.com/workflow-cli-master/drycc-latest-linux-amd64)
- [32 Bit Linux](https://storage.googleapis.com/workflow-cli-master/drycc-latest-linux-386)
- [64 Bit Mac OS X](https://storage.googleapis.com/workflow-cli-master/drycc-latest-darwin-amd64)
- [32 Bit Max OS X](https://storage.googleapis.com/workflow-cli-master/drycc-latest-darwin-386)
- [64 Bit Windows](https://storage.googleapis.com/workflow-cli-master/drycc-latest-windows-amd64.exe)
- [32 Bit Windows](https://storage.googleapis.com/workflow-cli-master/drycc-latest-windows-386.exe)

`drycc` is a command line utility used to interact with the [Drycc](http://drycc.cc) open source PaaS.

Please add any [issues](https://github.com/drycc/workflow-cli/issues) you find with this software to the [Drycc Workflow CLI Project](https://github.com/drycc/workflow-cli).

## Installation

### Pre-built Binary

See the appropriate sub-section below for your system to download and install the latest build of this software.

#### 64 Bit Linux

Master:

```console
curl -o drycc https://storage.googleapis.com/workflow-cli-master/drycc-latest-linux-amd64 && chmod +x drycc
```

Latest stable release:

```
curl -o drycc https://storage.googleapis.com/workflow-cli-release/drycc-stable-linux-amd64 && chmod +x drycc
```

#### 32 Bit Linux

Master:

```console
curl -o drycc https://storage.googleapis.com/workflow-cli-master/drycc-latest-linux-386 && chmod +x drycc
```

Latest stable release:

```
curl -o drycc https://storage.googleapis.com/workflow-cli-release/drycc-stable-linux-386 && chmod +x drycc
```

#### 64 Bit Mac OS X

Master:

```console
curl -o drycc https://storage.googleapis.com/workflow-cli-master/drycc-latest-darwin-amd64 && chmod +x drycc
```

Latest stable release:

```
curl -o drycc https://storage.googleapis.com/workflow-cli-release/drycc-stable-darwin-amd64 && chmod +x drycc
```

#### 32 Bit Max OS X

Master:

```console
curl -o drycc https://storage.googleapis.com/workflow-cli-master/drycc-latest-darwin-386 && chmod +x drycc
```

Latest stable release:

```
curl -o drycc https://storage.googleapis.com/workflow-cli-release/drycc-stable-darwin-386 && chmod +x drycc
```

#### 64 Bit Windows

Master:

```console
powershell -NoProfile -ExecutionPolicy Bypass -Command "(new-object net.webclient).DownloadFile('https://storage.googleapis.com/workflow-cli-master/drycc-latest-windows-amd64.exe', 'drycc.exe')"
```

Latest stable release:

```
powershell -NoProfile -ExecutionPolicy Bypass -Command "(new-object net.webclient).DownloadFile('https://storage.googleapis.com/workflow-cli-release/drycc-stable-windows-amd64.exe', 'drycc.exe')"
```

#### 32 Bit Windows

Master:

```console
powershell -NoProfile -ExecutionPolicy Bypass -Command "(new-object net.webclient).DownloadFile('https://storage.googleapis.com/workflow-cli-master/drycc-latest-windows-386.exe', 'drycc.exe')"
```

Latest stable release:

```
powershell -NoProfile -ExecutionPolicy Bypass -Command "(new-object net.webclient).DownloadFile('https://storage.googleapis.com/workflow-cli-release/drycc-stable-windows-386.exe', 'drycc.exe')"
```


After you execute the appropriate command for your system, you'll have a `drycc` binary in the current directory.

Run the following to see the version:

```console
$ ./drycc --version
```

You can then move it anywhere in your path:

```console
$ mv drycc /usr/local/bin
```

### From Scratch on OS X and Linux

To compile the client from scratch, ensure you have Docker installed and run

    $ make

### From Scratch on Windows

To compile the client from scratch, open PowerShell and execute the following commands in the source directory.

    $ .\make bootstrap
    $ .\make build

`.\make bootstrap` will fetch all required dependencies, while `.\make build` will compile and install
the client in the current directory.

    $ .\drycc --version

## Usage

Running `drycc help` will give you a up to date list of `drycc` commands.
To learn more about a command run `drycc help <command>`.

## License

see [LICENSE](https://github.com/drycc/workflow-cli/blob/master/LICENSE)

[v2.18]: https://github.com/drycc/workflow/releases/tag/v2.18.0
[k8s-home]: http://kubernetes.io
[install-k8s]: http://kubernetes.io/gettingstarted/
[mkdocs]: http://www.mkdocs.org/
[issues]: https://github.com/drycc/workflow/issues
[prs]: https://github.com/drycc/workflow/pulls
[Drycc website]: http://drycc.com/
[blog]: https://blog.drycc.info/blog/
[#community slack]: https://slack.drycc.cc/
[slack community]: https://slack.drycc.com/
[v2.18]: https://github.com/drycc/workflow/releases/tag/v2.18.0
[v2.19]: https://web.drycc.com
[v2.19.0]: https://gist.github.com/Cryptophobia/24c204583b18b9fc74c629fb2b62dfa3/revisions
[v2.19.1]: https://github.com/drycc/workflow/releases/tag/v2.19.1
[v2.19.2]: https://github.com/drycc/workflow/releases/tag/v2.19.2
[v2.19.3]: https://github.com/drycc/workflow/releases/tag/v2.19.3
[v2.19.4]: https://github.com/drycc/workflow/releases/tag/v2.19.4
