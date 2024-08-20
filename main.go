/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"LetsGoToSpace/cmd"

	"github.com/lpernett/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	cmd.Execute()
}
