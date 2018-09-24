
|![](https://upload.wikimedia.org/wikipedia/commons/thumb/4/4c/Anchor_pictogram_yellow.svg/156px-Anchor_pictogram_yellow.svg.png) | Hephy Workflow is the open source fork of Deis Workflow.<br />Please [go here](https://www.teamhephy.com/) for more detail. |
|---:|---|
| 08/27/2018 | Team Hephy [blog][] comes online |
| 08/20/2018 | Deis [#community slack][] goes dark |
| 08/10/2018 | Hephy Workflow [v2.19.4][] fourth patch release |
| 08/08/2018 | [Deis website][] goes dark, then redirects to Azure Kubernetes Service |
| 08/01/2018 | Hephy Workflow [v2.19.3][] third patch release |
| 07/17/2018 | Hephy Workflow [v2.19.2][] second patch release |
| 07/12/2018 | Hephy Workflow [v2.19.1][] first patch release |
| 06/29/2018 | Hephy Workflow [v2.19.0][] first release in the open source fork of Deis |
| 06/16/2018 | Hephy Workflow [v2.19][] series is announced |
| 03/01/2018 | End of Deis Workflow maintenance: critical patches no longer merged |
| 12/11/2017 | Team Hephy [slack community][] invites first volunteers |
| 09/07/2017 | Deis Workflow [v2.18][] final release before entering maintenance mode |
| 09/06/2017 | Team Hephy [slack community][] comes online |

# Deis Client

[![Build Status](https://ci.deis.io/buildStatus/icon?job=Deis/workflow-cli/master)](https://ci.deis.io/job/Deis/job/workflow-cli/job/master/)
[![Go Report Card](https://goreportcard.com/badge/github.com/deis/workflow-cli)](https://goreportcard.com/report/github.com/deis/workflow-cli)
[![codebeat badge](https://codebeat.co/badges/05d314a8-ca61-4211-b69e-e7a3033662c8)](https://codebeat.co/projects/github-com-deis-workflow-cli)
[![codecov](https://codecov.io/gh/deis/workflow-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/deis/workflow-cli)

Download Links:

- [64 Bit Linux](https://storage.googleapis.com/workflow-cli-master/deis-latest-linux-amd64)
- [32 Bit Linux](https://storage.googleapis.com/workflow-cli-master/deis-latest-linux-386)
- [64 Bit Mac OS X](https://storage.googleapis.com/workflow-cli-master/deis-latest-darwin-amd64)
- [32 Bit Max OS X](https://storage.googleapis.com/workflow-cli-master/deis-latest-darwin-386)
- [64 Bit Windows](https://storage.googleapis.com/workflow-cli-master/deis-latest-windows-amd64.exe)
- [32 Bit Windows](https://storage.googleapis.com/workflow-cli-master/deis-latest-windows-386.exe)

`deis` is a command line utility used to interact with the [Deis](http://deis.io) open source PaaS.

Please add any [issues](https://github.com/teamhephy/workflow-cli/issues) you find with this software to the [Deis Workflow CLI Project](https://github.com/teamhephy/workflow-cli).

## Installation

### Pre-built Binary

See the appropriate sub-section below for your system to download and install the latest build of this software.

#### 64 Bit Linux

Master:

```console
curl -o deis https://storage.googleapis.com/workflow-cli-master/deis-latest-linux-amd64 && chmod +x deis
```

Latest stable release:

```
curl -o deis https://storage.googleapis.com/workflow-cli-release/deis-stable-linux-amd64 && chmod +x deis
```

#### 32 Bit Linux

Master:

```console
curl -o deis https://storage.googleapis.com/workflow-cli-master/deis-latest-linux-386 && chmod +x deis
```

Latest stable release:

```
curl -o deis https://storage.googleapis.com/workflow-cli-release/deis-stable-linux-386 && chmod +x deis
```

#### 64 Bit Mac OS X

Master:

```console
curl -o deis https://storage.googleapis.com/workflow-cli-master/deis-latest-darwin-amd64 && chmod +x deis
```

Latest stable release:

```
curl -o deis https://storage.googleapis.com/workflow-cli-release/deis-stable-darwin-amd64 && chmod +x deis
```

#### 32 Bit Max OS X

Master:

```console
curl -o deis https://storage.googleapis.com/workflow-cli-master/deis-latest-darwin-386 && chmod +x deis
```

Latest stable release:

```
curl -o deis https://storage.googleapis.com/workflow-cli-release/deis-stable-darwin-386 && chmod +x deis
```

#### 64 Bit Windows

Master:

```console
powershell -NoProfile -ExecutionPolicy Bypass -Command "(new-object net.webclient).DownloadFile('https://storage.googleapis.com/workflow-cli-master/deis-latest-windows-amd64.exe', 'deis.exe')"
```

Latest stable release:

```
powershell -NoProfile -ExecutionPolicy Bypass -Command "(new-object net.webclient).DownloadFile('https://storage.googleapis.com/workflow-cli-release/deis-stable-windows-amd64.exe', 'deis.exe')"
```

#### 32 Bit Windows

Master:

```console
powershell -NoProfile -ExecutionPolicy Bypass -Command "(new-object net.webclient).DownloadFile('https://storage.googleapis.com/workflow-cli-master/deis-latest-windows-386.exe', 'deis.exe')"
```

Latest stable release:

```
powershell -NoProfile -ExecutionPolicy Bypass -Command "(new-object net.webclient).DownloadFile('https://storage.googleapis.com/workflow-cli-release/deis-stable-windows-386.exe', 'deis.exe')"
```


After you execute the appropriate command for your system, you'll have a `deis` binary in the current directory.

Run the following to see the version:

```console
$ ./deis --version
```

You can then move it anywhere in your path:

```console
$ mv deis /usr/local/bin
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

    $ .\deis --version

## Usage

Running `deis help` will give you a up to date list of `deis` commands.
To learn more about a command run `deis help <command>`.

## License

see [LICENSE](https://github.com/teamhephy/workflow-cli/blob/master/LICENSE)

[v2.18]: https://github.com/teamhephy/workflow/releases/tag/v2.18.0
[k8s-home]: http://kubernetes.io
[install-k8s]: http://kubernetes.io/gettingstarted/
[mkdocs]: http://www.mkdocs.org/
[issues]: https://github.com/teamhephy/workflow/issues
[prs]: https://github.com/teamhephy/workflow/pulls
[Deis website]: http://deis.com/
[blog]: https://blog.teamhephy.info/blog/
[#community slack]: https://slack.deis.io/
[slack community]: https://slack.teamhephy.com/
[v2.18]: https://github.com/teamhephy/workflow/releases/tag/v2.18.0
[v2.19]: https://web.teamhephy.com
[v2.19.0]: https://gist.github.com/Cryptophobia/24c204583b18b9fc74c629fb2b62dfa3/revisions
[v2.19.1]: https://github.com/teamhephy/workflow/releases/tag/v2.19.1
[v2.19.2]: https://github.com/teamhephy/workflow/releases/tag/v2.19.2
[v2.19.3]: https://github.com/teamhephy/workflow/releases/tag/v2.19.3
[v2.19.4]: https://github.com/teamhephy/workflow/releases/tag/v2.19.4
