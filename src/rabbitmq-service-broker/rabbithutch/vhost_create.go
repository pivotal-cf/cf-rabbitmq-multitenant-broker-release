package rabbithutch

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func (r *rabbitHutch) VHostCreate(vhost string) error {
	resp, err := r.client.PutVhost(vhost, rabbithole.VhostSettings{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return validateResponse(resp)
}
