// Copyright (c) 2010 Joseph D Poirier
// Distributable under the terms of The New BSD License
// that can be found in the LICENSE file.

package main

import (
	"cpu"
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("\nAn OS views physical cores and hardware threads as logical processors \nin a multi-package, multi-core, multi-threading environment.")
	fmt.Println("error                      : ", cpu.Error)
	fmt.Println("cpuid present              : ", cpu.CpuidPresent)
	fmt.Println("cpuid restricted           : ", cpu.CpuidRestricted)
	fmt.Println("hwt supported              : ", cpu.Hwt)          // hardware multi-threading
	fmt.Println("htt enabled                : ", cpu.Htt)          // hyper-threading
	fmt.Println("on line logical processors : ", cpu.OnlnProcCnt)  // also via cpu.Onln()
	fmt.Println("maximum logical processors : ", cpu.MaxProcCnt)   // also via cpu.Conf()
//	fmt.Println("package count                : ", cpu.PkgCnt)
	fmt.Println("physical core count        : ", cpu.PhyCoreCnt)
	fmt.Println("logical processor count    : ", cpu.LogProcCnt)
	fmt.Println("htt logical processor count: ", cpu.HttProcCnt)
	fmt.Println("vendor                     : ", cpu.Vendor)
	// show the exported fnctions for completeness
	fmt.Println("")
	fmt.Println("on line logical processors : ", cpu.Onln())
	fmt.Println("maximum logical processors : ", cpu.Conf())
	// set Go's runtime processor count
	runtime.GOMAXPROCS(int(cpu.MaxProcCnt))
}
