package io.buoyant.http.identifiers

import io.buoyant.linkerd.IdentifierInitializer


class HelloWorldIdentifierInitializer extends IdentifierInitializer {
  override def configId: String = "io.buoyant.helloWorldIdentifier"

  override def configClass: Class[_] = return classOf[HelloWorldIdentifierConfig]
}
