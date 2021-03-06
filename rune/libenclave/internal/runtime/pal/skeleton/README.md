The files in this directory are used to implement a skeleton enclave runtime in order to help to write your own enclave runtime.

Note that this code base is inspired by [v28 SGX in-tree driver](https://patchwork.kernel.org/patch/11418925/).

---

# Run skeleton with Docker
## Install sgx-tools
Refer to [this guide](https://github.com/alibaba/inclavare-containers/tree/master/sgx-tools/README.md).

Note that this step is only required when using SGX out-of-tree driver.

## Build liberpal-skeleton
```shell
cd "${path_to_inclavare_containers}/rune/libenclave/internal/runtime/pal/skeleton"
make
cp liberpal-skeleton-v*.so /usr/lib
```

## Build skeleton container image
```shell
cd "${path_to_inclavare_containers}/rune/libenclave/internal/runtime/pal/skeleton"
cat >Dockerfile <<EOF
FROM centos:7.5.1804

RUN mkdir -p /run/rune
WORKDIR /run/rune

COPY encl.bin .
COPY encl.elf .
COPY encl.ss .
# if any
COPY encl.token .
EOF
docker build . -t skeleton-enclave
```

## Build and install rune
Please refer to [this guide](https://github.com/alibaba/inclavare-containers#rune) to build `rune` from scratch.

---

# Run skeleton container image
## Configure OCI runtime
Add the `rune` OCI runtime configuration in dockerd config file, e.g, `/etc/docker/daemon.json`, in your system.

```json
{
	"runtimes": {
		"rune": {
			"path": "/usr/local/sbin/rune",
			"runtimeArgs": []
		}
	}
}
```

then restart dockerd on your system.
> e.g. `sudo systemctl restart docker` for CentOS, or `sudo service docker restart` for Ubuntu

You can check whether `rune` is correctly picked as supported OCI runtime or not with
```shell
docker info | grep rune
Runtimes: rune runc
```

## Run skeleton container image with rune
Note that replace `${SKELETON_PAL_VERSION}` with the actual version number. Currently skeleton supports PAL API v1 and v2.

```shell
docker run -it --rm --runtime=rune \
  -e ENCLAVE_TYPE=intelSgx \
  -e ENCLAVE_RUNTIME_PATH=/usr/lib/liberpal-skeleton-v${SKELETON_PAL_VERSION}.so \
  -e ENCLAVE_RUNTIME_ARGS="debug" \
  skeleton-enclave
```

where:
- @ENCLAVE_TYPE: specify the type of enclave hardware to use, such as `intelSgx`.
- @ENCLAVE_PATH: specify the path to enclave runtime to launch.
- @ENCLAVE_ARGS: specify the specific arguments to enclave runtime, seperated by the comma.

---

# Run skeleton OCI bundle
Note: The following method to launch skeleton with `rune` is usually provided for developmemt purpose.

## Create skeleton bundle
In order to use `rune` you must have your container image in the format of an OCI bundle. If you have Docker installed you can use its `export` method to acquire a root filesystem from an existing skeleton Docker container image.

```shell
# create the top most bundle directory
cd "$HOME/rune_workdir"
mkdir rune-container
cd rune-container

# create the rootfs directory
mkdir rootfs

# export skeleton image via Docker into the rootfs directory
docker export $(docker create skeleton-enclave) | sudo tar -C rootfs -xvf -
```

After a root filesystem is populated you just generate a spec in the format of a config.json file inside your bundle. `rune` provides a spec command which is similar to `runc` to generate a template file that you are then able to edit.

```shell
rune spec
```

To find features and documentation for fields in the spec please refer to the [specs](https://github.com/opencontainers/runtime-spec) repository.

In order to run the skeleton bundle with `rune`, you need to configure enclave runtime as following:
```json
  "annotations": {
      "enclave.type": "intelSgx",
      "enclave.runtime.path": "/usr/lib/liberpal-skeleton-v${SKELETON_PAL_VERSION}.so",
      "enclave.runtime.args": "debug"
  }
```

where:
- @enclave.type: specify the type of enclave hardware to use, such as intelSgx.
- @enclave.runtime.path: specify the path to enclave runtime to launch.
- @enclave.runtime.args: specify the specific arguments to enclave runtime, seperated by the comma.

## Run skeleton
Assuming you have an OCI bundle from the previous step you can execute the container in this way.

```shell
cd "$HOME/rune_workdir/rune-container"
sudo rune run skeleton-enclave-container
```
