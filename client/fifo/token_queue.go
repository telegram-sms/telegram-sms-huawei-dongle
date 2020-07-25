package fifo

import "container/list"

type TokenQueue struct {
	tokens *list.List
}

func (t *TokenQueue) Init() {
	t.tokens = list.New()
}

func (t *TokenQueue) Peek() string {
	front := t.tokens.Front()
	if front != nil {
		return front.Value.(string)
	}

	// no fifo available
	return ""
}

func (t *TokenQueue) Consume() string {
	front := t.tokens.Front()
	if front != nil {
		result := front.Value.(string)
		t.tokens.Remove(front)
		return result
	}

	// no fifo available
	return ""
}

func (t *TokenQueue) Add(token string) {
	if len(token) > 0 {
		t.tokens.PushBack(token)
	}
}

func (t *TokenQueue) HasAny() bool {
	return t.Len() > 0
}

func (t *TokenQueue) Len() int {
	return t.tokens.Len()
}

func (t *TokenQueue) Reset() {
	t.tokens = list.New()
}
