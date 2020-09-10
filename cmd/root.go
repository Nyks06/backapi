package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nyks06/backapi"
	"github.com/nyks06/backapi/http"
	"github.com/nyks06/backapi/pg"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"github.com/spf13/cobra"
)

var (
	configPath string
)

func init() {
	serverCmd.Flags().StringVar(&configPath, "configpath", "./", "Defines here the config path where global conf file can be found")

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(versionCmd)
}

func initStoreManager() (*backapi.StoreManager, error) {
	store, err := pg.NewStore(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	userStore := &pg.UserStore{
		DB: store.DB,
	}
	userFinder := &pg.UserFinder{
		DB: store.DB,
	}

	sessionStore := &pg.SessionStore{
		DB: store.DB,
	}
	sessionFinder := &pg.SessionFinder{
		DB: store.DB,
	}

	return &backapi.StoreManager{
		UserStore:  userStore,
		UserFinder: userFinder,

		SessionStore:  sessionStore,
		SessionFinder: sessionFinder,
	}, nil
}

func runServer(cmd *cobra.Command, args []string) {
	//config setup
	cfg := backapi.Config{}
	if err := confita.NewLoader(file.NewBackend(configPath)).Load(context.Background(), &cfg); err != nil {
		panic(err)
	}

	storeManager, err := initStoreManager()
	if err != nil {
		panic(err)
	}

	userService := &backapi.UserService{
		StoreManager: storeManager,
	}

	sessionService := &backapi.SessionService{
		StoreManager: storeManager,
	}

	// Load Server with services
	s := http.NewServer(userService, sessionService, cfg.HTTP.Port)
	if err := s.Setup(); err != nil {
		panic(err)
	}

	if err := s.Start(); err != nil {
		panic(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "live",
	Short: "The solution running both the backend and the frontend of the website",
	Long:  "Live is the binary used to run the whole platform - being both the backend and the frontend",
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Use this command to run the server.",
	Long:  "This command will run the server using given args and the config given as command line if any.",
	Run:   runServer,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of the core.",
	Long:  "Print the current version used to run the server in the given binary",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Core Server Version : %s\n", backapi.Version)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
