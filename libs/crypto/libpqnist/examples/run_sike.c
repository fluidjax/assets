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

/**
 * @file run_sike.c
 * @author Kealan McCusker
 * @brief Encapsulate and decapsulate a secret using SIKE
 */

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <amcl/pqnist.h>
#include <oqs/oqs.h>

// Set to have fixed values for the key pair
#define FIXED

/* Print encoded binary string in hex */
void qredo_print_hex(uint8_t *src, int src_len)
{
    int i;
    for (i = 0; i < src_len; i++)
    {
        printf("%02x", (unsigned char) src[i]);
    }
    printf("\n");
}

void cleanup_heap(uint8_t *secret_key, uint8_t *shared_secret_e,
                  uint8_t *shared_secret_d, uint8_t *public_key,
                  uint8_t *ciphertext, OQS_KEM *kem)
{
    if (kem != NULL)
    {
        OQS_MEM_secure_free(secret_key, kem->length_secret_key);
        OQS_MEM_secure_free(shared_secret_e, kem->length_shared_secret);
        OQS_MEM_secure_free(shared_secret_d, kem->length_shared_secret);
    }
    OQS_MEM_insecure_free(public_key);
    OQS_MEM_insecure_free(ciphertext);
    OQS_KEM_free(kem);
}

OQS_STATUS test(const char *method_name)
{

    OQS_KEM *kem = NULL;
    uint8_t *public_key = NULL;
    uint8_t *secret_key = NULL;
    uint8_t *ciphertext = NULL;
    uint8_t *shared_secret_e = NULL;
    uint8_t *shared_secret_d = NULL;
    OQS_STATUS rc;
    int rv;

    kem = OQS_KEM_new(method_name);
    if (kem == NULL)
    {
        return OQS_ERROR;
    }

    printf("================================================================================\n");
    printf("Sample computation for KEM %s\n", kem->method_name);
    printf("================================================================================\n");

    public_key = malloc(kem->length_public_key);
    secret_key = malloc(kem->length_secret_key);
    ciphertext = malloc(kem->length_ciphertext);
    shared_secret_e = malloc(kem->length_shared_secret);
    shared_secret_d = malloc(kem->length_shared_secret);

    if ((public_key == NULL) || (secret_key == NULL) || (ciphertext == NULL) || (shared_secret_e == NULL) || (shared_secret_d == NULL))
    {
        fprintf(stderr, "ERROR: malloc failed\n");
        cleanup_heap(secret_key, shared_secret_e, shared_secret_d, public_key, ciphertext, kem);
        return OQS_ERROR;
    }

    rc = OQS_KEM_keypair(kem, public_key, secret_key);
    if (rc != OQS_SUCCESS)
    {
        fprintf(stderr, "ERROR: OQS_KEM_keypair failed\n");
        cleanup_heap(secret_key, shared_secret_e, shared_secret_d, public_key, ciphertext, kem);
        return OQS_ERROR;
    }
    printf("pk: ");
    qredo_print_hex(public_key, kem->length_public_key);
    printf("sk: ");
    qredo_print_hex(secret_key, kem->length_secret_key);

    rc = OQS_KEM_encaps(kem, ciphertext, shared_secret_e, public_key);
    if (rc != OQS_SUCCESS)
    {
        fprintf(stderr, "ERROR: OQS_KEM_encaps failed\n");
        cleanup_heap(secret_key, shared_secret_e, shared_secret_d, public_key, ciphertext, kem);
        return OQS_ERROR;
    }
    printf("ciphertext: ");
    qredo_print_hex(ciphertext, kem->length_ciphertext);
    printf("shared_secret_e: ");
    qredo_print_hex(shared_secret_e, kem->length_shared_secret);

    rc = OQS_KEM_decaps(kem, shared_secret_d, ciphertext, secret_key);
    if (rc != OQS_SUCCESS)
    {
        fprintf(stderr, "ERROR: OQS_KEM_decaps failed\n");
        cleanup_heap(secret_key, shared_secret_e, shared_secret_d, public_key, ciphertext, kem);
        return OQS_ERROR;
    }
    printf("shared_secret_d: ");
    qredo_print_hex(shared_secret_d, kem->length_shared_secret);

    rv = memcmp(shared_secret_e, shared_secret_d, kem->length_shared_secret);
    if (rv != 0)
    {
        fprintf(stderr, "ERROR: shared secrets are not equal\n");
        cleanup_heap(secret_key, shared_secret_e, shared_secret_d, public_key, ciphertext, kem);
        return OQS_ERROR;
    }
    else
    {
        printf("shared secrets are equal\n");
    }

    printf("public_key len %ld\n", kem->length_public_key);
    printf("secret_key len %ld\n", kem->length_secret_key);
    printf("ciphertext len %ld\n", kem->length_ciphertext);
    printf("shared_secret len %ld\n", kem->length_shared_secret);

    cleanup_heap(secret_key, shared_secret_e, shared_secret_d, public_key, ciphertext, kem);

    return OQS_SUCCESS;
}

int main()
{
    OQS_STATUS rc;

#ifdef FIXED
    // Set RNG to known value
    uint8_t entropy_input[48];

    for (size_t i = 0; i < 48; i++)
    {
        entropy_input[i] = i;
    }

    rc = OQS_randombytes_switch_algorithm(OQS_RAND_alg_nist_kat);
    if (rc != OQS_SUCCESS)
    {
        return EXIT_FAILURE;
    }
    OQS_randombytes_nist_kat_init(entropy_input, NULL, 256);
#else
    // Use system RNG
    OQS_randombytes_switch_algorithm(OQS_RAND_alg_system);
#endif

    char *alg_name = ALG_NAME;
    //if (!OQS_KEM_alg_is_enabled(alg_name))
    rc = !OQS_KEM_alg_is_enabled(alg_name);
    if (rc != OQS_SUCCESS)
    {
        printf("error algorithm not defined\n");
        return EXIT_FAILURE;
    }

    rc = test(alg_name);
    if (rc != OQS_SUCCESS)
    {
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}
