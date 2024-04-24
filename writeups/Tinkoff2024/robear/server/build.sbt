ThisBuild / scalaVersion     := "2.13.8"
ThisBuild / version          := "1.0"
ThisBuild / organization     := "ru.spbctf"
ThisBuild / organizationName := "SPbCTF"

libraryDependencies ++= Seq(
  "com.github.finagle" %% "finch-core" % "0.34.0",
  "com.github.finagle" %% "finch-circe" % "0.34.0",
  "org.playframework" %% "play-json" % "3.0.2",
  "commons-codec" % "commons-codec" % "1.16.1",
  "org.bouncycastle" % "bcprov-jdk16" % "1.46",
  "javax.xml.bind" % "jaxb-api" % "2.3.0",
  "org.scalaj" %% "scalaj-http" % "2.4.2",
)

dependencyOverrides ++= Seq(
  "com.fasterxml.jackson.core" % "jackson-databind" % "2.13.3"
)
