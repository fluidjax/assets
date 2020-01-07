# CMake generated Testfile for 
# Source directory: /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/test/unit
# Build directory: /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/test/unit
# 
# This file includes the relevant testing commands required for 
# testing this directory and lists subdirectories to be tested as well.
add_test(test_aes_encrypt_CBC_256 "/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/test/unit/test_aes_encrypt_CBC_256" "aes/CBCMMT256.rsp")
set_tests_properties(test_aes_encrypt_CBC_256 PROPERTIES  PASS_REGULAR_EXPRESSION "SUCCESS" WORKING_DIRECTORY "/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/testVectors" _BACKTRACE_TRIPLES "/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/test/unit/CMakeLists.txt;30;add_test;/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/test/unit/CMakeLists.txt;40;amcl_test;/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/test/unit/CMakeLists.txt;0;")
add_test(test_aes_decrypt_CBC_256 "/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/test/unit/test_aes_decrypt_CBC_256" "aes/CBCMMT256.rsp")
set_tests_properties(test_aes_decrypt_CBC_256 PROPERTIES  PASS_REGULAR_EXPRESSION "SUCCESS" WORKING_DIRECTORY "/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/testVectors" _BACKTRACE_TRIPLES "/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/test/unit/CMakeLists.txt;30;add_test;/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/test/unit/CMakeLists.txt;41;amcl_test;/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/test/unit/CMakeLists.txt;0;")
