/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helm-api",
	Short: "Rest API Web Server to manage the install and uninstall of Hel Charts",
	Long:  `Rest API Web Server to manage the install and uninstall of Hel Charts`,
}

func GenerateDoc() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gen_docs",
		Short: "This command will build the documents for the cli",
		Long:  `This command will build the documents for the cli`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := doc.GenMarkdownTree(rootCmd, "./documents")
			if err != nil {
				return err
			}
			fmt.Println("Documents Generated")
			return nil
		},
	}
	return cmd
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(GenerateDoc())
}
