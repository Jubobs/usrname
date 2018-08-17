package usrname

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"sync"
	"unicode"
)

type Status string

const (
	UnknownStatus Status = "unknown"
	Invalid       Status = "invalid"
	Unavailable   Status = "unavailable"
	Available     Status = "available"
)

type Result struct {
	Username string
	Checker  Checker
	Status   Status
	Message  string
}

type Site interface {
	Name() string
	Link(username string) string
}

type Validator interface {
	Site
	Validate(username string) []Violation
	IllegalPattern() *regexp.Regexp
	Whitelist() *unicode.RangeTable
}

type Checker interface {
	Validator
	Check(client Client) func(string) Result
}

var (
	checkersMu sync.RWMutex
	checkers   = make(map[string]Checker)
)

func Register(name string, checker Checker) error {
	checkersMu.Lock()
	defer checkersMu.Unlock()
	if checker == nil {
		return errors.New("usrname: Register checker is nil")
	}
	if _, dup := checkers[name]; dup {
		return fmt.Errorf("usrname: Register called twice for checker %s", name)
	}
	checkers[name] = checker
	return nil
}

func Checkers() []string {
	checkersMu.RLock()
	defer checkersMu.RUnlock()
	var list []string
	for name := range checkers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func CheckerFor(name string) (Checker, error) {
	checkersMu.RLock()
	defer checkersMu.RUnlock()
	checker, exists := checkers[name]
	if !exists {
		err := fmt.Errorf("usrname: Checker not found for %s", name)
		return nil, err
	}
	return checker, nil
}
