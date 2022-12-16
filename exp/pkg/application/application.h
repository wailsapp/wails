//go:build darwin

#ifndef application_h
#define application_h

static void init(void);
static void run(void);
static void setActivationPolicy(int policy);
static char *getAppName(void);

#endif