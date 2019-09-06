package cmd

import (
	cli "github.com/spf13/cobra"
	config "github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/snowzach/gogrpcapi/conf"
	"github.com/snowzach/gogrpcapi/thingrpc/thingrpcserver"
	"github.com/snowzach/gogrpcapi/thingrpc"
	"github.com/snowzach/gogrpcapi/server"
	"github.com/snowzach/gogrpcapi/store/postgres"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var (
	apiCmd = &cli.Command{
		Use:   "api",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) { // Initialize the databse

			var thingStore thingrpc.ThingStore
			var err error
			switch config.GetString("storage.type") {
			case "postgres":
				thingStore, err = postgres.New()
			}
			if err != nil {
				logger.Fatalw("Database Error", "error", err)
			}

			// Create the GRPC/HTTP server
			s, err := server.New()
			if err != nil {
				logger.Fatalw("Could not create grpc/http server",
					"error", err,
				)
			}

			// Create the rpcserver
			rpcserver, err := thingrpcserver.New(thingStore)
			if err != nil {
				logger.Fatalw("Could not create thing rpcserver",
					"error", err,
				)
			}

			// Register the Thing RPC server to the GRPC Server
			thingrpc.RegisterThingRPCServer(s.GRPCServer(), rpcserver)
			s.GwReg(thingrpc.RegisterThingRPCHandlerFromEndpoint)

			err = s.ListenAndServe()
			if err != nil {
				logger.Fatalw("Could not start server",
					"error", err,
				)
			}

			<-conf.Stop.Chan() // Wait until StopChan
			conf.Stop.Wait()   // Wait until everyone cleans up
			zap.L().Sync()     // Flush the logger

		},
	}
)
