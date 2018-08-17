package usrname_test

import (
	"reflect"
	"testing"

	"github.com/fortytw2/leaktest"
	"github.com/jubobs/usrname"
	_ "github.com/jubobs/usrname/disqus"
	_ "github.com/jubobs/usrname/facebook"
	_ "github.com/jubobs/usrname/github"
	_ "github.com/jubobs/usrname/instagram"
	_ "github.com/jubobs/usrname/medium"
	_ "github.com/jubobs/usrname/pinterest"
	_ "github.com/jubobs/usrname/reddit"
	_ "github.com/jubobs/usrname/twitter"
)

const template = "Checkers(), got %q, want %q"

func TestCheckers(t *testing.T) {
	defer leaktest.Check(t)()
	expected := []string{
		"Disqus",
		"GitHub",
		"Instagram",
		"Medium",
		"Pinterest",
		"Twitter",
		"facebook",
		"reddit",
	}
	if actual := usrname.Checkers(); !reflect.DeepEqual(actual, expected) {
		t.Errorf(template, actual, expected)
	}
}
