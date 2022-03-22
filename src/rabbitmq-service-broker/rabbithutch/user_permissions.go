package rabbithutch

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func (r *rabbitHutch) AssignPermissionsTo(vhost, username string) error {
	permissions := rabbithole.Permissions{Configure: ".*", Write: ".*", Read: ".*"}
	resp, err := r.client.UpdatePermissionsIn(vhost, username, permissions)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return validateResponse(resp)
}
