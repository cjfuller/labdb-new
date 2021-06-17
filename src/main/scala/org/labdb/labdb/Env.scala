package org.labdb.labdb

object Env:
  // STOPSHIP
  def dev() = true

  def prod() = !dev()

  def proxyTarget(): Option[String] = sys.env.get("PROXY_TARGET")
