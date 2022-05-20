# Changelog

## [1.14.0](https://github.com/mesosphere/konvoy-image-builder/compare/v1.13.2...v1.14.0) (2022-05-20)


### Features

* **gpu:** add image tags ([0013788](https://github.com/mesosphere/konvoy-image-builder/commit/001378809917a73b795e6c83f7bef789201f1811))


### Bug Fixes

* remove extra "release"  keyword from vsphere template name ([#336](https://github.com/mesosphere/konvoy-image-builder/issues/336)) ([a14f6ef](https://github.com/mesosphere/konvoy-image-builder/commit/a14f6ef6e4e1f2f9306f9138f947c281979b27ab))

### [1.13.2](https://github.com/mesosphere/konvoy-image-builder/compare/v1.13.1...v1.13.2) (2022-05-11)


### Bug Fixes

* README remove old  test status ([67102d1](https://github.com/mesosphere/konvoy-image-builder/commit/67102d10fed6449af281b5b8474d27ff5023b63f))

### [1.13.1](https://github.com/mesosphere/konvoy-image-builder/compare/v1.13.0...v1.13.1) (2022-05-11)


### Bug Fixes

* adds a v1 ([ac6e72f](https://github.com/mesosphere/konvoy-image-builder/commit/ac6e72fb197e73050f56b941a1f850dfe5338f0a))

## [1.13.0](https://github.com/mesosphere/konvoy-image-builder/compare/v1.12.0...v1.13.0) (2022-05-09)


### Features

* **azure:** add azure `--instance-type` ([f5e1940](https://github.com/mesosphere/konvoy-image-builder/commit/f5e1940cfd776c5cab3c660af1c28ac17af681b3))


### Bug Fixes

* **aws:** deprecate `--aws-instance-type` ([e0dcc56](https://github.com/mesosphere/konvoy-image-builder/commit/e0dcc561e0286db15ac480f55e8a3291e6e5f544))
* **azure:** append build name to image sku ([#326](https://github.com/mesosphere/konvoy-image-builder/issues/326)) ([b921f42](https://github.com/mesosphere/konvoy-image-builder/commit/b921f42a9df41a721d3a571a03a2c3c3d961d4ec))
* **azure:** fix rhel 8 build name ([9e8ec95](https://github.com/mesosphere/konvoy-image-builder/commit/9e8ec952b789bb7ceb4c469c5a1b3512402d640c))

## [1.12.0](https://github.com/mesosphere/konvoy-image-builder/compare/v1.11.0...v1.12.0) (2022-04-14)


### Features

* use containerd with flatcar ([#300](https://github.com/mesosphere/konvoy-image-builder/issues/300)) ([b96f8bc](https://github.com/mesosphere/konvoy-image-builder/commit/b96f8bc65fa63cd047fe8d2ae1802005e2fe37c4))


### Bug Fixes

* **flatcar:** fix no update settings ([#308](https://github.com/mesosphere/konvoy-image-builder/issues/308)) ([03a618c](https://github.com/mesosphere/konvoy-image-builder/commit/03a618cf8b4901fbcf66572185c55ea77094cc16))
* use bastion in offline fips ova rhel builds ([#307](https://github.com/mesosphere/konvoy-image-builder/issues/307)) ([8d3e338](https://github.com/mesosphere/konvoy-image-builder/commit/8d3e3387779df6117e65ed98783ae1bc1194a69d))

## [1.11.0](https://github.com/mesosphere/konvoy-image-builder/compare/v1.10.0...v1.11.0) (2022-04-07)


### Features

* **azure:** support flatcar images ([6773d1c](https://github.com/mesosphere/konvoy-image-builder/commit/6773d1cf609eb6f3a6e20ef538d38b67665128b7))
* **azure:** support oracle images ([c11780a](https://github.com/mesosphere/konvoy-image-builder/commit/c11780a1d9fa48089bbb8f35b483a1a9202f2612))
* **azure:** support rhel images ([cdef472](https://github.com/mesosphere/konvoy-image-builder/commit/cdef4724278ce6e7860ca431d699c07cd9d565da))
* **azure:** support sles images ([49a7745](https://github.com/mesosphere/konvoy-image-builder/commit/49a774523ebaab6b7a26a229330f9dea5173b7b3))


### Bug Fixes

* support flatcar 3033.2.0 ([#299](https://github.com/mesosphere/konvoy-image-builder/issues/299)) ([43bdfd2](https://github.com/mesosphere/konvoy-image-builder/commit/43bdfd27298414130c4b3703636377c2c64fd8c7))

## [1.10.0](https://github.com/mesosphere/konvoy-image-builder/compare/v1.9.1...v1.10.0) (2022-03-31)


### Features

* **azure:** support ubuntu images ([91dac7f](https://github.com/mesosphere/konvoy-image-builder/commit/91dac7ffc7edf6251898c96143b11a2406b34f3d))


### Bug Fixes

* cleanup vsphere VM when building vSphere template in dry run ([#283](https://github.com/mesosphere/konvoy-image-builder/issues/283)) ([44b1a94](https://github.com/mesosphere/konvoy-image-builder/commit/44b1a9423b34ab444d095834a4537fafe7ca10ea))
* fixes version lock to set fact ([ebaba83](https://github.com/mesosphere/konvoy-image-builder/commit/ebaba83028e79adf8fb295b26961a18a15e9be50))
* hardcode v3.4.x etcd version ([8d9c8e9](https://github.com/mesosphere/konvoy-image-builder/commit/8d9c8e924906cf55ddbab161b410fe60eee0d804))
* makefile targets for NVIDIA GPU support ([#285](https://github.com/mesosphere/konvoy-image-builder/issues/285)) ([b56c5af](https://github.com/mesosphere/konvoy-image-builder/commit/b56c5afd51bad97844112854521a6aae8d7ff305))
* **packer:** fix spacing typo in packer template ([92fa950](https://github.com/mesosphere/konvoy-image-builder/commit/92fa950a29d8694dfd83e1817b643cc4c23095f0))

### [1.9.1](https://github.com/mesosphere/konvoy-image-builder/compare/v1.9.0...v1.9.1) (2022-03-24)


### Bug Fixes

* ova packer template ([ebcd1da](https://github.com/mesosphere/konvoy-image-builder/commit/ebcd1da234e4fd8d7a13030b4b49b77ee74a093b))

## [1.9.0](https://github.com/mesosphere/konvoy-image-builder/compare/v1.8.0...v1.9.0) (2022-03-24)


### Features

* bump k8s to 1.22.4 ([d36a81d](https://github.com/mesosphere/konvoy-image-builder/commit/d36a81dd2646e72e247ecca1aedf336bd449cb95))
* bump kubernetes to 1.22.8 and use the new automated repos ([809bbd9](https://github.com/mesosphere/konvoy-image-builder/commit/809bbd9baf29df21e44512444da3795a8195cca6))
* support azure ([#230](https://github.com/mesosphere/konvoy-image-builder/issues/230)) ([016481a](https://github.com/mesosphere/konvoy-image-builder/commit/016481af838a878dcbfa7e7f94c5be958e35364d))


### Bug Fixes

* **ansible:** allow rsa public keys ([#271](https://github.com/mesosphere/konvoy-image-builder/issues/271)) ([291e922](https://github.com/mesosphere/konvoy-image-builder/commit/291e9220b843560e7c905b48dd0c1b63ca8a7ab2))
* go-mod tidy ([7ccfaa7](https://github.com/mesosphere/konvoy-image-builder/commit/7ccfaa77ee44e2ef775c8c4f2c52390fdfa57eb8))
* **lint:** don't lint CHANGELOG.md ([b8401b2](https://github.com/mesosphere/konvoy-image-builder/commit/b8401b2b00d1751da70a60789daf87613c611014))
* remove note to add promotion job ([d9cd670](https://github.com/mesosphere/konvoy-image-builder/commit/d9cd6705fe0ef14eeff9bb3d6a29c7bd4d54ca03))

## [1.8.0](https://github.com/mesosphere/konvoy-image-builder/compare/v1.7.0...v1.8.0) (2022-03-17)


### Features

* allow user to run ssh communicator through bastion ([#251](https://github.com/mesosphere/konvoy-image-builder/issues/251)) ([84e9674](https://github.com/mesosphere/konvoy-image-builder/commit/84e967431ac9b1dc8b7563536006c8b15e574562))
* bulild vsphere template in air gapped environment ([#246](https://github.com/mesosphere/konvoy-image-builder/issues/246)) ([5a2c62f](https://github.com/mesosphere/konvoy-image-builder/commit/5a2c62f95a0293a20dda1ebfd78f7c5fcbb4915e))
* create vsphere template image for RedHat 8.4 and 7.9 ([#239](https://github.com/mesosphere/konvoy-image-builder/issues/239)) ([b5e7abe](https://github.com/mesosphere/konvoy-image-builder/commit/b5e7abe50acf824e7244a8fb63440164b8ec03ac))


### Bug Fixes

* **app:** remove unused `gen.go` ([8bf34b3](https://github.com/mesosphere/konvoy-image-builder/commit/8bf34b35eb56e30f6226b3092bed9a5db20fba53))
* **cmd:** add subcommads to `build` and `generate` ([4fb3798](https://github.com/mesosphere/konvoy-image-builder/commit/4fb3798804ca403ae547566046a7052f6dfacdf9))
* **lint:** fix markdown rules ([b052bff](https://github.com/mesosphere/konvoy-image-builder/commit/b052bff48bd215d032c364518b94d54015e0d617))
* **lint:** fix textlint rules ([731e192](https://github.com/mesosphere/konvoy-image-builder/commit/731e192937f5761da0d1fda01d9b8b3f76c2cf9f))
* move goreleaser to where it really is ([82992d5](https://github.com/mesosphere/konvoy-image-builder/commit/82992d53d1c1b5a75baa15f5d14120000e58b706))
* **pkg:** remove unused `config` package ([1c6509a](https://github.com/mesosphere/konvoy-image-builder/commit/1c6509a97fe057b93b2d6414bb6b303df44adaea))
* use crictl to pull images and supports mirrors ([#252](https://github.com/mesosphere/konvoy-image-builder/issues/252)) ([f14f841](https://github.com/mesosphere/konvoy-image-builder/commit/f14f841c5d325d9677349c18f5e717109e661e7f))

## [1.7.0](https://github.com/mesosphere/konvoy-image-builder/compare/v1.6.0...v1.7.0) (2022-02-14)


### Features

* use published os-packages-bundles ([680d56e](https://github.com/mesosphere/konvoy-image-builder/commit/680d56e9035980b47a40fc3a532c443d71db173e))


### Bug Fixes

* create systemd drop-in to disable NetworkManager-cloud-setup service ([2f6011a](https://github.com/mesosphere/konvoy-image-builder/commit/2f6011aeef8770802c5fd5db7ccb18fdda3ae1c2))
* disable nm-cloud-setup only for AWS provider ([9da50ce](https://github.com/mesosphere/konvoy-image-builder/commit/9da50ce65525199efdf33e647f81751484fb1968))
* linting errors in changelog ([4eb2f96](https://github.com/mesosphere/konvoy-image-builder/commit/4eb2f9613953743b4bb007f97ccdd81bc6acee64))
* **release:** run goreleaser on release publish ([f1218a1](https://github.com/mesosphere/konvoy-image-builder/commit/f1218a13167678c7e49ebc71347eb9f4f7f869a9))

## [1.6.0](https://github.com/mesosphere/konvoy-image-builder/compare/v1.5.0...v1.6.0) (2022-02-08)


### Features

* Add "dry run" build mode ([#228](https://github.com/mesosphere/konvoy-image-builder/issues/228)) ([e56fdea](https://github.com/mesosphere/konvoy-image-builder/commit/e56fdea37b217cb4c218ddd366fbe6bb6a203879))
* adds an upload artifacts command to konvoy image builder ([#214](https://github.com/mesosphere/konvoy-image-builder/issues/214)) ([9ed1806](https://github.com/mesosphere/konvoy-image-builder/commit/9ed18066608bf2570ce4ed76b5559d82ade78e93))
* convert centos 7 minimal iso to docker image ([#195](https://github.com/mesosphere/konvoy-image-builder/issues/195)) ([b8ecfc5](https://github.com/mesosphere/konvoy-image-builder/commit/b8ecfc57bf0444471e8882ebf27e8e6df6f981bb))
* create RHEL 8.2 and RHEL 8.4 fips image for air-gapped installations ([#208](https://github.com/mesosphere/konvoy-image-builder/issues/208)) ([51af272](https://github.com/mesosphere/konvoy-image-builder/commit/51af272d1c039c564551dbe7f70218629a426edf))
* gather images dynamically ([3fe415f](https://github.com/mesosphere/konvoy-image-builder/commit/3fe415f0b6fa176c4b9baa0f363303b867f08ff3))
* rhel82 FIPS ([#200](https://github.com/mesosphere/konvoy-image-builder/issues/200)) ([1688a02](https://github.com/mesosphere/konvoy-image-builder/commit/1688a028fc8a350eb04c3bd1355f40c06b90a2f4))


### Bug Fixes

* add testify to go.mod ([91c38d5](https://github.com/mesosphere/konvoy-image-builder/commit/91c38d500238fae2f996e6164bf65e5505986587))
* Allow user to provide a subset of registry configuration fields ([3571d3a](https://github.com/mesosphere/konvoy-image-builder/commit/3571d3ade91131fe9cbd439788500d733a1a1613))
* allows users to set kubernetes version through flag in build command ([04681f1](https://github.com/mesosphere/konvoy-image-builder/commit/04681f120d929103253e8efd02f03ada63311f3a))
* **ansible:** reuse roles for image saving ([e298986](https://github.com/mesosphere/konvoy-image-builder/commit/e298986a6c16e6b2ba82619bab46045346a9f097))
* configure NetworkManager to prevent interfering with interfaces ([#231](https://github.com/mesosphere/konvoy-image-builder/issues/231)) ([36de19f](https://github.com/mesosphere/konvoy-image-builder/commit/36de19f3ec5b6401240d5ee0082ae59b3efacc2c))
* have extra-vars work ([a9b962b](https://github.com/mesosphere/konvoy-image-builder/commit/a9b962ba6071d68e2b45f22b43b68de1a47581ff))
* lint errors ([46a6b19](https://github.com/mesosphere/konvoy-image-builder/commit/46a6b195ca8528b19fa171e21f739d1f5cc8e951))
* **os-packages:** prevent clean error ([cda2e50](https://github.com/mesosphere/konvoy-image-builder/commit/cda2e50d6b4da98d15209168b19791bb3e44cd1a))
* **os-packages:** set defaults for targets ([3bfc439](https://github.com/mesosphere/konvoy-image-builder/commit/3bfc439a2daf17baae37fe364a1d48c5588c574c))
* remove execute bits from playbook ([8a44a81](https://github.com/mesosphere/konvoy-image-builder/commit/8a44a8127ee6911f3976d22ef146de6a46508bf9))
* remove unused 'global' playbook ([a3356c0](https://github.com/mesosphere/konvoy-image-builder/commit/a3356c0680449d4108c809f1d9f2e3b7e3bea24f))
* remove unused 'images' group vars ([31eb405](https://github.com/mesosphere/konvoy-image-builder/commit/31eb4058eb4c9d278bbaecd7f6578722ad4799a1))
* replace broken centos 8 appstream repository with alma linux repository ([#227](https://github.com/mesosphere/konvoy-image-builder/issues/227)) ([ada2ca9](https://github.com/mesosphere/konvoy-image-builder/commit/ada2ca94bfd842526a4af9878a1b67ff80e2afa3))
* set the correct KIB version ami tag ([d92dd74](https://github.com/mesosphere/konvoy-image-builder/commit/d92dd744a92edcd8b3d0aae10a0cba44d78b5dd0))
