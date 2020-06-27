package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nyks06/backapi/pg"
	webcore "github.com/nyks06/backapi"
	"github.com/nyks06/backapi/http"
	"github.com/nyks06/backapi/stripe"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"github.com/spf13/cobra"

	mailer "github.com/nyks06/backapi/mailjet"

	"github.com/stripe/stripe-go/client"
)

var (
	configPath string
)

func init() {
	serverCmd.Flags().StringVar(&configPath, "configpath", "./", "Defines here the config path where global conf file can be found")

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(versionCmd)
}

func initStoreManager() (*webcore.StoreManager, error) {
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

	ticketStore := &pg.TicketStore{
		DB: store.DB,
	}
	ticketFinder := &pg.TicketFinder{
		DB: store.DB,
	}

	pronoStore := &pg.PronosticStore{
		DB: store.DB,
	}
	pronoFinder := &pg.PronosticFinder{
		DB: store.DB,
	}

	sportStore := &pg.SportStore{
		DB: store.DB,
	}
	sportFinder := &pg.SportFinder{
		DB: store.DB,
	}

	competitionStore := &pg.CompetitionStore{
		DB: store.DB,
	}
	competitionFinder := &pg.CompetitionFinder{
		DB: store.DB,
	}

	return &webcore.StoreManager{
		UserStore:  userStore,
		UserFinder: userFinder,

		SessionStore:  sessionStore,
		SessionFinder: sessionFinder,

		TicketStore:  ticketStore,
		TicketFinder: ticketFinder,

		PronosticStore:  pronoStore,
		PronosticFinder: pronoFinder,

		SportStore:  sportStore,
		SportFinder: sportFinder,

		CompetitionStore:  competitionStore,
		competitionFinder: competitionFinder,
	}, nil
}

func runServer(cmd *cobra.Command, args []string) {
	//config setup
	cfg := webcore.Config{}
	if err := confita.NewLoader(file.NewBackend(configPath)).Load(context.Background(), &cfg); err != nil {
		panic(err)
	}

	storeManager, err := initStoreManager()
	if err != nil {
		panic(err)
	}

	mailer := mailer.NewMailer(cfg.Mailer.FromMail, cfg.Mailer.FromUser, cfg.Mailer.MailjetAPIKey, cfg.Mailer.MailjetAPISecret)
	sc := &client.API{}
	sc.Init(cfg.Stripe.SecretKey, nil)
	paymentClient := stripe.NewPaymentClient(sc)

	// Create Services with store
	APIService := &webcore.APIService{
		UserStore:  userStore,
		UserFinder: userFinder,

		TicketStore:  ticketStore,
		TicketFinder: ticketFinder,

		PronosticsStore:  pronoStore,
		PronosticsFinder: pronoFinder,

		SportsStore:  sportStore,
		SportsFinder: sportFinder,

		PlayersStore:  playerStore,
		PlayersFinder: playersFinder,

		NewsStore:  newsStore,
		NewsFinder: newsFinder,

		CompetitionsStore:  competitionStore,
		CompetitionsFinder: competitionFinder,

		SessionFinder: sessionFinder,
		SessionStore:  sessionStore,

		Mailer:        mailer,
		PaymentClient: paymentClient,
	}

	// Load Server with services
	s := http.NewServer(APIService)
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
		fmt.Printf("Core Server Version : %s\n", webcore.Version)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
