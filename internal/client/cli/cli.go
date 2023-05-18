package cli

import (
	"fmt"
	"github.com/Spear5030/yagophkeeper/internal/client/cli/tui"
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"time"
)

type usecase interface {
	GetLoginsPasswords() []domain.LoginPassword
	ListSecrets() []string
	AddLoginPassword(domain.LoginPassword) error
	AddTextData(domain.TextData) error
	AddBinaryData(domain.BinaryData) error
	AddCardData(domain.CardData) error
	RegisterUser(user domain.User) error
	LoginUser(user domain.User) error
	CheckSync() (time.Time, error)
	GetLocalSyncTime() time.Time
	SyncData() error
	GetVersion() string
	GetBuildTime() string
}

type CLI struct {
	logger  *zap.Logger
	usecase usecase
}

func New(logger *zap.Logger, usecase usecase) *CLI {
	c := CLI{logger: logger, usecase: usecase}
	c.ListSecrets()
	c.RegisterUser()
	c.LoginUser()
	c.CheckSync()
	c.Sync()
	c.AddLPCmd()
	c.AddCardCmd()
	c.AddTextCmd()
	c.AddBinaryCmd()
	c.Version()
	rootCmd.AddCommand(addCmd)
	rootCmd.Run = func(cmd *cobra.Command, args []string) { tui.StartTUI(usecase) }
	return &c
}

func (cli *CLI) ListSecrets() {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "Print secrets. ",
		Long:  `Print all local secrets`,
		Run: func(cmd *cobra.Command, args []string) {
			secrets := cli.usecase.ListSecrets()
			for _, secret := range secrets {
				fmt.Println(secret)
			}
		},
	}
	rootCmd.AddCommand(listCmd)
}

func (cli *CLI) AddLPCmd() {
	var lp = &domain.LoginPassword{}
	var addLPCmd = &cobra.Command{
		Use:   "login",
		Short: "add login-password secret",
		Long:  `add login-password secret`,
		Run: func(cmd *cobra.Command, args []string) {
			err := cli.usecase.AddLoginPassword(*lp)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
	addLPCmd.Flags().StringVarP(&lp.Login, "login", "l", "", "login (required)")
	addLPCmd.MarkFlagRequired("login")
	addLPCmd.Flags().StringVarP(&lp.Password, "password", "p", "", "password (required)")
	addLPCmd.MarkFlagRequired("password")
	addLPCmd.Flags().StringVarP(&lp.Meta, "meta", "m", "", "meta field")
	addCmd.AddCommand(addLPCmd)
}

func (cli *CLI) AddTextCmd() {
	var td = &domain.TextData{}
	var addTextCmd = &cobra.Command{
		Use:   "text",
		Short: "add text secret",
		Long:  `add text secret`,
		Run: func(cmd *cobra.Command, args []string) {
			err := cli.usecase.AddTextData(*td)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
	addTextCmd.Flags().StringVarP(&td.Text, "text", "t", "", "text (required)")
	addTextCmd.MarkFlagRequired("text")
	addTextCmd.Flags().StringVarP(&td.Meta, "meta", "m", "", "meta field")
	addCmd.AddCommand(addTextCmd)
}

func (cli *CLI) AddBinaryCmd() {
	var bd = &domain.BinaryData{}
	var path string
	var err error
	var AddBinaryCmd = &cobra.Command{
		Use:   "binary",
		Short: "add binary secret",
		Long:  `add binary secret`,
		Run: func(cmd *cobra.Command, args []string) {
			bd.BinaryData, err = os.ReadFile(path)
			if err != nil {
				fmt.Println(err)
			}
			err = cli.usecase.AddBinaryData(*bd)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
	AddBinaryCmd.Flags().StringVarP(&path, "path", "p", "", "path to binary file (required)")
	AddBinaryCmd.MarkFlagRequired("path")
	AddBinaryCmd.Flags().StringVarP(&bd.Meta, "meta", "m", "", "meta field")

	addCmd.AddCommand(AddBinaryCmd)
}

func (cli *CLI) AddCardCmd() {
	var card = &domain.CardData{}
	var addCardCmd = &cobra.Command{
		Use:   "card",
		Short: "add card secret",
		Long:  `add card secret`,
	}
	addCardCmd.Flags().StringVarP(&card.Number, "number", "n", "", "number (required)")
	addCardCmd.MarkFlagRequired("number")
	addCardCmd.Flags().StringVarP(&card.CVC, "cvc", "v", "", "cvc (required)")
	addCardCmd.MarkFlagRequired("cvc")
	addCardCmd.Flags().StringVarP(&card.CardHolder, "cardholder", "", "", "card holder")
	addCardCmd.Flags().StringVarP(&card.Meta, "meta", "m", "", "meta field")
	addCmd.AddCommand(addCardCmd)
}

func (cli *CLI) RegisterUser() {
	var user = domain.User{}
	var regUserCmd = &cobra.Command{
		Use:   "register",
		Short: "register account",
		Long:  `register account`,
		Run: func(cmd *cobra.Command, args []string) {
			cli.logger.Debug("RegisterUser")
			err := cli.usecase.RegisterUser(user)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
	regUserCmd.Flags().StringVarP(&user.Email, "email", "l", "", "email (required)")
	regUserCmd.MarkFlagRequired("email")
	regUserCmd.Flags().StringVarP(&user.Password, "password", "p", "", "password (required)")
	regUserCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(regUserCmd)
}

func (cli *CLI) Sync() {
	var syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "sync secrets",
		Long:  `sync secrets with server`,
		Run: func(cmd *cobra.Command, args []string) {
			err := cli.usecase.SyncData()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("All secrets synced") //todo return last sync time
		},
	}
	rootCmd.AddCommand(syncCmd)
}

func (cli *CLI) CheckSync() {
	var checkSyncCmd = &cobra.Command{
		Use:   "checksync",
		Short: "get last sync time from server",
		Long:  `get last sync time from server`,
		Run: func(cmd *cobra.Command, args []string) {
			t, err := cli.usecase.CheckSync()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Secrets on server last sync time:", t)
			t = cli.usecase.GetLocalSyncTime()
			fmt.Println("Local secrets last sync time:", t)
		},
	}
	rootCmd.AddCommand(checkSyncCmd)
}

func (cli *CLI) LoginUser() {
	var user = domain.User{}
	var logUserCmd = &cobra.Command{
		Use:   "login",
		Short: "login account",
		Long:  `login account`,
		Run: func(cmd *cobra.Command, args []string) {
			err := cli.usecase.LoginUser(user)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
	logUserCmd.Flags().StringVarP(&user.Email, "email", "l", "", "email (required)")
	logUserCmd.MarkFlagRequired("email")
	logUserCmd.Flags().StringVarP(&user.Password, "password", "p", "", "password (required)")
	logUserCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(logUserCmd)
}

func (cli *CLI) Version() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "get version",
		Long:  `get version and build time`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version:" + cli.usecase.GetVersion())
			fmt.Println("Build time:" + cli.usecase.GetBuildTime())
		},
	}
	rootCmd.AddCommand(versionCmd)
}
