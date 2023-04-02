package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddSubscription(t *testing.T) {
	teardown := setup()
	defer teardown()

	testCases := []struct {
		sub Subscription
		err error
	}{
		{
			Subscription{
				URL:         "blog.tumblr.com",
				BlogName:    "blog",
				Source:      "tumblr",
				Title:       "some title",
				Description: "useful description",
			},
			nil,
		},
	}

	for _, testCase := range testCases {
		err := repo.AddSubscription(&testCase.sub)
		checkErrors(t, testCase.err, err)
		if err == nil {
			var dbSub Subscription
			DB.Get(
				&dbSub,
				"SELECT id, created_at, url, blog_name, source, title, description FROM subscriptions WHERE url = ? LIMIT 1",
				testCase.sub.URL,
			)
			assert.NotEqual(t, uint(0), testCase.sub.ID, "ID is not set for saved subscription")
			assert.NotEqual(t, uint(0), dbSub.ID, "ID is not fetched from database")
			assert.Equal(t, testCase.sub.URL, dbSub.URL)
			assert.Equal(t, testCase.sub.BlogName, dbSub.BlogName)
			assert.Equal(t, testCase.sub.Source, dbSub.Source)
			assert.Equal(t, testCase.sub.Title, dbSub.Title)
			assert.Equal(t, testCase.sub.Description, dbSub.Description)
		}
	}
}

func TestRemoveAllSubscriptions(t *testing.T) {
	teardown := setup()
	defer teardown()

	repo.AddSubscription(&Subscription{URL: "blog.tumblr.com", BlogName: "blog", Source: "tumblr"})
	repo.AddSubscription(&Subscription{URL: "a.tumblr.com", BlogName: "a", Source: "tumblr"})

	var subsCount int
	DB.Get(&subsCount, "SELECT count(*) FROM subscriptions")
	if subsCount != 2 {
		t.Errorf("Expected subscriptions count to be [2], got [%d]", subsCount)
	}

	repo.RemoveAllSubscriptions()
	DB.Get(&subsCount, "SELECT count(*) FROM subscriptions")

	if subsCount != 0 {
		t.Errorf("Expected subscriptions count to be [0], got [%d]", subsCount)
	}
}

func TestListSubscriptions(t *testing.T) {
	teardown := setup()
	defer teardown()
	repo.AddSubscription(&Subscription{URL: "blog.tumblr.com", BlogName: "blog", Source: "tumblr"})
	repo.AddSubscription(&Subscription{URL: "a.tumblr.com", BlogName: "a", Source: "tumblr"})

	subs, err := repo.ListSubscriptions()
	if err != nil {
		t.Errorf("Unexpected error on listing subscriptions: [%#v]", err)
	}

	assert.NotNil(t, subs[0].ID, "ID is nil")
	assert.Equal(t, "blog.tumblr.com", subs[0].URL)
	assert.Equal(t, "blog", subs[0].BlogName)
	assert.Equal(t, "tumblr", subs[0].Source)

	assert.NotNil(t, subs[1].ID, "ID is nil")
	assert.Equal(t, "a.tumblr.com", subs[1].URL)
	assert.Equal(t, "a", subs[1].BlogName)
	assert.Equal(t, "tumblr", subs[1].Source)
}
