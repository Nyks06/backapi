package webcore

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// APIService is the struct responsible of handling the API methods and apply the logic code.
// It is connected directly to the third party services and the stores in order to r/w information required.
type APIService struct {
	UserStore  UserStore
	UserFinder UserFinder

	SessionStore  SessionStore
	SessionFinder SessionFinder

	TicketStore  TicketStore
	TicketFinder TicketFinder

	PronosticsStore  PronosticsStore
	PronosticsFinder PronosticsFinder

	SportsStore  SportsStore
	SportsFinder SportsFinder

	CompetitionsStore  CompetitionsStore
	CompetitionsFinder CompetitionsFinder

	Mailer        Mailer
	PaymentClient PaymentClient
}

// ---------------------
// User related stuff
// _____________________

// CreateUser takes a user and stores it on our db.
// The user parameter is considered to be well formatted and filled.
func (a *APIService) CreateUser(ctx context.Context, user *User) (*User, error) {
	u, err := a.UserFinder.ByEmail(user.Email)
	// TODO : Handle error WTF ?

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, NewInternalServerError("the password can't be encrypted")
	}
	user.Password = string(password)
	user.SponsorshipID = fmt.Sprintf("%v%v%v", strings.ToUpper(user.Username), rand.Intn(99), rand.Intn(99))

	// Check if there is an internal server error or not
	if err != nil {
		if IsInternalServerError(err) {
			return nil, err
		}
	}
	// Check if a user has been found or not
	if u != nil {
		return nil, NewResourceAlreadyCreatedError("user_already_created")
	}

	customerID, err := a.PaymentClient.NewCustomer(user.Email)
	if err != nil {
		return nil, err
	}
	user.CustomerID = customerID

	createdUser, err := a.UserStore.Create(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// RemoveUser not implemented yet
func (a *APIService) RemoveUser(ctx context.Context, ID string) error {
	return nil
}

func (a *APIService) GetCurrentUser(ctx context.Context, CUser *User, UserID string) (*User, error) {
	if CUser == nil || CUser.ID != UserID {
		return nil, NewUnauthorizedError("user_mismatch", "user_mismatch")
	}

	return a.UserFinder.ByID(UserID)
}

func (a *APIService) GetUserByID(ctx context.Context, UserID string) (*User, error) {
	u, err := a.UserFinder.ByID(UserID)
	if err != nil {
		return nil, err
	}

	subs, err := a.SubscriptionGet(ctx, u.CustomerID)
	if err != nil {
		return nil, err
	}
	// for _, sub := range subs {
	// 	switch sub.ProductID {
	// 	// Fun
	// 	case "prod_GOm07o55qZvlln":
	// 		u.SubIDFun = sub.ID
	// 	// Prive
	// 	case "prod_GOm1ffYFrUjvBc":
	// 		u.SubIDPrive = sub.ID
	// 	// Champion
	// 	case "prod_GOm1yxzdPh1lI7":
	// 		u.SubIDChampion = sub.ID
	// 	}
	// }

	return u, nil
}

func (a *APIService) ListUsers(ctx context.Context) ([]User, error) {
	return a.UserFinder.List()
}

func (a *APIService) GetUserByEmail(ctx context.Context, Email string) (*User, error) {
	return a.UserFinder.ByEmail(Email)
}

func (a *APIService) ChangePassword(ctx context.Context, user *User, currentPassword, newPassword string) (*User, error) {
	u, err := a.UserFinder.ByEmail(user.Email)
	if err != nil {
		return nil, NewInternalServerError("the user can't be found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(currentPassword))
	if err != nil {
		return nil, NewUnauthorizedError("password_mismatch", "password_mismatch")
	}

	newPasswd, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, NewInternalServerError("the password can't be encrypted")
	}
	u.Password = string(newPasswd)

	err = a.UserStore.SetPassword(u.ID, u.Password)
	if err != nil {
		return nil, NewInternalServerError("the password can't be updated")
	}

	return u, nil
}

// func (a *APIService) UpdateDetails(ctx context.Context, user *User, Firstname string, Lastname string) (*User, error) {
// 	u, err := a.UserFinder.ByEmail(user.Email)
// 	if err != nil {
// 		return nil, NewInternalServerError("the user can't be found")
// 	}

// 	u.Firstname = Firstname
// 	u.Lastname = Lastname
// 	err = a.UserStore.UpdateDetails(u.ID, u.Firstname, u.Lastname)
// 	if err != nil {
// 		return nil, NewInternalServerError("the user details can't be updated")
// 	}

// 	return u, nil
// }

// func (a *APIService) UpdateContactSettings(ctx context.Context, user *User, email string, phone string, telegram string) (*User, error) {
// 	u, err := a.UserFinder.ByEmail(user.Email)
// 	if err != nil {
// 		return nil, NewInternalServerError("the user can't be found")
// 	}

// 	u.Email = email
// 	u.PhoneNumber = phone
// 	u.Telegram = telegram
// 	err = a.UserStore.UpdateContactSettings(u.ID, u.Email, u.PhoneNumber, u.Telegram)
// 	if err != nil {
// 		return nil, NewInternalServerError("the user contact settings can't be updated")
// 	}

// 	return u, nil
// }

// ---------------------
// Session related stuff
// _____________________

func (a *APIService) GetSessionByID(ctx context.Context, ID string) (*Session, error) {
	return a.SessionFinder.ByID(ID)
}

func (a *APIService) CreateSession(ctx context.Context, Email, Password string) (*Session, *User, error) {
	user, err := a.UserFinder.ByEmail(Email)
	if err != nil {
		return nil, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password))
	if err != nil {
		return nil, nil, NewUnauthorizedError("password_mismatch")
	}

	subs, err := a.SubscriptionGet(ctx, user.CustomerID)
	if err != nil {
		return nil, nil, err
	}
	// for _, sub := range subs {
	// 	switch sub.ProductID {
	// 	// Fun
	// 	case "prod_FbmLHcI2x4KvY3":
	// 		user.SubIDFun = sub.ID
	// 	// Prive
	// 	case "prod_FbmOtiMhXCPx8b":
	// 		user.SubIDPrive = sub.ID
	// 	// Champion
	// 	case "prod_FbmNbyjYyKBLI4":
	// 		user.SubIDChampion = sub.ID
	// 	}
	// }

	sess, err := a.SessionStore.Create(&Session{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(48 * time.Hour),
	})
	return sess, user, err
}

func (a *APIService) RemoveSession(ctx context.Context, ID string) error {
	return a.SessionStore.Delete(ID)
}

// ---------------------
// Subscription related stuff
// _____________________

func (a *APIService) SubscriptionCreate(ctx context.Context, CUser *User, PlanID string, SourceID string) error {
	if CUser.CustomerID == "" {
		customerID, err := a.PaymentClient.NewCustomer(CUser.Email)
		if err != nil {
			return err
		}
		CUser.CustomerID = customerID
		if err = a.UserStore.UpdateCustomerID(CUser.ID, CUser.CustomerID); err != nil {
			return err
		}
	}

	err := a.PaymentClient.SetPaymentCard(CUser.CustomerID, SourceID)
	if err != nil {
		return err
	}

	return a.PaymentClient.CreateSubscription(CUser.CustomerID, PlanID)
}

func (a *APIService) SubscriptionGet(ctx context.Context, CustomerID string) ([]Subscription, error) {
	return a.PaymentClient.GetSubscription(CustomerID)
}

// func (a *APIService) SubscriptionCancel(ctx context.Context, CustomerID string, id string) error {
// 	return a.PaymentClient.CancelSubscription(CustomerID, id)
// }

// ---------------------
// Ticket related stuff
// _____________________

func (a *APIService) TicketCreate(ctx context.Context, pGroup *Ticket) (*Ticket, error) {
	return a.TicketStore.Create(pGroup)
}

func (a *APIService) TicketUpdate(ctx context.Context, ID string, ticket *Ticket) error {
	return a.TicketStore.Update(ID, ticket)
}

func (a *APIService) TicketGetByID(ctx context.Context, ID string) (*Ticket, error) {
	ticket, err := a.TicketFinder.ByID(ID)
	if err != nil {
		return nil, err
	}

	pronos, err := a.PronosticsFinder.ListByTicketID(ID)
	if err != nil {
		return nil, err
	}

	for idx, prono := range pronos {
		sport, err := a.SportsFinder.ByID(prono.Sport.ID)
		if err != nil {
			return nil, err
		}

		compet, err := a.CompetitionsFinder.ByID(prono.Competition.ID)
		if err != nil {
			return nil, err
		}

		pronos[idx].Sport = sport
		pronos[idx].Competition = compet
	}

	ticket.Pronostics = pronos

	// Update ticket overall odd and status
	ticket.Odd = 1
	ticket.Status = StatusWin
	for _, p := range ticket.Pronostics {
		if p.Status != StatusCanceled {
			ticket.Odd *= p.Odd
		}
		if p.Status == StatusLose {
			ticket.Status = StatusLose
		} else if p.Status == StatusInProgress && ticket.Status == StatusWin {
			ticket.Status = StatusInProgress
		}
	}

	ticket.Odd = math.Round(ticket.Odd*100) / 100

	return ticket, nil
}

func (a *APIService) TicketDelete(ctx context.Context, ID string) error {
	return a.TicketStore.Delete(ID)
}
func (a *APIService) TicketList(ctx context.Context) ([]Ticket, error) {
	tickets, err := a.TicketFinder.List()
	if err != nil {
		return nil, err
	}

	for id := range tickets {
		pronos, err := a.PronosticsFinder.ListByTicketID(tickets[id].ID)
		if err != nil {
			return nil, err
		}

		for idx, prono := range pronos {
			sport, err := a.SportsFinder.ByID(prono.Sport.ID)
			if err != nil {
				return nil, err
			}

			compet, err := a.CompetitionsFinder.ByID(prono.Competition.ID)
			if err != nil {
				return nil, err
			}

			pronos[idx].Sport = sport
			pronos[idx].Competition = compet
		}

		tickets[id].Pronostics = pronos

		// Update ticket overall odd and status
		tickets[id].Odd = 1
		tickets[id].Status = StatusWin
		for _, p := range tickets[id].Pronostics {
			if p.Status != StatusCanceled {
				tickets[id].Odd *= p.Odd
			}
			if p.Status == StatusLose {
				tickets[id].Status = StatusLose
			} else if p.Status == StatusInProgress && tickets[id].Status == StatusWin {
				tickets[id].Status = StatusInProgress
			}
		}
		// Done in order to round the overall odd
		tickets[id].Odd = math.Round(tickets[id].Odd*100) / 100
	}

	return tickets, nil
}

// ---------------------
// Pronostic related stuff
// _____________________

func (a *APIService) PronosticCreate(ctx context.Context, prono *Pronostic) (*Pronostic, error) {
	return a.PronosticsStore.Create(prono)
}

func (a *APIService) PronosticGetByID(ctx context.Context, ID string) (*Pronostic, error) {
	prono, err := a.PronosticsFinder.ByID(ID)
	if err != nil {
		return nil, err
	}

	sport, err := a.SportsFinder.ByID(prono.Sport.ID)
	if err != nil {
		return nil, err
	}
	compet, err := a.CompetitionsFinder.ByID(prono.Competition.ID)
	if err != nil {
		return nil, err
	}

	prono.Sport = sport
	prono.Competition = compet

	return prono, nil
}

func (a *APIService) PronosticDelete(ctx context.Context, ID string) error {
	return a.PronosticsStore.Delete(ID)
}

func (a *APIService) PronosticUpdate(ctx context.Context, ID string, Status string) error {
	return a.PronosticsStore.Update(ID, Status)
}

func (a *APIService) PronosticList(ctx context.Context) ([]Pronostic, error) {
	pronos, err := a.PronosticsFinder.List()
	if err != nil {
		return nil, err
	}

	for idx, prono := range pronos {
		sport, err := a.SportsFinder.ByID(prono.Sport.ID)
		if err != nil {
			return nil, err
		}

		compet, err := a.CompetitionsFinder.ByID(prono.Competition.ID)
		if err != nil {
			return nil, err
		}

		pronos[idx].Sport = sport
		pronos[idx].Competition = compet
	}

	return pronos, nil
}

// ---------------------
// Sport related stuff
// _____________________

func (a *APIService) SportCreate(ctx context.Context, sport *Sport) (*Sport, error) {
	return a.SportsStore.Create(sport)
}

func (a *APIService) SportGetByID(ctx context.Context, ID string) (*Sport, error) {
	return a.SportsFinder.ByID(ID)
}

func (a *APIService) SportDelete(ctx context.Context, ID string) error {
	return a.SportsStore.Delete(ID)
}

func (a *APIService) SportList(ctx context.Context) ([]Sport, error) {
	return a.SportsFinder.List()
}

// ---------------------
// Competition related stuff
// _____________________

func (a *APIService) CompetitionCreate(ctx context.Context, competition *Competition) (*Competition, error) {
	return a.CompetitionsStore.Create(competition)
}

func (a *APIService) CompetitionGetByID(ctx context.Context, ID string) (*Competition, error) {
	compet, err := a.CompetitionsFinder.ByID(ID)
	if err != nil {
		return nil, err
	}

	sport, err := a.SportsFinder.ByID(compet.Sport.ID)
	if err != nil {
		return nil, err
	}

	compet.Sport = sport
	return compet, nil
}

func (a *APIService) CompetitionDelete(ctx context.Context, ID string) error {
	return a.CompetitionsStore.Delete(ID)
}

func (a *APIService) CompetitionList(ctx context.Context) ([]Competition, error) {
	compets, err := a.CompetitionsFinder.List()
	if err != nil {
		return nil, err
	}

	for idx, compet := range compets {
		sport, err := a.SportsFinder.ByID(compet.Sport.ID)
		if err != nil {
			return nil, err
		}

		compets[idx].Sport = sport
	}

	return compets, nil
}

// ---------------------
// Message related stuff
// _____________________

func (a *APIService) MessageCreate(ctx context.Context, Name string, Phone string, Email string, Message string) error {
	return nil
}
