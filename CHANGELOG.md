### v2.0.0-beta3 -> v2.0.0

#### Features

 - [`203fd0d`](https://github.com/deis/workflow-cli/commit/203fd0dda982aa8b200e96c94b1fcaa59e09ef5e) config: ignore comments in config:push

#### Fixes

 - [`f22a325`](https://github.com/deis/workflow-cli/commit/f22a32584b41e53f1d3a805ad3bae0b78cc7cb35) cmd: return error on bad limit parse
 - [`848b3a6`](https://github.com/deis/workflow-cli/commit/848b3a6176606f798a81bdcb04f0715b63889266) cmd: split only twice at most
 - [`cb24ddb`](https://github.com/deis/workflow-cli/commit/cb24ddb7426b3c69f034f33bd84774a7468b4c2d) cmd: strip printLogs with only one newline
 - [`aaf5e78`](https://github.com/deis/workflow-cli/commit/aaf5e78d9f2ec38c48fc4e8b6d86f4b8659dfa63) apps:create: Better output when --no-remote used
 - [`54e365b`](https://github.com/deis/workflow-cli/commit/54e365b3a4cfedf1b726a1f0ec03b5c700af16d4) (all): fix missing error checks caught by ineffassign
 - [`1e64709`](https://github.com/deis/workflow-cli/commit/1e647091a7b90355ed3cfddb813a0b56f5d06441) (all): correct spelling mistakes and linting problems caught by misspell and golint

#### Documentation

 - [`7283e7c`](https://github.com/deis/workflow-cli/commit/7283e7c775938e0f61e9670515ef62d1fdab0076) README: add note to move binary elsewhere in PATH
 - [`348f2b2`](https://github.com/deis/workflow-cli/commit/348f2b291f3e085cfec1e448beae9fc15f487608) badge: added code-beat badge
 - [`ee40d59`](https://github.com/deis/workflow-cli/commit/ee40d59364c93175d7f24a55bfca14f934d89bc4) CHANGELOG.md: update for v2.0.0-beta3

#### Maintenance

 - [`23e1970`](https://github.com/deis/workflow-cli/commit/23e197093bb55ee178ad150101d45af5c67cae62) version: bump client version to 2.0.0
 - [`2ff93ba`](https://github.com/deis/workflow-cli/commit/2ff93ba0956e87da626317c91c983865969b6c16) cmd/config.go: DEIS_RELEASE -> WORKFLOW_RELEASE
 - [`c48c568`](https://github.com/deis/workflow-cli/commit/c48c568448932b6893a102c956fdb29b3bdac2a7) git: Make builder git remote use different hostname

### v2.0.0-beta2 -> v2.0.0-beta3

#### Features

 - [`6fb32bf`](https://github.com/deis/workflow-cli/commit/6fb32bf26873a1cdcba35fce13b57284263ba400) registry: add support for setting private registry information per application

#### Fixes

 - [`80ce942`](https://github.com/deis/workflow-cli/commit/80ce942bfb4c70952aa1a0472df1817dff72830b) logger: Split on the double new lines when printing log messages
 - [`c673253`](https://github.com/deis/workflow-cli/commit/c673253f1da1d1cff48b6551c8fec9748ff34115) controller: use django HttpResponse for logs
 - [`852b0b0`](https://github.com/deis/workflow-cli/commit/852b0b00e029adafab32e59fc7f688917d067b1e) registry: add a missing case statement for registry in deis.go
