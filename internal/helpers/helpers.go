package helpers

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"github.com/zrougamed/tgCli/internal/models"
)

func CreateDefaultConfig(configFile string) error {
	defaultConfig := models.Config{
		TGCloud: models.TGCloudConfig{
			User:     "mail@domain.com",
			Password: "",
		},
		Machines: make(map[string]models.MachineConfig),
		Default:  "",
	}

	viper.Set("tgcloud", defaultConfig.TGCloud)
	viper.Set("machines", defaultConfig.Machines)
	viper.Set("default", defaultConfig.Default)

	if err := viper.WriteConfigAs(configFile); err != nil {
		log.Printf("Error creating default config: %v", err)
		return fmt.Errorf("unable to create default config file: %w", err)
	}
	return nil
}

func SaveConfig() error {
	return viper.WriteConfig()
}

func GracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nTerminating tgcli, Good Bye!")
		os.Exit(0)
	}()
}

func CheckForUpdates() (string, error) {
	// For now, just placeholder
	return "N/A", nil
}
