package rabbithutch

import (
	"net/http"

	"github.com/pivotal-cf/brokerapi"
)

func (r *rabbitHutch) DeleteUserAndConnections(username string) error {
	defer func() {
		conns, _ := r.client.ListConnections()
		for _, conn := range conns {
			if conn.User == username {
				resp, err := r.client.CloseConnection(conn.Name)
				if err == nil {
					resp.Body.Close()
				}
			}
		}
	}()

	return r.DeleteUser(username)
}

func (r *rabbitHutch) DeleteUser(username string) error {
	resp, err := r.client.DeleteUser(username)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return brokerapi.ErrBindingDoesNotExist
	}

	return validateResponse(resp)
}
