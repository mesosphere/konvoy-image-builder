# gon.hcl
#
source = ["./dist/konvoy-image-wrapper-osx_darwin_amd64/konvoy-image"]
bundle_id = "com.d2iq.dkp.konvoy-image-wrapper"
apple_id {
  password = "@env:AC_PASSWORD"
}
sign {
  application_identity = "@env:AC_APPLICATION_IDENTITY
}