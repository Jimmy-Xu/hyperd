#!/bin/bash +x
set -e
#set -x

if [ $# -ne 1 ];then
  echo "./start_nano.sh <pod-name>"
  exit 1
fi

BASE_DIR=$(dirname $(readlink -f $0))
ISO_DIR="/home/osboxes/iso"
BOOT_DISK="${ISO_DIR}/NanoBoot.raw"
POD_ID=$1
DM_DEV=`${BASE_DIR}/../common/get_pod_dm.py ${POD_ID} device`
if [ $? -ne 0 ];then
  echo "${POD_ID} doesn't exist"
  exit 1
fi

# umount old device
${BASE_DIR}/umount_dm.sh ${POD_ID} >/dev/null 2>&1
if [ $? -ne 0 ];then
   echo "umount ${POD_ID} failed"
   exit 1
fi


VM_ID=`${BASE_DIR}/../common/get_pod_dm.py ${POD_ID} vmid`

#UEFI boot, with network adaptor and seiral port
cd ${ISO_DIR}
#echo "Start NanoServer VM... (vnc port: 5908)"
#cat <<EOF
#qemu-system-x86_64 -enable-kvm -smp 1 -m 1024 
#  -bios /usr/share/edk2.git/ovmf-x64/OVMF-pure-efi.fd 
#  -netdev tap,id=network0,ifname=tap1,script=no,downscript=no
#  -device e1000,netdev=network0,mac=00:16:35:AF:94:4B
#  -machine vmport=off 
#  -boot order=c,menu=off 
#  -drive file=${BOOT_DISK},format=raw 
#  -drive file=${DM_DEV},format=raw 
#  -vnc :8 
#  -chardev stdio,id=mon0 
#  -mon chardev=mon0 
#  -serial unix:/var/run/hyper/${VM_ID}/win_ctl.sock,server,nowait 
#  -serial unix:/var/run/hyper/${VM_ID}/win_tty.sock,server,nowait
#EOF

echo ${POD_ID}
qemu-system-x86_64 -enable-kvm -smp 1 -m 1024 \
  -bios /usr/share/edk2.git/ovmf-x64/OVMF-pure-efi.fd \
  -machine vmport=off \
  -boot order=c,menu=off \
  -drive file=${BOOT_DISK},format=raw \
  -drive file=${DM_DEV},format=raw \
  -vnc :8 \
  -qmp unix:/var/run/hyper/${VM_ID}/win_qmp.sock,server,nowait \
  -serial unix:/var/run/hyper/${VM_ID}/win_ctl.sock,server,nowait \
  -serial unix:/var/run/hyper/${VM_ID}/win_tty.sock,server,nowait &

