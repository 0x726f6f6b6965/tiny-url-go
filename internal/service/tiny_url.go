package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/0x726f6f6b6965/tiny-url-go/internal/config"
	"github.com/0x726f6f6b6965/tiny-url-go/protos"
	"github.com/0x726f6f6b6965/tiny-url-go/utils"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ShortedURLService interface {
	// ShortURL - create new short URLs
	//
	// @ owner - A registered user account’s unique identifier.
	//
	// @ originalURL - The original long URL that is needed to be shortened.
	//
	// @ expiryDate - The optional expiration for the shortened URL.
	ShortURL(ctx context.Context, owner, originalURL string, expiryDate time.Time) (string, error)
	// RedirectURL - redirect a short URL
	//
	// @ urlKey - The shortened URL against which we need to fetch the long URL from the database.
	RedirectURL(ctx context.Context, urlKey string) (string, error)
	// DeleteURL - delete a short URL
	//
	// @ owner - A registered user account’s unique identifier.
	//
	// @ urlKey - The shortened URL against which we need to fetch the long URL from the database.
	DeleteURL(ctx context.Context, owner, urlKey string) error
	// UpdateURL - update a short URL
	//
	// @ owner - A registered user account’s unique identifier.
	//
	// @ short - The shortened URL against which we need to fetch the long URL from the database.
	//
	// @ originalURL - The original long URL that is needed to be shortened.
	//
	// @ expiry - The optional expiration date for the shortened URL.
	UpdateURL(ctx context.Context, owner, short, originalURL string, expiry time.Time) error
}

type shortenURLService struct {
	sequencer utils.Sequencer
	dynamodb  utils.Storage
	urlTable  string
	expire    time.Duration
}

var (
	ErrSequencer = errors.New("sequencer error")
	ErrStorage   = errors.New("storage error")
	ErrUnmarshal = errors.New("unmarshal error")
	ErrEmpty     = errors.New("empty")
)

// DeleteURL implements TinyURLService.
func (t *shortenURLService) DeleteURL(ctx context.Context, owner string, urlKey string) error {
	if utils.IsEmpty(owner) {
		return errors.Join(ErrEmpty, errors.New("owner is empyt"))
	}
	if utils.IsEmpty(urlKey) {
		return errors.Join(ErrEmpty, errors.New("urlKey is empyt"))
	}
	err := t.dynamodb.Delete(ctx, fmt.Sprintf("%s;URL#%s;USER#%s", t.urlTable, urlKey, owner))
	if err != nil {
		return errors.Join(ErrStorage, err)
	}
	return nil
}

// RedirectURL implements TinyURLService.
func (t *shortenURLService) RedirectURL(ctx context.Context, urlKey string) (string, error) {
	if utils.IsEmpty(urlKey) {
		return "", errors.Join(ErrEmpty, errors.New("urlKey is empyt"))
	}
	data, err := t.dynamodb.Get(ctx, fmt.Sprintf("%s;URL#%s;BeginWith USER#;query", t.urlTable, urlKey))
	if err != nil {
		return "", errors.Join(ErrStorage, err)
	}
	paginator, ok := data.(*dynamodb.QueryPaginator)
	if !ok {
		return "", ErrUnmarshal
	}
	var pages []protos.ShortenedURL
	if paginator.HasMorePages() {
		response, err := paginator.NextPage(ctx)
		if err != nil {
			return "", errors.Join(ErrStorage, err)
		}
		err = attributevalue.UnmarshalListOfMaps(response.Items, &pages)
		if err != nil {
			return "", errors.Join(ErrStorage, err)
		}
	}

	return pages[0].Original, nil
}

// ShortURL implements TinyURLService.
func (t *shortenURLService) ShortURL(ctx context.Context, owner string, originalURL string, expiryDate time.Time) (string, error) {
	if utils.IsEmpty(owner) {
		return "", errors.Join(ErrEmpty, errors.New("owner is empyt"))
	}
	if utils.IsEmpty(originalURL) {
		return "", errors.Join(ErrEmpty, errors.New("originalURL is empyt"))
	}
	seq, err := t.sequencer.Next()
	if err != nil {
		return "", errors.Join(ErrSequencer, err)
	}

	encoded := base64.RawURLEncoding.EncodeToString(seq.Bytes())
	now := time.Now().UTC()
	data := &protos.ShortenedURL{
		Shorten:   encoded,
		Original:  originalURL,
		ExpiresAt: now.Add(t.expire).UTC().Unix(),
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
		Owner:     owner,
	}
	if expiryDate.After(now) {
		data.ExpiresAt = expiryDate.UTC().Unix()
	}
	err = t.dynamodb.Save(ctx, fmt.Sprintf("%s;URL#%s;USER#%s", t.urlTable, encoded, owner), data)
	if err != nil {
		return "", errors.Join(ErrStorage, err)
	}
	return encoded, nil
}

// UpdateURL implements TinyURLService.
func (t *shortenURLService) UpdateURL(ctx context.Context, owner string, short string, originalURL string, expiry time.Time) error {
	if utils.IsEmpty(owner) {
		return errors.Join(ErrEmpty, errors.New("owner is empyt"))
	}
	if utils.IsEmpty(short) {
		return errors.Join(ErrEmpty, errors.New("short is empyt"))
	}
	data := &protos.ShortenedURL{}
	mask := []string{}
	if !utils.IsEmpty(originalURL) {
		data.Original = originalURL
		mask = append(mask, "Original")
	}
	if !expiry.IsZero() {
		data.ExpiresAt = expiry.UTC().Unix()
		mask = append(mask, "ExpiresAt")
	}
	if len(mask) == 0 {
		return nil
	}
	mask = append(mask, "UpdatedAt")
	data.UpdatedAt = time.Now().UTC().Unix()
	err := t.dynamodb.Update(ctx, fmt.Sprintf("%s;URL#%s;USER#%s", t.urlTable, short, owner), data, mask)
	if err != nil {
		return errors.Join(ErrStorage, err)
	}
	return nil
}

func NewTinyURLService(cfg *config.AppConfig, sequencer utils.Sequencer, dynamodb utils.Storage) ShortedURLService {
	return &shortenURLService{
		urlTable:  cfg.TableName,
		sequencer: sequencer,
		dynamodb:  dynamodb,
		expire:    cfg.Expire,
	}
}
