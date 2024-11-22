package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/serpentine"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "example [args]",
		Short: "A simple example program!",
		Long: `A simple example program!

This is an example program we made for serpentine!`,
		Example: `
# Run it:
example

# Run it with some argument:
example --foo=xyz
		`,
		Args: cobra.ArbitraryArgs,
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var foo string
	var bar int
	var zaz float64
	var d time.Duration
	cmd.PersistentFlags().StringVarP(&foo, "surname", "s", "doe", "Your surname")
	cmd.Flags().StringVar(&foo, "name", "john", "Your name")
	cmd.Flags().DurationVar(&d, "duration", 0, "Time since your last commit")
	cmd.Flags().IntVar(&bar, "age", 0, "Your age")
	cmd.Flags().Float64Var(&zaz, "idk", 0.0, "I don't know")

	cmd.Flags().MarkHidden("age")
	cmd.Flags().MarkHidden("duration")
	cmd.Flags().MarkHidden("idk")

	serpentine.Setup(cmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "sub [command] [flags] [args]",
		Short: "An example sub command",
	})
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
