package testing

import "testing"

type TestHandler interface {
	RunTest(*testing.T)
}

// The TestHandlerFunc type is an adapter to allow the use of
// ordinary functions as test handlers. If f is a function
// with the appropriate signature, TestHandlerFunc(f) is a
// Handler that calls f.
type TestHandlerFunc func(*testing.T)

// RunTest calls f(t).
func (f TestHandlerFunc) RunTest(t *testing.T) {
	f(t)
}

// A constructor for a piece of test middleware.
type Constructor func(TestHandler) TestHandler

// Chain acts as a list of TestHandler constructors.
// Chain is effectively immutable:
// once created, it will always hold
// the same set of constructors in the same order.
type Chain struct {
	constructors []Constructor
}

// New creates a new chain,
// memorizing the given list of test middleware constructors.
// New serves no other function,
// constructors are only called upon a call to Then().
func New(constructors ...Constructor) Chain {
	return Chain{append(([]Constructor)(nil), constructors...)}
}

// Then chains the middleware and returns the final http.Handler.
//     New(m1, m2, m3).Then(h)

// Then() treats nil as panic
func (c *Chain) Then(h TestHandler) TestHandler {
	if h == nil {
		panic("test handler is null")
	}
	for i := range c.constructors {
		h = c.constructors[len(c.constructors)-1-i](h)
	}
	return h
}

// ThenFunc works identically to Then, but takes
// a TestHandlerFunc instead of a TestHandler.
//
// The following two statements are equivalent:
//     c.Then(TestHandlerFunc(fn))
//     c.ThenFunc(fn)
//
// ThenFunc() treats nil as panic
func (c *Chain) ThenFunc(fn TestHandlerFunc) TestHandler {
	if fn == nil {
		panic("test handler function is null")
	}
	return c.Then(fn)
}

// Append extends a chain, adding the specified constructors
// as the last ones in the request flow.
//
// Append returns a new chain, leaving the original one untouched.
//
//     stdChain := alice.New(m1, m2)
//     extChain := stdChain.Append(m3, m4)
//     // requests in stdChain go m1 -> m2
//     // requests in extChain go m1 -> m2 -> m3 -> m4
func (c *Chain) Append(constructors ...Constructor) Chain {
	newOne := make([]Constructor, len(c.constructors)+len(constructors))
	newOne = append(newOne, c.constructors...) //move old constructors into new set
	newOne = append(newOne, constructors...)   //append new constructors into new set
	return Chain{newOne}
}

// Extend extends a chain by adding the specified chain
// as the last one in the request flow.
//
// Extend returns a new chain, leaving the original one untouched.
//
//     stdChain := alice.New(m1, m2)
//     ext1Chain := alice.New(m3, m4)
//     ext2Chain := stdChain.Extend(ext1Chain)
//     // requests in stdChain go  m1 -> m2
//     // requests in ext1Chain go m3 -> m4
//     // requests in ext2Chain go m1 -> m2 -> m3 -> m4

func (c *Chain) Extend(chain Chain) Chain {
	return c.Append(chain.constructors...)
}
