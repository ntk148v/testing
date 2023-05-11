package main

import (
	"fmt"
	garbler "github.com/michaelbironneau/garbler/lib"
)

func main() {
	//use defaults
	p, _ := garbler.NewPassword(nil)
	fmt.Println(p)

	//use Strong preset (Insecure/Easy/Medium/Strong/Paranoid are available)
	p, _ = garbler.NewPassword(&garbler.Strong)
	fmt.Println(p)

	//guess requirements from existing password
	reqs := garbler.MakeRequirements("asdfGG11!")
	p, _ = garbler.NewPassword(&reqs)
	fmt.Println(p)

	//specify requirements explicitly:
	//if specifying requirements you should not ignore error return,
	//in case the requirements are impossible to satisfy (eg. minimum length is
	//greater than maximum length)
	reqs = garbler.PasswordStrengthRequirements{MinimumTotalLength: 100}
	p, e := garbler.NewPassword(&reqs)
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(p)
}
