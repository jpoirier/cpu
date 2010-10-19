// Copyright (c) 2010 Joseph D Poirier
// Distributable under the terms of The New BSD License
// that can be found in the LICENSE file.

package main

import (
	"proctopo"
	"fmt"
)

func main() {
	fmt.Println("\nAn OS views threads as logical processors in a hardware multi-threading environment.")
	fmt.Println("error             : ", proctopo.Error)
	fmt.Println("cpuid present     : ", proctopo.CpuidPresent)
	fmt.Println("cpuid restricted  : ", proctopo.CpuidRestricted)
	fmt.Println("htt supported     : ", proctopo.HttSupported)
	fmt.Println("OS processors     : ", proctopo.OSProcCnt) // also via proctopo.Onln()
	fmt.Println("max processors    : ", proctopo.MaxProc)   // also via proctopo.Conf()
	fmt.Println("package count     : ", proctopo.PkgCntEnum)
	fmt.Println("core count        : ", proctopo.CoreCntEnum)
	fmt.Println("thread count      : ", proctopo.ThreadCntEnum)
	fmt.Println("htt/smt per core  : ", proctopo.HttSmtPerCore)
	fmt.Println("htt/smt per pkg   : ", proctopo.HttSmtPerPkg)
	fmt.Println("processors on line: ", proctopo.Onln())
	fmt.Println("vendor            : ", proctopo.Vendor)
}
