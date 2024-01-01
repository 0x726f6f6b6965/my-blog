package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/0x726f6f6b6965/my-blog/db/migrations"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/drivers"
)

var (
	cfgPath      string
	driverConfig drivers.Config
)

func main() {
	cobra.OnInitialize(initConfig)

	var rootCmd = &cobra.Command{
		Use:           "migrate [flags] <driver>",
		Short:         "Event DB migration tool.",
		Example:       `migrate psql`,
		PreRun:        getConfig,
		RunE:          run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "Filename of config file to override default lookup")

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
}

func getConfig(cmd *cobra.Command, args []string) {
	driverName := "psql"

	// Configure the driver
	driverConfig = map[string]interface{}{
		"whitelist": viper.GetStringSlice(driverName + ".whitelist"),
		"blacklist": viper.GetStringSlice(driverName + ".blacklist"),
	}

	for _, key := range viper.AllKeys() {
		if strings.HasPrefix(key, fmt.Sprintf("%s.", driverName)) {

			driverConfig[strings.ReplaceAll(key, fmt.Sprintf("%s.", driverName), "")] = viper.Get(key)
		}
	}
}

func run(cmd *cobra.Command, args []string) error {
	user := driverConfig.MustString(drivers.ConfigUser)
	pass, _ := driverConfig.String(drivers.ConfigPass)
	dbname := driverConfig.MustString(drivers.ConfigDBName)
	host := driverConfig.MustString(drivers.ConfigHost)
	port := driverConfig.DefaultInt(drivers.ConfigPort, 5432)
	sslmode := driverConfig.DefaultString(drivers.ConfigSSLMode, "require")

	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", user, pass, host, port, dbname, sslmode)
	err := migrations.MigrateSchema(context.Background(), uri)
	if err != nil {
		return fmt.Errorf("migrate failed to migrate to the latest schema, %+v/n", err)
	}

	return nil
}

func initConfig() {
	if cfgPath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgPath)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
