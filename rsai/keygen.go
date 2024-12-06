package rsai

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var (
	e = big.NewInt(3)
)

type PrivateKey struct {
	N *big.Int `json:"n"`
	D *big.Int `json:"d"`
}

/*
Creates a new private key

Parameters:

	n: The modulus, a part of both the public and private keys.
	d: The private exponent used for decryption.

Returns:
p: The new private key
*/
func newPrivateKey(n *big.Int, d *big.Int) *PrivateKey {
	p := PrivateKey{
		N: n,
		D: d,
	}
	return &p
}

type PublicKey struct {
	N *big.Int
	E *big.Int
}

/*
Creates a new public key

Parameters:

	n: The modulus, a part of both the public and private keys.
	e: The public exponent used for encryption.

Returns:
p: The new public key
*/
func NewPublicKey(n *big.Int, e *big.Int) *PublicKey {
	p := PublicKey{
		N: n,
		E: e,
	}
	return &p
}

/*
KeyGen generates RSA key pairs (public and private keys) based on the specified
bit length for the modulus.

Parameters:

	k : The key size.

Returns:

	n: The modulus, a part of both the public and private keys.
	d: The private exponent used for decryption.
	e: The public exponent used for encryption.
*/
func KeyGen(k int) (*big.Int, *big.Int, *big.Int) {
	p, q := generatePrimePair(e, k)
	d := computeD(p, q)
	n := new(big.Int).Mul(p, q)
	return n, d, e
}

func KeyGenStruct(k int) (*PrivateKey, *PublicKey) {
	n, d, e := KeyGen(k)
	sk := newPrivateKey(n, d)
	pk := NewPublicKey(n, e)
	return sk, pk
}

/*
generateRandomPrime generates a random prime number of half the specified key size.

Parameters:

	keySize : The bit length of the key size

Returns:

	p : A random prime number (probably)
	  https://yourbasic.org/golang/check-prime/
*/
func generateRandomPrime(keySize int) *big.Int {
	p, err := rand.Prime(rand.Reader, keySize/2)

	if err != nil {
		fmt.Printf("Error generating random prime: %v", err)
	}
	return p
}

/*
gcdS computes the greatest common divisor (gcd) of two given big.Int values and checks if it equals 1.

Parameters:

	a : The first value.
	b : The second value.

Returns:

	bool: True if the gcd of a and b is equal to 1; otherwise, false.
*/
func gcdS(b *big.Int) bool {
	got := new(big.Int).GCD(nil, nil, e, big.NewInt(0).Sub(b, big.NewInt(1)))
	want := big.NewInt(1)
	comparison := got.Cmp(want)
	if comparison == 0 {
		return true
	} else {
		return false
	}
}

/*
generatePrimePair generates a pair of random prime numbers ensuring gcd(e, p, q) = 1 for the generated primes.

Parameters:

	e : The public exponent (e) used for encryption.
	keySize : The bit length of the desired prime numbers, which determines the key size.

Returns:

	p, q : A pair of random prime numbers (p and q) satisfying the gcd(e, p, q) = 1 condition.

Notes:

	The function will make sure that gcd(e, p, q) = 1 is fulfilled on the primes that are returned.
*/
func generatePrimePair(e *big.Int, keySize int) (*big.Int, *big.Int) {
	lookingForPair := true

	for lookingForPair {
		p := generateRandomPrime(keySize)
		q := generateRandomPrime(keySize)

		// We need to remember to check that p & q are not alike
		gcdPAndE := gcdS(p)
		gcdQAndE := gcdS(q)

		bitlen := big.NewInt(0).Mul(p, q).BitLen()
		bitLenIsK := bitlen == keySize

		if gcdPAndE && gcdQAndE && bitLenIsK {
			lookingForPair = false
			// If we get here, everything went as planned
			return p, q
		}
	}

	// If we get here it means something went wrong
	err1 := big.NewInt(0)
	err2 := big.NewInt(0)
	return err1, err2
}

/*
FromBigIntToFloat converts a big.Int to a float64.
Global function

Parameters:

	myBigInt : The big.Int number to be converted.

Returns:

	result : The converted float64 value.
*/
func FromBigIntToFloat(myBigInt *big.Int) float64 {
	tempfloat := new(big.Float).SetInt(myBigInt)
	result, _ := tempfloat.Float64() // _ is the accuracy
	return result
}

/*
computeD calculates the private exponent (d) for RSA encryption using the provided prime numbers p and q.

Parameters:

	p : The first prime number.
	q : The second prime number.

Returns:

	d : The computed private exponent.
*/
func computeD(p *big.Int, q *big.Int) *big.Int {
	n := new(big.Int)
	n.Mul(p, q)

	// Calculate the first part of the equation
	part1 := new(big.Int)
	part1.Sub(p, big.NewInt(1))
	part1.Mul(part1, new(big.Int).Sub(q, big.NewInt(1)))

	// Calculate d using the modular multiplicative inverse
	d := new(big.Int)
	d.ModInverse(e, part1)
	return d
}
