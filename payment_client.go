package webcore

import "time"

type Subscription struct {
	ID          string
	PlanID      string
	ProviderID  string
	ProductID   string
	UserID      string
	Provider    string
	AutoRenewal bool
	CreatedAt   time.Time
	StartAt     time.Time
	ExpiresAt   time.Time
}

const (
	SubscriptionProviderStripe = "Stripe"
)

// PaymentClient interface implements all the required stuff to handle customers
type PaymentClient interface {
	NewCustomer(string) (string, error)
	SetPaymentCard(string, string) error
	CreateSubscription(string, string) error
	GetSubscription(string) ([]Subscription, error)
	CancelSubscription(string, string) error
}
