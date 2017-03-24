
# Build

```bash
cd ~/gopath/src/github.com/hyperhq/hyperd
make
```

# Run hyperd with ntfs support

```bash
sudo ./hyperd --v=3 --host tcp://0.0.0.0:2375 --host unix:///var/run/hyperd.sock --storage-driver=devicemapper --storage-opt dm.fs=ntfs-3g --storage-opt dm.mkfsarg="-U" --storage-opt dm.mkfsarg="-p 2048" --storage-opt dm.mkfsarg="-f" --storage-opt dm.mountopt="offset=$((0x0e400000))" 
```
