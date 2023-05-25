module github.com/drycc/workflow-cli

go 1.20

require (
	github.com/containerd/console v1.0.3
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/drycc/controller-sdk-go v0.0.0-20221102055456-444409b1af71
	github.com/drycc/pkg v0.0.0-20220830031116-26c11ff8667c
	github.com/ghodss/yaml v1.0.0
	github.com/gorilla/websocket v1.5.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/stretchr/testify v1.8.0
	golang.org/x/exp v0.0.0-20220827204233-334a2380cb91
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/goware/urlx v0.3.2 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/drycc/controller-sdk-go => github.com/jianxiaoguo/controller-sdk-go v0.0.0-20230524075627-dcfa3f6f1df1
