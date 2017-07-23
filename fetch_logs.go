/*
   Usage : <fetch_logs.go> <link to download> <pattern of download links>
   This program helps to download a set of href download files in parallel
*/



package main 

import (
	"bytes"
	"fmt"
	"os/exec"
	"os"
	"strings"
	"time"
    "golang.org/x/net/html"
)

var file_list [100] string;
var index int32 = 0;
var print_status_list map[string]string;
var link,pattern string;

func main () {
	start := time.Now ();
	link = os.Args[1];
	pattern = os.Args[2];
	print_status_list = make (map[string]string);
	html_list := download_log (link,"retout");
	//fmt.Printf("Output from log listing%q\n", html_list);
	doc, err := html.Parse(strings.NewReader(html_list))
	if err != nil {
		fmt.Println(err)
		return;
	} 
	index = 0;
    //parse_html (doc,"205018505");
    //parse_html (doc,"205006911"); 
    parse_html (doc,pattern); 
    //fmt.Println (file_list)
    var url, out string;
    go print_status ();
    for _,file := range file_list {
    	if (file == "") {
    		break;
    	}
		fmt.Println (file);
		if (file != "") { 
			url = link + file;
			out = "./logs/" + file;
			go download_log (url,out);
	    } else {
	    	break;
	    }
	}
    var input string;
	fmt.Scanln (&input);
	fmt.Printf ("Time taken by program : %s",time.Since(start));
}

func print_status () {
	var file,status string;
	var nooflines int;
	var row int;

	for { 
		nooflines = len(print_status_list);
		if (nooflines != 0) { 
			row = 0;
			for file,status = range print_status_list {
				fmt.Printf ("\033[%d;0H",row+15);
				fmt.Printf ("%s : %s              ",file,status);
				row++;
			}
		}
		time.Sleep (time.Second * 2);
	} 
}


func download_log (url string, output string) string {
	var cmd * exec.Cmd;
	if (output == "retout") {
            cmd = exec.Command ("curl",url);
		} else {
			cmd = exec.Command("curl", url, "-o",output);
			print_status_list[url] = "In progress";
		}	
	cmd.Stdin = strings.NewReader("some input");
	var out bytes.Buffer;
	cmd.Stdout = &out;

	err := cmd.Run();
	for err != nil {
		fmt.Println ("Error : ",err);
        time.Sleep (time.Second * 5);
        if (output == "retout") {
        		cmd = exec.Command ("curl",url);
        } else {
        		cmd = exec.Command("curl", url, "-o",output);
        }
        cmd = exec.Command("curl", url, "-o",output);
		cmd.Stdin = strings.NewReader("some input");
		var out bytes.Buffer;
		cmd.Stdout = &out;
        err = cmd.Run ();
	}
	//fmt.Printf("in all caps: %q\n", out.String())	
	if (output != "retout") {
		print_status_list[url] = "Finished";
	}
	return out.String ();
}


func parse_html (n *html.Node,s string) {
    if n.Type == html.ElementNode && n.Data == "a" {
        for _, a := range n.Attr {
            if a.Key == "href" && strings.Contains(a.Val,s) {
                //fmt.Println(a.Val)
                file_list[index] = a.Val;
                index++;
                //fmt.Println (index);
                break
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        parse_html(c,s)
    }
}


