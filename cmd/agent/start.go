package agent

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/faja/tornimo-agent/pkg/aggregator"
	"github.com/faja/tornimo-agent/pkg/collector"
	"github.com/faja/tornimo-agent/pkg/forwarder"
	"github.com/faja/tornimo-agent/pkg/serializer"
	"github.com/faja/tornimo-agent/pkg/statsd"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the agent in the foreground",
	Long:  "",
	RunE:  start,
}

func start(cmd *cobra.Command, args []string) error {
	if err := readConfig(); err != nil {
		return err
	}

	// TODO: defer stopAgent()

	// stop channel
	stopCh := make(chan error)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// "listen" for all the signals
	go func() {
		select {
		// TODO: stop command
		// TODO: catch all the errors
		case sig := <-signalCh:
			log.Printf("Receiving signal '%s', shutting down...", sig)
			stopCh <- nil
		}
	}()

	// start the agent
	if err := startAgent(); err != nil {
		return err
	}

	log.Println("[agent] agent started")

	// run forever until stopped
	select {
	case err := <-stopCh:
		return err
	}
}

func startAgent() error {
	// TODO: fix fowarder
	f := forwarder.NewDefaultForwarder(globalConfig["tornimo_put_address"])
	// TODO move start to New
	f.Start()

	s := serializer.NewSerializer(f, globalConfig["tornimo_token"])

	a := aggregator.InitAggregator(globalConfig["hostname"], s)

	collector.NewCollector()

	// TODO statsd_port is a string:(
	statsd.NewServer(globalConfig["statsd_port"], a.GetMetricsChan())

	return nil
}
