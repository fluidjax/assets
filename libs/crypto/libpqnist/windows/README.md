# Windows build

Instructions for how to build the library on Windows OS.

There are three shells used; git, bash and windows admin

## Dependencies

To build the C library you need to install the following.

* install [MinGW 64bit](http://mingw-w64.org/doku.php/download/mingw-builds) Please select the **Architecture** - *x86_64*
Add *C\<install_path>\bin* to the PATH variable
* install [CMake](http://www.cmake.org)
* install Git

### AMCL

[AMCL](https://github.com/apache/incubator-milagro-crypto-c) is required

Build the AMCL library

```sh
mkdir lib
git clone https://github.com/apache/incubator-milagro-crypto-c.git
cd incubator-milagro-crypto-c
mkdir build
cd build
cmake -G "MinGW Makefiles" -D WORD_SIZE=64 -D BUILD_SHARED_LIBS=OFF -D AMCL_CHUNK=64 -D AMCL_CURVE="BLS381,SECP256K1" -D AMCL_RSA="" -D BUILD_PYTHON=OFF -D BUILD_BLS=ON -D BUILD_WCC=OFF -D BUILD_MPIN=OFF -D BUILD_X509=OFF ..
mingw32-make
mingw32-make test
cp  ./lib/*.a ../../lib
```

### liboqs

[liboqs](https://github.com/open-quantum-safe/liboqs) is a C library for
quantum-resistant cryptographic algorithms. It is a API level on top of the
NIST round two submissions.

```sh
git clone https://github.com/open-quantum-safe/liboqs.git
cp -r include ./liboqs
cp build.bat ./liboqs\
cp common.c ./liboqs/src/common
cd liboqs
build.bat
cp liboqs.a ../lib
```

## Compiling

Build and test code. 

```sh
git checkout nogolang
cp -r ./windows/include/* ./include
mkdir build
cd build
cmake -G "MinGW Makefiles" -D CMAKE_BUILD_TYPE=Release -D BUILD_SHARED_LIBS=OFF -D BUILD_SIKE_COMPRESS=OFF ..
cp ../windows/lib/* ./src
mingw32-make
mingw32-make test
```
