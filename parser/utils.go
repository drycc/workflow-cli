package parser

import (
	"fmt"
	"log"
	"strconv"

	"github.com/drycc/workflow-cli/cmd"
)

func safeGetString(args map[string]interface{}, key string) string {
	return safeGetValue(args, key, "")
}

func safeGetBool(args map[string]interface{}, key string) bool {
	return safeGetValue(args, key, false)
}

func safeGetInt(args map[string]interface{}, key string) int {
	if args[key] == nil {
		return 0
	}
	retVal, err := strconv.Atoi(args[key].(string))
	if err != nil {
		log.Fatalf("could not convert %s to int: %v", args[key], err)
	}
	return retVal
}

func safeGetValue[T any](args map[string]interface{}, key string, defaultValue T) T {
	if args[key] == nil {
		return defaultValue
	}
	return args[key].(T)
}

func responseLimit(limit string) (int, error) {
	if limit == "" {
		return -1, nil
	}

	return strconv.Atoi(limit)
}

// PrintUsage runs if no matching command is found.
func PrintUsage(cmdr cmd.Commander) {
	cmdr.PrintErrln("Found no matching command, try 'drycc help'")
	cmdr.PrintErrln("Usage: drycc <command> [<args>...]")
}

func printHelp(argv []string, usage string) bool {
	if len(argv) > 1 {
		if argv[1] == "--help" || argv[1] == "-h" {
			fmt.Print(usage)
			return true
		}
	}

	return false
}
