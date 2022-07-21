source = ["./konvoy-image-bundle_darwin_amd64/konvoy-image"]
bundle_id = "com.d2iq.dkp.konvoy-image-wrapper"

sign {
  application_identity = "Developer ID Application: Mesosphere Inc. (JQJDUUPXFN)"
}

dmg {
  output_path = "konvoy-image.dmg"
  volume_name = "konvoy-image"
}
