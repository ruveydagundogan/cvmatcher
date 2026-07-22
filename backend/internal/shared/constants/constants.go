package constants

const AnonymousUserID = "00000000-0000-0000-0000-000000000000"

func NormalizeUserID(userID string) string {
	if userID == "" {
		return AnonymousUserID
	}
	return userID
}
