package application

import (
	"encoding/json"
	"sync"
)

// Pool for CallOptions - most critical allocation
var callOptionsPool = sync.Pool{
	New: func() interface{} {
		return &CallOptions{
			Args: make([]json.RawMessage, 0, 4), // Pre-allocate common size
		}
	},
}

// Pool for Args - frequent allocation for parameter parsing
var argsPool = sync.Pool{
	New: func() interface{} {
		return &Args{
			data: make(map[string]any, 8), // Pre-allocate common size
		}
	},
}

// Pool for QueryParams - frequent allocation in HTTP requests
var queryParamsPool = sync.Pool{
	New: func() interface{} {
		return make(QueryParams, 8) // Pre-allocate common size
	},
}

// Pool for CallError - frequent allocation during error handling
var callErrorPool = sync.Pool{
	New: func() interface{} {
		return &CallError{}
	},
}

// Pool for Parameter slices - used in method binding
var parameterSlicePool = sync.Pool{
	New: func() interface{} {
		return make([]*Parameter, 0, 8) // Pre-allocate common size
	},
}

// GetCallOptions retrieves a CallOptions from the pool
func GetCallOptions() *CallOptions {
	opts := callOptionsPool.Get().(*CallOptions)
	opts.Reset()
	return opts
}

// PutCallOptions returns a CallOptions to the pool
func PutCallOptions(opts *CallOptions) {
	if opts != nil {
		callOptionsPool.Put(opts)
	}
}

// Reset clears CallOptions for reuse
func (co *CallOptions) Reset() {
	co.MethodID = 0
	co.MethodName = ""
	co.Args = co.Args[:0] // Keep underlying capacity
}

// GetArgs retrieves an Args from the pool
func GetArgs() *Args {
	args := argsPool.Get().(*Args)
	args.Reset()
	return args
}

// PutArgs returns an Args to the pool
func PutArgs(args *Args) {
	if args != nil {
		argsPool.Put(args)
	}
}

// Reset clears Args for reuse
func (a *Args) Reset() {
	if a.data != nil {
		// Clear map but keep underlying capacity
		for k := range a.data {
			delete(a.data, k)
		}
	} else {
		a.data = make(map[string]any, 8)
	}
}

// GetQueryParams retrieves a QueryParams from the pool
func GetQueryParams() QueryParams {
	qp := queryParamsPool.Get().(QueryParams)
	// Clear map but keep underlying capacity
	for k := range qp {
		delete(qp, k)
	}
	return qp
}

// PutQueryParams returns a QueryParams to the pool
func PutQueryParams(qp QueryParams) {
	if qp != nil {
		queryParamsPool.Put(qp)
	}
}

// GetCallError retrieves a CallError from the pool
func GetCallError() *CallError {
	err := callErrorPool.Get().(*CallError)
	err.Reset()
	return err
}

// PutCallError returns a CallError to the pool
func PutCallError(err *CallError) {
	if err != nil {
		callErrorPool.Put(err)
	}
}

// Reset clears CallError for reuse
func (ce *CallError) Reset() {
	ce.Kind = ""
	ce.Message = ""
	ce.Cause = nil
}

// GetParameterSlice retrieves a Parameter slice from the pool
func GetParameterSlice() []*Parameter {
	slice := parameterSlicePool.Get().([]*Parameter)
	return slice[:0] // Reset length but keep capacity
}

// PutParameterSlice returns a Parameter slice to the pool
func PutParameterSlice(slice []*Parameter) {
	if slice != nil {
		// Clear references to avoid memory leaks
		for i := range slice {
			slice[i] = nil
		}
		parameterSlicePool.Put(slice)
	}
}

// ArgsFromQueryParams creates Args from QueryParams using pooled objects
func ArgsFromQueryParams(qp QueryParams) (*Args, error) {
	args := GetArgs()
	
	argData := qp["args"]
	if len(argData) == 1 {
		err := json.Unmarshal([]byte(argData[0]), &args.data)
		if err != nil {
			PutArgs(args) // Return to pool on error
			return nil, err
		}
	}
	
	return args, nil
}