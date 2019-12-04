package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func main() {
	root := newRoot()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// Return a new root command.
func newRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dqlite-demo",
		Short: "Demo application using dqlite",
		Long: `This demo shows how to integrate a Go application with dqlite.

Complete documentation is available at https://github.com/canonical/go-dqlite`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented")
		},
	}
	cmd.AddCommand(newStart())
	cmd.AddCommand(newAdd())
	cmd.AddCommand(newUpdate())
	cmd.AddCommand(newQuery())
	//	cmd.AddCommand(newBenchmark())
	cmd.AddCommand(newCluster())
	cmd.AddCommand(newAdhoc())
	cmd.AddCommand(newWeb())
	cmd.AddCommand(newDumper())

	return cmd
}

// Start a web server for remote clients.
func newWeb() *cobra.Command {
	var cluster *[]string
	var port *int

	cmd := &cobra.Command{
		Use:   "webserver",
		Short: "Start a web server.",
		//Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			WebServer(*port, defaultDatabase, *cluster)
			return nil
		},
	}

	flags := cmd.Flags()
	cluster = flags.StringSliceP("cluster", "c", defaultCluster, "addresses of existing cluster nodes")
	port = flags.IntP("port", "p", 4001, "port to serve traffic on")

	return cmd
}

// Return a new query key command.
func newQuery() *cobra.Command {
	var cluster *[]string

	cmd := &cobra.Command{
		Use:   "query <key>",
		Short: "Get a key from the demo table.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dbQuery(args[0], *cluster)
		},
	}

	flags := cmd.Flags()
	cluster = flags.StringSliceP("cluster", "c", defaultCluster, "addresses of existing cluster nodes")

	return cmd
}

// Return a new start command.
func newStart() *cobra.Command {
	var dir string
	var address string

	cmd := &cobra.Command{
		Use:   "start <id>",
		Short: "Start a dqlite-demo node.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.Wrapf(err, "%s is not a number", args[0])
			}
			return dbStart(id, dir, address)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&dir, "dir", "d", "/tmp/dqlite-demo", "base demo directory")
	flags.StringVarP(&address, "address", "a", "", "address of the node (default is 127.0.0.1:918<ID>)")

	return cmd
}

// Return a new update key command.
func newUpdate() *cobra.Command {
	var cluster *[]string

	cmd := &cobra.Command{
		Use:   "update <key> <value>",
		Short: "Insert or update a key in the demo table.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dbUpdate(args[0], args[1], *cluster)
		},
	}
	flags := cmd.Flags()
	cluster = flags.StringSliceP("cluster", "c", defaultCluster, "addresses of existing cluster nodes")

	return cmd
}

// Return a cluster nodes command.
func newCluster() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "display cluster nodes.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return dbCluster()
		},
	}
	return cmd
}

// Return a new add command.
func newAdd() *cobra.Command {
	var address string
	var cluster *[]string

	cmd := &cobra.Command{
		Use:   "add <id>",
		Short: "Add a node to the dqlite-demo cluster.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.Wrapf(err, "%s is not a number", args[0])
			}
			return dbAdd(id, address, *cluster)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&address, "address", "a", "", "address of the node to add (default is 127.0.0.1:918<ID>)")
	cluster = flags.StringSliceP("cluster", "c", defaultCluster, "addresses of existing cluster nodes")

	return cmd
}

const defaultDatabase = "demo.db"

// Return a new update key command.
func newAdhoc() *cobra.Command {
	var cluster *[]string

	cmd := &cobra.Command{
		Use:   "adhoc <statment>...",
		Short: "execute a statement against the demo database.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dbExec(defaultDatabase, *cluster, args...)
		},
	}
	flags := cmd.Flags()
	cluster = flags.StringSliceP("cluster", "c", defaultCluster, "addresses of existing cluster nodes")

	return cmd
}

// Return a new dump command.
func newDumper() *cobra.Command {
	var cluster *[]string

	cmd := &cobra.Command{
		Use:   "dump database",
		Short: "dump the Database and its associated WAL file.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dbDump(args[0], defaultDatabase, *cluster)
		},
	}
	flags := cmd.Flags()
	cluster = flags.StringSliceP("cluster", "c", defaultCluster, "addresses of existing cluster nodes")

	return cmd
}
