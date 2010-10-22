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

// CpuidPresent indicates whether the cpuid instruction is present.
var CpuidPresent bool

// CpuidRestricted indicates whether the cpuid instruction is restricted,
// i.e. not executable.
var CpuidRestricted bool

// HardwareThreading indicates whether hardware multi-threading is supported,
// can be hyper-threading and/or multiple physical cores.
var HardwareThreading bool

// HyperThreading indicates whether hyper-threading is enabled.
var HyperThreading bool

// Vendor is the package vendor's name.
var Vendor string

// MaxProcs is the maximum number of logical processors supported by the OS.
// This may include logical processors from packages outside of the one being
// reported on.
var MaxProcs uint32

// OnlnProcs is the number of logical processors that are on line.
var OnlnProcs uint32

// PhysicalCores is the number of physical cores in the package.
var PhysicalCores uint32

// LogicalProcs is the maximum addressable logical processors in the package,
// but not necessarily occupied by a logical processors
var LogicalProcs uint32

// HyperThreadingProcs is the number of hyper-threading logical processors in the package.
var HyperThreadingProcs uint32

// Error reports if an error occurred during the information gathering process.
// TODO: Needs to be fine grained so the caller knows where the error occurred
var Error bool

var Pkgs uint32

// PkgCnt is the number of physical packages in the system.
//var PkgCnt          uint32

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
	b[0] = byte(a >>  0); b[1] = byte(a >>  8);
	b[2] = byte(a >> 16); b[3] = byte(a >> 24)
	return fmt.Sprintf("%s", b)
}

func mask_width(x uint32) uint32 {
	if x == 0 { return 0 }
	return ^(0xFFFFFFFF<<x)
}

// CpuParams
func CpuParams() bool {
	Pkgs = 1
	HardwareThreading = false
	HyperThreading = false
	PhysicalCores = 1
	LogicalProcs = 1
	HyperThreadingProcs = 0
	MaxProcs = ConfProcs()
	OnlnProcs = OnlineProcs()
	// cpuid check
	CpuidPresent = have_cpuid()
	if !CpuidPresent { return false }
	// vendor name
	var r regs
	cpuid(&r, 0, 0)
	maxCpuid := r.eax
	Vendor = utos(r.ebx) + utos(r.edx) + utos(r.ecx)
	// restricted cpuid execution
	CpuidRestricted = false
	cpuid(&r, 0x80000000, 0)
	if maxCpuid<=4 && r.eax>0x80000004 {
		CpuidRestricted = true
		return false
	}
	// HardwareThreading enabled
	cpuid(&r, 1, 0)
	if r.edx>>28&1 != 0 {
		HardwareThreading = true
	}
	if !HardwareThreading { return false } // single core and no HyperThreading
	LogicalProcs = r.ebx >> 16 & 0xFF
	apicid := r.ebx >> 24 & 0xFF
	if Vendor == "GenuineIntel" {
		cpuid(&r, 4, 0)
		PhysicalCores = (r.eax >> 26 & 0x3F) + 1
	} else if Vendor == "AuthenticAMD" {
		cpuid(&r, 0x80000008, 0)
		PhysicalCores = (r.ecx & 0xFF) + 1
	} else {
		Error = true
		return false
	}
	// HyperThreading enabled and HardwareThreading logical processors
	smtid_mask := mask_width(LogicalProcs-PhysicalCores)
	if smtid_mask > 0 {
		HyperThreading = true
		HyperThreadingProcs = PhysicalCores * (apicid & smtid_mask)
	}
	return false
}

func init() {
	CpuParams()
}
