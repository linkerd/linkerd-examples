# Derived from https://github.com/istio/istio/blob/master/install/updateVersion.sh

function merge_files() {
  ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  ISTIO=$ROOT/istio-linkerd.yml

  echo "# GENERATED FILE. Use with Kubernetes 1.5+" > $ISTIO
  echo "# TO UPDATE, modify files in istio/ and run ./updateVersion.sh" >> $ISTIO
  cat $ROOT/mixer-pilot.yml >> $ISTIO
  cat $ROOT/istio-daemonset.yml >> $ISTIO
  cat $ROOT/istio-ingress.yml >> $ISTIO
  cat $ROOT/istio-egress.yml >> $ISTIO
}

merge_files
