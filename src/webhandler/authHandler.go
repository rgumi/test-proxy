package webhandler

import (
	"fmt"
	"log"
	"net/http"
)

func checkAuth(r *http.Request) (err error) {

	if auth := r.Header.Get("auth"); auth != "" {
		if auth != lockOwner {
			return fmt.Errorf("Lock has already been acquired by someone else")
		}
		return nil
	}
	return fmt.Errorf("Authtoken not set in header")
}

func checkAuthLock(r *http.Request) (err error) {
	// check if authtoken is set in header of request
	if auth := r.Header.Get("auth"); auth != "" {

		for key, value := range authMap {
			if auth == value {
				// auth token is valid

				// try to get lock
				err = manageLock(key, auth)
				return
			}
		}
		// auth token is invalid
		return fmt.Errorf("Authtoken not found")
	}

	// authtoken not set
	return fmt.Errorf("Authtoken not set in header")
}

func manageLock(name, requestor string) (err error) {
	if lockOwner == "" {
		log.Printf("setting lock owner to %s", name)
		lockOwner = requestor
		return nil
	}
	if lockOwner == requestor {
		log.Printf("releasing lock of owner %s", name)
		lockOwner = ""
		return nil
	}
	return fmt.Errorf("Lock has already been acquired by someone else")
}
