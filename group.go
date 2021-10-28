// Copyright (C) 2021 Alexander Sowitzki
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

// Package rungroup allows goroutines to share the same lifecycle.
//
// Like https://pkg.go.dev/golang.org/x/sync/errgroup, rungroup allows to create a handler that allows to run given
// functions in parallel with goroutines. It tries to be as close to the interface of errgroup but handles the
// cancelation of its routines and treatment of errors differently.
//
// Differences to errgroup
//
// Goroutine function gets a context as first argument, a new group does not get handled to the caller with the groups
// context. By default all routines get canceled as soon as one routine return, regardless if the error is nil or not.
// This can be overridden per routine. errgroup only cancels on the first non nil error. All non nil errors are
// returned to the creator of the group. errgroup only returns the first error and drops the rest.
package rungroup

import (
	"context"
	"sync"
)

// Group represents a set of goroutines which lifecycles are bound to each other.
type Group struct {
	ctx    context.Context // Passed to spawned goroutines.
	cancel func()          // Cancels ctx.
	wg     sync.WaitGroup  // WaitGroup that completes when all routines spawned from Group have returned.
	mtx    sync.Mutex      // Lock for access to errs.
	errs   []error         // Slice of errors that were returned by spawned goroutines.
}

type (
	// optionSet contains settings regarding spawned go routines. Only to be used with Option.
	optionSet struct{ noCancelOnSuccess bool }
	// Option modifies the settings of a spawned routine with Group.Go.
	Option func(o *optionSet)
)

// NoCancelOnSuccess prevents goroutines spawned with Group.Go to cancel the group context then they return
// a non nil error. Default is to cancel the group context on return regardless of the returned error.
func NoCancelOnSuccess(o *optionSet) { o.noCancelOnSuccess = true }

// New creates group for goroutine management. The context passed as parameter ctx with be taken as parent for the
// group context. Canceling it will cancel all spawned goroutines. ctx must not be nil.
func New(ctx context.Context) *Group {
	ctx, cancel := context.WithCancel(ctx)

	return &Group{ctx, cancel, sync.WaitGroup{}, sync.Mutex{}, []error{}}
}

// Wait block until all goroutines of the group have returned to it and returns a slice containing the errors that were
// returned by the routines. The slice is never nil and never contains nil values. The lowest index contains the error
// that was returned by the chronologically first finishing go routine (don't rely on order, there is no anti-scheduler
// magic here). A zero length slice indicates that no routine returned a non nil error. This method must not be called
// by multiple goroutines at the same time. After this call returnes, the group may not be reused.
func (g *Group) Wait() []error {
	g.wg.Wait()

	g.mtx.Lock()
	errs := g.errs
	g.errs = nil
	g.mtx.Unlock()

	return errs
}

// Go spawns a goroutine and calls the function fn with it. The context of the group is passed as the first
// argument to it.
//
// When any routine spawned by Go return, the following things happen: Panics are not recovered. If the returned error
// value of fn is non nil it is stored for retrieval by Wait. If NoCancelOnSuccess is part of opts the group context
// will be canceled if the returned error is not nil. If NoCancelOnSuccess is NOT part of opts the group context will
// be canceled regardless of the returned error.
//
// As long as no call from Wait has returned, Go may be called by any goroutines at the same time. Passing nil as fn or
// part of opts is not allowed.
func (g *Group) Go(fn func(context.Context) error, opts ...Option) {
	g.wg.Add(1)

	os := &optionSet{}
	for _, o := range opts {
		o(os)
	}

	go func() {
		defer g.wg.Done()

		if !os.noCancelOnSuccess {
			defer g.cancel()
		}

		if err := fn(g.ctx); err != nil {
			if os.noCancelOnSuccess {
				g.cancel()
			}

			g.mtx.Lock()
			g.errs = append(g.errs, err)
			g.mtx.Unlock()
		}
	}()
}
