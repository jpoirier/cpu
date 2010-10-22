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
var CpuidPresent    bool

// CpuidRestricted indicates whether the cpuid instruction is restricted,
// i.e. not executable.
var CpuidRestricted bool

// Hwt indicates whether hardware multi-threading is supported,
// can be hyper-threading and/or multiple physical cores.
var Hwt             bool

// Vendor is the package vendor's name.
var Vendor          string

// MaxProcCnt is the maximum number of logical processors supported by the OS.
var MaxProcCnt      uint32

// OnlnProcCnt is the number of logical processors that are on line.
var OnlnProcCnt     uint32

// PhyCoreCnt is the number of physical cores in the package.
var PhyCoreCnt      uint32

// LogProcCnt is the maximum addressable logical processors in the package,
// but not necessarily occupied by a logical processors
var LogProcCnt      uint32

// HttProcCnt is the number of hyper-threading logical processors in the package.
var HttProcCnt      uint32

// Error reports if an error occurred during the information gathering process.
// TODO: Needs to be fine grained so the caller knows where the error occurred
var Error           bool

// PkgCnt is the number of physical packages in the system.
//var PkgCnt          uint32

type regs struct {
	eax uint32
	ebx uint32
	ecx uint32
	edx uint32
}

// Onln returns the number of logical processors that are on line.
func Onln() uint32 {
	return uint32(C.onln())
}

// Conf returns the maximum number of logical processors supported by the OS.
func Conf() uint32 {
	return uint32(C.conf())
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

// CpuParams
func CpuParams() bool {
//	PkgCnt        = 1
	Hwt           = false
	PhyCoreCnt    = 1
	LogProcCnt    = 1
	MaxProcCnt    = Conf()
	OnlnProcCnt   = Onln()
	// cpuid check
	CpuidPresent = have_cpuid()
	if !CpuidPresent { return false }
	// vendor name
	var info regs
	cpuid(&info, 0, 0)
	maxCpuid := info.eax
	Vendor = utos(info.ebx) + utos(info.edx) + utos(info.ecx)
	// restricted cpuid execution
	var r regs
	CpuidRestricted = false
	cpuid(&r, 0x80000000, 0)
	if maxCpuid<=4 && r.eax>0x80000004 {
		CpuidRestricted = true
		return false
	}
	// htt enabled
	cpuid(&r, 1, 0)
	if r.edx>>28&1 != 0 {
		Hwt = true
	}
	if !Hwt { return false } // single core and no htt
	LogProcCnt = r.ebx >> 16 & 0xFF
	if Vendor == "GenuineIntel" {
		cpuid(&r, 4, 0)
		PhyCoreCnt = (r.eax >> 26 & 0x3F) + 1
	} else if Vendor == "AuthenticAMD" {
		cpuid(&r, 0x80000008, 0)
		PhyCoreCnt = (r.ecx & 0xFF) + 1
	} else {
		Error = true
		return false
	}
	if MaxProcCnt > PhyCoreCnt { HttProcCnt = MaxProcCnt - PhyCoreCnt }
	return true
}

func init() {
	CpuParams()
}
