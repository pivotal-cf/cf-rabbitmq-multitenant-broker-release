package rabbithutch

func (r *rabbitHutch) VHostDelete(vhost string) error {
	resp, err := r.client.DeleteVhost(vhost)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return validateResponse(resp)
}
