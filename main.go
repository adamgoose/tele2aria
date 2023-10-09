package main

import (
	"log"
	"path/filepath"

	"github.com/adamgoose/tele2aria/cmd"
	"github.com/adamgoose/tele2aria/lib"
	"github.com/defval/di"
	"github.com/spf13/viper"
	"github.com/zelenin/go-tdlib/client"
)

func main() {
	if err := lib.App.Apply(

		di.Provide(func() (*client.Client, error) {
			authorizer := client.ClientAuthorizer()
			go client.CliInteractor(authorizer)

			authorizer.TdlibParameters <- &client.SetTdlibParametersRequest{
				UseTestDc:              false,
				DatabaseDirectory:      filepath.Join(viper.GetString("TDLIB"), "database"),
				FilesDirectory:         filepath.Join(viper.GetString("TDLIB"), "files"),
				UseFileDatabase:        true,
				UseChatInfoDatabase:    true,
				UseMessageDatabase:     true,
				UseSecretChats:         false,
				ApiId:                  viper.GetInt32("TELEGRAM_APP_ID"),
				ApiHash:                viper.GetString("TELEGRAM_API_HASH"),
				SystemLanguageCode:     "en",
				DeviceModel:            "Server",
				SystemVersion:          cmd.Version,
				ApplicationVersion:     cmd.Version,
				EnableStorageOptimizer: true,
				IgnoreFileNames:        false,
			}

			_, err := client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
				NewVerbosityLevel: viper.GetInt32("TELEGRAM_VERBOSITY"),
			})
			if err != nil {
				log.Fatalf("SetLogVerbosityLevel error: %s", err)
			}

			return client.NewClient(authorizer)
		}),
	); err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}

func init() {
	viper.SetEnvPrefix("TELE2ARIA")
	viper.AutomaticEnv()

	viper.SetDefault("TELEGRAM_VERBOSITY", 0)
	viper.SetDefault("TDLIB", ".tdlib")
}
