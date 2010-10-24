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
	fmt.Println("\nAn OS views physical cores and hyper-threads as logical processors \nin a multi-package, multi-core, multi-threading environment.")
	fmt.Println("Note that the term package refers to a physical processor\nand system refers to multiple packages\n")

	fmt.Println("physical processors (aka packages)                        : ", cpu.Processors)
	fmt.Println("on line logical processors in the system                  : ", cpu.OnlnProcs) // also via cpu.OnlineProcs()
	fmt.Println("maximum logical processors in the system                  : ", cpu.MaxProcs) // also via cpu.ConfProcs()

	fmt.Println("cpuid present                                             : ", cpu.CpuidPresent)
	fmt.Println("cpuid restricted                                          : ", cpu.CpuidRestricted)
	fmt.Println("hardware-threading supported                              : ", cpu.HardwareThreading)
	fmt.Println("hyper-threading capable                                   : ", cpu.HyperThreadingCapable)
	fmt.Println("vendor name                                               : ", cpu.Vendor)

	fmt.Println("\n    --- processor hardware capability  ---")
	fmt.Println("logical processors per physical processor                 : ", cpu.LogicalProcsPkg)
	fmt.Println("physical cores per physical processor                     : ", cpu.PhysicalCoresPkg)
	fmt.Println("hyper-threading enabled                                   : ", cpu.HyperThreadingEnabled)
	fmt.Println("hyper-threading logical processors per physical processor : ", cpu.HyperThreadingProcsPkg)

	fmt.Println("\n    --- processor hardware configuration ---")
	fmt.Println("logical processors configured                            : ", cpu.LogicalProcsConf)
	fmt.Println("physical cores configured                                : ", cpu.PhysicalCoresConf)
	fmt.Println("hyper-threading processors configured                    : ", cpu.HyperThreadingProcsConf)

	fmt.Println("")
	fmt.Println("errors during system interrogation                       : ", cpu.Error)

	// show the exported fnctions for completeness
	fmt.Println("")
	fmt.Println("on line logical processors in the system                  : ", cpu.OnlineProcs())
	fmt.Println("maximum logical processors in the system                  : ", cpu.ConfProcs())

	// set Go's runtime processor count
	runtime.GOMAXPROCS(int(cpu.MaxProcs))
}
