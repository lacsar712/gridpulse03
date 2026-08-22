package kvbag

type Bag struct {
	m map[string]string
}

func New() *Bag {
	return &Bag{m: map[string]string{}}
}

func (b *Bag) Set(k, v string) {
	b.m[k] = v
}

func (b *Bag) Get(k string) (string, bool) {
	if b == nil || b.m == nil {
		return "", false
	}
	v, ok := b.m[k]
	return v, ok
}

func (b *Bag) Len() int {
	if b == nil || b.m == nil {
		return 0
	}
	return len(b.m)
}
