package wirepusher

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	AES_KEY_SIZE       int = 16 // WirePusher uses AES-128 for encryption
	ClearSingleMessage     = "wirepusher_clear_notification"
	ClearAllMessages       = "wirepusher_clear_all_notifications"
	baseURL                = "https://wirepusher.com/send"
)

// WirePusher service definition
type Service struct {
	client *http.Client
}

// New returns a new WirePusher service with the http DefaultClient
func New() *Service {
	return &Service{client: http.DefaultClient}
}

// NewWithClient returns a new WirePusher service with a custom http.Client
func NewWithClient(client *http.Client) *Service {
	return &Service{client: client}
}

// Send the message. If the provided user has a password, the message title and body will be encrypted
func (s *Service) Send(msg *Message, user *User) error {
	var (
		iv    []byte
		title string
		body  string
		err   error
	)
	if user.CanEncrypt() {
		iv, err = randomIV()
		if err != nil {
			return errors.Wrap(err, "error retrieving random iv during Send")
		}
		key, err := hex2bin(user.Password())
		if err != nil {
			return errors.Wrap(err, "error reading password of user during Send")
		}

		title, err = encrypt(msg.Title, iv, key)
		if err != nil {
			return errors.Wrap(err, "error encrypting title during Send")
		}

		body, err = encrypt(msg.Body, iv, key)
		if err != nil {
			return errors.Wrap(err, "error encrypting body during Send")
		}
	} else {
		title = msg.Title
		body = msg.Body
	}

	link := fmt.Sprintf("%s?id=%s&title=%s&message=%s&type=%s", baseURL, user.ID, url.QueryEscape(title), url.QueryEscape(body), msg.Type)
	if user.CanEncrypt() {
		link += fmt.Sprintf("&iv=%s", url.QueryEscape(bin2hex(iv)))
	}

	if msg.ID != 0 {
		link += fmt.Sprintf("&message_id=%d", msg.ID)
	}

	if len(msg.ActionURL) > 0 {
		link += fmt.Sprintf("&action=%s", msg.ActionURL)
	}

	if len(msg.ImageURL) > 0 {
		link += fmt.Sprintf("&image_url=%s", msg.ImageURL)
	}

	resp, err := s.client.Get(link)
	if err != nil {
		return errors.Wrapf(err, "error sending message with url %q", link)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http response code for url %q: %s", link, resp.Status)
	}

	return nil
}

// ClearMsg clears the message with the specified id from the notification tray of the specified user
func (s *Service) ClearMsg(msgid int, user *User) error {
	url := fmt.Sprintf("%s?id=%s&type=%s&message_id=%d", baseURL, user.ID, ClearSingleMessage, msgid)

	resp, err := s.client.Get(url)
	if err != nil {
		return errors.Wrapf(err, "error sending clear msg message with url %q", url)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected http response code for url %q: %s", url, resp.Status)
	}
	return nil
}

// ClearAllMsg clears all messages from the notification tray of the specified user
func (s *Service) ClearAllMsg(user *User) error {
	url := fmt.Sprintf("%s?id=%s&type=%s", baseURL, user.ID, ClearAllMessages)

	resp, err := s.client.Get(url)
	if err != nil {
		return errors.Wrapf(err, "error sending clear all msg message with url %q", url)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected http response code for url %q: %s", url, resp.Status)
	}
	return nil
}
