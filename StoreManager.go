package webcore

type StoreManager struct {
	UserStore  UserStore
	UserFinder UserFinder

	SessionStore  SessionStore
	SessionFinder SessionFinder

	TicketStore  TicketStore
	TicketFinder TicketFinder

	PronosticStore  PronosticsStore
	PronosticFinder PronosticsFinder

	SportStore  SportsStore
	SportFinder SportsFinder

	CompetitionStore  CompetitionsStore
	CompetitionFinder CompetitionsFinder
}
