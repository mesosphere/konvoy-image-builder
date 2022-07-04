source = ["./dist/konvoy-image-wrapper-universal-darwin_darwin_all/konvoy-image"]
bundle_id = "com.d2iq.dkp.konvoy-image-wrapper"
apple_id {
  username = "@env:AC_USERNAME"
  password = "@env:AC_PASSWORD"
}
sign {
  application_identity = "Developer ID Application: Mesosphere Inc. (JQJDUUPXFN)"
}
dmg{
  output_path = "./dist/konvoy-image-wrapper-universal-darwin_darwin_all/konvoy-image.dmg"
  volume_name = "konvoy-image"
}