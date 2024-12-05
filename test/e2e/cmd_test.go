package e2e

import (
	xcmd "go-infra/internal/cmd"
	"go-infra/internal/util/utilhttp"
	"os"
	"strings"
	"testing"
	"time"
)

// TestHealthController_Check_Stats tests the ?cmd=stats case in the Check method
func TestCmd(t *testing.T) {
	// Setup Echo context

	//
	os.Setenv("APP_ENV", "testing")
	os.Setenv("APP_SMS_GW_STDOUT", "1")
	os.Setenv("APP_EMAIL_GW_STDOUT", "1")

	cmd := xcmd.Command{}

	go cmd.Exec()

	time.Sleep(3 * time.Second)

	urls := []struct {
		title  string
		url    string
		form   map[string]string
		search string
	}{
		{title: "test email-secret-code", search: "123456789", url: "http://127.0.0.1:30780/sys/api/messenger/email-secret-code", form: map[string]string{"to": "test@example.com", "secret_code": "123456789", "lang": "en"}},
		{title: "test sms-secret-code", search: "123456789", url: "http://127.0.0.1:30780/sys/api/messenger/sms-secret-code", form: map[string]string{"to": "123456", "secret_code": "123456789", "lang": "en"}},
		{title: "test email-html", search: "123456789", url: "http://127.0.0.1:30780/sys/api/messenger/email-html", form: map[string]string{"to": "test@example.com", "html": "123456789"}},
		{title: "test sms-text", search: "123456789", url: "http://127.0.0.1:30780/sys/api/messenger/sms-text", form: map[string]string{"to": "123456", "text": "123456789"}},
	}

	for _, itm := range urls {

		t.Run(itm.title, func(t *testing.T) {

			t.Logf("url %v", itm.url)
			arr, err := utilhttp.PostFormURL(itm.url,
				nil, itm.form, nil,
			)

			if err != nil {
				t.Errorf("Error : %v", err)
			}

			if !strings.Contains(string(arr), itm.search) {
				t.Errorf("Error on %v", itm.url)
			}

		})

	}

	cmd.Stop()

	time.Sleep(1 * time.Second)

}
