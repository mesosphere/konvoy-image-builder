package _self

import _self.buildTypes.*
import jetbrains.buildServer.configs.kotlin.v2019_2.Project
import jetbrains.buildServer.configs.kotlin.v2019_2.RelativeId
import jetbrains.buildServer.configs.kotlin.v2019_2.projectFeatures.buildReportTab


object KIBProject: Project({
   val nonFIPSbuildTargets = listOf("centos7", "rhel79", "rhel82", "rhel84")
   var nonFIPSUpstreamProjectId = "MesosphereOnly_ClosedSource_SecureSigning_BuildNokmemRpmRepos"
   for (target in nonFIPSbuildTargets) {
       buildType(BuildChainAMI(target, nonFIPSUpstreamProjectId))
   }
})
