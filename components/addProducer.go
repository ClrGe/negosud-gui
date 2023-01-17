package components

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"net/http"
)

var newProducer []Producer

func postNewProducer(newProducer Producer) error {
	env, err := LoadConfig(".")
	if err != nil {
		return err
	}
	// convert producer struct to json
	producerJSON, err := json.Marshal(newProducer)
	if err != nil {
		return err
	}
	// create http client and request
	client := &http.Client{}
	req, err := http.NewRequest("POST", env.SERVER+"/api/producer", bytes.NewBuffer(producerJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	// make request
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		return fmt.Errorf("error posting new producer, status code: %d", res.StatusCode)
	}
	return nil
}

func producerForm(_ fyne.Window) fyne.CanvasObject {
	form := &widget.Form{}
	form.Append("Nom:", widget.NewEntry())
	form.Append("Details:", widget.NewMultiLineEntry())
	form.Append("Region:", widget.NewEntry())

	return form
}
