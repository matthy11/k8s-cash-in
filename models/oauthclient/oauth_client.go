package oauthclient

type OAuthClient struct {
	FailedAttempts int `firestore:"failedAttempts"`
}

type PrivateKey struct {
	Hash   string                            `firestore:"hash"`
	Status string                            `firestore:"status"`
	Claims map[string]map[string]interface{} `firestore:"claims"`
}
