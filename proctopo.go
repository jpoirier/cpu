// Copyright (c) 2010 Joseph D Poirier
// Distributable under the terms of The New BSD License
// that can be found in the LICENSE file.

package proctopo

// #include "proctopo.h"
import "C"

import (
	"fmt"
	"unsafe"
)

type ProcTopo_t struct {
	CpuidPresent    bool
	CpuidRestricted bool
	HttSupported    bool
	Vendor          string
	MaxCpus         int
	OSProcCnt       uint32
	PkgCntEnum      uint32
	CoreCntEnum     uint32
	ThreadCntEnum   uint32
	HttSmtPerCore   uint32
	HttSmtPerPkg    uint32
	Error           bool
}

type regs struct {
	eax uint32
	ebx uint32
	ecx uint32
	edx uint32
}

// Onln returns the number of logical processors that are on line.
func Onln() int {
	return int(C.onln())
}

// Conf returns the maximum number of logical processors supported by the OS.
func Conf() int {
	return int(C.conf())
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

// cpuParams
func cpuParams(s *ProcTopo_t) bool {
	s.MaxCpus = Conf()

	//----------------------------
	// cpuid check
	//----------------------------
	s.CpuidPresent = have_cpuid()
	if !s.CpuidPresent {
		return false
	}

	//----------------------------
	// vendor name
	//----------------------------
	var info regs
	cpuid(&info, 0, 0)
	maxCpuid := info.eax
	s.Vendor = utos(info.ebx) + utos(info.edx) + utos(info.ecx)

	//----------------------------
	// restricted cpuid execution
	//----------------------------
	var r regs
	s.CpuidRestricted = false
	cpuid(&r, 0x80000000, 0)
	if maxCpuid<=4 && r.eax>0x80000004 {
		s.CpuidRestricted = true
		return false
	}

	//----------------------------
	// htt enabled
	//----------------------------
	s.HttSupported = false
	cpuid(&r, 1, 0)
	if r.edx>>28&1 != 0 {
		s.HttSupported = true
	}

	return true
}

// ProcTopo
func ProcTopo(s *ProcTopo_t) {
	stat := cpuParams(s)
}
