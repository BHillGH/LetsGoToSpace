package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/lpernett/godotenv"
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
	fmt.Println(os.Getenv("symbol"))
	fmt.Println("Enter your unique symbol:")
	fmt.Scanln(&symbol)
	os.Setenv("symbol", symbol)

	fmt.Println("Your symbol is ", os.Getenv("symbol"))

}

func registerAgent() {
	url := "https://api.spacetraders.io/v2/register"

	// Create body in JSON to send in request
	body := &registerStruct{
		Symbol:  os.Getenv("symbol"),
		Faction: "COSMIC",
	}

	// Creating a new client
	client := &http.Client{}

	// Create the buffer and encode the body to json and create the request
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	req, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	// Do the actual client request
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		panic(res.Status)
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Step 2: Unmarshal JSON into the struct
	var responseVar responseStruct
	err = json.Unmarshal(responseBody, &responseVar)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Now you can access the response fields
	fmt.Printf("Contract Type: %s\n", responseVar.Data.Token)

	// Optionally, you can still pretty-print the JSON for debugging
	var out bytes.Buffer
	err = json.Indent(&out, responseBody, "", "    ")
	if err != nil {
		log.Fatalf("Error formatting JSON: %v", err)
	}
	fmt.Println("Be sure to keep this token safe!")

}

func writeTokenToDotenv(token string) {
	envMap := map[string]string{
		"token": token,
	}
	godotenv.Write(envMap, ".env")

}

type responseStruct struct {
	Data struct {
		Token string `json:"token"`
		Agent struct {
			AccountID       string `json:"accountId"`
			Symbol          string `json:"symbol"`
			Headquarters    string `json:"headquarters"`
			Credits         int    `json:"credits"`
			StartingFaction string `json:"startingFaction"`
			ShipCount       int    `json:"shipCount"`
		} `json:"agent"`
	} `json:"data"`
}
