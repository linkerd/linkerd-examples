package io.buoyant.http.classifiers;

import io.buoyant.linkerd.ResponseClassifierInitializer;

/**
 * This config initializer is loaded by linkerd at startup and registers the
 * `HeaderClassifierConfig` class under the id "io.buoyant.headerClassifier".
 * This tells linkerd's config system to deserialize response classifier config
 * blocks to `HeaderClassifierConfig` if the kind is
 * "io.buoyant.headerClassifer".
 *
 * In order for linkerd to load this class, it must be listed in the
 * `META-INF/services/io.l5d.linkerd.ResponseClassifierInitializer` file.
 */
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
