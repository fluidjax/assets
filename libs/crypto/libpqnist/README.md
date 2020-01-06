# libpqnist - Crypto Document functions

This is a library that implements the cryptographic functions required for the DTA

The code is quantum safe; it uses the round two algorithms from the [NIST](https://csrc.nist.gov/projects/post-quantum-cryptography/round-2-submissions)
competition.

The algorithms that are used are;

* [AES-CBC](https://csrc.nist.gov/projects/block-cipher-techniques/bcm/current-modes) for encryption using a 256 bit key
* [SIKE](https://sike.org) for secret encapsulation
* [BLS](https://datatracker.ietf.org/doc/draft-boneh-bls-signature) for digital signature

## Dependencies

To correctly build the C library you need to install the following.

```
sudo add-apt-repository universe
sudo apt-get update
sudo apt-get install -y gcc g++ git cmake doxygen autoconf automake libtool curl make unzip wget libssl-dev xsltproc lcov emacs
```

### liboqs

[liboqs](https://github.com/open-quantum-safe/liboqs) is a C library for
quantum-resistant cryptographic algorithms. It is a API level on top of the
NIST round two submissions.

```
git clone https://github.com/open-quantum-safe/liboqs.git
cd liboqs
autoreconf -i
./configure --without-openssl --disable-kem-kyber  --disable-kem-newhope   --disable-kem-ntru --disable-kem-saber --disable-sig-dilithium  --disable-sig-mqdss --disable-sig-sphincs --disable-kem-bike --disable-kem-frodokem --disable-sig-picnic --disable-sig-qtesla
make clean
make -j
sudo make install
```

### AMCL

[AMCL](https://github.com/apache/incubator-milagro-crypto-c) is required

Build and install the AMCL library

```
git clone https://github.com/apache/incubator-milagro-crypto-c.git
cd incubator-milagro-crypto-c
mkdir build
cd build
cmake -D CMAKE_BUILD_TYPE=Release -D BUILD_SHARED_LIBS=ON -D AMCL_CHUNK=64 -D AMCL_CURVE="BLS381,SECP256K1" -D AMCL_RSA="" -D BUILD_PYTHON=OFF -D BUILD_BLS=ON -D BUILD_WCC=OFF -D BUILD_MPIN=OFF -D BUILD_X509=OFF -D CMAKE_INSTALL_PREFIX=/usr/local ..
make
make test
sudo make install
```

### golang

There is a golang wrapper in the ./go directory

```
wget https://dl.google.com/go/go1.12.linux-amd64.tar.gz
tar -xzf go1.12.linux-amd64.tar.gz
sudo cp -r go /usr/local
echo 'export PATH=$PATH:/usr/local/go/bin' >> ${HOME}/.bashrc
```

#### configure GO

```
mkdir -p ${HOME}/go/bin 
mkdir -p ${HOME}/go/pkg 
mkdir -p ${HOME}/go/src 
echo 'export GOPATH=${HOME}/go' >> ${HOME}/.bashrc 
echo 'export PATH=$GOPATH/bin:$PATH' >> ${HOME}/.bashrc
```

This package is needed for testing.

```
go get github.com/stretchr/testify/assert
```

## Compiling

Build and test code. 

```sh
export GOPATH=${HOME}/src
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:./
mkdir build
cd build
cmake -D CMAKE_BUILD_TYPE=Release -D BUILD_SHARED_LIBS=ON -D BUILD_SIKE_COMPRESS=OFF ..
make
make doc
make test
sudo make install
```

## Windows

See ./windows/README.md

## Documentation

The documentation is generated using doxygen and can accessed (post build)
via the file

```
./build/doxygen/html/index.html
```

## Docker

Build and run tests using docker

```
docker build -t libpqnist:builder .
docker run --cap-add SYS_PTRACE --rm libpqnist:builder
```

Generate coverage figures

```
docker run --rm libpqnist:builder ./scripts/coverage.sh
```

Build release image

```
docker build -t libpqnist -f Dockerfile.Release .
```

## Virtual machine

In "./vagrant" there are configuration files to run the software on a VM