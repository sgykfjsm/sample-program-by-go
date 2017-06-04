package main

// This code is to learn https://upgear.io/blog/simple-golang-retry-function

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

type stop struct {
	error
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetWithRetry(requestURL string, attempts int) error {
	// Build the request
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("unable to make request: %s", err.Error())
	}

	fmt.Printf("Start to request %s\n", requestURL)
	return retry(attempts, time.Second, func() error {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			// This error will result in a retry
			return err
		}
		defer res.Body.Close()

		s := res.StatusCode
		switch s {
		case 500:
			// Retry
			return fmt.Errorf("server error: %d", s)
		case 401:
			return stop{fmt.Errorf("client error: %d", s)}
		default:
			fmt.Printf("Success: %d", s)
			return nil
		}
	})
}

func retry(attemts int, sleepDuration time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attemts--; attemts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleepDuration)))
			// exponential backoff
			sleepDuration = sleepDuration + jitter/2

			fmt.Printf("Sleeping for %f seconds to retry (remaining attempts %d)\n", sleepDuration.Seconds(), attemts)
			time.Sleep(sleepDuration)
			return retry(attemts, 2*sleepDuration, f)
		}
	}

	return nil
}

func main() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := rand.Int() % 100
		if i == 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200 OK"))
		} else if i < 10 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 UNAUTHORIZED"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 INTERNAL SERVER ERROR"))
		}
	}))
	defer ts.Close()

	err := GetWithRetry(ts.URL, 10)
	if err != nil {
		fmt.Print(err.Error())
	}
}
