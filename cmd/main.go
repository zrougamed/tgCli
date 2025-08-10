package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zrougamed/tgCli/internal/cloud"
	"github.com/zrougamed/tgCli/internal/config"
	"github.com/zrougamed/tgCli/internal/helpers"
	"github.com/zrougamed/tgCli/internal/models"
	"github.com/zrougamed/tgCli/internal/server"
	"github.com/zrougamed/tgCli/pkg/constants"
)

func init() {
	var err error
	constants.HomeDir, err = os.UserHomeDir()
	if err != nil {
		log.Fatal("Unable to get user home directory:", err)
	}

	constants.ConfigDir = filepath.Join(constants.HomeDir, ".tgcli")
	constants.ConfigFile = filepath.Join(constants.ConfigDir, "config.yml")
	constants.CredsFile = filepath.Join(constants.ConfigDir, "creds.bank")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(constants.ConfigDir, 0755); err != nil {
		log.Fatal("Unable to create config directory:", err)
	}

	// Initialize viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(constants.ConfigDir)

	// Set defaults
	viper.SetDefault("tgcloud.user", "mail@domain.com")
	viper.SetDefault("tgcloud.password", "")
	viper.SetDefault("machines", make(map[string]models.MachineConfig))
	viper.SetDefault("default", "")

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create default
			helpers.CreateDefaultConfig(constants.ConfigFile)
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	}
}

func main() {
	helpers.GracefulShutdown()
	availableVersion, err := helpers.CheckForUpdates()
	if err != nil {
		log.Printf("Error checking for updates: %v", err)
		availableVersion = "N/A"
	}
	var rootCmd = &cobra.Command{
		Use:   "tg",
		Short: "TigerGraph CLI tool for cloud and server management",
		Long:  `A comprehensive CLI tool for managing TigerGraph cloud instances and server operations.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add global flags
	rootCmd.PersistentFlags().BoolVarP(&constants.Debug, "debug", "d", false, "Enable debug mode")

	// Add version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("TigerGraph CLI\n")
			fmt.Printf("  Version Installed: %s\n", constants.VERSION_CLI)
			fmt.Printf("  Version Available: %s\n", availableVersion)
			fmt.Printf("Support:\n")
			fmt.Printf("   TigerGraph Community: https://community.tigergraph.com\n")
			fmt.Printf("   TigerGraph Discord: https://discord.gg/GkEmvDqB\n")
			fmt.Printf("Copyright (c) 2014-2024 TigerGraph. All rights reserved.\n")
		},
	}

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(createCloudCmd())
	rootCmd.AddCommand(createServerCmd())
	rootCmd.AddCommand(createConfCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createCloudCmd() *cobra.Command {
	var cloudCmd = &cobra.Command{
		Use:   "cloud",
		Short: "TigerGraph Cloud operations",
		Long:  `Manage TigerGraph Cloud instances including login, start, stop, terminate, and list operations.`,
	}

	// Login command
	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to tgcloud.io",
		Run:   cloud.RunLogin,
	}
	loginCmd.Flags().StringP("email", "e", "", "Email address for tgcloud.io")
	loginCmd.Flags().StringP("password", "p", "", "Password for tgcloud.io")
	loginCmd.Flags().StringP("save", "s", "n", "Save credentials (y/n)")
	loginCmd.Flags().StringP("output", "o", "stdout", "Output format (stdout/json)")

	// Start command
	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start a tgcloud instance",
		Run:   cloud.RunStart,
	}
	startCmd.Flags().StringP("id", "i", "", "TGCloud Machine ID")
	startCmd.MarkFlagRequired("id")

	// Stop command
	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop a tgcloud instance",
		Run:   cloud.RunStop,
	}
	stopCmd.Flags().StringP("id", "i", "", "TGCloud Machine ID")
	stopCmd.MarkFlagRequired("id")

	// Terminate command
	var terminateCmd = &cobra.Command{
		Use:   "terminate",
		Short: "Terminate a tgcloud instance",
		Run:   cloud.RunTerminate,
	}
	terminateCmd.Flags().StringP("id", "i", "", "TGCloud Machine ID")
	terminateCmd.MarkFlagRequired("id")

	// Archive command
	var archiveCmd = &cobra.Command{
		Use:   "archive",
		Short: "Archive a tgcloud instance",
		Run:   cloud.RunArchive,
	}
	archiveCmd.Flags().StringP("id", "i", "", "TGCloud Machine ID")
	archiveCmd.MarkFlagRequired("id")

	// List command
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all tgcloud instances",
		Run:   cloud.RunList,
	}
	listCmd.Flags().StringP("activeonly", "a", "y", "Hide terminated servers (y/n)")
	listCmd.Flags().StringP("output", "o", "stdout", "Output format (stdout/json)")

	// Create command
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a tgcloud instance",
		Run:   cloud.RunCreate,
	}
	createCmd.Flags().StringP("id", "i", "", "TGCloud Starter Kit")

	cloudCmd.AddCommand(loginCmd, startCmd, stopCmd, terminateCmd, archiveCmd, listCmd, createCmd)
	return cloudCmd
}

func createServerCmd() *cobra.Command {
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "TigerGraph Server operations",
		Long:  `Manage TigerGraph server operations including GSQL, demos, algorithms, and services.`,
	}

	// GSQL command
	var gsqlCmd = &cobra.Command{
		Use:   "gsql",
		Short: "Execute a GSQL terminal",
		Run:   server.RunGSQL,
	}
	gsqlCmd.Flags().StringP("alias", "a", "", "TigerGraph server alias to use")
	gsqlCmd.Flags().StringP("user", "u", "tigergraph", "TigerGraph user")
	gsqlCmd.Flags().StringP("password", "p", "tigergraph", "TigerGraph password")
	gsqlCmd.Flags().String("host", "http://127.0.0.1", "TigerGraph host")
	gsqlCmd.Flags().String("gsPort", "14240", "GSQL Port")

	// Backup command
	var backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup a TigerGraph server",
		Run:   server.RunBackup,
	}
	backupCmd.Flags().StringP("alias", "a", "", "TigerGraph server alias to use")
	backupCmd.Flags().StringP("user", "u", "tigergraph", "TigerGraph user")
	backupCmd.Flags().StringP("password", "p", "tigergraph", "TigerGraph password")
	backupCmd.Flags().String("host", "http://127.0.0.1", "TigerGraph host")
	backupCmd.Flags().String("gsPort", "14240", "GSQL Port")
	backupCmd.Flags().String("restPort", "9000", "REST Port")
	backupCmd.Flags().StringP("type", "t", "ALL", "Backup type (ALL/SCHEMA/DATA)")

	// Services command
	var servicesCmd = &cobra.Command{
		Use:   "services",
		Short: "Start/Stop GPE/GSE/RESTPP Services",
		Run:   server.RunServices,
	}
	servicesCmd.Flags().StringP("user", "u", "tigergraph", "TigerGraph user")
	servicesCmd.Flags().StringP("password", "p", "tigergraph", "TigerGraph password")
	servicesCmd.Flags().String("host", "http://127.0.0.1", "TigerGraph host")
	servicesCmd.Flags().String("gsPort", "14240", "GSQL Port")
	servicesCmd.Flags().String("ops", "start", "Operation (start/stop)")

	serverCmd.AddCommand(gsqlCmd, backupCmd, servicesCmd)
	return serverCmd
}

func createConfCmd() *cobra.Command {
	var confCmd = &cobra.Command{
		Use:   "conf",
		Short: "Configuration management",
		Long:  `Manage TigerGraph CLI configuration including server aliases and credentials.`,
	}

	// Add command
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add server configuration",
		Run:   config.RunConfAdd,
	}
	addCmd.Flags().StringP("alias", "a", "", "Server alias name")
	addCmd.Flags().StringP("user", "u", "tigergraph", "TigerGraph user")
	addCmd.Flags().StringP("password", "p", "tigergraph", "TigerGraph password")
	addCmd.Flags().String("host", "http://127.0.0.1", "TigerGraph host")
	addCmd.Flags().String("gsPort", "14240", "GSQL Port")
	addCmd.Flags().String("restPort", "9000", "REST Port")
	addCmd.Flags().StringP("default", "d", "n", "Set as default alias (y/n)")

	// Delete command
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete server configuration",
		Run:   config.RunConfDelete,
	}
	deleteCmd.Flags().StringP("alias", "a", "", "Server alias to delete")
	deleteCmd.MarkFlagRequired("alias")

	// List command
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all configurations",
		Run:   config.RunConfList,
	}

	// TGCloud command
	var tgcloudCmd = &cobra.Command{
		Use:   "tgcloud",
		Short: "Configure TGCloud credentials",
		Run:   config.RunConfTGCloud,
	}
	tgcloudCmd.Flags().StringP("email", "e", "", "TGCloud email")
	tgcloudCmd.Flags().StringP("password", "p", "", "TGCloud password")

	confCmd.AddCommand(addCmd, deleteCmd, listCmd, tgcloudCmd)
	return confCmd
}
