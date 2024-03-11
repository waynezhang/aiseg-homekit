package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/waynezhang/aiseg-homekit/internal/aisegmanager"
	"github.com/waynezhang/aiseg-homekit/internal/log"
)

func Execute() {
	var verbose bool
	var rootCmd = &cobra.Command{
		Use:   "aiseg",
		Short: "AiSEG controller",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetVerbose(verbose)
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(ListCmd)
	rootCmd.AddCommand(ToggleCmd)
	rootCmd.AddCommand(TurnOnCmd)
	rootCmd.AddCommand(TurnOffCmd)
	rootCmd.AddCommand(TurnAllOnCmd)
	rootCmd.AddCommand(TurnAllOffCmd)

	rootCmd.AddCommand(HKServeCmd)

	rootCmd.AddCommand(VersionCmd)

	_ = rootCmd.Execute()
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all devices",
	Run: func(cmd *cobra.Command, args []string) {
		mgr := aisegmanager.DiscoverNewAiSEGManager()
		for _, d := range mgr.Devices {
			state := " "
			if d.IsOn {
				state = "x"
			}
			fmt.Printf("[%s] [%s] %s\n", d.NodeId, state, d.Name)
		}
	},
}

var ToggleCmd = func() *cobra.Command {
	fn := func(cmd *cobra.Command, args []string) {
		nodeId := args[0]
		log.D("Toggle device %s", nodeId)

		mgr := aisegmanager.DiscoverNewAiSEGManager()
		if err := mgr.ToggleDevice(nodeId); err != nil {
			log.E("Failed to toggle device due to %s", err.Error())
		}
	}

	cmd := &cobra.Command{
		Use:   "toggle",
		Short: "Toggle device",
		Args:  cobra.ExactArgs(1),
		Run:   fn,
	}

	return cmd
}()

var TurnOnCmd = func() *cobra.Command {
	fn := func(cmd *cobra.Command, args []string) {
		nodeId := args[0]
		log.D("Turn device %s on", nodeId)

		mgr := aisegmanager.DiscoverNewAiSEGManager()
		if err := mgr.TurnDevice(nodeId, true); err != nil {
			log.E("Failed to turn device on due to %s", err.Error())
		}
	}

	cmd := &cobra.Command{
		Use:   "on",
		Short: "Turn device on",
		Args:  cobra.ExactArgs(1),
		Run:   fn,
	}

	return cmd
}()

var TurnOffCmd = func() *cobra.Command {
	fn := func(cmd *cobra.Command, args []string) {
		nodeId := args[0]
		log.D("Turn device %s off", nodeId)

		mgr := aisegmanager.DiscoverNewAiSEGManager()
		if err := mgr.TurnDevice(nodeId, false); err != nil {
			log.E("Failed to turn device off due to %s", err.Error())
		}
	}

	cmd := &cobra.Command{
		Use:   "off",
		Short: "Turn device off",
		Args:  cobra.ExactArgs(1),
		Run:   fn,
	}

	return cmd
}()

var TurnAllOnCmd = func() *cobra.Command {
	fn := func(cmd *cobra.Command, args []string) {
		argType := args[0]
		log.D("Turn device %s on", argType)

		var deviceType aisegmanager.DeviceType
		if argType == "light" {
			deviceType = aisegmanager.DeviceTypeLight
		} else if argType == "floorheating" {
			deviceType = aisegmanager.DeviceTypeFloorHeating
		} else {
			log.E("Unsupported device type %d", deviceType)
			return
		}

		mgr := aisegmanager.DiscoverNewAiSEGManager()
		if err := mgr.TurnAllDevices(deviceType, true); err != nil {
			log.E("Failed to turn device on due to %s", err.Error())
		}
	}

	cmd := &cobra.Command{
		Use:   "allon",
		Short: "Turn all devices on by type",
		Args:  cobra.ExactArgs(1),
		Run:   fn,
	}

	return cmd
}()

var TurnAllOffCmd = func() *cobra.Command {
	fn := func(cmd *cobra.Command, args []string) {
		argType := args[0]
		log.D("Turn device %s off", argType)

		var deviceType aisegmanager.DeviceType
		if argType == "light" {
			deviceType = aisegmanager.DeviceTypeLight
		} else if argType == "floorheating" {
			deviceType = aisegmanager.DeviceTypeFloorHeating
		} else {
			log.E("Unsupported device type %d", deviceType)
			return
		}

		mgr := aisegmanager.DiscoverNewAiSEGManager()
		if err := mgr.TurnAllDevices(deviceType, false); err != nil {
			log.E("Failed to turn device off due to %s", err.Error())
		}
	}

	cmd := &cobra.Command{
		Use:   "alloff",
		Short: "Turn all devices off by type",
		Args:  cobra.ExactArgs(1),
		Run:   fn,
	}

	return cmd
}()
