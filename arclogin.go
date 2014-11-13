/*
	arclogin [username] [password]

Copyright (c) <YEAR>, <OWNER>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer. 
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/
package main

import (
	"fmt"
	"bufio"
	"os"
	"os/user"
	"os/exec"
	"code.google.com/p/gopass"
	"flag"
	"strings"
)

func usage() {
	/* generate usage */
	fmt.Fprint(os.Stderr, "usage: arclogin [username] [password]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return false, err
}

func main() {
	var username string
	var passwd string
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if (len(args)<2){
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the username: ")
		u, _ := reader.ReadString('\n')
		pwd,err := gopass.GetPass("Enter the passphrase: ")
		if err!=nil {
			fmt.Println("Error entering the passphrase")
			os.Exit(1)
		}
		username=strings.TrimSpace(u)
		passwd=strings.TrimSpace(pwd)
	} else {
		username = args[0]
		passwd = args[1]
	}
	u, err := user.Current()
	if err!=nil {
		fmt.Println("Error could not get current user.")
		os.Exit(1)
	}
	home := u.HomeDir
	if (home == ""){
		fmt.Println("Error cannot find home directory for current user")
		os.Exit(1)
	}
	stat,err := exists(home+"/.archiver.photo")
	if err!=nil {
		fmt.Println("Error checking home directory")
		os.Exit(1)
	}
	if (!stat) {
		err := os.Mkdir(string(home) + "/.archiver.photo",0700)
		if err!=nil {
			fmt.Println("Could not create "+
				home+"/.archiver.photo");
			os.Exit(1)
		}
		fmt.Println("Created: "+string(home)+
			"/.archiver.photo");
	}
	curlargs := []string{"--data","user="+username+"&pass="+passwd, 
			"https://archiver.photo/cgi-bin/arcauth.cgi"};

	cookie,err := exec.Command("/usr/local/bin/curl",curlargs...).Output()
	if err!=nil {
		fmt.Println("Could not connect to auth server.")
		fmt.Println(err)
		os.Exit(1)
	}
	if (strings.TrimSpace(string(cookie)) == "401 Unauthorized") {
		fmt.Println("Error, invalid username and/or passphrase")
		os.Exit(1)
	}
	f, err := os.Create(home+"/.archiver.photo/cookie")
	if err!=nil {
		fmt.Println("Error creating cookie file.")
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	f.Write(cookie)
	f.Sync()
	fmt.Println("Auth cookie retrieved and saved.")
}
