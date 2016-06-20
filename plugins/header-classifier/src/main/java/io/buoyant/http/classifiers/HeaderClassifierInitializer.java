package io.buoyant.http.classifiers;

import io.buoyant.linkerd.ResponseClassifierInitializer;

public class HeaderClassifierInitializer extends ResponseClassifierInitializer {
    @Override
    public String configId() {
        return "io.buoyant.headerClassifier";
    }

    @Override
    public Class<?> configClass() {
        return HeaderClassifierConfig.class;
    }
}
