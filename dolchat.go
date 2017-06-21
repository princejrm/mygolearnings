// This program has been created by Richard John (princejrm@yahoo.com)
// It provides simple chat using go lang
// Syntax : go run dolchat <id> <dialport x.x.x.x:port> <listening port :port>
// example : go run dolchat Sam "127.0.0.1:9999" ":9998" 
//           go run dolchat David. "127.0.0.1:9998" ":9999"



package main 

import "fmt"
import "net"
import "encoding/gob"
import "time"
import "io"
import "os"
import "bufio"


var listenport string;
var dialport string;
var id string;

func server (){
	ln, err := net.Listen ("tcp",listenport);
	if err != nil {
		fmt.Println ("Listening error : ",err);
		return;
	}
	for {
		c, err := ln.Accept ();
		if err != nil {
			fmt.Println ("Connection Acceptance Error : ",err);
			continue;
		}
		go handleServerConnection (c);
	}
}

func handleServerConnection (c net.Conn) {
	for { 
		var msg string;
		err := gob.NewDecoder(c).Decode(&msg);
		if err != nil {
		    if err ==  io.EOF {
               fmt.Println ("Other side hung up connection");
               os.Exit(0);
            } else {
            	fmt.Printf ("Error is decoding the message : %s\n",err);
            }
		} else {
			fmt.Printf ("\t\t\t\t\t%s",msg);
		} 
	}
}


func client () {
	var c net.Conn;
	var err error;
	var progress int = 0;
	for {
		c, err = net.Dial ("tcp",dialport);
		if err != nil {
			switch progress { 
				case 0:
					fmt.Print("Trying to connect -");
					progress = 1;
					break;
				case 1:
					fmt.Printf("\b\\");
					progress = 2;
					break;
				case 2:
					fmt.Printf ("\b|");
					progress = 3;
					break;
				case 3:
					fmt.Printf ("\b/");
					progress = 4;
					break;
				case 4:
					fmt.Printf ("\b-");
					progress = 1;
					break;
			}
			time.Sleep (time.Millisecond * 125);
			continue;
		}
		fmt.Printf ("\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b");
		break;
	}
	fmt.Println ("Connection established, start chating...");
	for {
		var msg string;
		reader := bufio.NewReader(os.Stdin);
		msg, _ = reader.ReadString('\n');
		if (msg == "quit\n" || msg == "Quit\n") {
			break;
		}
		msg = id + ": " + msg;
		var err error
		err = gob.NewEncoder(c).Encode(msg);
		if err != nil {
			fmt.Println ("Error is sending information : ",err);
		}
	}
}

func main () {
	if (len(os.Args) != 4) {
		fmt.Println ("Invalid commandline arguments");
		fmt.Printf ("%s <id> <dialport x.x.x.x:port> <listening port :port>",os.Args[0]);
		os.Exit (0);
	}
	id = os.Args[1];
	dialport = os.Args[2];
	listenport = os.Args[3];
	go server ();
	client ();
}