package readeef

import (
	"database/sql"
	"errors"
	"net/url"

	"github.com/urandom/readeef/parser"
)

type Feed struct {
	parser.Feed

	User           User
	Id             int64
	Link           string
	SiteLink       string `db:"site_link"`
	HubLink        string `db:"hub_link"`
	UpdateError    string `db:"update_error"`
	SubscribeError string `db:"subscribe_error"`
	Articles       []Article
	Tags           []string

	lastUpdatedArticleLinks map[string]bool
}

type Article struct {
	parser.Article

	Id       int64
	Guid     sql.NullString
	FeedId   int64 `db:"feed_id"`
	Read     bool
	Favorite bool
	Score    int64
}

type ArticleScores struct {
	ArticleId int64
	Score     int64
	Score1    int64
	Score2    int64
	Score3    int64
	Score4    int64
	Score5    int64
}

var (
	EmptyArticleScores = ArticleScores{}
)

func (f Feed) UpdateFromParsed(pf parser.Feed) Feed {
	f.Feed = pf
	f.HubLink = pf.HubLink
	f.SiteLink = pf.SiteLink

	newArticles := make([]Article, len(pf.Articles))

	for i, pa := range pf.Articles {
		a := Article{Article: pa}
		a.FeedId = f.Id

		if pa.Guid != "" {
			a.Guid.Valid = true
			a.Guid.String = pa.Guid
		}

		newArticles[i] = a
	}

	f.Articles = newArticles

	return f
}

func (f Feed) Validate() error {
	if u, err := url.Parse(f.Link); err != nil || !u.IsAbs() {
		return ValidationError{errors.New("Feed has no link")}
	}

	return nil
}

func (a Article) Validate() error {
	if a.FeedId == 0 {
		return ValidationError{errors.New("Article has no feed id")}
	}

	return nil
}

func (asc ArticleScores) Validate() error {
	if asc.ArticleId == 0 {
		return ValidationError{errors.New("Article scores has no article id")}
	}

	return nil
}

func (asc ArticleScores) CalculateScore() int64 {
	return asc.Score1 + int64(0.1*float64(asc.Score2)) + int64(0.01*float64(asc.Score3)) + int64(0.001*float64(asc.Score4)) + int64(0.0001*float64(asc.Score5))
}
