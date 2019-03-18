package rabbithutch

import (
	"net/http"

	rabbithole "github.com/michaelklishin/rabbit-hole"
)

//go:generate counterfeiter -o ./fakes/api_client_fake.go $FILE APIClient

type APIClient interface {
	GetVhost(string) (*rabbithole.VhostInfo, error)
	PutVhost(string, rabbithole.VhostSettings) (*http.Response, error)
	UpdatePermissionsIn(vhost, username string, permissions rabbithole.Permissions) (res *http.Response, err error)
	PutPolicy(vhost, name string, policy rabbithole.Policy) (res *http.Response, err error)
	DeleteVhost(vhostname string) (res *http.Response, err error)
	DeleteUser(username string) (res *http.Response, err error)
	ListUsers() (users []rabbithole.UserInfo, err error)
	PutUser(string, rabbithole.UserSettings) (*http.Response, error)
	ProtocolPorts() (res map[string]rabbithole.Port, err error)
}

//go:generate counterfeiter -o ./fakes/rabbithutch_fake.go $FILE RabbitHutch

type RabbitHutch interface {
	EnsureVHostExists(string) error
	CreateUser(string, string, string) (string, error)
}

type rabbitHutch struct {
	client APIClient
}

func New(client APIClient) RabbitHutch {
	return &rabbitHutch{client}
}
