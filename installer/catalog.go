package main

// Everything the installer downloads, pinned by SHA-256. Any byte difference aborts.
// Keep the app list in sync with scripts/install-apps.sh (catalog_test.go checks it).

type download struct {
	Name   string // file name in the cache
	URLs   []string
	SHA256 string
	Size   int64 // bytes, for the progress bar; 0 if unknown
}

// Android SDK platform-tools r37.0.1 (adb, fastboot). SHA-1s match Google's repository2-3.xml.
var platformTools = map[string]download{
	"darwin": {
		Name:   "platform-tools_r37.0.1-darwin.zip",
		URLs:   []string{"https://dl.google.com/android/repository/platform-tools_r37.0.1-darwin.zip"},
		SHA256: "ee39ad5967e95c2a07f04dbcbde96b1a0c916ba376096db5d2f498b7727a5d1d",
		Size:   16110554,
	},
	"linux": {
		Name:   "platform-tools_r37.0.1-linux.zip",
		URLs:   []string{"https://dl.google.com/android/repository/platform-tools_r37.0.1-linux.zip"},
		SHA256: "d230f13842f60f782a8645f9c813f8f845bf36089ea7289f28c48f17979313f1",
		Size:   9054187,
	},
	"windows": {
		Name:   "platform-tools_r37.0.1-win.zip",
		URLs:   []string{"https://dl.google.com/android/repository/platform-tools_r37.0.1-win.zip"},
		SHA256: "45f4d63113e895ebde0c90f194099a4676b6ac653bd28d54314a9e022bbc1a99",
		Size:   8044989,
	},
}

const sfBase = "https://sourceforge.net/projects/andyyan-gsi/files/lineage-21-pre-qpr2-td/"

func sourceforge(file string) []string {
	return []string{
		sfBase + file + "/download?use_mirror=cfhcable",
		sfBase + file + "/download",
	}
}

// LineageOS 21 GSIs by AndyYan, build 20260918 (security patch 2026-09-01), "vndklite" variants:
// their ext4 has no shared_blocks, so /system can be remounted read-write for the camera
// overlay. MD5s match SourceForge's published ones; tested on the device.
var (
	imageNoGoogle = download{
		Name:   "lineage-21.0-20260918-UNOFFICIAL-arm64_bvN-vndklite.img.gz",
		URLs:   sourceforge("lineage-21.0-20260918-UNOFFICIAL-arm64_bvN-vndklite.img.gz"),
		SHA256: "4444245ac00c480c819d997f204c2603c8b94fd298e40bfc3fecf3808d2c0579",
		Size:   1143986059,
	}
	imageGoogle = download{
		Name:   "lineage-21.0-20260918-UNOFFICIAL-arm64_bgN-vndklite-signed.img.gz",
		URLs:   sourceforge("lineage-21.0-20260918-UNOFFICIAL-arm64_bgN-vndklite-signed.img.gz"),
		SHA256: "7baba0781c7733a8879c2eb7cbf34704056d3e4f9c7a0c185867733eec483885",
		Size:   1417963524,
	}
)

func fdroid(file, sha string) download {
	return download{
		Name:   file,
		URLs:   []string{"https://f-droid.org/repo/" + file, "https://f-droid.org/archive/" + file},
		SHA256: sha,
	}
}

// App bundle, from the official F-Droid repo (hashes from f-droid.org/repo/index-v2.json).
var appBundle = []download{
	fdroid("org.fdroid.fdroid_1023052.apk", "985f5181d48bb6bafd54083a048b391271e0ab28385881cc41294fb01a222762"),
	fdroid("com.termux_1002.apk", "e6265a57eb5ca363808488e3b01955958bed93bc0c8a0d281849b363b11027ec"),
	fdroid("org.localsend.localsend_app_643.apk", "82ec3568fba2aa5295b9aae8b76f701d7a4703d86b9f8bad749472038fbaeab3"),
	fdroid("me.zhanghai.android.files_40.apk", "2fe900bf43d725b655008d438f5ca46d0da0f301e3784403b5b96c3fc2df6e8a"),
	fdroid("org.mozilla.fennec_fdroid_1560020.apk", "27f2951376ca1085e0933066c902fdfe4d260916d381e233475c5fd6964f5745"),
	fdroid("org.primftpd_71.apk", "d36566e774e2d0446dcbb41438d22e62d198688aac2711ba7afaa6ea2fdd577d"),
}

var auroraStore = fdroid("com.aurora.store_76.apk", "fd9c75d90d0f4a7c132b9b4a5a2cf1992a45e03b8d8ff988b7dcfbc0db2c4d11")

const (
	expectedBoard = "d39g_4m_bml_s26ultra_mini_pt"
	firefoxPkg    = "org.mozilla.fennec_fdroid"
)
