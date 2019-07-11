// Package branch provides use-cases for creating or updating a branch.
package branch

import (
	"github.com/google/wire"
	"github.com/int128/ghcp/usecases"
)

var Set = wire.NewSet(
	wire.Struct(new(CreateBranch), "*"),
	wire.Bind(new(usecases.CreateBranch), new(*CreateBranch)),
	wire.Struct(new(UpdateBranch), "*"),
	wire.Bind(new(usecases.UpdateBranch), new(*UpdateBranch)),
)
