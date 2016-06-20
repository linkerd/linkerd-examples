package io.buoyant.http.classifiers;

import com.twitter.finagle.http.Response;
import com.twitter.finagle.service.ReqRep;
import com.twitter.finagle.service.ResponseClass;
import com.twitter.finagle.service.ResponseClass$;
import com.twitter.util.Function;

public class HeaderClassifier extends Function<ReqRep, ResponseClass> {
    private String headerName;

    public HeaderClassifier(String headerName) {
        this.headerName = headerName;
    }

    @Override
    public boolean isDefinedAt(ReqRep reqRep) {
        if (reqRep.response().isThrow()) return false;
        if (!(reqRep.response().get() instanceof Response)) return false;
        Response rep = (Response) reqRep.response().get();
        return rep.headerMap().get(headerName).isDefined();
    }

    @Override
    public ResponseClass apply(ReqRep reqRep) {
        Response rep = (Response) reqRep.response().get();
        String status = rep.headerMap().get(headerName).get();
        switch (status) {
            case "success": return ResponseClass$.MODULE$.Success();
            case "retry": return ResponseClass$.MODULE$.RetryableFailure();
            default: return ResponseClass$.MODULE$.NonRetryableFailure();
        }
    }
}
