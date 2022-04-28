package a

// +kanopy:builder=true
type UnspecifiedReceiver struct {
}

// +kanopy:builder=true
// +kanopy:receiver=ptr
type PtrReceiver struct {
}

// +kanopy:builder=false
// +kanopy:receiver=value
type ValueReceiver struct {
}
