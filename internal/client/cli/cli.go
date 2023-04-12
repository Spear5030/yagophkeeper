package cli

import (
	"fmt"
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"time"
)

type usecase interface {
	ListSecrets() []domain.LoginPassword
	AddLoginPassword(domain.LoginPassword) error
	RegisterUser(user domain.User) error
	LoginUser(user domain.User) error
	CheckSync() (time.Time, error)
	GetLocalSyncTime() time.Time
	SyncData() error
}

type CLI struct {
	logger  *zap.Logger
	usecase usecase
}

func New(logger *zap.Logger, usecase usecase) *CLI {
	c := CLI{logger: logger, usecase: usecase}
	c.ListSecrets()
	c.AddLoginPassword()
	c.RegisterUser()
	c.LoginUser()
	c.CheckSync()
	c.Sync()
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

func (cli *CLI) AddLoginPassword() {
	var lp = &domain.LoginPassword{}
	var addLPCmd = &cobra.Command{
		Use:   "add",
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
	rootCmd.AddCommand(addLPCmd)
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
