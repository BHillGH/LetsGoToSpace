package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type registerStruct struct {
	Symbol  string `json:"symbol"`
	Faction string `json:"faction"`
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Sets up your environment",
	Long:  `Creates a .env file with your callsign and stuff`,
	Run: func(cmd *cobra.Command, args []string) {
		val, ok := os.LookupEnv("symbol")
		if !ok {
			setupSymbol()
			registerAgent()
		} else {
			fmt.Println("var set: ", val)
		}

	},
}

func setupSymbol() {
	var symbol string

	fmt.Println("Enter your unique symbol:")
	fmt.Scanln(&symbol)
	os.Setenv("symbol", symbol)

	fmt.Println("Your symbol is ", os.Getenv("symbol"))

}

func registerAgent() {
	url := "https://api.spacetraders.io/v2/register"

	body := &registerStruct{
		Symbol:  os.Getenv("symbol"),
		Faction: "COSMIC",
	}

	client := &http.Client{}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	req, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		panic(res.Status)
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var out bytes.Buffer

	err = json.Indent(&out, responseBody, "", "    ")
	if err != nil {
		log.Fatalf("Error formatting JSON: %v", err)
	}

	// Print the indented JSON
	fmt.Println("Be sure to keep this token safe! \n", out.String())

}
