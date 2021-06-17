package org.labdb.labdb

object Proxy:
  val devProxyTarget = "http://localhost:3001"
  val proxySuffix = "-backend.labdb.io"

  def backendHost(requestHost: String): String =
    if Env.dev() then
      Env.proxyTarget() match {
        case Some(target) => target
        case None         => devProxyTarget
      }
    else "https://" + requestHost.replace(".labdb.io", proxySuffix)
