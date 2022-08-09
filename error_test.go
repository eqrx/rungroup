// Copyright (C) 2022 Alexander Sowitzki
//
// This program is free software: you can redistribute it and/or modify it under the terms of the
// GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
// warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
// details.
//
// You should have received a copy of the GNU Affero General Public License along with this program.
// If not, see <https://www.gnu.org/licenses/>.

package rungroup_test

import (
	"errors"
	"testing"

	"eqrx.net/rungroup"
)

func TestErrorsNil(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("empty err did not panic")
		}
	}()

	err := &rungroup.Error{Errs: nil}
	_ = err.Error()
}

func TestErrorsEmpty(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("empty err did not panic")
		}
	}()

	err := &rungroup.Error{Errs: []error{}}
	_ = err.Error()
}

func TestErrorsOne(t *testing.T) {
	t.Parallel()

	expectedMsg := "A B  C"

	actualMsg := (&rungroup.Error{Errs: []error{errors.New(expectedMsg)}}).Error()

	if actualMsg != expectedMsg {
		t.Fatalf("expected error message to be \"%s\", was %s", expectedMsg, actualMsg)
	}
}

func TestErrorsMany(t *testing.T) {
	t.Parallel()

	errs := &rungroup.Error{Errs: []error{
		errors.New("A "),
		errors.New("B"),
		errors.New(" C"),
	}}

	actualMsg := errs.Error()
	expectedMsg := "[ A , B,  C ]"

	if actualMsg != expectedMsg {
		t.Fatalf("expected error message to be \"%s\", was %s", expectedMsg, actualMsg)
	}
}

func TestErrorsNested(t *testing.T) {
	t.Parallel()

	errs := &rungroup.Error{Errs: []error{
		errors.New("A"),
		&rungroup.Error{Errs: []error{errors.New("B"), errors.New("C")}},
		&rungroup.Error{Errs: []error{
			&rungroup.Error{Errs: []error{errors.New("D")}},
		}},
	}}

	actualMsg := errs.Error()
	expectedMsg := "[ A, [ B, C ], D ]"

	if actualMsg != expectedMsg {
		t.Fatalf("expected error message to be \"%s\", was %s", expectedMsg, actualMsg)
	}
}
