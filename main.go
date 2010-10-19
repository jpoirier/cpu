// Copyright (c) 2010 Joseph D Poirier
// Distributable under the terms of The New BSD License
// that can be found in the LICENSE file.

package main

import (
	"proctopo"
	"fmt"
)

func main() {
	var s proctopo.ProcTopo_t
	proctopo.ProcTopo(&s)
	fmt.Println("cpuid present     : ", s.CpuidPresent)
	fmt.Println("cpuid restricted  : ", s.CpuidRestricted)
	fmt.Println("htt supported     : ", s.HttSupported)
	fmt.Println("Max CPUs          : ", s.MaxCpus)
	fmt.Println("vendor            : ", s.Vendor)
}
