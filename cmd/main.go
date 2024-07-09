package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/DVKunion/SeaMoon/cmd/client"
	"github.com/DVKunion/SeaMoon/cmd/server"
	"github.com/DVKunion/SeaMoon/pkg/api/database/drivers"
	"github.com/DVKunion/SeaMoon/pkg/system/version"
)

var (
	debug bool

	// server params
	addr  string
	port  string
	proto string

	rootCommand = &cobra.Command{}

	serverCommand = &cobra.Command{
		Use:   "server",
		Short: "SeaMoon server mod",
		RunE:  Server,
	}

	clientCommand = &cobra.Command{
		Use:   "client",
		Short: "SeaMoon client mod",
	}

	clientWebCommand = &cobra.Command{
		Use:   "web",
		Short: "SeaMoon client web mod",
		Run:   Client,
	}

	generateCommand = &cobra.Command{
		Use:   "generate",
		Short: "SeaMoon generate devs code",
		RunE:  drivers.Drive().Generate,
	}

	versionCommand = &cobra.Command{
		Use:   "version",
		Short: "SeaMoon version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("SeaMoon Powered By DVK")
			fmt.Printf("Version: %s\n", version.Version)
			fmt.Printf("Commit: %s\n", version.Commit)
			fmt.Printf("V2rayCoreVersion: %s\n", version.V2rayCoreVersion)
		},
	}
)

func Client(cmd *cobra.Command, args []string) {
	drivers.Init()
	client.Serve(cmd.Context(), debug)
}

func Server(cmd *cobra.Command, args []string) error {
	s, err := server.New(
		server.WithHost("0.0.0.0"),
		server.WithPort(port),
		server.WithProto(proto),
	)

	if err != nil {
		return err
	}

	return s.Serve(cmd.Context())
}

func init() {
	clientWebCommand.Flags().BoolVarP(&debug, "debug", "d", false, "proxy detail log")

	serverCommand.Flags().StringVarP(&addr, "addr", "a", "0.0.0.0", "server listen addr")
	serverCommand.Flags().StringVarP(&port, "port", "p", "9000", "server listen port")
	serverCommand.Flags().StringVarP(&proto, "proto", "t", "websocket", "server listen proto: (websocket/grpc)")

	rootCommand.AddCommand(versionCommand)
	rootCommand.AddCommand(clientCommand)
	rootCommand.AddCommand(serverCommand)
	rootCommand.AddCommand(generateCommand)

	clientCommand.AddCommand(clientWebCommand)
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
