package tweetnacl

/*
#include "tweetnacl.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// The number of bytes in a crypto_box public key
const crypto_box_PUBLICKEYBYTES int = 32

// The number of bytes in a crypto_box secret key
const crypto_box_SECRETKEYBYTES int = 32

// The number of zero padding bytes for a crypto_box message
const crypto_box_ZEROBYTES int = 32

// The number of zero padding bytes for a crypto_box ciphertext
const crypto_box_BOXZEROBYTES int = 16

// Constant zero-filled byte array used for padding messages
var crypto_box_PADDING = []byte{0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00}

type KeyPair struct {
	PublicKey []byte
	SecretKey []byte
}

// Wrapper function for crypto_box_keypair.
//
// Randomly generates a secret key and a corresponding public key. It guarantees that the secret key
// has crypto_box_PUBLICKEYBYTES bytes and that the public key has crypto_box_SECRETKEYBYTES bytes,
// returns a KeyPair initialised with a crypto_box public/private key pair.
//
// Ref. http://nacl.cr.yp.to/box.html
func CryptoBoxKeyPair() (*KeyPair, error) {
	pk := make([]byte, crypto_box_PUBLICKEYBYTES)
	sk := make([]byte, crypto_box_SECRETKEYBYTES)
	rc := C.crypto_box_keypair((*C.uchar)(unsafe.Pointer(&pk[0])), (*C.uchar)(unsafe.Pointer(&sk[0])))

	if rc == 0 {
		return &KeyPair{PublicKey: pk, SecretKey: sk}, nil
	}

	return nil, fmt.Errorf("Error generating key pair (error code %v)", rc)
}

// Wrapper function for crypto_box.
//
// Encrypts and authenticates the message using the secretKey, publicKey and nonce. The zero padding
// required by the crypto_box C API is added internally and should not be included in the supplied
// message. Likewise the zero padding that prefixes the ciphertext returned by the crypto_box C API
// is stripped from the returned ciphertext.
//
// Ref. http://nacl.cr.yp.to/box.html
func CryptoBox(message, nonce, publicKey, secretKey []byte) ([]byte, error) {
	buffer := make([]byte, len(message)+crypto_box_ZEROBYTES)

	copy(buffer[0:crypto_box_ZEROBYTES], crypto_box_PADDING)
	copy(buffer[crypto_box_ZEROBYTES:], message)

	rc := C.crypto_box((*C.uchar)(unsafe.Pointer(&buffer[0])),
		(*C.uchar)(unsafe.Pointer(&buffer[0])),
		(C.ulonglong)(len(buffer)),
		(*C.uchar)(unsafe.Pointer(&nonce[0])),
		(*C.uchar)(unsafe.Pointer(&publicKey[0])),
		(*C.uchar)(unsafe.Pointer(&secretKey[0])))

	if rc == 0 {
		return buffer[crypto_box_BOXZEROBYTES:], nil
	}

	return nil, fmt.Errorf("Error encrypting message (error code %v)", rc)
}

// Wrapper function for crypto_box_open.
//
// Verifies and decrypts the ciphertext using the secretKey, publicKey and nonce. The zero padding
// required by the crypto_box C API is added internally and should not be included in the supplied
// message. Likewise the zero padding that prefixes the plaintext returned by the crypto_box C API
// is stripped from the returned plaintext.
//
// Ref. http://nacl.cr.yp.to/box.html
func CryptoBoxOpen(ciphertext, nonce, publicKey, secretKey []byte) ([]byte, error) {
	buffer := make([]byte, len(ciphertext)+crypto_box_BOXZEROBYTES)

	copy(buffer[0:crypto_box_BOXZEROBYTES], crypto_box_PADDING)
	copy(buffer[crypto_box_BOXZEROBYTES:], ciphertext)

	rc := C.crypto_box_open((*C.uchar)(unsafe.Pointer(&buffer[0])),
		(*C.uchar)(unsafe.Pointer(&buffer[0])),
		(C.ulonglong)(len(buffer)),
		(*C.uchar)(unsafe.Pointer(&nonce[0])),
		(*C.uchar)(unsafe.Pointer(&publicKey[0])),
		(*C.uchar)(unsafe.Pointer(&secretKey[0])))

	if rc == 0 {
		return buffer[crypto_box_ZEROBYTES:], nil
	}

	return nil, fmt.Errorf("Error decrypting message (error code %v)", rc)
}
