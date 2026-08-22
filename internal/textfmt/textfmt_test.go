package textfmt

import "testing"

func TestExcerptAndWrap(t *testing.T) {
	if Excerpt("abcdef", 3) != "abc..." {
		t.Fatalf("%q", Excerpt("abcdef", 3))
	}
	if WordCount("a b  c") != 3 {
		t.Fatalf("%d", WordCount("a b  c"))
	}
	if Wrap("one two three", 8) == "" {
		t.Fatal("wrap empty")
	}
}

func TestCSV(t *testing.T) {
	got := SplitCSVLine(`a,"b,c",d`)
	if len(got) != 3 || got[1] != "b,c" {
		t.Fatalf("%v", got)
	}
	if JoinCSV([]string{"a", "b,c"}) != `a,"b,c"` {
		t.Fatalf("%q", JoinCSV([]string{"a", "b,c"}))
	}
}

func TestMarkdown(t *testing.T) {
	if Heading(2, "Hi") != "## Hi" {
		t.Fatal(Heading(2, "Hi"))
	}
	if !IsIdent("ab_1") || IsIdent("1ab") {
		t.Fatal("ident")
	}
}
