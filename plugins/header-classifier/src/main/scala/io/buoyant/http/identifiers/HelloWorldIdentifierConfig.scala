package io.buoyant.http.identifiers

import com.fasterxml.jackson.annotation.JsonIgnore
import com.twitter.finagle.{Dtab, Path}
import com.twitter.finagle.http.Request
import io.buoyant.linkerd.protocol.HttpIdentifierConfig
import io.buoyant.router.RoutingFactory.Identifier


class HelloWorldIdentifierConfig extends HttpIdentifierConfig{
  /* This public member is populated by the json property of the same name. */
  var name: String = null

  @JsonIgnore
  override def newIdentifier(prefix: Path, baseDtab: () => Dtab): Identifier[Request] = {
    new HelloWorldIdentifier(prefix, name, baseDtab)
  }
}
