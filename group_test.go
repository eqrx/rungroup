// Copyright 2021 Alexander Sowitzki.
// GNU Affero General Public License version 3 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     https://opensource.org/licenses/AGPL-3.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rungroup_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/eqrx/rungroup"
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

	// Wait for all routines to finish and get the error slice. Normally you would go through the list and use the
	// errors package to identify and handle them. Go supports wrapping only one error so just wrapping and returning
	// is not an option ;). You have to decide how you want to handle them.
	//
	// ... and yes, this chicken-egg-solver heavily favours the egg.
	if errs := group.Wait(); len(errs) != 0 {
		fmt.Printf("the error group finished with errors: %#v", errs)
	}
}
