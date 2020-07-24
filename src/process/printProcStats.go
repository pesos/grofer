package process

import (
	"fmt"
)

func (p *Process) PrintStats() {
	fmt.Println(p.Name, p.CPUPercent)
	// fmt.Println(p) //not sure of format in which we'll print
}
