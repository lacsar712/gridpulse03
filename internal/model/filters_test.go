package model

import "testing"

func TestSlugify(t *testing.T) {
	if got := Slugify("Hello World"); got != "hello-world" {
		t.Fatalf("%q", got)
	}
	if got := Slugify("   "); got != "item" {
		t.Fatalf("%q", got)
	}
}

func TestTagsRoundtrip(t *testing.T) {
	tags := SplitTags("a, b, a,")
	if len(tags) != 2 || tags[0] != "a" || tags[1] != "b" {
		t.Fatalf("%v", tags)
	}
	if JoinTags([]string{"x", " y "}) != "x,y" {
		t.Fatalf("%q", JoinTags([]string{"x", " y "}))
	}
}

func TestListFilterNormalized(t *testing.T) {
	f := ListFilter{Limit: 0, Offset: -1}.Normalized(20)
	if f.Limit != 20 || f.Offset != 0 {
		t.Fatalf("%+v", f)
	}
}
