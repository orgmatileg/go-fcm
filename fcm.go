package fcm

import (
	"context"
	"errors"
	"fmt"

	firebase "firebase.google.com/go"
	message "firebase.google.com/go/messaging"
)

// Message contain all information message
type Message struct {
	Token string
	Title string
	Body  string
	Topic string
	Data  map[string]string

	// Apple / iOS
	Subtitle  string
	ActioniOS string

	// Android
	ActionAndroid string

	// Webpush
	ActionWebPush  string
	CustomData     map[string]interface{}
	BadgeIconImage string
	SoundWebPush   string
	TagsCategory   string
}

// Initiation Firebase Cloud Messaging
// Require private key json from firebase
// to get private key -> https://console.firebase.google.com/u/0/project/_/settings/serviceaccounts/adminsdk
// Dont forget set environtment variable GOOGLE_APPLICATION_CREDENTIALS to location private key
// example -> os.Getenv("/app/generated-private-key.json")
func initFirebaseCloudMessaging() (*firebase.App, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing new app: %s", err.Error())
	}
	return app, nil
}

// SendMessageOne send message to one token/device
// return
func SendMessageOne(msg *Message) (string, error) {
	if msg.Token == "" {
		return "", errors.New("token cannot be empty")
	}
	app, err := initFirebaseCloudMessaging()
	if err != nil {
		return "", err
	}
	c, err := app.Messaging(context.Background())
	if err != nil {
		return "", fmt.Errorf("error initializing messaging: %s", err.Error())
	}
	m := copyMyMessageToFcmMessage(msg)
	s, err := c.Send(context.Background(), m)
	if err != nil {
		return "", fmt.Errorf("error sending messaging: %s", err.Error())
	}
	return s, nil
}

// ValidateToken validate valid token if not return an error
func ValidateToken(msg *Message) (string, error) {
	if msg.Token == "" {
		return "", errors.New("token cannot be empty")
	}
	app, err := initFirebaseCloudMessaging()
	if err != nil {
		return "", err
	}
	c, err := app.Messaging(context.Background())
	if err != nil {
		return "", fmt.Errorf("error initializing messaging: %s", err.Error())
	}
	m := copyMyMessageToFcmMessage(msg)
	s, err := c.SendDryRun(context.Background(), m)
	if err != nil {
		return "", fmt.Errorf("error sending messaging: %s", err.Error())
	}
	return s, nil
}

func copyMyMessageToFcmMessage(msg *Message) *message.Message {
	return &message.Message{
		Topic: msg.Topic,
		Token: msg.Token,
		Data:  msg.Data,
		APNS: &message.APNSConfig{
			Payload: &message.APNSPayload{
				Aps: &message.Aps{
					Category:   msg.TagsCategory,
					CustomData: msg.CustomData,
					Alert: &message.ApsAlert{
						ActionLocKey: msg.ActioniOS,
						Title:        msg.Title,
						SubTitle:     msg.Subtitle,
						Body:         msg.Body,
					},
				},
				CustomData: msg.CustomData,
			},
		},
		Android: &message.AndroidConfig{
			Data: msg.Data,
			Notification: &message.AndroidNotification{
				Title:       msg.Title,
				Icon:        msg.BadgeIconImage,
				Tag:         msg.TagsCategory,
				Body:        msg.Body,
				ClickAction: msg.ActioniOS,
			},
			Priority: "high",
		},
		Webpush: &message.WebpushConfig{
			Data: msg.Data,
			Notification: &message.WebpushNotification{
				Title:      msg.Title,
				Body:       msg.Body,
				Data:       msg.CustomData,
				CustomData: msg.CustomData,
				Icon:       msg.BadgeIconImage,
				Badge:      msg.BadgeIconImage,
				Image:      msg.BadgeIconImage,
				Vibrate:    []int{500, 500, 200, 200, 200},
				Tag:        msg.TagsCategory,
			},
			FcmOptions: &message.WebpushFcmOptions{
				Link: msg.ActionWebPush,
			},
		},
		Notification: &message.Notification{
			Title: msg.Title,
			Body:  msg.Body,
		},
	}
}

// TODO: add register topic
// TODO: add unregister topic

/*	USEFUL SOURCE

	https://developers.google.com/web/fundamentals/push-notifications/display-a-notification
	https://developer.mozilla.org/en-US/docs/Web/API/notification/Notification
	https://godoc.org/firebase.google.com/go/messaging
	https://firebase.google.com/docs/reference/fcm/rest/v1/projects.messages
	https://firebase.google.com/docs/cloud-messaging/send-message
	https://firebase.google.com/docs/admin/setup/


*/
