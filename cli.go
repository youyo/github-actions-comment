package ghc

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Run
func Run(cmd *cobra.Command, args []string) error {
	title := viper.GetString("title")
	body := viper.GetString("body")
	failure := viper.GetBool("failure")

	g := New()

	if err := g.GetCommentUrl(); err != nil {
		return err
	}

	comment, err := g.GenerateComment(title, body, failure)
	if err != nil {
		return err
	}

	err = g.PostComment(comment)

	return err
}
