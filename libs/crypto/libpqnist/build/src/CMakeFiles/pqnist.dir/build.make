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
include src/CMakeFiles/pqnist.dir/depend.make

# Include the progress variables for this target.
include src/CMakeFiles/pqnist.dir/progress.make

# Include the compile flags for this target's objects.
include src/CMakeFiles/pqnist.dir/flags.make

src/CMakeFiles/pqnist.dir/pqnist.c.o: src/CMakeFiles/pqnist.dir/flags.make
src/CMakeFiles/pqnist.dir/pqnist.c.o: ../src/pqnist.c
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --progress-dir=/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/CMakeFiles --progress-num=$(CMAKE_PROGRESS_1) "Building C object src/CMakeFiles/pqnist.dir/pqnist.c.o"
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/src && /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -o CMakeFiles/pqnist.dir/pqnist.c.o   -c /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/src/pqnist.c

src/CMakeFiles/pqnist.dir/pqnist.c.i: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Preprocessing C source to CMakeFiles/pqnist.dir/pqnist.c.i"
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/src && /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -E /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/src/pqnist.c > CMakeFiles/pqnist.dir/pqnist.c.i

src/CMakeFiles/pqnist.dir/pqnist.c.s: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Compiling C source to assembly CMakeFiles/pqnist.dir/pqnist.c.s"
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/src && /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -S /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/src/pqnist.c -o CMakeFiles/pqnist.dir/pqnist.c.s

# Object files for target pqnist
pqnist_OBJECTS = \
"CMakeFiles/pqnist.dir/pqnist.c.o"

# External object files for target pqnist
pqnist_EXTERNAL_OBJECTS =

src/libpqnist.2.0.0.dylib: src/CMakeFiles/pqnist.dir/pqnist.c.o
src/libpqnist.2.0.0.dylib: src/CMakeFiles/pqnist.dir/build.make
src/libpqnist.2.0.0.dylib: src/CMakeFiles/pqnist.dir/link.txt
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --bold --progress-dir=/Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/CMakeFiles --progress-num=$(CMAKE_PROGRESS_2) "Linking C shared library libpqnist.dylib"
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/src && $(CMAKE_COMMAND) -E cmake_link_script CMakeFiles/pqnist.dir/link.txt --verbose=$(VERBOSE)
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/src && $(CMAKE_COMMAND) -E cmake_symlink_library libpqnist.2.0.0.dylib libpqnist.2.dylib libpqnist.dylib
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/src && /usr/local/Cellar/cmake/3.15.5/bin/cmake -E copy /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/src/lib* /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/go/

src/libpqnist.2.dylib: src/libpqnist.2.0.0.dylib
	@$(CMAKE_COMMAND) -E touch_nocreate src/libpqnist.2.dylib

src/libpqnist.dylib: src/libpqnist.2.0.0.dylib
	@$(CMAKE_COMMAND) -E touch_nocreate src/libpqnist.dylib

# Rule to build all files generated by this target.
src/CMakeFiles/pqnist.dir/build: src/libpqnist.dylib

.PHONY : src/CMakeFiles/pqnist.dir/build

src/CMakeFiles/pqnist.dir/clean:
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/src && $(CMAKE_COMMAND) -P CMakeFiles/pqnist.dir/cmake_clean.cmake
.PHONY : src/CMakeFiles/pqnist.dir/clean

src/CMakeFiles/pqnist.dir/depend:
	cd /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build && $(CMAKE_COMMAND) -E cmake_depends "Unix Makefiles" /Users/chris/dev/qredo/assets/libs/crypto/libpqnist /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/src /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/src /Users/chris/dev/qredo/assets/libs/crypto/libpqnist/build/src/CMakeFiles/pqnist.dir/DependInfo.cmake --color=$(COLOR)
.PHONY : src/CMakeFiles/pqnist.dir/depend

