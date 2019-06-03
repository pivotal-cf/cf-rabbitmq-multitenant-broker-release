package rabbithutch

func (b *rabbitHutch) VHostDelete(vhost string) error {
	resp, err := b.client.DeleteVhost(vhost)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return validateResponse(resp)
}
