package main

import (
	"net"
)

func main(){
	
 listener,err := net.Listen("tcp",":5557")
 
 if err != nil{
	panic("failed to start server "+err.Error())
 }
 
 defer listener.Close()
 
 for{
 
 	con,err :=	listener.Accept()

	if err != nil{
		panic(err.Error())
		continue
	}
	
	go handleConnection(con)
 
 }
 
}

func handleConnection(con net.Conn){
	defer con.Close()
	
	var payloadSize int = 1024*5;
	buffer := make([]byte,payloadSize)
	
	for{
		
		n,err  := con.Read(buffer)
		
		if err != nil{
			print(err)
			panic(err.Error())
			return
		}
		
		print(buffer[:n])
		
	}
	
}



