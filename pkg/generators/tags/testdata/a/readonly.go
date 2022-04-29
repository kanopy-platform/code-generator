package a

type StructWithMemberComments struct {
	// This member can be set.
	//
	// Comments are gibberish...
	// Name must be unique within a namespace.
	// Cannot be updated.
	// +optional
	Name string

	// This member cannot be set.
	//
	// Populated by the system.
	// Read-only.
	UID string
}
