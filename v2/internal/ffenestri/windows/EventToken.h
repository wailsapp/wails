

/* this ALWAYS GENERATED file contains the definitions for the interfaces */


 /* File created by MIDL compiler version 8.01.0622 */
/* @@MIDL_FILE_HEADING(  ) */



/* verify that the <rpcndr.h> version is high enough to compile this file*/
#ifndef __REQUIRED_RPCNDR_H_VERSION__
#define __REQUIRED_RPCNDR_H_VERSION__ 500
#endif

/* verify that the <rpcsal.h> version is high enough to compile this file*/
#ifndef __REQUIRED_RPCSAL_H_VERSION__
#define __REQUIRED_RPCSAL_H_VERSION__ 100
#endif

#include "rpc.h"
#include "rpcndr.h"

#ifndef __RPCNDR_H_VERSION__
#error this stub requires an updated version of <rpcndr.h>
#endif /* __RPCNDR_H_VERSION__ */


#ifndef __eventtoken_h__
#define __eventtoken_h__

#if defined(_MSC_VER) && (_MSC_VER >= 1020)
#pragma once
#endif

/* Forward Declarations */

#ifdef __cplusplus
extern "C"{
#endif


/* interface __MIDL_itf_eventtoken_0000_0000 */
/* [local] */

//  Microsoft Windows
//  Copyright (c) Microsoft Corporation. All rights reserved.
#pragma once
typedef struct EventRegistrationToken
    {
    __int64 value;
    } 	EventRegistrationToken;



extern RPC_IF_HANDLE __MIDL_itf_eventtoken_0000_0000_v0_0_c_ifspec;
extern RPC_IF_HANDLE __MIDL_itf_eventtoken_0000_0000_v0_0_s_ifspec;

/* Additional Prototypes for ALL interfaces */

/* end of Additional Prototypes */

#ifdef __cplusplus
}
#endif

#endif

