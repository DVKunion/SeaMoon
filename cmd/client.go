package main

import (
	"SeaMoon/pkg/proxy"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	proxyAddr          string
	listenAddr         string
	clientMod          string
	verbose            bool
	rootClientCommand  = &cobra.Command{}
	proxyClientCommand = &cobra.Command{
		Use:   "proxy",
		Short: "SeaMoon Proxy Client",
		Run: func(cmd *cobra.Command, args []string) {
			Proxy()
		},
	}
	versionCommand = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("V1.0.0-BETA")
		},
	}
)

func init() {
	rootClientCommand.AddCommand(versionCommand)
	proxyClientCommand.Flags().StringVarP(&clientMod, "mod", "m", "http", "mod of SeaMoon client")
	proxyClientCommand.Flags().StringVarP(&listenAddr, "laddr", "l", ":9000", "local client address like : 0.0.0.0:9000")
	proxyClientCommand.Flags().StringVarP(&proxyAddr, "paddr", "p", "", "proxy server address")
	proxyClientCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "proxy detail log")

	rootClientCommand.AddCommand(proxyClientCommand)
}

func Proxy() {
	switch clientMod {
	case "http":
		proxy.NewHttpClient(listenAddr, proxyAddr, verbose)
		break
	case "socks5":
		proxy.NewSocks5Client(listenAddr, proxyAddr, verbose)
		break
	}
}

func main() {
	if err := rootClientCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
