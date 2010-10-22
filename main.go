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
	fmt.Println("HardwareThreading supported: ", cpu.HardwareThreading) // hardware multi-threading
	fmt.Println("HyperThreading enabled     : ", cpu.HyperThreading)    // hyper-threading
	fmt.Println("on line logical processors : ", cpu.OnlnProcs)         // also via cpu.OnlineProcs()
	fmt.Println("maximum logical processors : ", cpu.MaxProcs)          // also via cpu.ConfProcs()
	fmt.Println("physical processors        : ", cpu.Pkgs)
	fmt.Println("physical cores             : ", cpu.PhysicalCores)
	fmt.Println("logical processors         : ", cpu.LogicalProcs)
	fmt.Println("HyperThreading processors  : ", cpu.HyperThreadingProcs)
	fmt.Println("vendor                     : ", cpu.Vendor)
	// show the exported fnctions for completeness
	fmt.Println("")
	fmt.Println("on line logical processors : ", cpu.OnlineProcs())
	fmt.Println("maximum logical processors : ", cpu.ConfProcs())
	// set Go's runtime processor count
	runtime.GOMAXPROCS(int(cpu.MaxProcs))
}
