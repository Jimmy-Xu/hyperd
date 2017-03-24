package mount

import (
	"time"

	"github.com/Sirupsen/logrus"
)

// GetMounts retrieves a list of mounts for the current running process.
func GetMounts() ([]*Info, error) {
	return parseMountTable()
}

// Mounted looks at /proc/self/mountinfo to determine of the specified
// mountpoint has been mounted
func Mounted(mountpoint string) (bool, error) {
	entries, err := parseMountTable()
	if err != nil {
		return false, err
	}

	// Search the table for the mountpoint
	for _, e := range entries {
		if e.Mountpoint == mountpoint {
			return true, nil
		}
	}
	return false, nil
}

// Mount will mount filesystem according to the specified configuration, on the
// condition that the target path is *not* already mounted. Options must be
// specified like the mount or fstab unix commands: "opt1=val1,opt2=val2". See
// flags.go for supported option flags.
func Mount(device, target, mType, options string) error {
	logrus.Infof("[Mount] begin: device:%v target:%v mType:%v options:%v", device, target, mType, options)
	flag, _ := parseOptions(options)
	if flag&REMOUNT != REMOUNT {
		logrus.Infof("[Mount] before Mounted: device:%v target:%v mType:%v options:%v", device, target, mType, options)
		if mounted, err := Mounted(target); err != nil || mounted {
			return err
		}
		logrus.Infof("[Mount] after Mounted: device:%v target:%v mType:%v options:%v", device, target, mType, options)
	}
	logrus.Infof("[Mount] end: device:%v target:%v mType:%v options:%v", device, target, mType, options)
	return ForceMount(device, target, mType, options)
}

// ForceMount will mount a filesystem according to the specified configuration,
// *regardless* if the target path is not already mounted. Options must be
// specified like the mount or fstab unix commands: "opt1=val1,opt2=val2". See
// flags.go for supported option flags.
func ForceMount(device, target, mType, options string) error {
	logrus.Debugf("[ForceMount] Begin device:%v, target:%v, mType:%v, options:%v", device, target, mType, options)
	logrus.Debugf("[ForceMount] Before parseOptions:%v", options)
	flag, data := parseOptions(options)
	logrus.Debugf("[ForceMount] After parseOptions():%v", options)
	logrus.Debugf("[ForceMount] Before mount(): device:%v, target:%v, mType:%v, uintptr(flag):%v, data:%v", device, target, mType, uintptr(flag), data)
	if err := mount(device, target, mType, uintptr(flag), data); err != nil {
		logrus.Debugf("[ForceMount] After mount() err:%v ", err)
		return err
	}
	logrus.Debugf("[ForceMount] After mount(): device:%v, target:%v, mType:%v, uintptr(flag):%v, data:%v", device, target, mType, uintptr(flag), data)
	logrus.Debugf("[ForceMount] End device:%v, target:%v, mType:%v, options:%v", device, target, mType, options)
	return nil
}

// Unmount will unmount the target filesystem, so long as it is mounted.
func Unmount(target string) error {
	if mounted, err := Mounted(target); err != nil || !mounted {
		return err
	}
	return ForceUnmount(target)
}

// ForceUnmount will force an unmount of the target filesystem, regardless if
// it is mounted or not.
func ForceUnmount(target string) (err error) {
	// Simple retry logic for unmount
	for i := 0; i < 10; i++ {
		if err = unmount(target, 0); err == nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}
