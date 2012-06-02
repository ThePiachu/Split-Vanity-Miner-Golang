package main

import(
	"fmt"
	"os"
	"bufio"
	"net/http"
	"io/ioutil"
	"mymath"
	"math/big"
	"bitelliptic"
	"bitecdsa"
	"crypto/rand"
)

var Address string

type Work struct{
	PublicKey string
	X *big.Int
	Y *big.Int
	Pattern string
	NetByte byte
	Reward float64
}

func main() {
	fmt.Println("Hello world!")
	file, err := os.Open("address.txt")
	if err != nil {
		fmt.Println("Can't find address.txt file. Terminating.")
        return
    }
    reader := bufio.NewReader(file)
    part, _, err := reader.ReadLine()
    if err!=nil{
    	fmt.Println("Problems reading address.txt file. Terminating.")
    	return
    }
    Address=string(part)
	fmt.Println("Your address is - ", Address)
	fmt.Println("Fetching work (vanitypooltest)...")
	
	
	
	
	response, err := http.Get("https://vanitypooltest.appspot.com/getWork")
	if err != nil {
		fmt.Println("Problems fetching work. Terminating.")
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Problems reading fetched work. Terminating.")
		return
	}
	workString:=string(body)
	fmt.Println("Fetched work:")
	fmt.Println(workString)
	if len(workString)==0{
		fmt.Println("No work. Terminating.")
		return
	}
	tasks:=mymath.SplitStrings(workString, "\n")
	
	work:=make([]Work, len(tasks)-1)
	for i:=0;i<len(tasks)-1;i++{
		task:=mymath.SplitStrings(tasks[i], ";")[0]
		parts:=mymath.SplitStrings(task, ":")
		//fmt.Println(parts)
		work[i].PublicKey=parts[0]
		work[i].Pattern=parts[1]
		work[i].NetByte=byte(mymath.Str2Int64(parts[2]))
		work[i].Reward=mymath.Str2Float(parts[3])
		
		work[i].X, work[i].Y=mymath.PublicKeyToPointCoordinates(work[i].PublicKey)
	}
	//fmt.Println(work)
	
	fmt.Println("")
	fmt.Println("Starting the work!")
	
	curve:=bitelliptic.S256()
	
	for count:=0;len(work)>0;count++{
		if count%100==99{
			fmt.Println("Checked 100 keys...")
		}
		key, err:=bitecdsa.GenerateKey(curve, rand.Reader)
		if err!=nil{
			fmt.Println("Error encountered while generating keys:")
			fmt.Println(err.Error())
			return
		}
		for i:=0;i<len(work);i++{
			x, y:=curve.Add(work[i].X, work[i].Y, key.X, key.Y)
			address:=string(mymath.NewFromPublicKey(work[i].NetByte,
				append(append([]byte{0x04}, mymath.Big2HexPadded(x, 32)...), mymath.Big2HexPadded(y, 32)...)).Base)
			ok:=true
			for j:=0;j<len(work[i].Pattern);j++{
				if work[i].Pattern[j]!=address[j]{
					ok=false
				}
			}
			if ok==true{
				fmt.Println("Solved a work!")
				fmt.Println(tasks[i])
				fmt.Printf("Solution - %X\n", mymath.Big2Hex(key.D))
				
				fmt.Println("Attempting to hand in work...")
				
				
				postAddress:="https://vanitypooltest.appspot.com/solveWork?key="+work[i].PublicKey+":"+work[i].Pattern+"&privateKey="+mymath.Hex2Str(mymath.Big2Hex(key.D))+"&bitcoinAddress="+Address
				response, err = http.Get(postAddress)
				if err != nil {
					fmt.Println("Problems fetching work. Terminating.")
					return
				}
				defer response.Body.Close()
				body, err = ioutil.ReadAll(response.Body)
				if err != nil {
					fmt.Println("Problems reading fetched work. Terminating.")
					return
				}
				
				fmt.Println("Server response:", string(body))
				if string(body)=="OK!"{
					fmt.Println("Work accepted, yay!")
					work=append(work[0:i], work[i+1:]...)
				} else {
					fmt.Println("Work did not get accepted...")
				}
			}
		}
	}
}
