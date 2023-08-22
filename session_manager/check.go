package sessionmanager

import (
	"time"
)

func ValidateAuthToken(pat string) (bool, error) {
	log.Debug("Entering ValidateAuthToken...\tReceived encrypted Auth Token ")

	prat, err := decrypt_session_ID(pat)
	if err != nil {
		log.WithError(err).Error("Failed to decrypt the Auth Token.")
		return false, err
	}

	expirationTime, err := val_session_time(prat)
	if err != nil {
		log.WithError(err).Error("Failed to validate the session.")
		return false, err
	}

	if time.Now().After(expirationTime) {
		log.Debug("Session is expired.")
		return false, nil
	}
	log.Debug("Session is valid.")
	return true, nil
}
