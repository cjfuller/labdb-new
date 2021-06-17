package org.labdb.labdb

import org.eclipse.jetty.server.{
  CustomRequestLog,
  Server,
  Slf4jRequestLogWriter
}
import org.eclipse.jetty.server.handler.{HandlerCollection, RequestLogHandler}
import org.eclipse.jetty.servlet.{DefaultServlet, ServletContextHandler}
import org.eclipse.jetty.webapp.WebAppContext
import org.scalatra.servlet.ScalatraListener

object JettyLauncher:
  @main def main() =
    val port =
      if (System.getenv("PORT") != null) System.getenv("PORT").toInt else 8080

    val server = Server(port)
    val context = WebAppContext()
    context setContextPath "/"
    context.setResourceBase("src/main/webapp")
    context.addEventListener(ScalatraListener())
    context.addServlet(classOf[DefaultServlet], "/")

    val requestLogHandler = RequestLogHandler()
    val requestLog = CustomRequestLog(
      Slf4jRequestLogWriter(),
      CustomRequestLog.EXTENDED_NCSA_FORMAT
    )

    requestLogHandler.setRequestLog(requestLog)

    val handlers = HandlerCollection()
    handlers.addHandler(context)
    handlers.addHandler(requestLogHandler)

    server.setHandler(handlers)

    server.start
    server.join
