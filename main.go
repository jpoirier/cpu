// Copyright (c) 2010 Joseph D Poirier
// Distributable under the terms of The New BSD License
// that can be found in the LICENSE file.

package main

import (
	"cpu"
	"fmt"
)

func main() {
	fmt.Println("\nAn OS views physical cores and hardware threads as logical processors \nin a multi-package, multi-core, multi-threading environment.")
	fmt.Println("error             : ", cpu.Error)
	fmt.Println("cpuid present     : ", cpu.CpuidPresent)
	fmt.Println("cpuid restricted  : ", cpu.CpuidRestricted)
	fmt.Println("htt supported     : ", cpu.HttSupported)
	fmt.Println("OS processors     : ", cpu.OSProcCnt) // also via cpu.Onln()
	fmt.Println("max processors    : ", cpu.MaxProc)   // also via cpu.Conf()
	fmt.Println("package count     : ", cpu.PkgCntEnum)
	fmt.Println("core count        : ", cpu.CoreCntEnum)
	fmt.Println("thread count      : ", cpu.ThreadCntEnum)
	fmt.Println("htt/smt per core  : ", cpu.HttSmtPerCore)
	fmt.Println("htt/smt per pkg   : ", cpu.HttSmtPerPkg)
	fmt.Println("processors on line: ", cpu.Onln())
	fmt.Println("vendor            : ", cpu.Vendor)
}
