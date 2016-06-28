package io.buoyant.http.classifiers;

import com.twitter.finagle.http.Response;
import com.twitter.finagle.service.ReqRep;
import com.twitter.finagle.service.ResponseClass;
import com.twitter.finagle.service.ResponseClass$;
import com.twitter.util.Function;

/**
 * HeaderClassifier is an HTTP response classifier that classifies responses
 * based on the value of a response header.  If the header value is "success",
 * the response is classified as a success.  If the header value is "retry",
 * the response is classified as a retryable failure.  Otherwise, the response
 * is classified as a non-retryable failure.
 */
public class HeaderClassifier extends Function<ReqRep, ResponseClass> {
    private String headerName;

    /**
     * @param headerName the name of the response header to use
     */
    public HeaderClassifier(String headerName) {
        this.headerName = headerName;
    }

    /**
     * This defines which inputs this classifier accepts.  We accept any
     * response where the header is present.
     */
    @Override
    public boolean isDefinedAt(ReqRep reqRep) {
        if (reqRep.response().isThrow()) return false;
        if (!(reqRep.response().get() instanceof Response)) return false;
        Response rep = (Response) reqRep.response().get();
        return rep.headerMap().get(headerName).isDefined();
    }

    /**
     * Classify responses based on the header.
     */
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
