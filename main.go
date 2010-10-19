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
	fmt.Println("error             : ", s.Error)
	fmt.Println("cpuid present     : ", s.CpuidPresent)
	fmt.Println("cpuid restricted  : ", s.CpuidRestricted)
	fmt.Println("htt supported     : ", s.HttSupported)
	fmt.Println("max processors    : ", s.MaxProc) // also via proctopo.Conf()
	fmt.Println("package count     : ", s.PkgCntEnum)
	fmt.Println("core count        : ", s.CoreCntEnum)
	fmt.Println("thread count      : ", s.ThreadCntEnum)
	fmt.Println("htt/smt per core  : ", s.HttSmtPerCore)
	fmt.Println("htt/smt per pkg   : ", s.HttSmtPerPkg)
	fmt.Println("processors on line: ", proctopo.Onln())
	fmt.Println("vendor            : ", s.Vendor)
}
