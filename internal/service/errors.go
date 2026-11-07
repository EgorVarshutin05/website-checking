package service

import "errors"

var (
	ErrInvalidSite = errors.New("некорректные данные сайта")
	ErrSiteExists  = errors.New("сайт уже существует")
)
