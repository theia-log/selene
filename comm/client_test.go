package comm

import "testing"

func TestEventFilterCreate(t *testing.T) {
	filter := Filter(11.2).
		MatchEnd(12.3).
		MatchContent("content").
		MatchTag("tag1", "tag2").
		OrderAsc()

	if filter.Content == nil || *filter.Content != "content" {
		t.Fatal("Content match improperly set")
	}

	if filter.End == nil || *filter.End != 12.3 {
		t.Fatal("End timestamp impropery set")
	}

	if filter.Start != 11.2 {
		t.Fatal("Start timestamp improperly set")
	}

	if filter.Order == nil || *filter.Order != OrderAsc {
		t.Fatal("Order improperly set")
	}

	if filter.Tags == nil || len(filter.Tags) != 2 {
		t.Fatal("Expected 2 tags filters")
	}

	// reset order

	filter.OrderDesc()

	if filter.Order == nil || *filter.Order != OrderDesc {
		t.Fatal("Order improperly re-set, expected to be reset to Desc")
	}
}
