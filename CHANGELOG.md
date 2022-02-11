# Changelog

## [1.6.0](https://github.com/mesosphere/konvoy-image-builder/compare/v1.5.0...v1.6.0) (2022-02-08)


### Features

* Add "dry run" build mode ([#228](https://github.com/mesosphere/konvoy-image-builder/issues/228)) ([e56fdea](https://github.com/mesosphere/konvoy-image-builder/commit/e56fdea37b217cb4c218ddd366fbe6bb6a203879))
* adds an upload artifacts command to konvoy image builder ([#214](https://github.com/mesosphere/konvoy-image-builder/issues/214)) ([9ed1806](https://github.com/mesosphere/konvoy-image-builder/commit/9ed18066608bf2570ce4ed76b5559d82ade78e93))
* convert centos 7 minimal iso to docker image ([#195](https://github.com/mesosphere/konvoy-image-builder/issues/195)) ([b8ecfc5](https://github.com/mesosphere/konvoy-image-builder/commit/b8ecfc57bf0444471e8882ebf27e8e6df6f981bb))
* create RHEL 8.2 and RHEL 8.4 fips image for air-gapped installations ([#208](https://github.com/mesosphere/konvoy-image-builder/issues/208)) ([51af272](https://github.com/mesosphere/konvoy-image-builder/commit/51af272d1c039c564551dbe7f70218629a426edf))
* gather images dynamically ([3fe415f](https://github.com/mesosphere/konvoy-image-builder/commit/3fe415f0b6fa176c4b9baa0f363303b867f08ff3))
* rhel82 FIPS ([#200](https://github.com/mesosphere/konvoy-image-builder/issues/200)) ([1688a02](https://github.com/mesosphere/konvoy-image-builder/commit/1688a028fc8a350eb04c3bd1355f40c06b90a2f4))


### Bugfixes

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
