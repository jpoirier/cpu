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

var CpuidPresent    bool
var CpuidRestricted bool
var HttSupported    bool
var Vendor          string
var MaxProc         uint32
var OSProcCnt       uint32
var PkgCntEnum      uint32
var CoreCntEnum     uint32
var ThreadCntEnum   uint32
var HttSmtPerCore   uint32
var HttSmtPerPkg    uint32
var Error           bool

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

// cpuid executes the function f1 and sub function f2 with the output
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

// CpuParams
func CpuParams() bool {
	PkgCntEnum    = 1
	CoreCntEnum   = 1
	ThreadCntEnum = 1
	HttSmtPerCore = 0
	HttSmtPerPkg  = 0
	MaxProc       = Conf()
	OSProcCnt     = Onln()
	// cpuid check
	CpuidPresent = have_cpuid()
	if !CpuidPresent {
		return false
	}
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
	HttSupported = false
	cpuid(&r, 1, 0)
	if r.edx>>28&1 != 0 {
		HttSupported = true
	}
	// physical and logical core cnt for this package
	var logCoreCnt uint32 = r.ebx >> 16 & 0xFF
	var phyCoreCnt uint32
	if Vendor == "GenuineIntel" {
		cpuid(&r, 4, 0)
		phyCoreCnt = (r.eax >> 26 & 0x3F) + 1
	} else if Vendor == "AuthenticAMD" {
		cpuid(&r, 0x80000008, 0)
		phyCoreCnt = (r.ecx & 0xFF) + 1
	} else {
		Error = true
		return false
	}
	PkgCntEnum    = MaxProc / logCoreCnt // wrong? symmetrical for multiple packages?
	CoreCntEnum   = phyCoreCnt
	ThreadCntEnum = logCoreCnt
	HttSmtPerPkg  = logCoreCnt - phyCoreCnt
	if logCoreCnt > phyCoreCnt && HttSupported { HttSmtPerCore = 1 /* always 1? */ }
	return true
}

func init() {
	CpuParams()
}
