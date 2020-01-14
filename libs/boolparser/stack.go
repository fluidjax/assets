/*
The MIT License (MIT)

Copyright (c) 2015 Marc Abi Khalil

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

Original Source : https://github.com/marcmak/calc			2015
Updated & Modified: Christopher Morris chris@qredo.com 		2020

*/

package boolparser

// Stack is a LIFO data structure
type Stack struct {
	Values []Token
}

// Pop removes the token at the top of the stack and returns its value
func (self *Stack) Pop() Token {
	if len(self.Values) == 0 {
		return Token{}
	}
	token := self.Values[len(self.Values)-1]
	self.Values = self.Values[:len(self.Values)-1]
	return token
}

// Push adds tokens to the top of the stack
func (self *Stack) Push(i ...Token) {
	self.Values = append(self.Values, i...)
}

// Peek returns the token at the top of the stack
func (self *Stack) Peek() Token {
	if len(self.Values) == 0 {
		return Token{}
	}
	return self.Values[len(self.Values)-1]
}

// EmptyInto dumps all tokens from one stack to another
func (self *Stack) EmptyInto(s *Stack) {
	if !self.IsEmpty() {
		for i := self.Length() - 1; i >= 0; i-- {
			s.Push(self.Pop())
		}
	}
}

// IsEmpty checks if there are any tokens in the stack
func (self *Stack) IsEmpty() bool {
	return len(self.Values) == 0
}

// Length returns the amount of tokens in the stack
func (self *Stack) Length() int {
	return len(self.Values)
}
