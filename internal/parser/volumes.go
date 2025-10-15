package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewVolumesCommand creates the volumes command
func NewVolumesCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volumes",
		Short: i18n.T("Manage volumes for your applications"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.VolumesList(app, results)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(volumesListCommand(cmdr))
	cmd.AddCommand(volumesAddCommand(cmdr))
	cmd.AddCommand(volumesExpandCommand(cmdr))
	cmd.AddCommand(volumesInfoCommand(cmdr))
	cmd.AddCommand(volumesRemoveCommand(cmdr))
	cmd.AddCommand(volumesClientCommand(cmdr))
	cmd.AddCommand(volumesMountCommand(cmdr))
	cmd.AddCommand(volumesUnmountCommand(cmdr))
	return cmd
}

func volumesAddCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		vtype        string
		name         string
		size         string
		nfsServer    string
		nfsPath      string
		ossServer    string
		ossBucket    string
		ossPathStyle bool
		ossAccessKey string
		ossSecretKey string
	}

	cmd := &cobra.Command{
		Use: "add <name> <size>",
		Example: template.CustomExample(
			"drycc volumes add myvolume 1G",
			map[string]string{
				"<name>": i18n.T("The volume name"),
				"<size>": i18n.T("The volume size, such as '500G'"),
			},
		),
		Short: i18n.T("Create a volume for the application"),
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			flags.name = args[0]
			flags.size = args[1]
			parameters := make(map[string]any)
			switch flags.vtype {
			case "nfs":
				parameters["nfs"] = map[string]any{
					"server": flags.nfsServer,
					"path":   flags.nfsPath,
				}
			case "oss":
				parameters["oss"] = map[string]any{
					"server":     flags.ossServer,
					"bucket":     flags.ossBucket,
					"access_key": flags.ossAccessKey,
					"secret_key": flags.ossSecretKey,
					"path_style": flags.ossPathStyle,
				}
			}
			return cmdr.VolumesCreate(app, flags.name, flags.vtype, flags.size, parameters)
		},
	}

	cmd.Flags().StringVarP(&flags.vtype, "type", "t", "csi", i18n.T("The volume type, such as csi, nfs, oss"))
	cmd.Flags().StringVar(&flags.nfsServer, "nfs-server", "", i18n.T("The hostname or ip address of the nfs server"))
	cmd.Flags().StringVar(&flags.nfsPath, "nfs-path", "", i18n.T("Path that is exported by the nfs server"))

	cmd.Flags().StringVar(&flags.ossServer, "oss-server", "", i18n.T("Endpoint url for object storage service"))
	cmd.Flags().StringVar(&flags.ossBucket, "oss-bucket", "", i18n.T("Bucket name in object storage"))
	cmd.Flags().StringVar(&flags.ossAccessKey, "oss-access-key", "", i18n.T("Access key id for authentication"))
	cmd.Flags().StringVar(&flags.ossSecretKey, "oss-secret-key", "", i18n.T("Secret access key for authentication"))
	cmd.Flags().BoolVar(&flags.ossPathStyle, "oss-path-style", false, i18n.T("Force a path-style endpoint to be used"))

	cmd.Flags().SortFlags = false
	cmd.MarkFlagsRequiredTogether("nfs-server", "nfs-path")
	cmd.MarkFlagsRequiredTogether("oss-server", "oss-bucket", "oss-access-key", "oss-secret-key")

	volumeTypeCompletion := completion.VolumeTypeCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("type", volumeTypeCompletion.CompletionFunc)
	return cmd
}

func volumesExpandCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		name string
		size string
	}

	volumeCompletion := completion.VolumeCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "expand <name> <size>",
		Example: template.CustomExample(
			"drycc volumes expand myvolume 2G",
			map[string]string{
				"<name>": i18n.T("The volume name"),
				"<size>": i18n.T("The volume size, such as '500G'"),
			},
		),
		Short:             i18n.T("Expand a volume capacity for the application"),
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: volumeCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.name = args[0]
			flags.size = args[1]
			return cmdr.VolumesExpand(app, flags.name, flags.size)
		},
	}

	return cmd
}

func volumesListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List volumes in the application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.VolumesList(app, results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}

func volumesInfoCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		name string
	}

	volumeCompletion := completion.VolumeCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "info <name>",
		Example: template.CustomExample(
			"drycc volumes info myvolume",
			map[string]string{
				"<name>": i18n.T("The volume name to be info"),
			},
		),
		Short:             i18n.T("Print information about a volume"),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: volumeCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.name = args[0]
			return cmdr.VolumesInfo(app, flags.name)
		},
	}

	return cmd
}

func volumesRemoveCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		name string
	}

	volumeCompletion := completion.VolumeCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "remove <name>",
		Example: template.CustomExample(
			"drycc volumes remove myvolume",
			map[string]string{
				"<name>": i18n.T("The volume name to be removed"),
			},
		),
		Short:             i18n.T("Delete a volume from the application"),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: volumeCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.name = args[0]
			return cmdr.VolumesDelete(app, flags.name)
		},
	}

	return cmd
}

func volumesClientCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: i18n.T("The client used to manage volume files"),
	}
	cmd.AddCommand(&cobra.Command{
		Use: "ls <target>",
		Example: template.CustomExample(
			"drycc volumes client ls vol://myvolume/tmp",
			map[string]string{
				"<target>": i18n.T("The target path of volume"),
			},
		),
		Short: i18n.T("List volume objects"),
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.VolumesClient(app, "ls", args...)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use: "cp <source> <target>",
		Example: template.CustomExample(
			"drycc volumes client cp vol://myvolume/tmp /tmp",
			map[string]string{
				"<source>": i18n.T("The volume or local source path"),
				"<target>": i18n.T("The volume or local target path"),
			},
		),
		Short: i18n.T("Copy volume objects"),
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.VolumesClient(app, "cp", args...)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use: "rm <target>",
		Example: template.CustomExample(
			"drycc volumes client rm vol://myvolume/tmp",
			map[string]string{
				"<target>": i18n.T("The target path of volume"),
			},
		),
		Short: i18n.T("Remove volume objects"),
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.VolumesClient(app, "rm", args...)
		},
	})

	return cmd
}

func volumesMountCommand(cmdr *commands.DryccCmd) *cobra.Command {
	volumesMountCompletion := completion.VolumesMountCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "mount <name> <ptype>=<path>...",
		Example: template.CustomExample(
			"drycc volumes mount myvolume web=/data",
			map[string]string{
				"<name>":  i18n.T("The volume name"),
				"<ptype>": i18n.T("The process name as defined in your Procfile"),
				"<path>":  i18n.T("The filesystem path"),
			},
		),
		Short:             i18n.T("Mount a volume to process of the application"),
		Args:              cobra.MinimumNArgs(2),
		ValidArgsFunction: volumesMountCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			mountSpecs := args[1:]
			return cmdr.VolumesMount(app, name, mountSpecs)
		},
	}

	return cmd
}

func volumesUnmountCommand(cmdr *commands.DryccCmd) *cobra.Command {
	volumesUnmountCompletion := completion.VolumesUnmountCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "unmount <name> <ptype>...",
		Example: template.CustomExample(
			"drycc volumes unmount myvolume web worker",
			map[string]string{
				"<name>":  i18n.T("The volume name"),
				"<ptype>": i18n.T("The process name as defined in your Procfile"),
			},
		),
		Short:             i18n.T("Unmount a volume from process of the application"),
		Args:              cobra.MinimumNArgs(2),
		ValidArgsFunction: volumesUnmountCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			ptypes := args[1:]
			return cmdr.VolumesUnmount(app, name, ptypes)
		},
	}

	return cmd
}
