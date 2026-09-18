package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type standardResponse struct {
	message    string
	statusCode int
}

func main() {

	listener, err := net.Listen("tcp", ":5557")

	if err != nil {
		panic("failed to start server " + err.Error())
	}

	fmt.Print("Server started on port 5557 \n")

	defer listener.Close()

	for {

		con, err := listener.Accept()

		if err != nil {
			fmt.Print(err.Error())
			break
		}

		go handleConnection(con)

	}

}

func handleConnection(con net.Conn) {
	defer con.Close()

	for {

		var payloadSize int = 1024 * 5
		buffer := make([]byte, payloadSize)

		_, err := con.Read(buffer)

		if err != nil {
			print(err.Error() + "\n")
			return
		}

		var data string = string(buffer)

		returnData, err := doHttpParsing(data)

		if err != nil {
			errMessage := standardResponse{
				message:    "error parsing http request",
				statusCode: 400,
			}

			Bs, err := json.Marshal(errMessage)

			_, err = con.Write(Bs)

			if err != nil {
				fmt.Print("error writing to client: " + err.Error())
			}
		}

		Bs, err := json.Marshal(returnData)

		_, err = con.Write(Bs)

		if err != nil {
			fmt.Print("error writing to client: " + err.Error())
		}

	}

}

func doHttpParsing(data string) (string, error) {

	// var header = make(map[string]string)
	// header["GET"] = ""

	var dataSeparator = "\r\n\r\n"
	var dataInLine = strings.Split(data, "\n")
	var dataGotten = strings.SplitN(data, dataSeparator, 2)
	var requestType string
	var httpVersion string
	var endpoint string
	var contentType string
	var contentLenght string

	var partHi = "/hi"     //post
	var partPing = "/ping" //get

	for eachLine := range dataInLine {
		//for first line the POST /api/upload?file=test.txt&debug=true HTTP/1.1
		// the first would be requestType eg Post,get,put second would be endpoint the three would be http version eg http/1.1
		if eachLine == 0 {

			var firstLine = strings.Split(dataInLine[eachLine], " ") //would not work for search because they might be space in it eg ?name=eze larry
			for splited := range firstLine {

				// not supported
				// strings.Contains("PATCH",strings.ToUpper(firstLine[splited])) ||
				// strings.Contains("HEAD",strings.ToUpper(firstLine[splited])) ||
				// strings.Contains("OPTIONS",strings.ToUpper(firstLine[splited]))

				if strings.Contains("POST", strings.ToUpper(firstLine[splited])) {
					requestType = "POST"
				} else if strings.Contains("GET", strings.ToUpper(firstLine[splited])) {
					requestType = "GET"

				} else if strings.Contains("PUT", strings.ToUpper(firstLine[splited])) {
					requestType = "PUT"
				} else if strings.Contains("DELETE", strings.ToUpper(firstLine[splited])) {
					requestType = "DELETE"
				}

				if strings.ContainsAny(firstLine[splited], "/") {
					endpoint = firstLine[splited]
				}

				if strings.Contains("HTTP", strings.ToUpper(firstLine[splited])) {
					httpVersion = firstLine[splited]
				}

			}

		}
		// print(dataInLine[eachLine] + "\n")

		var twoOfThem = strings.SplitN(dataInLine[eachLine], ":", 2)
		if twoOfThem[0] == "Content-Type" {
			contentType = strings.TrimSpace(twoOfThem[1])
		}
		if twoOfThem[0] == "Content-Length" {
			contentLenght = strings.TrimSpace(twoOfThem[1])
		}

		if endpoint == partHi && requestType == "POST" {
			return "hello", nil

		} else if endpoint == partPing && requestType == "GET" {

			var d = "pong \n here is your data \n " + dataGotten[1]
			return d, nil

		} else if requestType == "/" {
			return "Whttp server", nil
		}
		fmt.Printf("content-type %g content-length %y http-version %p", contentType, contentLenght, httpVersion)
		// print(twoOfThem[0])

		// for oneOfTheme := range twoOfThem{

		// }

	}
	return "nil", nil
}
