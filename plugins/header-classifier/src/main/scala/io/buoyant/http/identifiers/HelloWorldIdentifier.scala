package io.buoyant.http.identifiers



import com.twitter.finagle.buoyant.Dst
import com.twitter.finagle.http.Request
import com.twitter.finagle.{Dtab, Path}
import com.twitter.util.Future
import io.buoyant.router.RoutingFactory
import io.buoyant.router.RoutingFactory.{IdentifiedRequest, RequestIdentification}

object HelloWorldIdentifier {
  def mk(
    prefix: Path,
    baseDtab: () => Dtab = () => Dtab.base
  ): RoutingFactory.Identifier[Request] = HelloWorldIdentifier(prefix, "hello", baseDtab)
}

case class HelloWorldIdentifier(
  prefix: Path,
  value: String,
  baseDtab: () => Dtab = () => Dtab.base
) extends RoutingFactory.Identifier[Request] {

  def apply(req: Request): Future[RequestIdentification[Request]] = {
      val dst = Dst.Path(Path.Utf8(req.path), baseDtab(), Dtab.local)
      if (req.headerMap.contains("l5d-hello")) {
        req.headerMap.set ("l5d-hello", "hello")
      }

      Future.value(new IdentifiedRequest[Request](dst, req))
  }
}
