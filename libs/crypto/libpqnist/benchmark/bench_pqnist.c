/*
	Licensed to the Apache Software Foundation (ASF) under one
	or more contributor license agreements.  See the NOTICE file
	distributed with this work for additional information
	regarding copyright ownership.  The ASF licenses this file
	to you under the Apache License, Version 2.0 (the
	"License"); you may not use this file except in compliance
	with the License.  You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing,
	software distributed under the License is distributed on an
	"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
	KIND, either express or implied.  See the License for the
	specific language governing permissions and limitations
	under the License.
*/

/*
   Time the run through the flow of encrypting, ecapsulating and signing a message
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <amcl/utils.h>
#include <amcl/randapi.h>
#include <amcl/bls_BLS381.h>
#include <oqs/oqs.h>
#include <amcl/pqnist.h>
#include <time.h>

#define G2LEN 4*BFS_BLS381
#define SIGLEN BFS_BLS381+1

#define MIN_TIME 5.0
#define MIN_ITERS 100

int main()
{
    int i,rc;

    int iterations;
    clock_t start;
    double elapsed;

    // Seed value for CSPRNG
    char seed[PQNIST_SEED_LENGTH];
    octet SEED = {sizeof(seed),sizeof(seed),seed};

    csprng RNG;

    // AES Key
    char k[PQNIST_AES_KEY_LENGTH];
    octet K= {0,sizeof(k),k};

    // Initialization vectors
    char iv[PQNIST_AES_IV_LENGTH];
    octet IV= {sizeof(iv),sizeof(iv),iv};
    char iv2[PQNIST_AES_IV_LENGTH];
    octet IV2= {sizeof(iv2),sizeof(iv2),iv2};

    // Message to be sent to Bob
    char p[256];
    octet P = {0, sizeof(p), p};
    OCT_jstring(&P,"Hello Bob! This is a message from Alice");

    // Pad message
    int l = 16 - (P.len % 16);
    if (l < 16)
    {
        OCT_jbyte(&P,0,l);
    }

    // AES CBC ciphertext
    char c[256];
    octet C = {0, sizeof(c), c};

    // non random seed value
    for (i=0; i<PQNIST_SEED_LENGTH; i++) SEED.val[i]=i+1;

    // initialise random number generator
    CREATE_CSPRNG(&RNG,&SEED);

    // Generate 256 bit AES Key
    K.len=PQNIST_AES_KEY_LENGTH;
    generateRandom(&RNG,&K);

    // Generate SIKE and BLS keys

    // Bob's SIKE keys
    uint8_t SIKEpk[OQS_KEM_sike_length_public_key];
    uint8_t SIKEsk[OQS_KEM_sike_length_secret_key];

    // Alice's BLS keys
    char BLSsk[BGS_BLS381];
    char BLSpk[G2LEN];

    rc = pqnist_sike_keys(seed, SIKEpk, SIKEsk);
    if (rc)
    {
        fprintf(stderr, "ERROR pqnist_sike_keys rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    rc = pqnist_bls_keys(seed, BLSpk, BLSsk);
    if (rc)
    {
        fprintf(stderr, "ERROR pqnist_bls_keys rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    // BLS signature
    char S[SIGLEN];

    // SIKE encapsulated key
    uint8_t ek[OQS_KEM_sike_length_ciphertext];

    // Alice

    // Random initialization value
    generateRandom(&RNG,&IV);

    // Copy plaintext
    OCT_copy(&C,&P);

    // Encrypt plaintext
    iterations=0;
    start=clock();
    do
    {
        pqnist_aes_cbc_encrypt(K.val, K.len, IV.val, C.val, C.len);
        iterations++;
        elapsed=(clock()-start)/(double)CLOCKS_PER_SEC;
    }
    while (elapsed<MIN_TIME || iterations<MIN_ITERS);
    elapsed=1000000.0*elapsed/iterations;
    printf("pqnist_aes_cbc_encrypt              - %8d iterations  ",iterations);
    printf(" %8.2lf us per iteration\n",elapsed);

    generateRandom(&RNG,&IV2);

    // Generate an AES which is ecapsulated using SIKE. Use this key to
    // AES encrypt the K parameter.
    iterations=0;
    start=clock();
    do
    {
        pqnist_encapsulate_encrypt(K.val, K.len, IV2.val, SIKEpk, ek);
        iterations++;
        elapsed=(clock()-start)/(double)CLOCKS_PER_SEC;
    }
    while (iterations<MIN_ITERS);
    elapsed=1000.0*elapsed/iterations;
    printf("pqnist_encapsulate_encrypt              - %8d iterations  ",iterations);
    printf(" %8.2lf ms per iteration\n",elapsed);

    // Obtain encapsulated AES key and decrypt K
    iterations=0;
    start=clock();
    do
    {
        pqnist_decapsulate_decrypt(K.val, K.len, IV2.val, SIKEsk, ek);
        iterations++;
        elapsed=(clock()-start)/(double)CLOCKS_PER_SEC;
    }
    while (iterations<MIN_ITERS);
    elapsed=1000.0*elapsed/iterations;
    printf("pqnist_decapsulate_decrypt              - %8d iterations  ",iterations);
    printf(" %8.2lf ms per iteration\n",elapsed);


    iterations=0;
    start=clock();
    do
    {
        pqnist_aes_cbc_decrypt(K.val, K.len, IV.val, C.val, C.len);
        iterations++;
        elapsed=(clock()-start)/(double)CLOCKS_PER_SEC;
    }
    while (elapsed<MIN_TIME || iterations<MIN_ITERS);
    elapsed=1000000.0*elapsed/iterations;
    printf("pqnist_aes_cbc_decrypt              - %8d iterations  ",iterations);
    printf(" %8.2lf us per iteration\n",elapsed);



    // Alice signs message
    iterations=0;
    start=clock();
    do
    {
        pqnist_bls_sign(P.val, P.len, BLSsk, S);
        iterations++;
        elapsed=(clock()-start)/(double)CLOCKS_PER_SEC;
    }
    while (elapsed<MIN_TIME || iterations<MIN_ITERS);
    elapsed=1000.0*elapsed/iterations;
    printf("pqnist_bls_sign              - %8d iterations  ",iterations);
    printf(" %8.2lf ms per iteration\n",elapsed);

    // Bob verifies message
    iterations=0;
    start=clock();
    do
    {
        pqnist_bls_verify(P.val, P.len, BLSpk, S);
        iterations++;
        elapsed=(clock()-start)/(double)CLOCKS_PER_SEC;
    }
    while (elapsed<MIN_TIME || iterations<MIN_ITERS);
    elapsed=1000.0*elapsed/iterations;
    printf("pqnist_bls_verify              - %8d iterations  ",iterations);
    printf(" %8.2lf ms per iteration\n",elapsed);

    printf("BLS sk len %d\n", BGS_BLS381);
    printf("BLS pk len %d\n", G2LEN);
    printf("BLS sig len %d\n", SIGLEN);

    exit(EXIT_SUCCESS);
}
