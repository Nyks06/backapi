package stripe

import (
	"time"

	webcore "github.com/nyks06/backapi"

	striperoot "github.com/stripe/stripe-go"
	stripecli "github.com/stripe/stripe-go/client"
)

type PaymentClient struct {
	Client *stripecli.API
}

func NewPaymentClient(cl *stripecli.API) *PaymentClient {
	return &PaymentClient{
		Client: cl,
	}
}

func (p *PaymentClient) NewCustomer(email string) (string, error) {
	customer, err := p.Client.Customers.New(&striperoot.CustomerParams{
		Email: striperoot.String(email),
	})
	if err != nil {
		return "", webcore.NewInternalServerError(err.Error())
	}

	return customer.ID, nil
}

func (p *PaymentClient) SetPaymentCard(customerID string, sourceID string) error {
	params := &striperoot.CustomerParams{}
	params.SetSource(sourceID)
	_, err := p.Client.Customers.Update(customerID, params)

	if err != nil {
		return webcore.NewInternalServerError(err.Error())
	}
	return nil
}

func (p *PaymentClient) CreateSubscription(customerID string, planID string) error {
	coupon := ""
	if planID == "plan_GROWklP0N5Kkbx" {
		coupon = "g4mnMZjk"
	}

	if coupon != "" {
		p.Client.Subscriptions.New(&striperoot.SubscriptionParams{
			Customer: striperoot.String(customerID),
			Coupon:   striperoot.String(coupon),
			Items: []*striperoot.SubscriptionItemsParams{
				{
					Plan: striperoot.String(planID),
				},
			},
		})
	} else {
		p.Client.Subscriptions.New(&striperoot.SubscriptionParams{
			Customer: striperoot.String(customerID),
			Items: []*striperoot.SubscriptionItemsParams{
				{
					Plan: striperoot.String(planID),
				},
			},
		})
	}

	return nil
}

func (p *PaymentClient) GetSubscription(customerID string) ([]webcore.Subscription, error) {
	subs := make([]webcore.Subscription, 0)

	i := p.Client.Subscriptions.List(&striperoot.SubscriptionListParams{
		Customer: customerID,
		Status:   "all",
	})

	for i.Next() {
		sub := webcore.Subscription{
			ID:         i.Subscription().ID,
			PlanID:     i.Subscription().Plan.ID,
			ProviderID: webcore.SubscriptionProviderStripe,
			UserID:     i.Subscription().Customer.ID,
			ProductID:  i.Subscription().Plan.Product.ID,
			Provider:   webcore.SubscriptionProviderStripe,
			CreatedAt:  time.Unix(i.Subscription().Created, 0),
			ExpiresAt:  time.Unix(i.Subscription().CurrentPeriodEnd, 0),
		}
		if i.Subscription().Status == striperoot.SubscriptionStatusActive {
			sub.AutoRenewal = true
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

func (p *PaymentClient) CancelSubscription(customerID string, ID string) error {
	subs, err := p.GetSubscription(customerID)
	if err != nil {
		return err
	}
	for _, s := range subs {
		if s.AutoRenewal == true {
			if err = p.cancelSubscription(s.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *PaymentClient) cancelSubscription(subID string) error {
	_, err := p.Client.Subscriptions.Cancel(subID, nil)
	if err != nil {
		return webcore.NewInternalServerError(err.Error())
	}

	return nil
}
