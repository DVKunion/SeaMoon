package main

import (
	"github.com/DVKunion/SeaMoon/pkg/client"
	"github.com/DVKunion/SeaMoon/pkg/consts"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	mod        string
	debug      bool
	verbose    bool
	listenAddr string
	proxyAddr  string

	clientMap = map[string]func(listenAddr string, proxyAddr string, verbose bool){
		"http":   client.NewHttpClient,
		"socks5": client.NewSocks5Client,
	}

	rootCommand = &cobra.Command{
		Use:   "client",
		Short: "SeaMoon Client",
		Run: func(cmd *cobra.Command, args []string) {
			Client()
		},
	}
	versionCommand = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info(consts.Version)
		},
	}
)

func init() {
	rootCommand.AddCommand(versionCommand)
	rootCommand.Flags().StringVarP(&mod, "mod", "m", "http", "mod of SeaMoon client")
	rootCommand.Flags().StringVarP(&listenAddr, "laddr", "l", ":9000", "local client address like : 0.0.0.0:9000")
	rootCommand.Flags().StringVarP(&proxyAddr, "paddr", "p", "", "proxy server address")
	rootCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "proxy detail log")
	rootCommand.Flags().BoolVarP(&debug, "debug", "d", false, "proxy detail log")
}

func Client() {
	handle := clientMap[mod]
	handle(listenAddr, proxyAddr, verbose)
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
