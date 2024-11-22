package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/serpentine"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "foo",
		Short: "foo is bar",
		Long: `Foo is bar!
This is a new thing we made!`,
		Example: `
# bla bla bla
example

# bla bla bla
example --foo=xyz
		`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// nothing
		},
	}
	var foo string
	cmd.PersistentFlags().StringVarP(&foo, "foo", "f", "abc", "foo does something")
	cmd.Flags().StringVar(&foo, "bar", "lalala", "something something something")

	serpentine.Setup(cmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "another",
		Short: "hello world",
	})
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
