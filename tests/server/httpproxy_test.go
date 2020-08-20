package server_test

import (
	"encoding/json"
	"github.com/seknox/trasa/server/api/auth/serviceauth"
	"github.com/seknox/trasa/server/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHTTPAccessProxy(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name        string
		args        args
		wantSuccess bool
	}{
		{
			"should fail when hostname is incorrect",
			args{getreqWithBody(t, serviceauth.NewSession{
				HostName:  "gitlab01.trasa.io",
				TfaMethod: "",
				TotpCode:  getTotpCode(totpSEC),
				ExtToken:  "cb6dd3f6-54c2-4cb0-b294-e22c2aa708e4",
			})},
			false,
		},

		{
			"should fail when ext token is incorrect",
			args{getreqWithBody(t, serviceauth.NewSession{
				HostName:  "gitlab01.trasa.io",
				TfaMethod: "",
				TotpCode:  getTotpCode(totpSEC),
				ExtToken:  "db6dd3f6-54c2-4cb0-b294-e22c2aa708e4",
			})},
			false,
		},

		{
			"should fail when totp is incorrect",
			args{getreqWithBody(t, serviceauth.NewSession{
				HostName:  "gitlab01.trasa.io",
				TfaMethod: "",
				TotpCode:  "123456",
				ExtToken:  "cb6dd3f6-54c2-4cb0-b294-e22c2aa708e4",
			})},
			false,
		},

		{
			"should fail if service is not authorised",
			args{getreqWithBody(t, serviceauth.NewSession{
				HostName:  "test00.trasa.io",
				TfaMethod: "",
				TotpCode:  getTotpCode(totpSEC),
				ExtToken:  "cb6dd3f6-54c2-4cb0-b294-e22c2aa708e4",
			})},
			false,
		},

		{
			"should pass",
			args{getreqWithBody(t, serviceauth.NewSession{
				HostName:  "gitlab01.trasa.io",
				TfaMethod: "U2F",
				TotpCode:  "",
				ExtToken:  "cb6dd3f6-54c2-4cb0-b294-e22c2aa708e4",
			})},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(serviceauth.AgentLogin)

			handler.ServeHTTP(rr, tt.args.req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			var resp models.TrasaResponseStruct
			err := json.Unmarshal(rr.Body.Bytes(), &resp)
			if err != nil {
				t.Fatal(err)
			}

			if tt.wantSuccess && resp.Status != "success" {
				t.Errorf("AgentLogin() wanted success, got:%s reason %s", resp.Status, resp.Reason)
				return
			}

		})
	}

}
