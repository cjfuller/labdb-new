val ScalatraVersion = "2.7.+"
ThisBuild / scalaVersion := "3.0.0"
ThisBuild / organization := "org.labdb"
name := "labdb"
version := "3.0.0"
libraryDependencies ++= Seq(
  "org.scalatest" %% "scalatest" % "3.2.9" % "test",
  ("org.scalatra" %% "scalatra" % ScalatraVersion)
    .cross(CrossVersion.for3Use2_13),
  ("org.scalatra" %% "scalatra-scalatest" % ScalatraVersion % "test")
    .cross(CrossVersion.for3Use2_13),
  "ch.qos.logback" % "logback-classic" % "1.2.3" % "runtime",
  "org.eclipse.jetty" % "jetty-webapp" % "9.4.35.v20201120",
  "javax.servlet" % "javax.servlet-api" % "3.1.0" % "provided",
  "com.softwaremill.sttp.client3" %% "core" % "3.3.6"
)

enablePlugins(JettyPlugin)
