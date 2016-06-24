package io.buoyant.http.classifiers;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.twitter.finagle.service.ReqRep;
import com.twitter.finagle.service.ResponseClass;
import io.buoyant.linkerd.ResponseClassifierConfig;
import scala.PartialFunction;

public class HeaderClassifierConfig implements ResponseClassifierConfig {
    public String headerName;

    @Override
    @JsonIgnore
    public PartialFunction<ReqRep, ResponseClass> mk() {
        String headerName = this.headerName;
        if (headerName == null) headerName = "status";
        return new HeaderClassifier(headerName);
    }
}
