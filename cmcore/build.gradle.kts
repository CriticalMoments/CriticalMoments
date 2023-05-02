// build.gradle.kts

import org.jetbrains.kotlin.gradle.plugin.mpp.apple.XCFramework

plugins {
    kotlin("multiplatform") version "1.8.21"
}

repositories {
    mavenCentral()
}

// Gradle issue requires unique attributes
// https://slack-chats.kotlinlang.org/t/10340543/after-upgrading-to-gradle-8-1-from-7-x-i-got-consumable-conf
val cmPlatformAttribute = Attribute.of("io.criticalmoments.platform", String::class.java)

kotlin {

    // Build a hello world binary from exeMain src. 
    // Not needed, just getting hands dirty on multi-target multi-package build system
    macosX64("exe") { // on macOS
    // linuxX64("native") // on Linux
    // mingwX64("native") // on Windows
        binaries {
            executable()
        }
    }

    // Build XCFramework for iOS. Builds the commonMain directory
    // TODO: add a sourceSet for appcoreMain, include in framework, and use that instead. 
    val xcf = XCFramework("appcore")
    listOf(
        iosX64 {
          // See comment above
	  attributes.attribute(cmPlatformAttribute, "iosX64")
	},
        iosArm64(),
        iosSimulatorArm64()
    ).forEach {
        it.binaries.framework {
	    baseName = "appcore"
            xcf.add(this)
        }
    }

}

tasks.withType<Wrapper> {
    gradleVersion = "7.6"
    distributionType = Wrapper.DistributionType.BIN
}
