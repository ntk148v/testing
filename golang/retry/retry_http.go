// DeleteThing attempts to delete a thing. It will try a maximum of three times.
func DeleteThing(id string) error {
	// Build the request
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("https://unreliable-api/things/%s", id),
		nil,
	)
	if err != nil {
		return fmt.Errorf("unable to make request: %s", err)
	}

	// Execute the request
	return retry(3, time.Second, func() error {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// This error will result in a retry
			return err
		}
		defer resp.Body.Close()

		s := resp.StatusCode
		switch {
		case s >= 500:
			// Retry
			return fmt.Errorf("server error: %v", s)
		case s >= 400:
			// Don't retry, it was client's fault
			return stop{fmt.Errorf("client error: %v", s)}
		default:
			// Happy
			return nil
		}
	})
}
