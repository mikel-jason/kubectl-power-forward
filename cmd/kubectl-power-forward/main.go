package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mikel-jason/kube-power-forward/pkg/cmd/powerforward"
	"github.com/mikel-jason/kube-power-forward/pkg/proxy"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var rootCmd = &cobra.Command{
	Use:     "kubectl-power-forward",
	Version: version,
	Short:   "kubectl-power-forward is kubectl port-forward on steroids",
	Long: `kubectl-power-forward allows port-fowarding to a Kubernetes service with auto-reconnect on connection lost.
It contains a SOCKS and HTTP reverse proxy to simulate accessing your workloads with real-word domains.

Example configuration file:
--------------------------------------------------
proxy:
  socks:
    enabled: true
    listenAddress: "0.0.0.0:1080" # this is the default
  httpReverse:
    enabled: true
    listenAddress: "0.0.0.0:80" # this is the default
forwards:
  - namespace: default
    serviceName: echoserver
    podPort: 5678
    localPort: 8080
    hostname: echo.example.com
--------------------------------------------------

For using privileged ports, you must bind the capability to the binary like this:
sudo setcap 'cap_net_bind_service+ep' <PATH_TO_BINARY>
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath := cmd.Flag("config").Value.String()

		configFileContents, err := os.ReadFile(configFilePath)
		if err != nil {
			log.Fatalf("Cannot open config file %s: %s", configFilePath, err.Error())
		}
		var cfg ConfigFile
		err = yaml.Unmarshal(configFileContents, &cfg)
		if err != nil {
			log.Fatalf("Cannot parse config file %s: %s", configFilePath, err.Error())
		}
		cfg.FillDefaults()
		pfConfig := &powerforward.Config{
			Forwards: cfg.Forwards,
		}

		if cfg.Proxy.Socks.Enabled {
			pfConfig.SocksListenAddr = cfg.Proxy.Socks.ListenAddr
		}

		if cfg.Proxy.HttpReverse.Enabled {
			pfConfig.HttpReverseListenAddr = cfg.Proxy.HttpReverse.ListenAddr
		}

		err = powerforward.Start(pfConfig)
		if err != nil {
			log.Fatalf("Cannot start: %s", err.Error())
		}

		signalForShutdownChan := make(chan os.Signal, 1)
		signal.Notify(signalForShutdownChan, os.Interrupt, syscall.SIGTERM)

		<-signalForShutdownChan
		log.Println("Received interrupt signal, gracefully shutting down")
		proxy.Stop()
		powerforward.Stop()
	},
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to identify current directory: %s", err.Error())
	}
	defaultConfigFile := filepath.Join(pwd, ".power-forward.yaml")
	rootCmd.Flags().StringP("config", "f", defaultConfigFile, "Path to config file")

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
