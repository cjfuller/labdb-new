package org.labdb.labdb

import org.scalatra._

class LabdbMain extends ScalatraServlet:
  get("/") {
    <body>
      <h1>Hello</h1>
      Hello, world!
    </body>
  }