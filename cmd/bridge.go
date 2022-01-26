package cmd

import (
	"fmt"

	"github.com/natekfl/maveth/mavtunnel"
	"github.com/natekfl/maveth/network/tap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Starts the bridge over MAVLink",
	Long: `
		Starts the bridge over MAVLink.
		
		The MAVLink endpoint to use is specified by the --mavendpoint (-m) flag and is required. Possible endpoints types are:
			udps:listen_ip:port (udp, server mode)
			udpc:dest_ip:port (udp, client mode)
			udpb:broadcast_ip:port (udp, broadcast mode)
			tcps:listen_ip:port (tcp, server mode)
			tcpc:dest_ip:port (tcp, client mode)
			serial:port:baudrate (serial)
		In general, only the serial type should be used. As all the others are network based, they have the possibility of creating an infinite loop if you don't know what you're doing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mavEndpoint, err := mavtunnel.Endpoint(viper.GetString("mavendpoint"))
		if err != nil {
			return err
		}
		forwardErr := make(chan error)
		go (func() {
			forwardErr <- mavtunnel.Connect(mavEndpoint, func(b []byte) {
				err := tap.SendPacket(b)
				if err != nil {
					forwardErr <- err
				}
			})
		})()
		go (func() {
			forwardErr <- tap.StartInterfaceRead(mavtunnel.SendPacket)
		})()
		fmt.Println("Bridge started")
		return fmt.Errorf("error when running: %s", <-forwardErr)
	},
}

func init() {
	rootCmd.AddCommand(bridgeCmd)

	bridgeCmd.Flags().StringP("mavendpoint", "m", "", "The MAVLink endpoint to use")
	viper.GetViper().BindPFlag("mavendpoint", bridgeCmd.Flags().Lookup("mavendpoint"))
}
