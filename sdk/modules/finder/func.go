package finder

import "net"

// Func is the function that is called for each robot that is found. The first
// boolean returned indicates if the given robot should be added to the list of
// found robots. The second boolean returned indicates if the search should
// continue or not.
//
// Note that the Text SDK protocol does not support serial numbers, so it will
// always be nil in that case.
type Func func(ip net.IP, serial []byte) (bool, bool)
