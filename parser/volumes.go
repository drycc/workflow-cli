package parser

import (
	"os"

	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
	"golang.org/x/exp/slices"
)

// Volumes commands to their specific function.
func Volumes(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for volumes:

volumes:create           create a volume for the application
volumes:expand           expand a volume for the application
volumes:list             list volumes in the application
volumes:info             print information about a volume
volumes:delete           delete a volume from the application
volumes:client           the client used to manage volume files
volumes:mount            mount a volume to process of the application
volumes:unmount          unmount a volume from process of the application

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "volumes:create":
		return volumesCreate(argv, cmdr)
	case "volumes:expand":
		return volumesExpand(argv, cmdr)
	case "volumes:list":
		return volumesList(argv, cmdr)
	case "volumes:info":
		return volumesInfo(argv, cmdr)
	case "volumes:delete":
		return volumesDelete(argv, cmdr)
	case "volumes:client":
		return volumesClient(argv, cmdr)
	case "volumes:mount":
		return volumesMount(argv, cmdr)
	case "volumes:unmount":
		return volumesUnmount(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "volumes" {
			argv[0] = "volumes:list"
			return volumesList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func volumesCreate(argv []string, cmdr cmd.Commander) error {
	usage := `
Create a volume for the application.

Usage: drycc volumes:create <name> <size> [options]

Arguments:
  <name>
    the volume name.
  <size>
    the volume size, such as '500G'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -t --type=<type>
    the volume type, such as csi, nfs, default is 'csi'.
  --nfs-server=<nfs-server>
    the hostname or ip address of the nfs server.
  --nfs-path=<nfs-path>
    path that is exported by the nfs server.
  --nfs-read-only
    a flag indicating whether the storage will be mounted as read only.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	vType := safeGetValue(args, "--type", "csi")
	name := safeGetString(args, "<name>")
	size := safeGetString(args, "<size>")
	parameters := map[string]interface{}{}
	if vType == "nfs" {
		server := safeGetString(args, "--nfs-server")
		path := safeGetString(args, "--nfs-path")
		readOnly := safeGetBool(args, "--nfs-read-only")
		parameters["nfs"] = map[string]interface{}{
			"server":   server,
			"path":     path,
			"readOnly": readOnly,
		}
	}

	return cmdr.VolumesCreate(app, name, vType, size, parameters)
}

func volumesExpand(argv []string, cmdr cmd.Commander) error {
	usage := `
Expand a volume for the application.

Usage: drycc volumes:expand <name> <size> [options]

Arguments:
  <name>
    the volume name.
  <size>
    the volume size, such as '500G'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	name := safeGetString(args, "<name>")
	size := safeGetString(args, "<size>")

	return cmdr.VolumesExpand(app, name, size)
}

func volumesList(argv []string, cmdr cmd.Commander) error {
	usage := `
List volumes in the application.

Usage: drycc volumes:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	results, err := responseLimit(safeGetString(args, "--limit"))

	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")

	return cmdr.VolumesList(app, results)
}

func volumesInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
Print information about a volume.

Usage: drycc volumes:info <name> [options]

Arguments:
  <name>
    the volume name to be info.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	name := safeGetString(args, "<name>")

	return cmdr.VolumesInfo(app, name)
}

func volumesDelete(argv []string, cmdr cmd.Commander) error {
	usage := `
Delete a volume from the application.

Usage: drycc volumes:delete <name> [options]

Arguments:
  <name>
    the volume name to be removed.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	name := safeGetString(args, "<name>")

	return cmdr.VolumesDelete(app, name)
}

func volumesClient(argv []string, cmdr cmd.Commander) error {
	usage := `
The client used to manage volume files.

Usage: drycc volumes:client <cmd> [options] -- <args>...

Arguments:
  <cmd>
    ls         list volume files
    cp         copy volume files
    rm         remove volume files
  <args>
    arguments for running commands

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	cmd := safeGetString(args, "<cmd>")
	index := slices.Index(os.Args, "--")
	arguments := os.Args[index+1:]
	return cmdr.VolumesClient(app, cmd, arguments...)
}

func volumesMount(argv []string, cmdr cmd.Commander) error {
	usage := `
Mount a volume for an application.

Usage: drycc volumes:mount <name> <type>=<path>... [options]

Arguments:
  <name>
    the volume name.
  <type>
    the process name as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.
  <path>
    the filesystem path.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	name := safeGetString(args, "<name>")

	return cmdr.VolumesMount(app, name, args["<type>=<path>"].([]string))
}

func volumesUnmount(argv []string, cmdr cmd.Commander) error {
	usage := `
Unmount a volume for an application.

Usage: drycc volumes:unmount <name> <type>... [options]

Arguments:
  <name>
    the volume name.
  <type>
    the process name as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`
	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	name := safeGetString(args, "<name>")

	return cmdr.VolumesUnmount(app, name, args["<type>"].([]string))
}
