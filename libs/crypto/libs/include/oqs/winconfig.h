#ifndef __WINCONFIG_H
#define __WINCONFIG_H

/* Enable schemes supported on Windows */
///// OQS_COPY_FROM_PQCLEAN_FRAGMENT_KEMS_START
#define OQS_ENABLE_KEM_sidh_p751
#define OQS_ENABLE_KEM_sidh_p751_compressed
#define OQS_ENABLE_KEM_sike_p751
#define OQS_ENABLE_KEM_sike_p751_compressed
#define OQS_KEM_DEFAULT OQS_KEM_alg_sike_p751
#define OQS_SIG_DEFAULT OQS_SIG_alg_picnic_L1_FS

#define OQS_MASTER_BRANCH /**/
#define OQS_VERSION_NUMBER 0x00201000L
#define OQS_VERSION_TEXT "0.2.1-dev"

#endif
