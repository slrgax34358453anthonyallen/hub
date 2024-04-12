package main

import (
	"context"
	"sync"

	"github.com/getAlby/nostr-wallet-connect/events"
	"github.com/getAlby/nostr-wallet-connect/nip47"
	"github.com/nbd-wtf/go-nostr"
	"github.com/sirupsen/logrus"
)

func (svc *Service) HandleMultiPayKeysendEvent(ctx context.Context, nip47Request *Nip47Request, requestEvent *RequestEvent, app *App, publishResponse func(*Nip47Response, nostr.Tags)) {

	multiPayParams := &Nip47MultiPayKeysendParams{}
	resp := svc.decodeNip47Request(nip47Request, requestEvent, app, multiPayParams)
	if resp != nil {
		publishResponse(resp, nostr.Tags{})
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, keysendInfo := range multiPayParams.Keysends {
		wg.Add(1)
		go func(keysendInfo Nip47MultiPayKeysendElement) {
			defer wg.Done()

			keysendDTagValue := keysendInfo.Id
			if keysendDTagValue == "" {
				keysendDTagValue = keysendInfo.Pubkey
			}
			dTag := []string{"d", keysendDTagValue}

			resp := svc.checkPermission(nip47Request, requestEvent.NostrId, app, keysendInfo.Amount)
			if resp != nil {
				publishResponse(resp, nostr.Tags{dTag})
				return
			}

			payment := Payment{App: *app, RequestEvent: *requestEvent, Amount: uint(keysendInfo.Amount / 1000)}
			mu.Lock()
			insertPaymentResult := svc.db.Create(&payment)
			mu.Unlock()
			if insertPaymentResult.Error != nil {
				svc.Logger.WithFields(logrus.Fields{
					"requestEventNostrId": requestEvent.NostrId,
					"recipientPubkey":     keysendInfo.Pubkey,
					"keysendId":           keysendInfo.Id,
				}).Errorf("Failed to process event: %v", insertPaymentResult.Error)
				return
			}

			svc.Logger.WithFields(logrus.Fields{
				"requestEventNostrId": requestEvent.NostrId,
				"appId":               app.ID,
				"recipientPubkey":     keysendInfo.Pubkey,
			}).Info("Sending payment")

			preimage, err := svc.lnClient.SendKeysend(ctx, keysendInfo.Amount, keysendInfo.Pubkey, keysendInfo.Preimage, keysendInfo.TLVRecords)
			if err != nil {
				svc.Logger.WithFields(logrus.Fields{
					"requestEventNostrId": requestEvent.NostrId,
					"appId":               app.ID,
					"recipientPubkey":     keysendInfo.Pubkey,
				}).Infof("Failed to send payment: %v", err)
				svc.EventPublisher.Publish(&events.Event{
					Event: "nwc_payment_failed",
					Properties: map[string]interface{}{
						// "error":   fmt.Sprintf("%v", err),
						"keysend": true,
						"multi":   true,
						"amount":  keysendInfo.Amount / 1000,
					},
				})

				publishResponse(&Nip47Response{
					ResultType: nip47Request.Method,
					Error: &Nip47Error{
						Code:    nip47.ERROR_INTERNAL,
						Message: err.Error(),
					},
				}, nostr.Tags{dTag})
				return
			}
			payment.Preimage = &preimage
			mu.Lock()
			svc.db.Save(&payment)
			mu.Unlock()
			svc.EventPublisher.Publish(&events.Event{
				Event: "nwc_payment_succeeded",
				Properties: map[string]interface{}{
					"keysend": true,
					"multi":   true,
					"amount":  keysendInfo.Amount / 1000,
				},
			})
			publishResponse(&Nip47Response{
				ResultType: nip47Request.Method,
				Result: Nip47PayResponse{
					Preimage: preimage,
				},
			}, nostr.Tags{dTag})
		}(keysendInfo)
	}

	wg.Wait()
}
