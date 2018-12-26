package main

import (
	"math/rand"
	"time"
)

// Quote holds a single quote and the person who said it
type Quote struct {
	Text   string `json:"text"`
	Person string `json:"person"`
}

// QuotesCollection holds a collection of quotes!
type QuotesCollection struct {
	quotes []*Quote
}

// AddQuote creates a Quote object with the given inputs and
// adds it to the Quotes collection
func (Q *QuotesCollection) AddQuote(text, person string) {
	Q.quotes = append(Q.quotes, &Quote{Text: text, Person: person})
}

// GetQuote returns a random Quote object from the Quotes collection
func (Q *QuotesCollection) GetQuote() *Quote {
	return Q.quotes[rand.Intn(len(Q.quotes))]
}

// newQuotesCollection creates a new QuotesCollection
func newQuotesCollection() *QuotesCollection {
	result := &QuotesCollection{}
	rand.Seed(time.Now().Unix())
	result.AddQuote("Age is an issue of mind over matter. If you don't mind, it doesn't matter", "Mark Twain")
	result.AddQuote("Anyone who stops learning is old, whether at twenty or eighty. Anyone who keeps learning stays young. The greatest thing in life is to keep your mind young", "Henry Ford")
	result.AddQuote("Wrinkles should merely indicate where smiles have been", "Mark Twain")
	result.AddQuote("True terror is to wake up one morning and discover that your high school class is running the country", "Kurt Vonnegut")
	result.AddQuote("A diplomat is a man who always remembers a woman's birthday but never remembers her age", "Robert Frost")
	result.AddQuote("As I grow older, I pay less attention to what men say. I just watch what they do", "Andrew Carnegie")
	result.AddQuote("How incessant and great are the ills with which a prolonged old age is replete", "C. S. Lewis")
	result.AddQuote("Old age, believe me, is a good and pleasant thing. It is true you are gently shouldered off the stage, but then you are given such a comfortable front stall as spectator", "Confucius")
	result.AddQuote("Old age has deformities enough of its own. It should never add to them the deformity of vice", "Eleanor Roosevelt")
	result.AddQuote("Nobody grows old merely by living a number of years. We grow old by deserting our ideals. Years may wrinkle the skin, but to give up enthusiasm wrinkles the soul", "Samuel Ullman")
	result.AddQuote("An archaeologist is the best husband a woman can have. The older she gets the more interested he is in her", "Agatha Christie")
	result.AddQuote("All diseases run into one, old age", "Ralph Waldo Emerson")
	result.AddQuote("Bashfulness is an ornament to youth, but a reproach to old age", "Aristotle")
	result.AddQuote("Like everyone else who makes the mistake of getting older, I begin each day with coffee and obituaries", "Bill Cosby")
	result.AddQuote("Age appears to be best in four things old wood best to burn, old wine to drink, old friends to trust, and old authors to read", "Francis Bacon")
	result.AddQuote("None are so old as those who have outlived enthusiasm", "Henry David Thoreau")
	result.AddQuote("Every man over forty is a scoundrel", "George Bernard Shaw")
	result.AddQuote("Forty is the old age of youth fifty the youth of old age", "Victor Hugo")
	result.AddQuote("You can't help getting older, but you don't have to get old", "George Burns")
	result.AddQuote("Alas, after a certain age every man is responsible for his face", "Albert Camus")
	result.AddQuote("Youth is when you're allowed to stay up late on New Year's Eve. Middle age is when you're forced to", "Bill Vaughan")
	result.AddQuote("Old age is like everything else. To make a success of it, you've got to start young", "Theodore Roosevelt")
	result.AddQuote("A comfortable old age is the reward of a well-spent youth. Instead of its bringing sad and melancholy prospects of decay, it would give us hopes of eternal youth in a better world", "Maurice Chevalier")
	result.AddQuote("A man growing old becomes a child again", "Sophocles")
	result.AddQuote("I will never be an old man. To me, old age is always 15 years older than I am", "Francis Bacon")
	result.AddQuote("Age considers youth ventures", "Rabindranath Tagore")
	return result
}
