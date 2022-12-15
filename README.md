# Kubernetes config generator
Generate kubernete configuration including staging and production deployment files. It consists of `nginx-ingress`, `kustomize`, etc. Generate production ready configuration in minutes.
# Installation
## For Ubuntu
* Download .deb file from [releases](https://github.com/pravinbanjade/k8s-config-generator/releases) and install it.

```
wget https://github.com/pravinbanjade/k8s-config-generator/releases/download/v0.0.6/k8s-config-generator_0.0.6_linux_386.deb

sudo dpkg -i k8s-config-generator_0.0.6_linux_386.deb
```
Note: Please replace replace the version with the lastest version from [releases page](https://github.com/pravinbanjade/k8s-config-generator/releases)

## For Mac OS
### M1 and M2 chip (M chip series)
Download the binary and install it into your `/usr/local/bin` directory
```
curl -L -o kcg.tar.gz https://github.com/pravinbanjade/k8s-config-generator/releases/download/v0.0.6/k8s-config-generator_0.0.6_Linux_arm64.tar.gz
```
Now extract the tar file
```
tar -xzvf kcg.tar.gz
```
move `kcg` binary to `/usr/local/bin` to use the command from anywhere
```
sudo mv kcg /usr/local/bin/kcg
```
### For intel series - mac os
Download the binary and install it into your `/usr/local/bin` directory
```
curl -L -o kcg.tar.gz https://github.com/pravinbanjade/k8s-config-generator/releases/download/v0.0.6/k8s-config-generator_0.0.6_Linux_x86_64.tar.gz
```
Now extract the tar file
```
tar -xzvf kcg.tar.gz
```
move `kcg` binary to `/usr/local/bin` to use the command from anywhere
```
sudo mv kcg /usr/local/bin/kcg
```
Note: Please replace replace the version with the lastest version from [releases page](https://github.com/pravinbanjade/k8s-config-generator/releases)

## For Windows

Why you using windows bro? Switch to linux

Haha just kidding!

* Download tar file from [Here](https://github.com/pravinbanjade/k8s-config-generator/releases/download/v0.0.6/k8s-config-generator_0.0.6_Windows_x86_64.tar.gz) and extract and run .exe file

# USAGE
To generate kubernetes config just run `kcg` command and follow the interactive terminal.
Your files will be generated in the working directory in your specified `appName` folder.
```bash
kcg
```
