package result

// Callback is the prototype for functions that need to handle the result
// associated with changes to keys.
type Callback func(result *Result)
