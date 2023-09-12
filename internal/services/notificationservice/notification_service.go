package notificationservice

import (
	"context"
	"fmt"
	"heypay-cash-in-server/internal/services/firebaseservice"
	"heypay-cash-in-server/utils"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/messaging"
)

type Payload struct {
	Data  map[string]string
	Title string
	Body  string
}

var (
	SendPushNotification func(ctx context.Context, userID string, payload Payload) error
)

func init() {
	SendPushNotification = sendPushNotification
}

func getFcmTokens(ctx context.Context, userID string) ([]string, error) {
	fcmTokenSnapshots, err := firebaseservice.Db.Collection("fcmTokens").Where("userId", "==", userID).Documents(ctx).GetAll()
	if len(fcmTokenSnapshots) == 0 && err == nil {
		utils.Warn(ctx, "No fcm tokens found")
		return nil, nil
	}
	if err != nil {
		utils.Errorf(ctx, "Error getting fcm tokens - reason: %v", err)
		return nil, err
	}
	tokens := make([]string, len(fcmTokenSnapshots))
	for i, doc := range fcmTokenSnapshots {
		tokens[i] = doc.Ref.ID
	}
	return tokens, nil
}

func handleNotificationResponse(ctx context.Context, response *messaging.BatchResponse, tokens []string) {
	tokensToDelete := make([]*firestore.DocumentRef, len(tokens))
	index := 0
	for i, result := range response.Responses {
		err := result.Error
		if err != nil {
			utils.Infof(ctx, "Sending notification to %v failed - reason: %v", tokens[i], err)
			if err.Error() == "messaging/invalid-registration-token" || err.Error() == "messaging/registration-token-not-registered" {
				tokensToDelete[index] = firebaseservice.Db.Doc(fmt.Sprintf("fcmTokens/%s", tokens[i]))
				index++
			}
		} else {
			utils.Infof(ctx, "Notification sent to %v", tokens[i])
		}
	}
	for i := 0; i < index; i++ {
		_, err := tokensToDelete[i].Delete(ctx)
		if err != nil {
			utils.Errorf(ctx, "Could not delete fcmToken %v - reason: %v", tokensToDelete[i].ID, err)
		}
	}
}

// sendPushNotification sends a push notification to an user by userID with payload
func sendPushNotification(ctx context.Context, userID string, payload Payload) error {
	tokens, err := getFcmTokens(ctx, userID)
	if err != nil {
		utils.Errorf(ctx, "Error getting tokens - reason: %v", err)
		return err
	}
	if tokens == nil {
		utils.Warn(ctx, "No fcm tokens found, no notification sent")
		return nil
	}
	messages := make([]*messaging.Message, len(tokens))
	for i, token := range tokens {
		messages[i] = &messaging.Message{
			Data:  payload.Data,
			Token: token,
			Notification: &messaging.Notification{
				Title: payload.Title,
				Body:  payload.Body,
			},
		}
	}
	response, err := firebaseservice.Msg.SendAll(ctx, messages)
	handleNotificationResponse(ctx, response, tokens)
	return nil
}
