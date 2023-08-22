package sessionmanager

func RemoveSession(token string) error {

	log.Debug("Entering RemoveSession() in sessionmanager / remove.go ...")

	prat, err := decrypt_session_ID(token)
	if err != nil {
		log.WithError(err).Error("Failed to decrypt the Auth Token.")
		return err
	}


	is_removed, err := remove_user_session(prat)
	if !is_removed || err != nil {
		log.WithError(err).Error("Failed to remove session from the database")
		return err
	}

	log.Debug("Successfully removed session from the database.")

	return nil
}