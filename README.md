[![Build Status](https://travis-ci.org/drycc/workflow-cli.svg?branch=master)](https://travis-ci.org/drycc/workflow-cli)

# Drycc Client

[![Build Status](https://ci.drycc.cc/buildStatus/icon?job=Drycc/workflow-cli/master)](https://ci.drycc.cc/job/Drycc/job/workflow-cli/job/master/)
[![Go Report Card](https://goreportcard.com/badge/github.com/drycc/workflow-cli)](https://goreportcard.com/report/github.com/drycc/workflow-cli)
[![codebeat badge](https://codebeat.co/badges/05d314a8-ca61-4211-b69e-e7a3033662c8)](https://codebeat.co/projects/github-com-drycc-workflow-cli)
[![codecov](https://codecov.io/gh/drycc/workflow-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/drycc/workflow-cli)

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

see [LICENSE](https://github.com/drycc/workflow-cli/blob/master/LICENSE)

[k8s-home]: http://kubernetes.io
[install-k8s]: http://kubernetes.io/gettingstarted/
[mkdocs]: http://www.mkdocs.org/
[issues]: https://github.com/drycc/workflow/issues
[prs]: https://github.com/drycc/workflow/pulls
[Drycc website]: http://drycc.com/
[blog]: https://blog.drycc.info/blog/
[slack community]: https://slack.drycc.com/
