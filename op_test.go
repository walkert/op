package op

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var item = `{"uuid":"randogoo","templateUuid":"001","trashed":"N","createdAt":"2019-04-09T13:20:52Z","updatedAt":"2019-04-17T00:48:26Z","changerUuid":"uusid1","itemVersion":2,"vaultUuid":"rando1","details":{"fields":[{"designation":"username","name":"username","type":"T","value":"user@bar.com"},{"designation":"password","name":"password","type":"P","value":"greatpass"}],"sections":[{"name":"linked items","title":"Related Items"},{"fields":[{"k":"concealed","n":"TOTP_foo","t":"one-time password","v":"otpauth://totp"}],"name":"Section_3"}]},"overview":{"URLs":[{"u":"https://rando.com"}],"ainfo":"doofg@bar.com","pbe":85.38914735493256,"pgrng":true,"ps":100,"title":"FOOBAR","url":"https://foo.com"}}`

var configData = `{"latest_signin": "my_team","accounts": [{"shorthand": "my_team","url": "https://my_team.1password.com","email": "user@myteam.com","accountKey": "some-key-here","userUUID": "somuuid"}]}`

// mockCmd passes the real args to the underlying test executable
// by running TestHelperProcess directly. Lifted from exec_test.go
func mockCmd(name string, args ...string) (cmd *exec.Cmd) {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, name)
	cs = append(cs, args...)
	cmd = exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// declare our mock implementation of the read interface
type mockConfiger struct{}

func (m mockConfiger) Read() ([]byte, error) {
	return []byte(configData), nil
}

func TestTotp(t *testing.T) {
	configImpl = mockConfiger{}
	o, err := New(withCmdFunc(mockCmd))
	if err != nil {
		t.Fatal(err)
	}
	want := "123456"
	got, err := o.GetTotp("foo")
	if string(want) != "123456" {
		t.Fatalf("Got: %s, want: %s\n", got, want)
	}
}

func TestGetUserPass(t *testing.T) {
	configImpl = mockConfiger{}
	tests := []struct {
		name     string
		item     string
		wantUser string
		wantPass string
		wantErr  bool
		err      string
	}{
		{
			"ValidItem",
			"FOOBAR",
			"user@bar.com",
			"greatpass",
			false,
			"",
		},
		{
			"InvalidItem",
			"invalid",
			"",
			"",
			true,
			`error running [get item invalid]: item not found
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := New(withCmdFunc(mockCmd))
			if err != nil {
				t.Fatal(err)
			}
			user, pass, err := o.GetUserPass(tt.item)
			if err != nil {
				if tt.wantErr {
					if err.Error() != tt.err {
						t.Fatalf("Expected error: %s, got: %v\n", tt.err, err)
					}
					return
				}
				t.Fatalf("Unexpected error: %v\n", err)
			}
			if user != tt.wantUser {
				t.Fatalf("Wanted user: %s, got %s\n", tt.wantUser, user)
			}
			if pass != tt.wantPass {
				t.Fatalf("Wanted password: %s, got %s\n", tt.wantPass, pass)
			}
		})
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	cmd, args := args[0], args[1:]
	switch cmd {
	case "op":
		switch args[0] {
		case "signin":
			fmt.Println(`export OP_SESSION_my_team="RANDO"`)
		case "get":
			switch args[1] {
			case "totp":
				fmt.Println("123456")
			case "item":
				if args[2] == "FOOBAR" {
					fmt.Println(item)
				} else {
					fmt.Println("item not found")
					os.Exit(1)
				}
			}
		}
	}
}
