package main

import (
	"fmt"
	"github.com/altmer/bellboy/context"
	"github.com/altmer/bellboy/media"
	"github.com/altmer/bellboy/tumblr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "0.0.1"

func main() {
	ensureBellboyDir()

	viper.SetConfigName("bellboy")
	viper.AddConfigPath(bellboyDirPath())

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	db := context.NewDBConnection(viper.GetString("db"))
	defer db.Close()

	fmt.Printf("Opened database at [%s]\n", viper.GetString("db"))
	fmt.Printf("Media folder is [%s]\n", viper.GetString("media_folder"))

	syncer := tumblr.Syncer{
		BlogName: viper.GetString("tumblr.blog"),
		Client:   tumblr.New(viper.GetStringMapString("tumblr")),
		Repo:     media.NewRepository(db),
	}

	var cmdSync = &cobra.Command{
		Use:   "sync",
		Short: "Sync tumblr blog posts and likes",
		Run: func(cmd *cobra.Command, args []string) {
			syncer.Sync()
		},
	}

	var cmdSubsDown = &cobra.Command{
		Use:   "subsdown",
		Short: "Loads tumblr subscriptions to local DB (destructive operation!)",
		Run: func(cmd *cobra.Command, args []string) {
			syncer.SubsDown()
		},
	}

	var cmdSubsUp = &cobra.Command{
		Use:   "export",
		Short: "Exports subscriptions from local DB to tumblr blog (follows all blogs)",
		Run: func(cmd *cobra.Command, args []string) {
			syncer.SubsUp()
		},
	}

	var rootCmd = &cobra.Command{Use: "bellboy"}
	rootCmd.AddCommand(cmdSync, cmdSubsDown, cmdSubsUp)
	rootCmd.Execute()
}
