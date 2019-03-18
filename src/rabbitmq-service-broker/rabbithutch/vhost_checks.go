package rabbithutch

import (
	"net/http"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/pivotal-cf/brokerapi"
)

func (r *rabbitHutch) EnsureVHostExists(vhost string) error {
	if _, err := r.client.GetVhost(vhost); err != nil {
		if rabbitErr, ok := err.(rabbithole.ErrorResponse); ok && rabbitErr.StatusCode == http.StatusNotFound {
			return brokerapi.ErrInstanceDoesNotExist
		}

		return err
	}

	return nil
}
