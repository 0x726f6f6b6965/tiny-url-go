package service

import (
	"net/url"
)

type TinyURLService interface {
	// ShortURL - create new short URLs
	//
	// @ apiDevKey - A registered user account’s unique identifier.
	//
	// @ originalURL - The original long URL that is needed to be shortened.
	//
	// @ customAlias - The optional key that the user defines as a customer short URL.
	//
	// @ expiryDate - The optional expiration date for the shortened URL.
	ShortURL(apiDevKey, originalURL, customAlias string, expiryDate int64) (*url.URL, error)
	// RedirectURL - redirect a short URL
	//
	// @ apiDevKey - A registered user account’s unique identifier.
	//
	// @ urlKey - The shortened URL against which we need to fetch the long URL from the database.
	RedirectURL(apiDevKey, urlKey string) (*url.URL, error)
	// DeleteURL - delete a short URL
	//
	// @ apiDevKey - A registered user account’s unique identifier.
	//
	// @ urlKey - The shortened URL against which we need to fetch the long URL from the database.
	DeleteURL(apiDevKey, urlKey string) error
}
