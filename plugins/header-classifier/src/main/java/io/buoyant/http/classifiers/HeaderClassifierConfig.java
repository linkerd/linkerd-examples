package io.buoyant.http.classifiers;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.twitter.finagle.service.ReqRep;
import com.twitter.finagle.service.ResponseClass;
import io.buoyant.linkerd.ResponseClassifierConfig;
import scala.PartialFunction;

/**
 * HeaderClassifierConfig defines the structure of the config block for this
 * plugin and constructs the response classifier.
 */
public class HeaderClassifierConfig extends ResponseClassifierConfig {

    /* This public member is populated by the json property of the same name. */
    public String headerName;

    /**
     * Construct the repsonse classifier.
     */
    @Override
    @JsonIgnore
    public PartialFunction<ReqRep, ResponseClass> mk() {
        String headerName = this.headerName;
        if (headerName == null) headerName = "status";
        return new HeaderClassifier(headerName);
    }
}
