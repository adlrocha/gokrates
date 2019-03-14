#!/bin/sh
echo "[*] Installing Zokrates base image"
git clone https://github.com/JacobEberhardt/ZoKrates
cd ZoKrates
docker build -t zokrates .