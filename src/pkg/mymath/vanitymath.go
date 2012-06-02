package mymath

import(
	"math/big"
	"log"
	"bitelliptic"
	"bitecdsa"
)

func CombinePrivateKeys(one, two string) *big.Int{
	a:=big.NewInt(0)
	b:=big.NewInt(0)
	
	a.SetString(one, 16)
	b.SetString(two, 16)
	
	one1:=big.NewInt(1)
	
	tmp:=a.Add(a, b)
	tmp=tmp.Sub(tmp, one1)
	//mod,_:=new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 16)
	//mod:=new(big.Int).Sub(bitelliptic.S256().P, one1)
	mod:=bitelliptic.S256().N
	tmp=tmp.Mod(tmp, mod)
	answer:=tmp.Add(tmp, one1)
	//answer:=tmp
	//answer:=a.Mod(a.Add(a, b), bitelliptic.S256().P)
	
	log.Printf("Sum 1 - %X", Big2Hex(answer))
	
	return answer
}

func CombinePublicKeys(one, two string) (*big.Int, *big.Int){
	a, b:=PublicKeyToPointCoordinates(one)
	c, d:=PublicKeyToPointCoordinates(two)
	e, f:=bitelliptic.S256().Add(a, b, c, d)
	
	log.Printf("04%X%X", Big2Hex(e), Big2Hex(f))
	
	return e, f
}

func PublicKeyToPointCoordinates(pubKey string)(*big.Int, *big.Int){
	if len(pubKey)!=130{
		log.Printf("PublicKeyToPointCoordinates - len(pubKey)!=130")
		log.Printf(pubKey)
		return nil, nil
	}
	if pubKey[0]!='0' || pubKey[1]!='4'{
		log.Printf("pubKey[0]!='0' || pubKey[1]!='4'")
		return nil, nil
	}
	a:=big.NewInt(0)
	b:=big.NewInt(0)
	
	a.SetString(pubKey[2:66], 16)
	b.SetString(pubKey[66:130], 16)
	
	//log.Printf("pubKey - %s", pubKey)
	//log.Printf("X - %v, Y - %s", a.String(), b.String())
	
	return a, b
}

func PointCoordinatesToPublicKey(x, y *big.Int)string{
	one:=Big2Hex(x)
	if len(one)<32{
		tmp:=make([]byte, 32-len(one))
		one=append(tmp, one...)
	}
	two:=Big2Hex(y)
	if len(two)<32{
		tmp:=make([]byte, 32-len(two))
		two=append(tmp, two...)
	}
	return "04"+Hex2Str(one)+Hex2Str(two)
}

func CheckSolution(pubKey, solution, pattern string, netByte byte) (string, string){

	d:=big.NewInt(0)
	d.SetString(solution, 16)
	
	private, err:=bitecdsa.GenerateFromPrivateKey(d, bitelliptic.S256())
	if err!=nil{
		log.Printf("vanitymath err - %s", err)
		return "", err.Error()
	}
	a, b:=PublicKeyToPointCoordinates(pubKey)
	
	x, y:=bitelliptic.S256().Add(a, b, private.PublicKey.X, private.PublicKey.Y)
	
	ba:=NewFromPublicKeyString(netByte, PointCoordinatesToPublicKey(x, y))
	address:=string(ba.Base)
	for i:=0;i<len(pattern);i++{
		if address[i]!=pattern[i]{
			return "", "Wrong pattern"
		}
	}
	return address, ""
}