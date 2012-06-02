// Copyright 2011 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mymath
//class for creating and handling bitcoin addresses

import(
	"crypto/rand"
	"bitecdsa"
	"bitelliptic"
	"log"
)

//structure to keep different encodings of the bitcoin adress
type BitAddress struct{
	PrivateKey []byte//private key, can be encrypted
	Encrypted bool//checked if the private key is encrypted
	PublicKey []byte//public key
	Hash160 []byte //RIPEMD hash
	Hash []byte //RIPEMD with network bit and checksum
	Base Base58//base58 encoded address
}

func NewRandomAddress() (*BitAddress, error){
	curve:=bitelliptic.S256()
	
	priv, err:=bitecdsa.GenerateKey(curve, rand.Reader)
	
	if err!=nil{
		return nil, err
	}
	ba:=new(BitAddress)
	ba.PrivateKey=Big2Hex(priv.D)
	ba.Encrypted=false
	ba.PublicKey=append(append([]byte{0x04}, Big2Hex(priv.PublicKey.X)...), Big2Hex(priv.PublicKey.Y)...)
	ba.Hash160=SHARipemd(ba.PublicKey)
	ba.Hash=append(append([]byte{0x00}, ba.Hash160...), DoubleSHA(append([]byte{0x00}, ba.Hash160...))[0:4]...)
	ba.Base=Hex2Base58(ba.Hash)
	
	return ba, nil
}

func NewRandomAddressOtherNets(netByte byte) (*BitAddress, error){
	curve:=bitelliptic.S256()
	
	priv, err:=bitecdsa.GenerateKey(curve, rand.Reader)
	
	if err!=nil{
		return nil, err
	}
	ba:=new(BitAddress)
	ba.PrivateKey=Big2Hex(priv.D)
	ba.Encrypted=false
	ba.PublicKey=append(append([]byte{netByte}, Big2Hex(priv.PublicKey.X)...), Big2Hex(priv.PublicKey.Y)...)
	ba.Hash160=SHARipemd(ba.PublicKey)
	ba.Hash=append(append([]byte{0x00}, ba.Hash160...), DoubleSHA(append([]byte{0x00}, ba.Hash160...))[0:4]...)
	ba.Base=Hex2Base58(ba.Hash)
	
	return ba, nil
}

func NewAddressFromPrivateKey(privateKey []byte) (*BitAddress, error){
	curve:=bitelliptic.S256()
	
	priv, err:=bitecdsa.GenerateFromPrivateKey(Hex2Big(privateKey), curve)
	
	if err!=nil{
		return nil, err
	}
	ba:=new(BitAddress)
	ba.PrivateKey=Big2Hex(priv.D)
	ba.Encrypted=false
	ba.PublicKey=append(append([]byte{0x04}, Big2Hex(priv.PublicKey.X)...), Big2Hex(priv.PublicKey.Y)...)
	ba.Hash160=SHARipemd(ba.PublicKey)
	ba.Hash=append(append([]byte{0x00}, ba.Hash160...), DoubleSHA(append([]byte{0x00}, ba.Hash160...))[0:4]...)
	ba.Base=Hex2Base58(ba.Hash)
	
	return ba, nil
}
func NewAddressFromPrivateKeyOtherNets(netByte byte, privateKey []byte) (*BitAddress, error){
	curve:=bitelliptic.S256()
	
	priv, err:=bitecdsa.GenerateFromPrivateKey(Hex2Big(privateKey), curve)
	
	if err!=nil{
		return nil, err
	}
	ba:=new(BitAddress)
	ba.PrivateKey=Big2Hex(priv.D)
	ba.Encrypted=false
	ba.PublicKey=append(append([]byte{netByte}, Big2Hex(priv.PublicKey.X)...), Big2Hex(priv.PublicKey.Y)...)
	ba.Hash160=SHARipemd(ba.PublicKey)
	ba.Hash=append(append([]byte{0x00}, ba.Hash160...), DoubleSHA(append([]byte{0x00}, ba.Hash160...))[0:4]...)
	ba.Base=Hex2Base58(ba.Hash)
	
	return ba, nil
}


























//creates the structure to hold the base58 address
func NewFromBase(b Base58) *BitAddress{
	ba:=new(BitAddress)
	ba.Base=b
	ba.Hash=b.BitHex()
	ba.Hash160=ba.Hash[1:21]
	return ba
}

func NewFromBaseString(s string) *BitAddress{
	return NewFromBase(String2Base58(s))
}

//TODO: do something for TESTnet?
//TODO: update and revise
//creates a bitcoin address from the public key
func NewFromPublicKey(netbyte byte, key []byte) *BitAddress{
	ba:=new(BitAddress)
	
	ba.PublicKey=key//store the public key
	
	ba.Hash=make([]byte, 25)
	ba.Hash[0]=netbyte
	
	tmp:=SHARipemd(key)
	
	ba.Hash160=tmp//stores RIPEMD hash
	
	for i:=0;i<len(tmp);i++{
		ba.Hash[i+1]=tmp[i]//copy RIPEMD hash to extended hash
	}
	tmp=DoubleSHA(ba.Hash[0:21])
	
	//hash checksum
	ba.Hash[21]=tmp[0]
	ba.Hash[22]=tmp[1]
	ba.Hash[23]=tmp[2]
	ba.Hash[24]=tmp[3]
	
	//encodes the address into base 58
	ba.Base=Hex2Base58(ba.Hash)
	
	return ba//return
}

//creates a new bitcoin address from a string representing a public key
func NewFromPublicKeyString(netByte byte, s string) *BitAddress{
	return NewFromPublicKey(netByte, String2Hex(s))
}


//TODO: do
//check validity of bitcoin address
func (ba BitAddress)CheckValidity() bool{
	tmp:=DoubleSHA(ba.Hash[0:21])
	if ba.Hash[21]==tmp[0] && ba.Hash[22]==tmp[1]&& ba.Hash[23]==tmp[2] && ba.Hash[24]==tmp[3] {
		return true
	}
	return false
}

func CheckAddressValidity(addr Base58)bool{
	return NewFromBase(addr).CheckValidity()
}

func CheckAddressStringValidity(addr string)bool{
	return CheckAddressValidity(String2Base58(addr))
}

func TestBitAddress(){
	addr, _:=NewAddressFromPrivateKey(Str2Hex("18E14A7B6A307F426A94F8114701E7C8E774E7F9A47E2C2035DB29A206321725"))
	log.Printf("Address - %X, %s", addr, addr.Base)
}