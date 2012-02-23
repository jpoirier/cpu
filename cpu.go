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

var PackageVersion string = "v0.13.3"

// Processors is the number of physical processors (that plug in to a socket).
var Processors uint32

// CpuidPresent indicates whether the cpuid instruction is present.
var CpuidPresent bool

// CpuidRestricted indicates whether the cpuid instruction is restricted,
// i.e. not executable.
//var CpuidRestricted bool

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

// LogicalProcsSharingCache
var LogicalProcsSharingCache uint32

// HyperThreadingProcsConf is the number of hyper-threading logical processors configured in the package.
var HyperThreadingProcsConf uint32

// HyperThreadingProcsPkg is the number of hyper-threading logical processors available in the package.
var HyperThreadingProcsPkg uint32

var ProcessorFamily string

var ProcessorL2Cache uint32

var ProcessorL2CacheLine uint32


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

// Cpuid executes the EAX function f1 and ECX sub function f2, the output
// registers from the operation returned in r.
func Cpuid(r *regs, f1, f2 uint32) {
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
	ProcessorL2Cache = 0
	ProcessorL2CacheLine = 0
	ProcessorFamily = "Unknown"
	Vendor = "Unknown"
//	CpuidRestricted = false
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
	Cpuid(&r, 0, 0)
	maxStdLevel := r.eax
	Vendor = utos(r.ebx) + utos(r.edx) + utos(r.ecx)

	Cpuid(&r, 0x80000000, 0)
	maxExtLevel := r.eax

	// this check includes some old processors (P4 & M, Old Xeon)
	// that we could report processor name on but probably not
	// worth the time
	if maxStdLevel < 5  {
		return
	}

	if maxExtLevel >= 0x80000004 {
		Cpuid(&r, 0x80000002, 0)
		ProcessorFamily = utos(r.eax) + utos(r.ebx) + utos(r.ecx) + utos(r.edx)
		Cpuid(&r, 0x80000003, 0)
		ProcessorFamily += utos(r.eax) + utos(r.ebx) + utos(r.ecx) + utos(r.edx)
		Cpuid(&r, 0x80000004, 0)
		ProcessorFamily += utos(r.eax) + utos(r.ebx) + utos(r.ecx) + utos(r.edx)
	}

	if maxExtLevel >= 0x80000006 {
		Cpuid(&r, 0x80000006, 0)
		ProcessorL2CacheLine = (r.ecx & 0xFF)
		ProcessorL2Cache = ((r.ecx >> 16) & 0xFFFF) * 1024
	}

	// The hardware capability of a package may be different from its configuration.
	// A package may be capable of addressing multiple logical processors,
	// in the case of multiple cores, but that's not a good indication that
	// core multi-processing is enabled. E.g. each core in an Athlon 64 X2
	// multi-core CPU is its own distinct processor and shares no resources with
	// other cores. Multi-core processors are distinguished by their level of
	// integration. Do AMD processors also have core multi-processing?

	Cpuid(&r, 1, 0)
	if r.edx >> 28 & 1 == 0 {
		return // single core and no hardware-threading
	}

	if Vendor == "GenuineIntel" {
		LogicalProcsPkg = r.ebx >> 16 & 0xFF
		Cpuid(&r, 4, 0)
		PhysicalCoresPkg = (r.eax >> 26 & 0x3F) + 1
		LogicalProcsSharingCache = (r.eax >> 14 & 0xFFF) + 1

		if PhysicalCoresPkg > 1 {
			HardwareThreading = true
		}
	} else if Vendor == "AuthenticAMD" {
		Cpuid(&r, 0x80000008, 0)
		apicid_sz := (r.ecx >> 12) & 0xF
		LogicalProcsPkg = (r.ecx & 0xFF) + 1

		if apicid_sz == 0 {
			PhysicalCoresPkg = LogicalProcsPkg // legacy mode check
		} else {
			PhysicalCoresPkg = pow(2, apicid_sz)
		}
	} else {
// TODO: abort? handle other vendors
		panic("Unknown processor...")
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
