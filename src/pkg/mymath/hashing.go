package mymath
//Subpackage used for all variations of hashing algorithms


//TODO: update this file to the other reposiroty? Added a new function, but don't know if it is useful
import(
    "hash"
    "crypto/sha256"
    "crypto/sha1"
    "ripemd160"
    //"code.google.com/p/go.crypto/ripemd160"
)

//TODO: test
func ReverseByteOrder(b []byte)([]byte){
	//reversing the order of bits 4 at a time
	//https://bitcointalk.org/index.php?topic=51281.msg742754#msg742754
	if len(b)%4!=0{
		return nil
	}
	answer:=make([]byte, len(b))
	for i:=0;i<len(b)/4;i++{
		answer[i*4+0]=b[i*4+3]
		answer[i*4+1]=b[i*4+2]
		answer[i*4+2]=b[i*4+1]
		answer[i*4+3]=b[i*4+0]
	}
	return answer
}

//TODO: test and add to tests
//double SHA-256 hashing of a single byte array
func DoubleSHA(b []byte)([]byte){
   var h hash.Hash = sha256.New()
   h.Write(b)
   var h2 hash.Hash = sha256.New()
   h2.Write(h.Sum(nil))
   
   return h2.Sum(nil)
}

//TODO: test and add to tests
//reverse double SHA-256 hashing of a single byte array
func DoubleSHARev(b []byte)([]byte){   
   return Rev(DoubleSHA(Rev(b)))
}

//TODO: test and add to tests
//Single SHA-256 hashing of a single byte array
func SingleSHA(b []byte)([]byte){
   var h hash.Hash = sha256.New()
   h.Write(b)
   
   return h.Sum(nil)
}

//TODO: test and add to tests
//Reversed single SHA-256 hashing of a single byte array
func SingleSHARev(b []byte)([]byte){   
   return Rev(SingleSHA(Rev(b)))
}

//TODO: test and add to tests
//Single SHA-1 hashing of a single byte array
func SingleSHA1(b []byte)([]byte){
   var h hash.Hash = sha1.New()//TODO: double check
   h.Write(b)
   
   return h.Sum(nil)
}

//Reversed SHA-1 hashing of a single byte array
//TODO: test and add to tests
func SingleSHA1Rev(b []byte)([]byte){   
   return Rev(SingleSHA1(Rev(b)))//TODO: double check
}

//TODO: test and add to tests
//double hash input bytes, double hash their concatanation
func DoubleDoubleSHA(a []byte, b []byte)([]byte){
	var tmp []byte
	
	//hash first input
    var h1 hash.Hash = sha256.New()
    h1.Write(a)
   
   //hash the first input the second time
    tmp=h1.Sum(nil)
    h1=sha256.New()
    h1.Write(tmp)
    
    //hash the second input
    var h2 hash.Hash = sha256.New()
    h2.Write(b)
   
	//hash the second input the second time
    tmp=h2.Sum(nil)
    h2=sha256.New()
    h2.Write(tmp)
    
    //hash the concatenation of the double hashes of both inputs
    var answer hash.Hash=sha256.New()
    answer.Write(append(h1.Sum(nil), h2.Sum(nil)...))
    
    //double hash the concatenation
    tmp=answer.Sum(nil)
    answer=sha256.New()
    answer.Write(tmp)
   
    return answer.Sum(nil)//return result
}

//TODO: test and add to tests
//double SHA-256 hash of a concatenation of the two inputs
func DoubleSHAPair(a []byte, b []byte)([]byte){

	var tmp []byte
   
    //hash the concatenation of the two inputs
    var answer hash.Hash=sha256.New()
    answer.Write(append(a, b...))
    
    //hash it the second time
    tmp=answer.Sum(nil)
    answer=sha256.New()
    answer.Write(tmp)
   
    return answer.Sum(nil)//return
}

//TODO: test and add to tests
//double SHA-256 hash of a concatenation of the reverse of two inputs
func DoubleSHAPairRev(a []byte, b []byte)([]byte){
	var tmp []byte

	//hash the concatenation of the reverse of both inputs
    var answer hash.Hash=sha256.New()
    //answer.Write(append(a[:]+b[:]))
    answer.Write(append(Rev(a), Rev(b)...))
    
    //hash it again
    tmp=answer.Sum(nil)
    answer=sha256.New()
    answer.Write(tmp)
   
    return Rev(answer.Sum(nil))//return the reverse of the output
}

//TODO: test and add to tests
//SHA-256 RIPEMD-160 operation for bitcoin address hashing
func SHARipemd(b []byte)([]byte){
	
	//sha hashing of the input
    var h hash.Hash = sha256.New()
    h.Write(b)

	//ripemd hashing of the sha hash
    var h2 hash.Hash = ripemd160.New()
    h2.Write(h.Sum(nil))

    return h2.Sum(nil)//return
}

//TODO: test and add to tests
//reverse SHA-256 RIPEMD-160 hash
func SHARipemdRev(b []byte)([]byte){
   return Rev(SHARipemd(Rev(b)))
}

//TODO: test and add to tests
//RIPEMD-160 operation for bitcoin address hashing
func Ripemd(b []byte)([]byte){
	//ripemd hashing of the sha hash
    var h hash.Hash = ripemd160.New()
    h.Write(b)

    return h.Sum(nil)//return
}

//TODO: test and add to tests
//reverse RIPEMD-160 hash
func RipemdRev(b []byte)([]byte){
   return Rev(Ripemd(Rev(b)))
}