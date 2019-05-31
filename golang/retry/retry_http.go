package httpretry

func init() {
	rand.Seed(time.Now().UnixNano())
}

func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

type stop struct {
	error
}

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
