package backapi

// StoreManager handle all the different stores accessible
type StoreManager struct {
	UserStore  UserStore
	UserFinder UserFinder

	SessionStore  SessionStore
	SessionFinder SessionFinder
}
