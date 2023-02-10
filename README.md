# Drycc Client

[![Build Status](https://woodpecker.drycc.cc/api/badges/drycc/workflow-cli/status.svg)](https://woodpecker.drycc.cc/drycc/workflow-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/drycc/workflow-cli)](https://goreportcard.com/report/github.com/drycc/workflow-cli)
[![codebeat badge](https://codebeat.co/badges/b609cb7f-7b42-4214-8787-09298f553176)](https://codebeat.co/projects/github-com-drycc-workflow-cli-main)
[![codecov](https://codecov.io/gh/drycc/workflow-cli/branch/main/graph/badge.svg)](https://codecov.io/gh/drycc/workflow-cli)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdrycc%2Fworkflow-cli.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdrycc%2Fworkflow-cli?ref=badge_shield)

Download Links: https://github.com/drycc/workflow-cli/releases

`drycc` is a command line utility used to interact with the [Drycc](http://drycc.cc) open source PaaS.

Please add any [issues](https://github.com/drycc/workflow-cli/issues) you find with this software to the [Drycc Workflow CLI Project](https://github.com/drycc/workflow-cli).

## Download and run the Drycc install script(not windows):

For example, install v1.0.1 version:

```console
sudo bash < <(curl -fsSL https://github.com/drycc/workflow-cli/releases/download/v1.0.1/install-drycc.sh)
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

see [LICENSE](https://github.com/drycc/workflow-cli/blob/main/LICENSE)

[k8s-home]: http://kubernetes.io
[install-k8s]: http://kubernetes.io/gettingstarted/
[mkdocs]: http://www.mkdocs.org/
[issues]: https://github.com/drycc/workflow/issues
[prs]: https://github.com/drycc/workflow/pulls
[Drycc website]: http://drycc.com/
[blog]: https://blog.drycc.info/blog/
[slack community]: https://slack.drycc.com/

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdrycc%2Fworkflow-cli.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdrycc%2Fworkflow-cli?ref=badge_large)
