# CMAKE generated file: DO NOT EDIT!
# Generated by "Unix Makefiles" Generator, CMake Version 3.15

# Delete rule output on recipe failure.
.DELETE_ON_ERROR:


#=============================================================================
# Special targets provided by cmake.

# Disable implicit rules so canonical targets will work.
.SUFFIXES:


# Remove some rules from gmake that .SUFFIXES does not remove.
SUFFIXES =

.SUFFIXES: .hpux_make_needs_suffix_list


# Suppress display of executed commands.
$(VERBOSE).SILENT:


# A target that is always out of date.
cmake_force:

.PHONY : cmake_force

#=============================================================================
# Set environment variables for the build.

# The shell in which to execute make rules.
SHELL = /bin/sh

# The CMake executable.
CMAKE_COMMAND = /usr/local/Cellar/cmake/3.15.5/bin/cmake

# The command to remove a file.
RM = /usr/local/Cellar/cmake/3.15.5/bin/cmake -E remove -f

# Escaping for special characters.
EQUALS = =

# The top-level source directory on which CMake was run.
CMAKE_SOURCE_DIR = /Users/chris/dev/qredo/assets/libs/crypto/libpqnist

# The top-level build directory on which CMake was run.
CMAKE_BINARY_DIR = /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build

# Include any dependencies generated for this target.
include test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/depend.make

# Include the progress variables for this target.
include test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/progress.make

# Include the compile flags for this target's objects.
include test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/flags.make

test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/test_aes_encrypt.c.o: test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/flags.make
test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/test_aes_encrypt.c.o: ../test/unit/test_aes_encrypt.c
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --progress-dir=/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/CMakeFiles --progress-num=$(CMAKE_PROGRESS_1) "Building C object test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/test_aes_encrypt.c.o"
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/test/unit && /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -o CMakeFiles/test_aes_encrypt_CBC_256.dir/test_aes_encrypt.c.o   -c /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/test/unit/test_aes_encrypt.c

test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/test_aes_encrypt.c.i: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Preprocessing C source to CMakeFiles/test_aes_encrypt_CBC_256.dir/test_aes_encrypt.c.i"
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/test/unit && /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -E /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/test/unit/test_aes_encrypt.c > CMakeFiles/test_aes_encrypt_CBC_256.dir/test_aes_encrypt.c.i

test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/test_aes_encrypt.c.s: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Compiling C source to assembly CMakeFiles/test_aes_encrypt_CBC_256.dir/test_aes_encrypt.c.s"
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/test/unit && /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -S /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/test/unit/test_aes_encrypt.c -o CMakeFiles/test_aes_encrypt_CBC_256.dir/test_aes_encrypt.c.s

# Object files for target test_aes_encrypt_CBC_256
test_aes_encrypt_CBC_256_OBJECTS = \
"CMakeFiles/test_aes_encrypt_CBC_256.dir/test_aes_encrypt.c.o"

# External object files for target test_aes_encrypt_CBC_256
test_aes_encrypt_CBC_256_EXTERNAL_OBJECTS =

test/unit/test_aes_encrypt_CBC_256: test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/test_aes_encrypt.c.o
test/unit/test_aes_encrypt_CBC_256: test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/build.make
test/unit/test_aes_encrypt_CBC_256: src/libpqnist.2.0.0.dylib
test/unit/test_aes_encrypt_CBC_256: test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/link.txt
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --bold --progress-dir=/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/CMakeFiles --progress-num=$(CMAKE_PROGRESS_2) "Linking C executable test_aes_encrypt_CBC_256"
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/test/unit && $(CMAKE_COMMAND) -E cmake_link_script CMakeFiles/test_aes_encrypt_CBC_256.dir/link.txt --verbose=$(VERBOSE)

# Rule to build all files generated by this target.
test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/build: test/unit/test_aes_encrypt_CBC_256

.PHONY : test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/build

test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/clean:
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/test/unit && $(CMAKE_COMMAND) -P CMakeFiles/test_aes_encrypt_CBC_256.dir/cmake_clean.cmake
.PHONY : test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/clean

test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/depend:
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build && $(CMAKE_COMMAND) -E cmake_depends "Unix Makefiles" /Users/chris/dev/qredo/assets/libs/crypto/libpqnist /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/test/unit /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/test/unit /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/DependInfo.cmake --color=$(COLOR)
.PHONY : test/unit/CMakeFiles/test_aes_encrypt_CBC_256.dir/depend

