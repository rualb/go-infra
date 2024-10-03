package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-infra/internal/config"
	"go-infra/internal/tool/toolhttp"
	xlog "go-infra/internal/tool/toollog"
	"maps"
	"slices"
)

type dataSender struct {
	QueryData   map[string]string
	BodyForm    map[string]string
	BodyJSON    map[string]string
	HeadersData map[string]string
}

func newDataSender() dataSender {
	return dataSender{
		QueryData:   map[string]string{},
		BodyForm:    map[string]string{},
		BodyJSON:    map[string]string{},
		HeadersData: map[string]string{},
	}
}
func (sd dataSender) fillQuery(
	gw config.AppConfigMessageGateway,
	extract func(string) (string, error),
) error {

	if gw.Query != "" {

		if err := json.Unmarshal([]byte(gw.Query), &sd.QueryData); err != nil {
			return fmt.Errorf("query param: %v", err)
		}

		for _, key := range slices.Collect(maps.Keys(sd.QueryData)) {
			val, err := extract(key)
			if err != nil {
				return err
			}
			sd.QueryData[key] = val
		}
	}

	return nil

}

func (sd dataSender) fillBody(gw config.AppConfigMessageGateway,
	extract func(string) (string, error),
) error {

	// email
	if gw.Body != "" {

		if err := json.Unmarshal([]byte(gw.Body), &sd.BodyJSON); err != nil {
			return fmt.Errorf("body param: %v", err)
		}

		for _, key := range slices.Collect(maps.Keys(sd.BodyJSON)) {
			val, err := extract(key)
			if err != nil {
				return err
			}

			sd.BodyForm[key] = val
		}
	}

	return nil

}

func (sd dataSender) sendData(gw config.AppConfigMessageGateway) error {

	if gw.User != "" {
		auth := gw.User + ":" + gw.Password
		sd.HeadersData["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	}

	if gw.URL == "" {

		return fmt.Errorf("error gateway URL is empty")

	}

	respBody, err := toolhttp.PostFormURL(gw.URL, sd.QueryData, sd.BodyForm, sd.HeadersData)

	if gw.Stdout {
		if err != nil && len(respBody) > 0 {
			xlog.Info("Resp: %s", string(respBody))
		}
	}

	return err
}
