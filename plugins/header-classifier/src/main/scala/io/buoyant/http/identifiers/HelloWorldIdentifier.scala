package io.buoyant.http.identifiers

import com.twitter.finagle.http.Request
import com.twitter.finagle.{Dtab, Path}
import com.twitter.util.Future
import io.buoyant.router.RoutingFactory
import io.buoyant.router.RoutingFactory.{RequestIdentification, UnidentifiedRequest}


case class HelloWorldIdentifier(
  prefix: Path,
  name: String,
  baseDtab: () => Dtab = () => Dtab.base
) extends RoutingFactory.Identifier[Request] {

  private[this] val MoveOn =
    Future.value(new UnidentifiedRequest[Request]("MoveOn to next identifier"))

  def apply(req: Request): Future[RequestIdentification[Request]] = {
    req.headerMap.set ("l5d-hello", name)
    MoveOn
  }
}
