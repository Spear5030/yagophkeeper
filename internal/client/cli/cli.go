package cli

import (
	"fmt"
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type usecase interface {
	ListSecrets() []domain.LoginPassword
	AddLoginPassword(domain.LoginPassword) error
	RegisterUser(user domain.User) (string, error)
	LoginUser(user domain.User) (string, error)
}

type CLI struct {
	logger  *zap.Logger
	usecase usecase
}

func New(logger *zap.Logger, usecase usecase) *CLI {
	c := CLI{logger: logger, usecase: usecase}
	c.AddListSecrets()
	c.AddLoginPassword()
	return &c
}

func (cli *CLI) AddListSecrets() {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "Print secrets",
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
	var user = &domain.User{}
	var regUserCmd = &cobra.Command{
		Use:   "register",
		Short: "register account",
		Long:  `register account`,
		Run: func(cmd *cobra.Command, args []string) {
			token, err := cli.usecase.RegisterUser(*user)
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
