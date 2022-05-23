package _self

import _self.buildTypes.BuildChainAMI
import jetbrains.buildServer.configs.kotlin.v2019_2.Project


object KIBProject: Project({
   val buildTargets = listOf("centos7", "rhel79", "rhel82", "rhel84")
   val upstreamProjectId = "MesosphereOnly_ClosedSource_SecureSigning_BuildNokmemRpmRepos"
   for (target in buildTargets) {
       buildType(BuildChainAMI(target, upstreamProjectId))
   }

    val fipsBuildTargets = listOf("centos7-fips", "rhel84-fips")
    val fipsUpstreamProjectId = "MesosphereOnly_ClosedSource_SecureSigning_BuildFipsRpmRepos"
    for (target in fipsBuildTargets) {
       buildType(BuildChainAMI(target, fipsUpstreamProjectId))
   }
})
