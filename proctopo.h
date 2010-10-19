// Copyright (c) 2010 Joseph D Poirier
// Distributable under the terms of The New BSD License
// that can be found in the LICENSE file.

#ifndef PROCTOPO_H
#define PROCTOPO_H

#include <stdbool.h>
#include <stdint.h>

typedef struct {
    uint32_t eax;
    uint32_t ebx;
    uint32_t ecx;
    uint32_t edx;
} regs_t;

extern bool have_cpuid(void);
extern void cpuid(regs_t* r, uint32_t f1, uint32_t f2);
extern int conf(void);
extern int onln(void);

#endif // PROCTOPO_H
