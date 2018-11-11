# Chroot
```
apt install qemu-arm-static
mount -o loop disk.img image
cp -r ./image/* ./root/
cp $(which qemu-arm-static) ./root/usr/bin
cp /etc/resolv.conf ./root/etc/resolv.conf
mount -t proc proc ./root/proc
mount -t sysfs sys ./root//sys
mount -o bind /dev ./root/dev
chroot ./root qemu-arm-static /bin/bash
```

# Alsa
```
apt update
apt install build-essential libtool curl automake
cd /tmp
curl ftp://ftp.alsa-project.org/pub/lib/alsa-lib-1.1.7.tar.bz2 | tar xvj
cd alsa-lib-1.1.7
libtoolize --force --copy --automake
aclocal
autoheader
automake --foreign --copy --add-missing
autoconf
./configure --enable-shared=no --enable-static=yes --host=arm-linux-gnueabihf
make -j12
export DESTDIR=${SRCPATH}/libs/arm
make install
```

# UPX
```
curl -L https://github.com/upx/upx/releases/download/v3.95/upx-3.95-amd64_linux.tar.xz | tar xJ
```