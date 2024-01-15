package commands

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/report"
	"github.com/photoprism/photoprism/pkg/rnd"
)

// AuthAddFlags specifies the "photoprism auth add" command flags.
var AuthAddFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "name, n",
		Usage: "arbitrary name to help identify the access `TOKEN`",
	},
	cli.StringFlag{
		Name:  "scope, s",
		Usage: "authorization `SCOPE` for the access token e.g. \"metrics\" or \"photos albums\" (\"*\" to allow all scopes)",
	},
	cli.Int64Flag{
		Name:  "expires, e",
		Usage: "access token lifetime in `SECONDS`, after which it expires and a new token must be created (-1 to disable)",
		Value: entity.UnixYear,
	},
}

// AuthAddCommand configures the command name, flags, and action.
var AuthAddCommand = cli.Command{
	Name:        "add",
	Usage:       "Creates a new client access token",
	Description: "Specifying a username as argument creates a personal access token for a registered user account.",
	ArgsUsage:   "[username]",
	Flags:       AuthAddFlags,
	Action:      authAddAction,
}

// authAddAction shows detailed session information.
func authAddAction(ctx *cli.Context) error {
	return CallWithDependencies(ctx, func(conf *config.Config) error {
		// Get username from command flag.
		userName := clean.Username(ctx.Args().First())

		// Find user account.
		user := entity.FindUserByName(userName)

		if user == nil && userName != "" {
			return fmt.Errorf("user %s not found", clean.LogQuote(userName))
		}

		// Get token name from command flag or ask for it.
		tokenName := ctx.String("name")

		if tokenName == "" {
			prompt := promptui.Prompt{
				Label:   "Token Name",
				Default: rnd.Name(),
			}

			res, err := prompt.Run()

			if err != nil {
				return err
			}

			tokenName = clean.Name(res)
		}

		// Get auth scope from command flag or ask for it.
		authScope := ctx.String("scope")

		if authScope == "" {
			prompt := promptui.Prompt{
				Label:   "Authorization Scope",
				Default: "*",
			}

			res, err := prompt.Run()

			if err != nil {
				return err
			}

			authScope = clean.Scope(res)
		}

		// Create client session.
		sess, err := entity.CreateClientAccessToken(tokenName, ctx.Int64("expires"), authScope, user)

		if err != nil {
			return fmt.Errorf("failed to create access token: %s", err)
		} else {
			// Show client authentication credentials.
			if sess.UserUID == "" {
				fmt.Printf("\nPLEASE WRITE DOWN THE FOLLOWING RANDOMLY GENERATED CLIENT ACCESS TOKEN, AS YOU WILL NOT BE ABLE TO SEE IT AGAIN:\n")
			} else {
				fmt.Printf("\nPLEASE WRITE DOWN THE FOLLOWING RANDOMLY GENERATED PERSONAL ACCESS TOKEN, AS YOU WILL NOT BE ABLE TO SEE IT AGAIN:\n")
			}

			result := report.Credentials("Access Token", sess.AuthToken(), "Authorization Scope", sess.Scope())

			fmt.Printf("\n%s\n", result)
		}

		return err
	})
}
