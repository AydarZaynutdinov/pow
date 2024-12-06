package repository

import (
	"errors"
	"math/rand"
	"time"
)

var EmptyQuoteListErr = errors.New("empty quote list")

type Quote struct {
	random    *rand.Rand
	quoteList []string
}

func NewQuote(quoteList []string) (*Quote, error) {
	if len(quoteList) == 0 {
		return nil, EmptyQuoteListErr
	}
	return &Quote{
		random:    rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
		quoteList: quoteList,
	}, nil
}

func (q *Quote) GetQuote() []byte {
	return []byte(q.quoteList[rand.Intn(len(q.quoteList))])
}
