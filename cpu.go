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

// CpuidPresent is the cpuid instruction present flag.
var CpuidPresent    bool
// CpuidRestricted is the cpuid instruction restricted flag.
var CpuidRestricted bool
// HttSupported is the Hyper-Threading Technology supported flag.
var HttSupported    bool
// Vendor is the package vendor's name.
var Vendor          string
// MaxProc is the maximum number of logical processors supported by the OS.
var MaxProc         uint32
// OSProcCnt is the number of logical processors that are on line.
var OSProcCnt       uint32
// PkgCnt is the number of physical packages in the system.
var PkgCnt          uint32
// CoreCnt is the number of physical cores in the system.
var CoreCnt         uint32
// ThreadCnt is the number of threads running in the system.
var ThreadCnt       uint32
// HttSmtPerCore is the logical core count associated to a physical core,
// excluding the physical core (always 1?)
var HttSmtPerCore   uint32
// HttSmtPerPkg is the logical core count minus the physical core count in a package.
var HttSmtPerPkg    uint32
// Error
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

// rev_u32 returns the bits of a in reverse order
func rev_u32(a uint32) uint32 {
	a = (((a&0XAAAAAAAA)>>1)|((a&0X55555555)<< 1))
	a = (((a&0XCCCCCCCC)>>2)|((a&0X33333333)<< 2))
	a = (((a&0XF0F0F0F0)>>4)|((a&0X0F0F0F0F)<< 4))
	a = (((a&0XFF00FF00)>>8)|((a&0X00FF00FF)<< 8))
	return ((a>>16)|(a<<16))
}

var l_table = [32]uint32 {
	0,9,1,10,13,21,2,29,11,14,16,18,22,25,3,30,
	8,12,20,28,15,17,24,7,19,27,23,6,26,5,4,31}

// high_bit returns the position of the highest order set bit in a
func high_bit(a uint32) uint32 {
	a |= a>>1; a |= a>>2; a |= a>>4; a |= a>>8; a |= a>>16
	return l_table[(a*0x07C4ACDD)>>27]
}

// CpuParams
func CpuParams() bool {
	PkgCnt        = 1
	CoreCnt       = 1
	ThreadCnt     = 1
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
	PkgCnt        = MaxProc / logCoreCnt // wrong? symmetrical for multiple packages?
	CoreCnt       = phyCoreCnt
	ThreadCnt     = logCoreCnt
	HttSmtPerPkg  = logCoreCnt - phyCoreCnt
	if logCoreCnt > phyCoreCnt && HttSupported { HttSmtPerCore = 1 /* always 1? */ }
	return true
}

func init() {
	CpuParams()
}
