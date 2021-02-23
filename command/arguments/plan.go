package arguments

import (
	"github.com/hashicorp/terraform/tfdiags"
)

// Plan represents the command-line arguments for the plan command.
type Plan struct {
	// State, Operation, and Vars are the common extended flags
	State     *State
	Operation *Operation
	Vars      *Vars

	// Destroy can be set to generate a plan to destroy all infrastructure.
	Destroy bool

	// DetailedExitCode enables different exit codes for error, success with
	// changes, and success with no changes.
	DetailedExitCode bool

	// InputEnabled is used to disable interactive input for unspecified
	// variable and backend config values. Default is true.
	InputEnabled bool

	// OutPath contains an optional path to store the plan file
	OutPath string

	// ViewType specifies which output format to use
	ViewType ViewType
}

// ParsePlan processes CLI arguments, returning a Plan value and errors.
// If errors are encountered, a Plan value is still returned representing
// the best effort interpretation of the arguments.
func ParsePlan(args []string) (*Plan, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics
	apply := &Plan{
		State:     &State{},
		Operation: &Operation{},
		Vars:      &Vars{},
	}

	cmdFlags := extendedFlagSet("plan", apply.State, apply.Operation, apply.Vars)
	cmdFlags.BoolVar(&apply.Destroy, "destroy", false, "destroy")
	cmdFlags.BoolVar(&apply.DetailedExitCode, "detailed-exitcode", false, "detailed-exitcode")
	cmdFlags.BoolVar(&apply.InputEnabled, "input", true, "input")
	cmdFlags.StringVar(&apply.OutPath, "out", "", "out")

	if err := cmdFlags.Parse(args); err != nil {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Failed to parse command-line flags",
			err.Error(),
		))
	}

	args = cmdFlags.Args()

	if len(args) > 0 {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Too many command line arguments",
			"To specify a working directory for the plan, use the global -chdir flag.",
		))
	}

	diags = diags.Append(apply.Operation.Parse())

	switch {
	default:
		apply.ViewType = ViewHuman
	}

	return apply, diags
}
