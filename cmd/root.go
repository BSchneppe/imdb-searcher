/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"time"

	imdb_seeder "github.com/BSchneppe/imdb-searcher/internal/pkg/imdb-seeder"

	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

var cfg meilisearch.ClientConfig

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "This is a CLI tool to seed the IMDB database into MeiliSearch.",
	Long:  `This should be run as a cron job to keep the MeiliSearch database up to date with the latest IMDB data.`,
	Run: func(cmd *cobra.Command, args []string) {
		debugEnabled, _ := cmd.Flags().GetBool("debug")
		logger, _ := zap.NewProduction()
		if debugEnabled {
			logger, _ = zap.NewDevelopment()
		}
		if cfg.Host == "" {
			logger.Fatal("meili-host is required")
		}
		imdb_seeder.Seed(cfg, logger)
	},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "imdb-searcher",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringVar(&cfg.Host, "meili-host", "", "Host of your Meilisearch database")
	rootCmd.PersistentFlags().StringVar(&cfg.APIKey, "meili-api-key", "", "API Key for accessing Meilisearch")
	rootCmd.PersistentFlags().DurationVar(&cfg.Timeout, "meili-timeout", 5*time.Second, "Timeout duration (e.g., 5s, 1m)")
	rootCmd.AddCommand(seedCmd)
	err := rootCmd.MarkPersistentFlagRequired("meili-host")
	if err != nil {
		os.Exit(1)
	}

}
