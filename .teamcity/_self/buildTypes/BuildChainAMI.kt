package _self.buildTypes

import jetbrains.buildServer.configs.kotlin.v10.toExtId
import jetbrains.buildServer.configs.kotlin.v2019_2.buildSteps.script
import jetbrains.buildServer.configs.kotlin.v2019_2.triggers.finishBuildTrigger
import jetbrains.buildServer.configs.kotlin.v2019_2.triggers.retryBuild
import jetbrains.buildServer.configs.kotlin.v2019_2.triggers.vcs

class BuildChainAMI(osTarget: String, upstreamDependency: String) : BuildType({
    id("BuildChainAMI$osTarget".toExtId())
    name = "Build Chain - AMI $osTarget"
    val k8s_version = %dep.$upstreamDependency.teamcity.build.branch%
    steps {
        script {
            
            name = "Build $osTarget AMI"
            criptContent = """
                cat <<EOF > buildchain-overrides.yaml
                build_name_extra: -buildchain
                EOF                
                make devkit.run WHAT="make $osTarget \
                ADDITIONAL_OVERRIDES=buildchain-overrides.yaml \
                ADDITIONAL_ARGS=\"--extra-vars kubernetes-version=$k8s_version \" \
                BUILD_DRY_RUN=true" """.trimIndent()
        }
    }
    vcs {
        root(DslContext.settingsRoot)
    }
    triggers {
        vcs {
            id = "vcsTrigger"
            enabled = false
            branchFilter = """
                +:*
                -:<default>
            """.trimIndent()
        }
        finishBuildTrigger {
            id = "TRIGGER_153"
            buildType = upstreamDependency
            successfulOnly = true
            branchFilter = """
                +:v*
                -:v1.21.0
            """.trimIndent()
        }
        retryBuild {
            id = "retryBuildTrigger"
            delaySeconds = 120
        }
    }

    dependencies {
        snapshot(AbsoluteId(upstreamDependency)) {
        }
    }

    requirements {
        exists("DOCKER_VERSION", "RQ_26")
    }
})
