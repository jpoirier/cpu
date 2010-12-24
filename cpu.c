// Copyright (c) 2010 Joseph D Poirier
// Distributable under the terms of The New BSD License
// that can be found in the LICENSE file.

#if defined(__WINDOWS__)
# define WIN32_LEAN_AND_MEAN
# include <windows.h>
#elif defined(__LINUX__) || defined(__DARWIN__) || defined(__FREEBSD__)
# include <stdlib.h>
# include <sys/param.h>
# include <sys/sysctl.h>
# include <unistd.h>
# ifndef _SC_NPROCESSORS_ONLN
#  define _SC_NPROCESSORS_ONLN (-1)
# endif
# ifndef _SC_NPROCESSORS_CONF
#  define _SC_NPROCESSORS_CONF (-1)
# endif
#else
# error "Invalid GOOS: must be darwin, freebsd, linux, or windows"
#endif

#include "cpu.h"

#if defined(__DARWIN__)
# define MIB_0   CTL_HW
# define MIB_1   HW_AVAILCPU
#elif defined(__LINUX__) || defined(__FREEBSD__)
# if defined(CTL_HW) && defined(HW_NCPU)
#  define MIB_0   CTL_HW
#  define MIB_1   HW_NCPU
# endif
#endif


bool have_cpuid(void) {
    uint32_t a = true;
#if defined(__386__) && !defined(__AMD64__)
    __asm__ __volatile__ (
        "pushfl\n\t"
        "popl %%eax\n\t"
        "movl %%eax, %%ecx\n\t"
        "xorl $0x200000, %%eax\n\t"
        "pushl %%eax\n\t"
        "popfl\n\t"
        "pushfl\n\t"
        "popl %%eax\n\t"
        "xorl %%ecx, %%eax\n\t"
        "shrl $21, %%eax\n\t"
        "andl $1, %%eax\n\t"
        "movl %%eax, %0\n\t"
        : "=r"(a)
        :
        : "eax", "ecx"
    );
#endif
    return a;
}

void cpuid(regs_t* r, uint32_t f1, uint32_t f2) {
    __asm__ __volatile__ (
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
        "popl %%ebx\n\t"
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

uint32_t onlineProcs(void) {
#if defined(__WINDOWS__)
	return (uint32_t) conf();
#else
    int x; uint32_t cnt; size_t sz = sizeof(cnt);
# if defined(MIB_0) && defined(MIB_1)
    int mib[2] = {MIB_0, MIB_1};
# endif
    if ((x = sysconf(_SC_NPROCESSORS_ONLN)) != -1) {
        return (uint32_t) x;
    }
# if defined(MIB_0) && defined(MIB_1)
    if ((x = sysctl(mib, 2, &cnt, &sz, NULL, 0)) != -1 ) {
        return (uint32_t) x;
    }
# endif
    if ((x = sysctlbyname("hw.ncpu", &cnt, &sz, NULL, 0)) != -1 ) {
        return (uint32_t) x;
    }
# if defined(MIB_0) && defined(MIB_1)
    if ((x = sysctlnametomib("hw.ncpu", mib, &sz)) != -1 ) {
        return (uint32_t) x;
    }
# endif
	return 0;
#endif
}

//  Number of OS configured processors
uint32_t confProcs(void) {
#if defined(__WINDOWS__)
	SYSTEM_INFO sysinfo;
	GetSystemInfo(&sysinfo);
	return (uint32_t) sysinfo.dwNumberOfProcessors;
#else
    int x;
    if ((x = sysconf(_SC_NPROCESSORS_CONF)) == -1) {
        x = onlineProcs();
    }
	return (uint32_t) x;
#endif
}
