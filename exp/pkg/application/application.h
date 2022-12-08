//go:build darwin

#ifndef application_h
#define application_h

void Init(void);
void Run(void);
void SetActivationPolicy(int policy);

#endif