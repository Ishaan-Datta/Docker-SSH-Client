package cmd

import (
	"fmt"
	"math/rand"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use: "generate",
	Short: "Generate random passwords",
	Long: `Generate random passwords with customizable options.
	For example:

	passwordgen generate -l 12 -d -s`,

	Run: generatePassword,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().IntP("length", "l", 8, "Length of the generated password")
	generateCmd.Flags().BoolP("digits", "d", false, "Include digitis in the generated password")
	generateCmd.Flags().BoolP("special_chars", "s", false, "Include special chars in generated password")
}

func generatePassword(cmd *cobra.Command, args []string) {
	length, _ := cmd.Flags().GetInt("length")
	isDigits, _ := cmd.Flags().GetBool("digits")
	isSpecialChars, _ := cmd.Flags().GetBool("special-chars")

	charset := "alphabet"

	if isDigits {
		charset += "0123"
	}

	if isSpecialChars {
		charset += "!A@%99231"
	}

	password := make([]byte, length)

	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}

	fmt.Println("Generating password")
	fmt.Println(password)
}