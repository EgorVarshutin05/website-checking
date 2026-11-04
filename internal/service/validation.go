package service

import (
	"errors"
	"net/url"
)

func ValidateSiteURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("разрешены только HTTP и HTTPS")
	}

	if parsed.Hostname() == "" {
		return errors.New("имя хоста не может быть пустым")
	}

	if parsed.User != nil {
		return errors.New("логин и пароль в URL запрещены")
	}

	if parsed.Fragment != "" || parsed.RawFragment != "" {
		return errors.New("фрагменты URL запрещены")
	}

	port := parsed.Port()

	if port != "" && port != "80" && port != "443" {
		return errors.New("разрешены только порты 80 и 443")
	}

	return nil
}
