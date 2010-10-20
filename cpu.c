// Copyright (c) 2010 Joseph D Poirier
// Distributable under the terms of The New BSD License
// that can be found in the LICENSE file.

#if defined(__WINDOWS__)
#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#elif defined(__LINUX__) || defined(__DARWIN__) || defined(__FREEBSD__)
#include <unistd.h>
#else
#error "Invalid GOOS: must be darwin, freebsd, linux, or windows"
#endif

#include "cpu.h"

/*
    - check if cpuid is executable
        Yes:
            - CPUID 0:
                    Vendor Information: EBX ECX EDX
                    GenuineIntel or AuthenticAMD
            - CPUID 1:
                    EDX[28]     : HTT enabled
                    EBX[23:16]  : logical core count

            - Logical core count:
                if GenuineIntel:
                    CPUID 4:
                        EAX[31:26] -> value + 1
                if AuthenticAMD
                    CPUID 0x80000008:
                        ECX[7:0]
        No:
            core cnt is 1

    r.ebx: 0x756E6547  uneG
    r.ecx: 0x6C65746E  letn
    r.edx: 0x49656E69  Ieni

    r.ebx: 0x68747541  htuA
    r.ecx: 0x444D4163  DMAc
    r.edx: 0x69746E65  itne
*/

#if defined(__386__) && !defined(__AMD64__)
/* eflag register checks, upper and lower boundaries */
#define CHK_386 (0x040000)
#define CHK_486 (0x200000)
uint32_t eflg_chks[2] = {CHK_386, CHK_486};
#endif

//
bool have_cpuid(void) {
#if defined(__386__) && !defined(__AMD64__)
    /*
        if can't flip ac bit (0x040000)
            386 or below, cpuid is not accessible
        else if can't flip cpuid (0x200000)
            older 486, cpuid is accessible but not executable
        else
            newer 486 or above, cpuid is executable
    */
    uint32_t a, b:
    int32_t j, i;
    for (i = 0; i < 2; i++) {
        j = eflg_chks[i];
        __asm__ __volatile__ (
            "pushfl\n\t"
            "popl %%eax\n\t"
            "movl %%eax, %0\n\t"
            "xorl %3, %%eax\n\t"
            "pushl %%eax\n\t"
            "popfl\n\t"
            "pushfl\n\t"
            "popl %%eax\n\t"
            "movl %%eax, %1\n\t"
            "pushl %0\n\t"
            "popfl\n"
            : "r="(a), "r="(b)
            : "r"(j)
            : "eax"
        );
        if ((a & j) != (b & j)
            return false;
    }
#endif
    return true;
}

//
void cpuid(regs_t* r, uint32_t f1, uint32_t f2) {
#if defined(__386__) && !defined(__AMD64__)
        "pushl %%ebx\n\t"
#endif
        "movl %4, %%eax\n\t"
        "movl %5, %%ecx\n\t"
        "cpuid\n\t"
        "movl %%eax, %0\n\t"
        "movl %%ebx, %1\n\t"
        "movl %%ecx, %2\n\t"
        "movl %%edx, %3\n\t"
#if defined(__386__) && !defined(__AMD64__)
        "popl %%ebx\n"
#endif
        : "=m"(r->eax), "=m"(r->ebx), "=m"(r->ecx), "=m"(r->edx)
        : "r"(f1), "r"(f2)
        : "eax",
#if defined(__AMD64__) && !defined(__386__)
          "ebx",
#endif
          "ecx", "edx", "cc", "memory"
    );
}

//  Number of online processors
uint32_t onln(void) {
#if defined(__WINDOWS__)
	return (uint32_t) conf();
#else
	return (uint32_t) sysconf(_SC_NPROCESSORS_ONLN);
#endif
}

//  Number of OS configured processors
uint32_t conf(void) {
#if defined(__WINDOWS__)
	SYSTEM_INFO sysinfo;
	GetSystemInfo(&sysinfo);
	return (uint32_t) sysinfo.dwNumberOfProcessors;
#else
	return (uint32_t) sysconf(_SC_NPROCESSORS_CONF);
#endif
}
