// Copyright (c) 2010 Joseph D Poirier
// Distributable under the terms of The New BSD License
// that can be found in the LICENSE file.

package cpu

// #include "cpu.h"
import "C"

import (
	"fmt"
	"unsafe"
)


// Processors is the number of physical processors (that plug in to a socket).
var Processors uint32

// CpuidPresent indicates whether the cpuid instruction is present.
var CpuidPresent bool

// CpuidRestricted indicates whether the cpuid instruction is restricted,
// i.e. not executable.
var CpuidRestricted bool

// HardwareThreading indicates whether hardware multi-threading is supported,
// can be hyper-threading and/or multiple physical cores within a package.
var HardwareThreading bool

// HyperThreadingEnabled indicates whether the package has hyper-threading enabled.
var HyperThreadingEnabled bool

// Vendor is the package vendor's name.
var Vendor string

// MaxProcs is the maximum number of logical processors supported by the OS.
// This may include logical processors from packages outside of the one being
// reported on.
var MaxProcs uint32

// OnlnProcs is the number of logical processors that are on line.
var OnlnProcs uint32

// PhysicalCoresConf is the number of physical cores configured in the package.
var PhysicalCoresConf uint32

// PhysicalCoresPkg is the number of physical cores in the package.
var PhysicalCoresPkg uint32

// LogicalProcsConf is the number of logical processors configured in the package.
var LogicalProcsConf uint32

// LogicalProcsPkg is the maximum number of addressable logical processors in the package,
// but not necessarily occupied by a logical processors
var LogicalProcsPkg uint32

var LogicalProcsSharingCache uint32

// HyperThreadingProcsConf is the number of hyper-threading logical processors configured in the package.
var HyperThreadingProcsConf uint32

// HyperThreadingProcsPkg is the number of hyper-threading logical processors available in the package.
var HyperThreadingProcsPkg uint32

type regs struct {
	eax uint32
	ebx uint32
	ecx uint32
	edx uint32
}

// OnlinenProcs returns the number of logical processors that are on line.
func OnlineProcs() uint32 {
	return uint32(C.onlineProcs())
}

// ConfProcs returns the maximum number of logical processors supported by the OS.
func ConfProcs() uint32 {
	return uint32(C.confProcs())
}

// have_cpuid returns whether or not the cpuid instruction exists
func have_cpuid() bool {
	return bool(C.have_cpuid())
}

// cpuid executes the function f1 and sub function f2, the output
// registers from the operation returned in r.
func cpuid(r *regs, f1, f2 uint32) {
	C.cpuid((*C.regs_t)(unsafe.Pointer(r)), C.uint32_t(f1), C.uint32_t(f2))
}

// utos returns a converted to a string.
func utos(a uint32) string {
	var b [4]byte
	b[0] = byte(a >> 0)
	b[1] = byte(a >> 8)
	b[2] = byte(a >> 16)
	b[3] = byte(a >> 24)
	return fmt.Sprintf("%s", b)
}

func pow(x, e uint32) uint32 {
	v := uint32(1)
	for {
		if e == 0 {
			break
		}
		if e & 1 != 0 {
			v *= x
		}
		x *= x
		e >>= 1
	}
	return v
}

// Params
func Params() {
	Processors = 1
	HardwareThreading = false
	HyperThreadingEnabled = false
	PhysicalCoresConf = 1
	PhysicalCoresPkg = 1
	LogicalProcsConf = 1
	LogicalProcsPkg = 1
	HyperThreadingProcsConf = 0
	HyperThreadingProcsPkg = 0
	Vendor = "Unknown"
	CpuidRestricted = false
	LogicalProcsSharingCache = 0

	MaxProcs = ConfProcs()
	OnlnProcs = OnlineProcs()

	// cpuid check, assumes 32-bit processor
	CpuidPresent = have_cpuid()
	if !CpuidPresent {
		return
	}

	// vendor name
	var r regs
	cpuid(&r, 0, 0)
//	maxStdLevel := r.eax
	Vendor = utos(r.ebx) + utos(r.edx) + utos(r.ecx)

	cpuid(&r, 0x80000000, 0)
//	maxExtLevel = r.eax

// XXX: validate this check!
	// restricted cpuid execution
//	if maxStdLevel <= 4 && maxExtLevel > 0x80000004 {
//		CpuidRestricted = true
//		return
//	}

	// The hardware capability of a package may be different from its configuration.
	// A package may be capable of addressing multiple logical processors,
	// in the case of multiple cores, but that's not a good indication that
	// core multi-processing is enabled. E.g. each core in an Athlon 64 X2
	// multi-core CPU is its own distinct processor and shares no esources with
	// other cores. Multi-core processors are distinguished by their level of
	// integration. Do AMD processors also do core multi-processing?

	cpuid(&r, 1, 0)
	if r.edx >> 28 & 1 == 0 {
		return // single core and no Hhardware-threading
	}

	if Vendor == "GenuineIntel" {
		LogicalProcsPkg = r.ebx >> 16 & 0xFF
		cpuid(&r, 4, 0)
		PhysicalCoresPkg = (r.eax >> 26 & 0x3F) + 1
		LogicalProcsSharingCache = (r.eax >> 14 & 0xFFF) + 1

		if PhysicalCoresPkg > 1 {
			HardwareThreading = true
		}
	} else if Vendor == "AuthenticAMD" {
		cpuid(&r, 0x80000008, 0)
		apicid_sz := (r.ecx >> 12) & 0xF
		LogicalProcsPkg = (r.ecx & 0xFF) + 1

		if apicid_sz == 0 {
			PhysicalCoresPkg = LogicalProcsPkg // legacy mode check
		} else {
			PhysicalCoresPkg = pow(2, apicid_sz)
		}
	} else {
// TODO: abort? handle other vendors
		return
	}

	if LogicalProcsPkg < PhysicalCoresPkg {
		LogicalProcsPkg = PhysicalCoresPkg // a hardware problem if this happens!
	} else if (LogicalProcsPkg - PhysicalCoresPkg) > 0 {
		HyperThreadingEnabled = true
//		HyperThreadingProcsPkg = LogicalProcsPkg - PhysicalCoresPkg
		HyperThreadingProcsPkg = PhysicalCoresPkg * (LogicalProcsSharingCache - 1)
	}

	// Intel supports only homogeneous MP
	if MaxProcs > LogicalProcsPkg {
		Processors = MaxProcs / LogicalProcsPkg
	}

	return
}

func init() {
	Params()
}
