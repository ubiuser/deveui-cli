package register

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/NickGowdy/deveui-cli/client"
)

type CodeRegister struct {
	HttpClient client.Client
	Code       string
}

func (cr CodeRegister) RegisterCode() (*http.Response, error) {
	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": cr.Code}

	err := json.NewEncoder(b).Encode(&reqBody)
	if err != nil {
		log.Print(err)
	}

	resp, err := cr.HttpClient.Post("sensor-onboarding-sample", "application/json", b)
	if err != nil {
		log.Print(err)
	}

	return resp, nil
}
