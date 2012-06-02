// Copyright 2011 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mymath
//Subpackage used for all Bitcoin-specific math calculations, operations and conversions

import(
	"log"
	"bytes"
	
    "math/big"
    "crypto/elliptic"
    "strconv"
    "BitSHA"
)
//Bitcoin-related math

//TODO: test and add to tests
func CalculateBitMidstate(data []byte)[]byte{
    var h BitSHA.Hash = BitSHA.New()
	h.Write(RevWords2(data[0:64]))
	return RevWords2(h.Midstate())
}

//TODO: test and add to tests
func Bits2Target(bits uint32) *big.Int{
	//log.Print("Bits2Target")
	answer:=big.NewInt(int64(bits%0x01000000))
	answer.Mul(answer, big.NewInt(2).Exp(big.NewInt(2), big.NewInt(8*(int64(bits/0x01000000)-3)), nil))
	
	return answer
}

//TODO: test and add to tests
func Bits2TargetHex(bits uint32) []byte{
	//log.Print("Bits2TargetHex")
	bigTarget:=Bits2Target(bits)
	
	array:=bigTarget.Bytes()
	answer:=make([]byte, 32)
	copy(answer[(32-len(array)):], array)
	return answer
}

//TODO: test and add to tests
func Bits2TargetHexRev(bits uint32) []byte{
	//log.Print("Bits2TargetHexRev")
	bigTarget:=Bits2Target(bits)
	
	array:=bigTarget.Bytes()
	answer:=make([]byte, 32)
	copy(answer[(32-len(array)):], array)
	answerrev:=make([]byte, 32)
	for i:=0;i<len(answer);i++{
		answerrev[i]=answer[len(answer)-1-i]
	}
	return answerrev
}

//TODO: test and add to tests
func BitsRev2TargetHexRev(bits uint32) []byte{
	//log.Print("BitsRev2TargetHexRev")
	var newbits uint32
	newbits=bits/0x01000000
	newbits+=((bits/0x010000)%0x0100)*0x0100
	newbits+=((bits/0x0100)%0x0100)*0x010000
	newbits+=(bits%0x0100)*0x01000000
	return Bits2TargetHexRev(newbits)
}

//TODO: test and add to tests
func Target2Bits(target *big.Int) uint32{
	tmp:=Big2Hex(target)
	
	properLen:=len(tmp)
	var bits uint32
	bits=0
	bits+=0x010000*uint32(tmp[0])
	bits+=0x0100*uint32(tmp[1])
	bits+=uint32(tmp[2])
	
	if(tmp[0]>0x80){
		properLen++
		bits/=0x0100
	}
	
	bits+=0x01000000*uint32(properLen)
	return bits
}

//TODO: test and add to tests
func Bits2Difficulty(bits uint32) (float64, error){
	return Target2Difficulty(Bits2Target(bits))
}

//TODO: test and add to tests
func Target2Difficulty(target *big.Int) (float64, error){
	a, err:=strconv.ParseFloat("26959535291011309493156476344723991336010898738574164086137773096960", 64)//decimal target equivalent of difficulty 1
	if err!=nil{
		return 0.0, err
	}
	b, err:=strconv.ParseFloat(target.String(), 64)
	if err!=nil{
		return 0.0, err
	}
	return a/b, nil
}

func BitsString2TargetHexRev(bits string) []byte{
	return Bits2TargetHexRev(Str2Uint32(bits))
}



//TODO: test and add to tests
func Makesecp256k1(){

	var p256 = new(elliptic.CurveParams)
	 //p=FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2
	 //a=00000000 00000000 00000000 00000000 00000000 00000000 0000000000000000
	 //b=00000000 00000000 00000000 00000000 00000000 00000000 0000000000000007
	 //G, compressed = 02 79BE667E F9DCBBAC 55A06295 CE870B07 029BFCDB 2DCE28D9 59F2815B 16F81798
	 //G, uncompressed = 04 79BE667E F9DCBBAC 55A06295 CE870B07 029BFCDB 2DCE28D9 59F2815B 16F81798 483ADA77 26A3C465 5DA4FBFC 0E1108A8 FD17B448 A6855419 9C47D08F FB10D4B8
    //n=FFFFFFFF FFFFFFFF FFFFFFFF FFFFFFFE BAAEDCE6 AF48A03B BFD25E8C D0364141
    //h= 01
    
    
    p256.P, _ = new(big.Int).SetString("115792089210356248762697446949407573530086143415290314195533631308867097853951", 10)
	p256.N, _ = new(big.Int).SetString("115792089210356248762697446949407573529996955224135760342422259061068512044369", 10)
	p256.B, _ = new(big.Int).SetString("5ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b", 16)
	p256.Gx, _ = new(big.Int).SetString("6b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c296", 16)
	p256.Gy, _ = new(big.Int).SetString("4fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f5", 16)
	p256.BitSize = 256

}

type VarInt uint64

func (vi VarInt)Len()int{
	if vi<0xfd{
		return 1
	}
	if vi<=0xffff{
		return 3
	}
	if vi<=0xffffffff{
		return 5
	}
	return 9
}

func (vi VarInt)Compile()[]byte{
	return VarInt2HexRev(vi)
}

//TODO: test
func DecodeVarInt(b []byte)VarInt{
	if len(b)==0{
		return 0
	}
	if b[0]<0xFD{
		return VarInt(b[0])
	}
	if b[0]==0xFD{
		return VarInt(HexRev2Uint64(b[1:3]))
	}
	if b[0]==0xFE{
		return VarInt(HexRev2Uint64(b[1:5]))
	}
	if b[0]==0xFF{
		return VarInt(HexRev2Uint64(b[1:9]))
	}
	return 0
}
//TODO: check and add to tests
func DecodeVarIntGiveRest(b []byte) (VarInt, []byte){
	if len(b)==0{
		return 0, b
	}
	if b[0]<0xFD{
		return VarInt(b[0]), b[1:]
	}
	if b[0]==0xFD{
		return VarInt(HexRev2Uint64(b[1:3])), b[3:]
	}
	if b[0]==0xFE{
		return VarInt(HexRev2Uint64(b[1:5])), b[5:]
	}
	if b[0]==0xFF{
		return VarInt(HexRev2Uint64(b[1:9])), b[9:]
	}
	return 0, b
}

//variable-length hex result
//https://en.bitcoin.it/wiki/Protocol_specification#Variable_length_integer
func VarInt2Hex(vi VarInt) []byte{
	answer:=make([]byte, 9)
	copy(answer[1:], Uint642Hex(uint64(vi)))
	if vi<0xfd {
		return answer[8:]
	}
	if vi<=0xffff {
		answer[6]=0xfd
		return answer[6:]
	}
	if vi<=0xffffffff {
		answer[4]=0xfe
		return answer[4:]
	}
	answer[0]=0xff
	return answer
}

//variable-length hex result
//https://en.bitcoin.it/wiki/Protocol_specification#Variable_length_integer
func VarInt2HexRev(vi VarInt) []byte{
	answer:=make([]byte, 9)
	copy(answer[1:], Uint642HexRev(uint64(vi)))
	
	if vi<0xfd {
		return answer[1:2]
	}
	if vi<=0xffff {
		answer[0]=0xfd
		return answer[0:3]
	}
	if vi<=0xffffffff {
		answer[0]=0xfe
		return answer[0:5]
	}
	answer[0]=0xff
	return answer
}

type VarStr struct{
	length VarInt
	str []byte
}

func (vs *VarStr)Len()int{
	return vs.length.Len()+len(vs.str)
}

func (vs *VarStr)Set(newString string){
	vs.str=make([]byte, len(newString))
	
	for i:=0;i<len(newString);i++{
		vs.str[i]=newString[i]
	}
	vs.length=VarInt(len(newString))
}

func (vs *VarStr)Compile()[]byte{
	answer:=make([]byte, vs.Len())
	copy(answer[:], vs.length.Compile())
	copy(answer[vs.length.Len():], vs.str)
	return answer
}

func (vs *VarStr)Read()string{
	return string(vs.str)
}
//TODO: test and add to tests
func GenerateMerkleTreeFromString(leafs []string)[]string{
	tmp:=make([][]byte, len(leafs))
	
	for i:=0;i<len(tmp);i++{
		tmp[i]=make([]byte, len(leafs[i])/2)
		copy(tmp[i][:], Str2Hex(leafs[i]))
	}
	merkletree:=GenerateMerkleTree(tmp)
	
	answer:=make([]string, len(merkletree))
	for i:=0;i<len(merkletree);i++{
		answer[i]=Hex2Str(merkletree[i])
	}
	return answer
}

//TODO: test and add to tests
func GenerateMerkleTree(leafs [][]byte)[][]byte{
	answer:=make([][]byte, len(leafs))
	for i:=0;i<len(answer);i++{
		answer[i]=make([]byte, len(leafs[i]))
		copy(answer[i][:], leafs[i][:])
	}
	level:=make([][]byte, len(leafs))
	for i:=0;i<len(answer);i++{
		level[i]=make([]byte, len(leafs[i]))
		copy(level[i][:], leafs[i][:])
	}
	//answer:=leafs[:]
	//level:=leafs[:]
	for;len(level)>1;{
		//currentlevel:=level[:int(math.Ceil(float64(len(level)/2.0)))]
		//currentlevel:=make([][]byte, int(math.Ceil(float64(1.0*len(level)/2.0))))
		currentlevel:=make([][]byte, len(level)/2+len(level)%2)
		log.Printf("len - %d", len(currentlevel))
		for i:=0;i<len(currentlevel)-1;i++{
			currentlevel[i]=DoubleSHAPair(level[2*i], level[2*i+1])
		}
		if(len(level)%2==1){
			currentlevel[len(currentlevel)-1]=DoubleSHAPair(level[2*(len(currentlevel)-1)], level[2*(len(currentlevel)-1)])
		}else{
			currentlevel[len(currentlevel)-1]=DoubleSHAPair(level[2*(len(currentlevel)-1)], level[2*(len(currentlevel)-1)+1])
		}
		//answer=append(answer, currentlevel)
		tmp:=make([][]byte, len(answer)+len(currentlevel))
		for i:=0;i<len(answer);i++{
			tmp[i]=make([]byte, len(answer[i]))
			copy(tmp[i][:], answer[i][:])
		}
		level=make([][]byte, len(currentlevel))
		for i:=0;i<len(currentlevel);i++{
			tmp[len(answer)+i]=make([]byte, len(currentlevel[i]))
			copy(tmp[len(answer)+i][:], currentlevel[i][:])
			level[i]=make([]byte, len(currentlevel[i]))
			copy(level[i][:], currentlevel[i][:])
		}
		//copy(tmp[:], answer)
		//copy(tmp[len(answer):], currentlevel)
		answer=tmp
		//level=currentlevel
	}
	return answer
}

/*
func GenerateMerkleTree(leafs [][]byte)[][]byte{
	answer:=make([][]byte, len(leafs))
	for i:=0;i<len(answer);i++{
		answer[i]=make([]byte, len(leafs[i]))
		copy(answer[i][:], leafs[i][:])
	}
	level:=make([][]byte, len(leafs))
	for i:=0;i<len(answer);i++{
		level[i]=make([]byte, len(leafs[i]))
		copy(level[i][:], leafs[i][:])
	}
	//answer:=leafs[:]
	//level:=leafs[:]
	for;len(level)>1;{
		//currentlevel:=level[:int(math.Ceil(float64(len(level)/2.0)))]
		//currentlevel:=make([][]byte, int(math.Ceil(float64(1.0*len(level)/2.0))))
		currentlevel:=make([][]byte, len(level)/2+len(level)%2)
		log.Printf("len - %d", len(currentlevel))
		for i:=0;i<len(currentlevel)-1;i++{
			currentlevel[i]=DoubleSHAPairRev(level[2*i], level[2*i+1])
		}
		if(len(level)%2==1){
			currentlevel[len(currentlevel)-1]=DoubleSHAPairRev(level[2*(len(currentlevel)-1)], level[2*(len(currentlevel)-1)])
		}else{
			currentlevel[len(currentlevel)-1]=DoubleSHAPairRev(level[2*(len(currentlevel)-1)], level[2*(len(currentlevel)-1)+1])
		}
		//answer=append(answer, currentlevel)
		tmp:=make([][]byte, len(answer)+len(currentlevel))
		for i:=0;i<len(answer);i++{
			tmp[i]=make([]byte, len(answer[i]))
			copy(tmp[i][:], answer[i][:])
		}
		level=make([][]byte, len(currentlevel))
		for i:=0;i<len(currentlevel);i++{
			tmp[len(answer)+i]=make([]byte, len(currentlevel[i]))
			copy(tmp[len(answer)+i][:], currentlevel[i][:])
			level[i]=make([]byte, len(currentlevel[i]))
			copy(level[i][:], currentlevel[i][:])
		}
		//copy(tmp[:], answer)
		//copy(tmp[len(answer):], currentlevel)
		answer=tmp
		//level=currentlevel
	}
	return answer
}

*/

//TODO: test and add to tests
func GenerateMerkleRoot(leafs [][]byte) []byte{
	tree:=GenerateMerkleTree(leafs)
	return tree[len(tree)-1]
}

func TestGenerateMerkleTree()bool{
	/* //http://blockexplorer.com/rawblock/0000000000000ae4e079775dfb20e0571ca5e24d7fc73489a8a48bb9758e17d6
	"b30c84a7e9a68873e40f807557b9af8d364e66f0cbfaa64e684a5ed662e5a197",
    "74182c272bd3b201a6afc1fcaf5fa60c603dbd52f2d3150ce7aa593d697fd01f",
    "48ba51a6a9a7ede4618e5cdeb8cd3a000ca9e369df96a2bdd9cd36563781f072",
    "665960a9ee63b2ab7c27eaff6f8e1cccbb35ef8b8a20a82f3bfd6ce3445057e6",
    "2ba0808f3cbbf21b45916585bc8b2c93c131abb79425fae92e3e0cb6968808ca",
    "58e13ec3fc70f745a2e6345f799f030ee6f26499eb4a59959b92155857325241",
    "96ec719655fa0a59a1570eab9172efb67c8de29d091486a10f0a1514dadc5e3a",
    "ee96e27544effdb27f19d90b1938b2e60c84f4091521a4082b0d91f16f27fdae",
    "92ea9d03ae8e51874030c69958734607246d645c1f047b11edf6be2d595e80b9",
    "6250804544c2010a86b147f7e00f1943794396f009ba2239728497c74cc2aeb8",
    
    "c84f67766b3cccca370a317d44cd042861629434f127715065e4a35d6cda32e9",
    "533638ea7a72c2ef1cf43cbd56861c7329c064f2830e380f04a686cd5a91db28",
    "af7aac7595207ad1cd5e80d1cb2c875dab6674b363cea87482912d514345dad4",
    "a3c61f0642c833f9a804bc92f824feed1927e903103df2d34a8d3987d52b2838",
    "dcff5d86211ea488879b9d495126371fa0eca322d8e14cbf09a46bf2ce7ba045",
    
    "25ee9f74d7e462e3432981ba231fbb630944575704dbfbf61a03f21371f04a80",
    "ea7894cfeb08bcd7a7de6d731cf036087a335221b0b4e3bf123252ffaff4b5f8",
    "e2d0b6819628e8f64d2e6ac632285c7cad9d5f2d78524317a108f26b66aebb0d",
    
    "fd918efdc92eac85717d23083c5d749d5cd2c43c153cb11d6299100fb743788d",
    "f47bf9fe04c7051d09461f0930c695413c51f2e9831e8c5efa6700e2f22161ab",
    
    "d6402a834394e961a8d07c1cd4d949badeefd371230f0f49fda9cb96158dab92"
    */
	merkletree:=make([][]byte, 10)
	merkletree[0]=String2Hex("b30c84a7e9a68873e40f807557b9af8d364e66f0cbfaa64e684a5ed662e5a197")
	merkletree[1]=String2Hex("74182c272bd3b201a6afc1fcaf5fa60c603dbd52f2d3150ce7aa593d697fd01f")
	merkletree[2]=String2Hex("48ba51a6a9a7ede4618e5cdeb8cd3a000ca9e369df96a2bdd9cd36563781f072")
	merkletree[3]=String2Hex("665960a9ee63b2ab7c27eaff6f8e1cccbb35ef8b8a20a82f3bfd6ce3445057e6")
	merkletree[4]=String2Hex("2ba0808f3cbbf21b45916585bc8b2c93c131abb79425fae92e3e0cb6968808ca")
	merkletree[5]=String2Hex("58e13ec3fc70f745a2e6345f799f030ee6f26499eb4a59959b92155857325241")
	merkletree[6]=String2Hex("96ec719655fa0a59a1570eab9172efb67c8de29d091486a10f0a1514dadc5e3a")
	merkletree[7]=String2Hex("ee96e27544effdb27f19d90b1938b2e60c84f4091521a4082b0d91f16f27fdae")
	merkletree[8]=String2Hex("92ea9d03ae8e51874030c69958734607246d645c1f047b11edf6be2d595e80b9")
	merkletree[9]=String2Hex("6250804544c2010a86b147f7e00f1943794396f009ba2239728497c74cc2aeb8")
	
	//log.Printf("merkletree - %dx%d", len(merkletree), len(merkletree[0]))
	//log.Printf("b30c84a7e9a68873e40f807557b9af8d364e66f0cbfaa64e684a5ed662e5a197 = %X", String2Hex("b30c84a7e9a68873e40f807557b9af8d364e66f0cbfaa64e684a5ed662e5a197"))
	
	/*one:=make([]byte, 2)
	two:=make([]byte, 2)
	one[0]=0x00
	one[1]=0x01
	two[0]=0x02
	two[1]=0x03
	
	log.Printf("%X", append(one[:], two[:]))
	
	three:=[]byte{0, 1}
	four:=[]byte{2, 3}
	
	five:=append(three, four)*/
	
	
	
	//log.Printf("aaaa+bbbb=%X", append(String2Hex("aaaa"), String2Hex("bbbb")))
	
	/*for i:=0;i<len(merkletree);i++{
		log.Printf("merkletree %d - %X", i, merkletree[i])
	}*/
	answer:=GenerateMerkleTree(merkletree)
	
	if(len(answer)==21){
		return true
	}
	
	/*for i:=0;i<len(merkletree);i++{
		log.Printf("merkletree %d - %X", i, merkletree[i])
	}
	for i:=0;i<len(answer);i++{
		log.Printf("answer %d - %X", i, answer[i])
	}*/
	return false
}




//Testing

func TestEverythingBitmath()bool{
	if(TestVarInt()==false){
		log.Print("TestVarInt()==false")
		return false
	}
	if(TestCompile()==false){
		log.Print("TestCompile()==false")
		return false
	}
	if(TestLen()==false){
		log.Print("TestLen()==false")
		return false
	}
	if(TestVarInt2Hex()==false){
		log.Print("TestVarInt2Hex()==false")
		return false
	}
	if(TestVarInt2HexRev()==false){
		log.Print("TestVarInt2HexRev()==false")
		return false
	}
	if(TestVarStr()==false){
		log.Print("TestVarStr()==false")
		return false
	}
	//log.Print("All Bitmath tests okay!")
	TestGenerateMerkleTree()
	return true
}



func TestVarInt()bool{
	var vi1 VarInt
	vi1=0
	
	var vi2 VarInt
	vi2=1
	
	var vi3 VarInt
	vi3=0xffffffff
	
	if(vi1!=0){
		return false
	}
	if(vi2!=1){
		return false
	}
	if(vi3!=0xffffffff){
		return false
	}
	
	return true
}

func TestCompile()bool{
	if(bytes.Compare(VarInt(0).Compile(), VarInt2HexRev(VarInt(0)))!=0){
		log.Printf("%X - %X", VarInt(0).Compile(), VarInt2Hex(VarInt(0)))
		return false
	}
	if(bytes.Compare(VarInt(1).Compile(), VarInt2HexRev(VarInt(1)))!=0){
		log.Printf("%X - %X", VarInt(1).Compile(), VarInt2Hex(VarInt(1)))
		return false
	}
	if(bytes.Compare(VarInt(0x010000).Compile(), VarInt2HexRev(VarInt(0x010000)))!=0){
		log.Printf("%X - %X", VarInt(0x010000).Compile(), VarInt2Hex(VarInt(0x010000)))
		return false
	}
	if(bytes.Compare(VarInt(0xffffffff).Compile(), VarInt2HexRev(VarInt(0xffffffff)))!=0){
		log.Printf("%X - %X", VarInt(0xffffffff).Compile(), VarInt2Hex(VarInt(0xffffffff)))
		return false
	}
	return true
}

func TestLen()bool{
	if(VarInt(0).Len()!=1){
		log.Print("TestLen(0)==false")
		return false
	}
	if(VarInt(1).Len()!=1){
		log.Print("TestLen(1)==false")
		return false
	}
	if(VarInt(0xFC).Len()!=1){
		log.Print("TestLen(fc)==false")
		return false
	}
	if(VarInt(0xFD).Len()!=3){
		log.Print("TestLen(fd)==false")
		return false
	}
	if(VarInt(0xFFFE).Len()!=3){
		log.Print("TestLen(fffe)==false")
		return false
	}
	if(VarInt(0xFFFF).Len()!=3){
		log.Print("TestLen(ffff)==false")
		return false
	}
	if(VarInt(0x010000).Len()!=5){
		log.Print("TestLen(0x010000)==false")
		return false
	}
	if(VarInt(0xFFFFFFFE).Len()!=5){
		log.Print("TestLen(fffffffe)==false")
		return false
	}
	if(VarInt(0xFFFFFFFF).Len()!=5){
		log.Print("TestLen(ffffffff)==false")
		return false
	}
	if(VarInt(0x0100000000).Len()!=9){
		log.Print("TestLen(0x0100000000)==false")
		return false
	}
	return true
}

func TestVarInt2Hex() bool{
	var byte1 []byte
	var byte3 []byte
	var byte5 []byte
	var byte9 []byte
	
	byte1=make([]byte, 1)
	byte3=make([]byte, 3)
	byte5=make([]byte, 5)
	byte9=make([]byte, 9)
	
	byte1[0]=0
	byte3[0]=0
	byte5[0]=0
	byte9[0]=0
	
	if(bytes.Compare(byte1, VarInt2Hex(VarInt(0)))!=0){
		log.Print("0")
		return false
	}
	byte1[0]=0xF0
	if(bytes.Compare(byte1, VarInt2Hex(VarInt(0xF0)))!=0){
		log.Print("F0")
		return false
	}
	byte1[0]=0xFC
	if(bytes.Compare(byte1, VarInt2Hex(VarInt(0xFC)))!=0){
		log.Print("FC")
		return false
	}
	byte3[0]=0xFD
	byte3[1]=0
	byte3[2]=0xFD
	if(bytes.Compare(byte3, VarInt2Hex(VarInt(0xFD)))!=0){
		log.Printf("FD - %X", VarInt2Hex(VarInt(0xFD)))
		return false
	}
	byte3[2]=0xFF
	if(bytes.Compare(byte3, VarInt2Hex(VarInt(0xFF)))!=0){
		log.Print("FF")
		return false
	}
	byte3[1]=0xFF
	byte3[2]=0xFE
	if(bytes.Compare(byte3, VarInt2Hex(VarInt(0xFFFE)))!=0){
		log.Print("FFFE")
		return false
	}
	byte3[2]=0xFF
	if(bytes.Compare(byte3, VarInt2Hex(VarInt(0xFFFF)))!=0){
		log.Print("FFFF")
		return false
	}
	
	byte5[0]=0xFE
	byte5[1]=0
	byte5[2]=0x01
	byte5[3]=0
	byte5[4]=0
	if(bytes.Compare(byte5, VarInt2Hex(VarInt(0x010000)))!=0){
		log.Print("0x010000")
		return false
	}
	byte5[1]=0xFF
	byte5[2]=0xFF
	byte5[3]=0xFF
	byte5[4]=0xFE
	if(bytes.Compare(byte5, VarInt2Hex(VarInt(0xFFFFFFFE)))!=0){
		log.Print("0xFFFFFFFE")
		return false
	}
	byte5[4]=0xFF
	if(bytes.Compare(byte5, VarInt2Hex(VarInt(0xFFFFFFFF)))!=0){
		log.Print("0xFFFFFFFF")
		return false
	}
	
	byte9[0]=0xFF
	byte9[1]=0
	byte9[2]=0
	byte9[3]=0
	byte9[4]=0x01
	byte9[5]=0
	byte9[6]=0
	byte9[7]=0
	byte9[8]=0
	
	if(bytes.Compare(byte9, VarInt2Hex(VarInt(0x0100000000)))!=0){
		log.Print("0x0100000000")
		return false
	}
	
	return true
}

func TestVarInt2HexRev() bool{
	var byte1 []byte
	var byte3 []byte
	var byte5 []byte
	var byte9 []byte
	
	byte1=make([]byte, 1)
	byte3=make([]byte, 3)
	byte5=make([]byte, 5)
	byte9=make([]byte, 9)
	
	byte1[0]=0
	byte3[0]=0
	byte5[0]=0
	byte9[0]=0
	
	if(bytes.Compare(byte1, VarInt2HexRev(VarInt(0)))!=0){
		log.Print("0")
		return false
	}
	byte1[0]=0xF0
	if(bytes.Compare(byte1, VarInt2HexRev(VarInt(0xF0)))!=0){
		log.Print("F0")
		return false
	}
	byte1[0]=0xFC
	if(bytes.Compare(byte1, VarInt2HexRev(VarInt(0xFC)))!=0){
		log.Print("FC")
		return false
	}
	byte3[0]=0xFD
	byte3[2]=0
	byte3[1]=0xFD
	if(bytes.Compare(byte3, VarInt2HexRev(VarInt(0xFD)))!=0){
		log.Printf("FD - %X", VarInt2HexRev(VarInt(0xFD)))
		return false
	}
	byte3[1]=0xFF
	if(bytes.Compare(byte3, VarInt2HexRev(VarInt(0xFF)))!=0){
		log.Print("FF")
		return false
	}
	byte3[2]=0xFF
	byte3[1]=0xFE
	if(bytes.Compare(byte3, VarInt2HexRev(VarInt(0xFFFE)))!=0){
		log.Print("FFFE")
		return false
	}
	byte3[1]=0xFF
	if(bytes.Compare(byte3, VarInt2HexRev(VarInt(0xFFFF)))!=0){
		log.Print("FFFF")
		return false
	}
	
	byte5[0]=0xFE
	byte5[4]=0
	byte5[3]=0x01
	byte5[2]=0
	byte5[1]=0
	if(bytes.Compare(byte5, VarInt2HexRev(VarInt(0x010000)))!=0){
		log.Print("0x010000")
		return false
	}
	byte5[4]=0xFF
	byte5[3]=0xFF
	byte5[2]=0xFF
	byte5[1]=0xFE
	if(bytes.Compare(byte5, VarInt2HexRev(VarInt(0xFFFFFFFE)))!=0){
		log.Print("0xFFFFFFFE")
		return false
	}
	byte5[1]=0xFF
	if(bytes.Compare(byte5, VarInt2HexRev(VarInt(0xFFFFFFFF)))!=0){
		log.Print("0xFFFFFFFF")
		return false
	}
	
	byte9[0]=0xFF
	byte9[8]=0
	byte9[7]=0
	byte9[6]=0
	byte9[5]=0x01
	byte9[4]=0
	byte9[3]=0
	byte9[2]=0
	byte9[1]=0
	
	if(bytes.Compare(byte9, VarInt2HexRev(VarInt(0x0100000000)))!=0){
		log.Print("0x0100000000")
		return false
	}
	return true
}

func TestVarStr()bool{
	vs1:=new(VarStr)
	vs1.Set("Hello!")
	
	if(vs1.Len()!=7){
		log.Print("s1.Len()")
		return false
	}
	if(vs1.Read()!="Hello!"){
		log.Print("s1.Read()")
		return false
	}
	
	compiled:=make([]byte, 7)
	compiled[0]=0x06
	compiled[1]=0x48
	compiled[2]=0x65
	compiled[3]=0x6C
	compiled[4]=0x6C
	compiled[5]=0x6F
	compiled[6]=0x21
	if(bytes.Compare(compiled, vs1.Compile())!=0){
		log.Print("vs1.Compile()")
		return false
	}
	
	return true
}

//TODO: expand and add to tests
func TestBitsTargetDifficultyConversions() bool{
	//1d00ffff
	testSucceeded:=true
	tmp:=big.NewInt(0)
	if AreStringsEqual(Hex2Str(Big2Hex(Bits2Target(0x1d00ffff))), "FFFF0000000000000000000000000000000000000000000000000000")!=true{
		log.Printf("%X", Big2Hex(Bits2Target(0x1d00ffff)))
		testSucceeded=false
	}
	ans, err:=Bits2Difficulty(0x1d00ffff)
	tmp, _=tmp.SetString("FFFF0000000000000000000000000000000000000000000000000000", 16)
	if Bits2Target(0x1d00ffff).Cmp(tmp)!=0{
		log.Printf("%v, %v", ans, err)
		testSucceeded=false
	}
	log.Printf("%X=?1d00ffff", Target2Bits(tmp))
	
	if AreStringsEqual(Hex2Str(Big2Hex(Bits2Target(0x1b0404cb))), "0404CB000000000000000000000000000000000000000000000000")!=true{
		log.Printf("%X", Big2Hex(Bits2Target(0x1b0404cb)))
		testSucceeded=false
	}
	ans, err=Bits2Difficulty(0x1b0404cb)
	tmp, _=tmp.SetString("0404CB000000000000000000000000000000000000000000000000", 16)
	if Bits2Target(0x1b0404cb).Cmp(tmp)!=0{
		log.Printf("%v, %v", ans, err)
		testSucceeded=false
	}
	
	if testSucceeded==false{
		log.Print("TestBitsTargetDifficultyConversions failed!")
	}
	log.Printf("%X=?0x1b0404cb", Target2Bits(tmp))
	
	for i:=0;i<10;i++{
		bits:=uint32(0x0000FFFF)
		bits+=uint32((i+10)*0x01000000)
		log.Printf("%X=?%X", Target2Bits(Bits2Target(bits)), bits)
	}
	
	return testSucceeded
}








func midstateStuff(){
 //log.Printf("00000000032e361b246e96c4e594523b6ff42bc0527e560f203fb20a08a86185, %X", mymath.RevWords(mymath.Str2Hex("00000000032e361b246e96c4e594523b6ff42bc0527e560f203fb20a08a86185")))
    
    
   /* var h BitSHA.Hash = BitSHA.New()
    //h.Write(mymath.Str2Hex("00000001c570c4764aadb3f09895619f549000b8b51a789e7f58ea750000709700000000103ca064f8c76c390683f8203043e91466a7fcc40e6ebc428fbcc2d8"))
    //h.Write(mymath.RevWords(mymath.Str2Hex("00000001c570c4764aadb3f09895619f549000b8b51a789e7f58ea750000709700000000103ca064f8c76c390683f8203043e91466a7fcc40e6ebc428fbcc2d89b574a864db8345b1b00b5ac00000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000080020000")))
    h.Write(mymath.Str2Hex("0100000076C470C5F0B3AD4A9F619598B80090549E781AB575EA587F977000000000000064A03C10396CC7F820F8830614E94330C4FCA76642BC6E0ED8C2BC8F"))
    //e772fc6964e7b06d8f855a6166353e48b2562de4ad037abc889294cea8ed1070
    //E772FC6964E7B06D8F855A6166353E48B2562DE4AD037ABC889294CEA8ED1070
    //69FC72E7 6DB0E764 615A858F 483E3566 E42D56B2 BC7A03AD CE949288 7010EDA8
    log.Printf("Mid - %X", h.Midstate())
    log.Printf("Copy - %X",h.Copy())
    log.Printf("Sum - %X",h.Sum())
    
    log.Printf("New midstate - %X", mymath.CalculateBitMidstate(mymath.Str2Hex("00000001c570c4764aadb3f09895619f549000b8b51a789e7f58ea750000709700000000103ca064f8c76c390683f8203043e91466a7fcc40e6ebc428fbcc2d89b574a864db8345b1b00b5ac00000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000080020000")))
    //log.Printf("Copy - %X",h.Copy())*/
    
    /*log.Printf("Reved words - %X", mymath.RevWords(mymath.Str2Hex("00000001c570c4764aadb3f09895619f549000b8b51a789e7f58ea750000709700000000103ca064f8c76c390683f8203043e91466a7fcc40e6ebc428fbcc2d8")))
    var h2 hash.Hash = sha256.New()
    h2.Write(mymath.RevWords(mymath.Str2Hex("00000001c570c4764aadb3f09895619f549000b8b51a789e7f58ea750000709700000000103ca064f8c76c390683f8203043e91466a7fcc40e6ebc428fbcc2d8")))
    //e772fc6964e7b06d8f855a6166353e48b2562de4ad037abc889294cea8ed1070
    log.Printf("Sum 2 - %X",h2.Sum())
    log.Printf("Reved words - %X", mymath.RevWords2(mymath.Str2Hex("00000001c570c4764aadb3f09895619f549000b8b51a789e7f58ea750000709700000000103ca064f8c76c390683f8203043e91466a7fcc40e6ebc428fbcc2d8")))
    var h3 hash.Hash = sha256.New()
    h3.Write(mymath.RevWords2(mymath.Str2Hex("00000001c570c4764aadb3f09895619f549000b8b51a789e7f58ea750000709700000000103ca064f8c76c390683f8203043e91466a7fcc40e6ebc428fbcc2d8")))
    //e772fc6964e7b06d8f855a6166353e48b2562de4ad037abc889294cea8ed1070
    log.Printf("Sum 2 - %X",h3.Sum())
    
    */
    
    
    
    //log.Printf("Reved words - %X", mymath.RevWords(mymath.Str2Hex("00000001c570c4764aadb3f09895619f549000b8b51a789e7f58ea750000709700000000103ca064f8c76c390683f8203043e91466a7fcc40e6ebc428fbcc2d89b574a864db8345b1b00b5ac00000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000080020000")))

}