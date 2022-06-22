source = ["./dist/konvoy-image-osx_darwin_amd64/konvoy-image"]
bundle_id = "com.d2iq.dkp.konvoy-image-wrapper"
apple_id {
  username = "@env:AC_USERNAME"
  password = "@env:AC_PASSWORD"
}
sign {
  application_identity = "@env:AC_APPLICATION_IDENTITY"
}