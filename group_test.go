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

package rungroup_test

import (
	"context"
	"errors"
	"fmt"

	"eqrx.net/rungroup"
)

var (
	ErrEgg     = errors.New("the egg came first")
	ErrChicken = errors.New("the chicken came first")
)

func ExampleGroup() {
	// Create a new group instance. Normally you would you current in here so the group gets canceled when it does.
	group := rungroup.New(context.Background())

	// This routine immediately return without error and causes the group context to be cancelled.
	// Since the error is nil it will not be added to the error list.
	group.Go(func(c context.Context) error { return nil })

	// Routine returns immediately with an error: The group gets cancelled and the error gets added to the list.
	group.Go(func(c context.Context) error { return ErrEgg })

	// Return immediately but do not cancel other routines. This is done by not returning an error
	// and starting with rungroup.NoCancelOnSucess set.
	group.Go(func(c context.Context) error { return nil }, rungroup.NoCancelOnSuccess)

	// When a routine returns with an error the context will be canceled,
	// regardless if rungroup.NoCancelOnSucess is set or not.
	group.Go(func(c context.Context) error { return ErrChicken }, rungroup.NoCancelOnSuccess)

	// Wait until the context of the group is canceled and return the context error.
	// Since the context is already done returning from here will not affect the context,
	// the error is stored though.
	group.Go(func(ctx context.Context) error {
		<-ctx.Done()

		fmt.Println("the question has been answered!")

		return nil
	})

	// Wait for all routines to finish and get all errors in a bundle. This bundled error does have a Unwrap method
	// since Go wrapping only supports for unwrap a single error and the code can't decide for you which one is the
	// most relevant. In case you want to handle all returned errors you may iterate over them. If you print or wrap
	// it, all wrapped errors will be printed.
	//
	// ... and yes, this chicken-egg-solver heavily favours the egg.
	if err := group.Wait(); err != nil {
		fmt.Printf("the error group failed: %v", err)
	}
}
