package user

type User struct {
	ID                     string // added on runtime
	PrimaryAccountID       string `firestore:"primaryAccountId"`
	PrimaryAccountCategory string `firestore:"primaryAccountCategory"`
}
