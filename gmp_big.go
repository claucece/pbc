/*
	Copyright © 2015 Nik Unger

	This file is part of The PBC Go Wrapper.

	The PBC Go Wrapper is free software: you can redistribute it and/or modify
	it under the terms of the GNU Lesser General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	The PBC Go Wrapper is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
	GNU Lesser General Public License for more details.

	You should have received a copy of the GNU Lesser General Public License
	along with The PBC Go Wrapper. If not, see <http://www.gnu.org/licenses/>.

	The PBC Go Wrapper makes use of The PBC library. The PBC Library and its
	use are covered under the terms of the GNU Lesser General Public License
	version 3, or (at your option) any later version.
*/

package pbc

/*
#include <gmp.h>
*/
import "C"

import (
	"math/big"
	"runtime"
	"unsafe"
)

var wordSize C.size_t
var bitsPerWord C.size_t

func clearMpz(x *C.mpz_t) {
	C.mpz_clear(&x[0])
}

func newMpz() *C.mpz_t {
	out := &C.mpz_t{}
	C.mpz_init(&out[0])
	runtime.SetFinalizer(out, clearMpz)
	return out
}

// big2thisMpz imports the value of num into out
func big2thisMpz(num *big.Int, out *C.mpz_t) {
	words := num.Bits()
	if len(words) > 0 {
		C.mpz_import(&out[0], C.size_t(len(words)), -1, wordSize, 0, 0, unsafe.Pointer(&words[0]))
	}
}

// big2mpz allocates a new mpz_t and imports a big.Int value
func big2mpz(num *big.Int) *C.mpz_t {
	out := newMpz()
	big2thisMpz(num, out)
	return out
}

// mpz2thisBig imports the value of num into out
func mpz2thisBig(num *C.mpz_t, out *big.Int) {
	wordsNeeded := (C.mpz_sizeinbase(&num[0], 2) + (bitsPerWord - 1)) / bitsPerWord
	words := make([]big.Word, wordsNeeded)
	var wordsWritten C.size_t
	C.mpz_export(unsafe.Pointer(&words[0]), &wordsWritten, -1, wordSize, 0, 0, &num[0])
	out.SetBits(words)
}

// mpz2big allocates a new big.Int and imports an mpz_t value
func mpz2big(num *C.mpz_t) (out *big.Int) {
	out = &big.Int{}
	mpz2thisBig(num, out)
	return
}

func init() {
	var oneWord big.Word
	size := unsafe.Sizeof(oneWord)
	wordSize = C.size_t(size)
	bitsPerWord = C.size_t(8 * size)
}
