package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/getAlby/nostr-wallet-connect/events"
	"github.com/getAlby/nostr-wallet-connect/nip47"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"github.com/sirupsen/logrus"
)

type Relay interface {
	Publish(ctx context.Context, event nostr.Event) (nostr.Status, error)
}

type Nip47Notifier struct {
	svc   *Service
	relay Relay
}

func NewNip47Notifier(svc *Service, relay Relay) *Nip47Notifier {
	return &Nip47Notifier{
		svc:   svc,
		relay: relay,
	}
}

func (notifier *Nip47Notifier) ConsumeEvent(ctx context.Context, event *events.Event) error {
	if event.Event != "nwc_payment_received" {
		return nil
	}

	if notifier.svc.lnClient == nil {
		return nil
	}

	paymentReceivedEventProperties, ok := event.Properties.(*events.PaymentReceivedEventProperties)
	if !ok {
		notifier.svc.Logger.WithField("event", event).Error("Failed to cast event")
		return errors.New("failed to cast event")
	}

	transaction, err := notifier.svc.lnClient.LookupInvoice(ctx, paymentReceivedEventProperties.PaymentHash)
	if err != nil {
		notifier.svc.Logger.
			WithField("paymentHash", paymentReceivedEventProperties.PaymentHash).
			WithError(err).
			Error("Failed to lookup invoice by payment hash")
		return err
	}

	notifier.notifySubscribers(ctx, &Nip47Notification{
		Notification:     transaction,
		NotificationType: nip47.PAYMENT_RECEIVED_NOTIFICATION,
	}, nostr.Tags{})
	return nil
}

func (notifier *Nip47Notifier) notifySubscribers(ctx context.Context, notification *Nip47Notification, tags nostr.Tags) {
	apps := []App{}

	// TODO: join apps and permissions
	notifier.svc.db.Find(&apps)

	for _, app := range apps {
		hasPermission, _, _ := notifier.svc.hasPermission(&app, nip47.NOTIFICATIONS_PERMISSION, 0)
		if !hasPermission {
			continue
		}
		notifier.notifySubscriber(ctx, &app, notification, tags)
	}
}

func (notifier *Nip47Notifier) notifySubscriber(ctx context.Context, app *App, notification *Nip47Notification, tags nostr.Tags) {
	notifier.svc.Logger.WithFields(logrus.Fields{
		"notification": notification,
		"appId":        app.ID,
	}).Info("Notifying subscriber")

	ss, err := nip04.ComputeSharedSecret(app.NostrPubkey, notifier.svc.cfg.NostrSecretKey)
	if err != nil {
		notifier.svc.Logger.WithFields(logrus.Fields{
			"notification": notification,
			"appId":        app.ID,
		}).WithError(err).Error("Failed to compute shared secret")
		return
	}

	payloadBytes, err := json.Marshal(notification)
	if err != nil {
		notifier.svc.Logger.WithFields(logrus.Fields{
			"notification": notification,
			"appId":        app.ID,
		}).WithError(err).Error("Failed to stringify notification")
		return
	}
	msg, err := nip04.Encrypt(string(payloadBytes), ss)
	if err != nil {
		notifier.svc.Logger.WithFields(logrus.Fields{
			"notification": notification,
			"appId":        app.ID,
		}).WithError(err).Error("Failed to encrypt notification payload")
		return
	}

	allTags := nostr.Tags{[]string{"p", app.NostrPubkey}}
	allTags = append(allTags, tags...)

	event := &nostr.Event{
		PubKey:    notifier.svc.cfg.NostrPublicKey,
		CreatedAt: nostr.Now(),
		Kind:      nip47.NOTIFICATION_KIND,
		Tags:      allTags,
		Content:   msg,
	}
	err = event.Sign(notifier.svc.cfg.NostrSecretKey)
	if err != nil {
		notifier.svc.Logger.WithFields(logrus.Fields{
			"notification": notification,
			"appId":        app.ID,
		}).WithError(err).Error("Failed to sign event")
		return
	}

	status, err := notifier.relay.Publish(ctx, *event)
	if err != nil {
		notifier.svc.Logger.WithFields(logrus.Fields{
			"notification": notification,
			"appId":        app.ID,
			"status":       status,
		}).WithError(err).Error("Failed to publish notification")
		return
	}
	notifier.svc.Logger.WithFields(logrus.Fields{
		"notification": notification,
		"appId":        app.ID,
		"status":       status,
	}).Info("Published notification event")

}
