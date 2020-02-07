package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/manifoldco/promptui"
)

type Artist struct {
	ID   string
	Name string
}

func main() {

	artist, err := connectNewArtist()
	if err != nil {
		fmt.Println("Error connecting new artist:", err)
		return
	}

	for {
		err := promptForCommand(artist)
		if err != nil {
			fmt.Println("Error", err)
		}
	}
}

func connectNewArtist() (Artist, error) {
	namePrompt := promptui.Prompt{
		Label: "Name?",
	}

	name, err := namePrompt.Run()
	if err != nil {
		return Artist{}, err
	}

	artist := Artist{
		Name: name,
	}

	artistData, err := json.Marshal(artist)
	if err != nil {
		return Artist{}, err
	}

	resp, err := http.Post("http://localhost:8080/api/artists", "application/json", bytes.NewBuffer(artistData))
	if err != nil {
		return Artist{}, err
	}
	if resp.StatusCode != 200 {
		return Artist{}, fmt.Errorf("Unexpected status code when creating new artist", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&artist)
	if err != nil {
		return Artist{}, err
	}

	return artist, nil
}

type Command struct {
	Description string
}

func promptForCommand(artist Artist) error {

	drawPrompt := promptui.Prompt{
		Label: "Draw something?",
	}

	commandTxt, err := drawPrompt.Run()
	if err != nil {
		return err
	}

	command := Command{
		Description: commandTxt,
	}

	commandData, err := json.Marshal(command)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8080/api/artists/"+artist.ID+"/moves", "application/json", bytes.NewBuffer(commandData))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Unexpected status code when creating new artist", resp.StatusCode)
	}

	return nil
}
