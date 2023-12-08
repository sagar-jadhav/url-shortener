package metrics

import (
	"testing"
)

func Test_Add(t *testing.T) {
	metrics := DomainMetrics{}
	domain := "localhost"

	metrics.Add(domain)
	if metrics.Get(domain) != 1 {
		t.Fatalf("Add() => expected 1 got %d", metrics.Get(domain))
	}

	metrics.Add(domain)
	if metrics.Get(domain) != 2 {
		t.Fatalf("Add() => expected 2 got %d", metrics.Get(domain))
	}
}

func Test_Sort(t *testing.T) {
	metrics := DomainMetrics{}
	localhostDomain := "localhost"
	githubDomain := "github"

	metrics.Add(localhostDomain)
	metrics.Add(localhostDomain)
	metrics.Add(localhostDomain)

	metrics.Add(githubDomain)
	metrics.Add(githubDomain)

	slice := metrics.Sort()

	if len(slice) != 2 {
		t.Fatalf("Sort() => expected length to be 2 but got %d", len(slice))
	}
	if slice[0].Domain != localhostDomain {
		t.Fatalf("Sort() => expected slice[0].domain to be %s but got %s", localhostDomain, slice[0].Domain)
	}
	if slice[0].ShortenedCount != 3 {
		t.Fatalf("Sort() => expected slice[0].shortenedCount to be %d but got %d", 3, slice[0].ShortenedCount)
	}
	if slice[1].Domain != githubDomain {
		t.Fatalf("Sort() => expected slice[1].domain to be %s but got %s", githubDomain, slice[1].Domain)
	}
	if slice[1].ShortenedCount != 2 {
		t.Fatalf("Sort() => expected slice[1].shortenedCount to be %d but got %d", 2, slice[1].ShortenedCount)
	}
}
