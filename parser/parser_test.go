package parser

import "github.com/drycc/workflow-cli/cmd"

// DryccCmd is an implementation of Commander.
type FakeDryccCmd cmd.DryccCmd

func (d FakeDryccCmd) Println(...interface{}) (int, error) {
	return 1, nil
}

func (d FakeDryccCmd) Print(...interface{}) (int, error) {
	return 1, nil
}

func (d FakeDryccCmd) Printf(string, ...interface{}) (int, error) {
	return 1, nil
}

func (d FakeDryccCmd) PrintErrln(...interface{}) (int, error) {
	return 1, nil
}

func (d FakeDryccCmd) PrintErr(...interface{}) (int, error) {
	return 1, nil
}

func (d FakeDryccCmd) PrintErrf(string, ...interface{}) (int, error) {
	return 1, nil
}

func (d FakeDryccCmd) ServicesAdd(string, string, string, string) error {
	return nil
}

func (d FakeDryccCmd) ServicesList(string) error {
	return nil
}

func (d FakeDryccCmd) ServicesRemove(string, string, string, int) error {
	return nil
}
