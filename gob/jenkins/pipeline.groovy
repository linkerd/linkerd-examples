// These constants could also be defined as [build parameters](https://wiki.jenkins-ci.org/display/JENKINS/Parameterized+Build)
def namerd = "100.100.100.10"        // The host used by namerctl to reach the namerd cluster.
def namerdNamespace = "default"      // The namerd namespace of the dtab we'll be updating.
def k8sNamespace = "sb-jenkins-test" // The k8s namespace that the services are running in.
def webFrontend = "100.100.100.11"   // The host used to run integration tests against the service mesh
def logicalName = "gen"              // The name used by other services for routing
def currentVersion = "gen-v0"        // The unique k8s svc name for the current production service
def newVersion = "gen-v1"            // The unique k8s svc name for the updated service we want to shift traffic to
def prefix = "/host/${logicalName}"  // The dtab prefix for the dentry we want to change

node {
    stage "auth"
    // Add the $GOPATH to $PATH if it's not there already.
    env.PATH = "${env.GOPATH}/bin:${env.PATH}"
    // Check out this repo in order to use the k8s yml in gob/k8s/gen.
    git branch: 'esbie/sb-jenkins-test', url: 'https://github.com/BuoyantIO/linkerd-examples.git'
    // Authenticate against gcloud with a service account.
    sh "gcloud auth activate-service-account --key-file credentials.json --project buoyant-hosted"
    sh "gcloud docker --authorize-only"
    sh "gcloud config set compute/zone us-central1-b"

    // Save the original dst in case we need to rollback.
    def originalDst = getDstForPrefix(prefix, namerd, namerdNamespace)

    stage "deploy"
    // Signal our intent to deploy, in case someone else is trying to edit the dtab at the same time.
    signalDeploy(prefix, namerd, namerdNamespace)
    // Deploy the new service!
    dir('gob') {
        sh "kubectl apply -f k8s/gen"
    }
    // give the instance some time to start
    sleep 5

    stage "integration testing"
    def dtabOverride = "l5d-dtab: ${prefix} => /tmp/${newVersion}"
    runIntegrationTests(webFrontend, dtabOverride)
    input message: "Integration tests successful! You can reach the service with `curl -H '${dtabOverride}' '${webFrontend}'`", ok: 'OK, done with manual testing step'

    stage "shift traffic (10%)"
    updateDentry(prefix, "1 * /tmp/${newVersion} & 9 * /tmp/${currentVersion}", false, namerd, namerdNamespace)
    try {
        input message: "Shifting 10% of traffic. To view, open: ${webFrontend}:9990", ok: 'OK, success rates look stable'
    } catch(err) {
        echo "reverting traffic back to ${originalDst}"
        updateDentry(prefix, originalDst, true, namerd, namerdNamespace)
        throw err
    }

    stage "shift traffic (100%)"
    updateDentry(prefix, "/tmp/${newVersion}", false, namerd, namerdNamespace)
    input message: "Ready to cleanup?", ok: 'Yep, everything looks good'

    stage "cleanup"
    updateDentry(prefix, "/srv/${newVersion}", true, namerd, namerdNamespace)
    sh "kubectl delete svc ${currentVersion} --namespace=${k8sNamespace}"
    sh "kubectl delete rc ${currentVersion} --namespace=${k8sNamespace}"
}

// The world's simplest integration test
def runIntegrationTests(webFrontend, dtabOverride) {
    def resp = sh(script: "curl -sL -w \"%{http_code}\" -H '${dtabOverride}' \"${webFrontend}/gob?limit=10\" -o /dev/null", returnStdout: true).trim()
    if (resp != "200") {
        error "could not reach"
    }
}

// Calls namerctl dtab get through the pipeline sh step and returns the result.
def getDtab(namerd, namerdNamespace) {
    return sh (script: "namerctl dtab get ${namerdNamespace} --base-url=http://${namerd} --json", returnStdout: true)
}

def getDstForPrefix(prefix, namerd, namerdNamespace) {
    def jsonResp = getDtab(namerd, namerdNamespace)
    return getDst(prefix, jsonResp)
}

// Stops deploy if dst is already marked with /tmp
def signalDeploy(prefix, namerd, namerdNamespace) {
    def jsonResp = getDtab(namerd, namerdNamespace)
    if (isTmp(prefix, jsonResp)) {
        error "dtab is already marked as being deployed! ${jsonResp}"
    }
    def updatedResp = markAsTmp(prefix, jsonResp)
    updateDtab(updatedResp, namerd, namerdNamespace)
}

def updateDentry(prefix, dst, replace, namerd, namerdNamespace) {
    def jsonResp = getDtab(namerd, namerdNamespace)
    def updatedResp = replace ? replaceDst(prefix, dst, jsonResp) : addToDst(prefix, dst, jsonResp)
    updateDtab(updatedResp, namerd, namerdNamespace)
}

// Calls namerctl dtab update through the pipeline sh step and echoes the resulting dtab
def updateDtab(serializedJson, namerd, namerdNamespace) {
    writeFile file: 'default.dtab', text: serializedJson
    sh "namerctl dtab update ${namerdNamespace} default.dtab --base-url=http://${namerd} --json"
    def dtab = getDtab(namerd, namerdNamespace)
    echo "dtab now updated to ${dtab}"
}

// Finds a matching prefix, and checks if it contains tmp
@NonCPS
def isTmp(prefix, jsonResp) {
    def json = new groovy.json.JsonSlurper().parseText(jsonResp)
    for (dentry in json.dtab) {
      if (dentry.prefix == prefix) {
        if (dentry.dst.contains("tmp")) {
            return true
        }
      }
    }
    return false
}

// Finds a matching prefix, and replaces srv with tmp
@NonCPS
def markAsTmp(prefix, jsonResp) {
    def json = new groovy.json.JsonSlurper().parseText(jsonResp)
    def dtab = json.dtab
    for (dentry in dtab) {
      if (dentry.prefix == prefix) {
        dentry.dst = dentry.dst.replaceAll("srv", "tmp")
      }
    }
    return groovy.json.JsonOutput.toJson(dtab)
}

// Finds a matching prefix, and replaces its dst
@NonCPS
def replaceDst(prefix, dst, jsonResp) {
    def json = new groovy.json.JsonSlurper().parseText(jsonResp)
    def dtab = json.dtab
    for (dentry in dtab) {
      if (dentry.prefix == prefix) {
        dentry.dst = dst
      }
    }
    return groovy.json.JsonOutput.toJson(dtab)
}

// Finds a matching prefix, and adds dst as an alternate
@NonCPS
def addToDst(prefix, dst, jsonResp) {
    def json = new groovy.json.JsonSlurper().parseText(jsonResp)
    def dtab = json.dtab
    for (dentry in dtab) {
      if (dentry.prefix == prefix) {
        dentry.dst = dst + " | " + dentry.dst
      }
    }
    return groovy.json.JsonOutput.toJson(dtab)
}

// Finds a matching prefix, and returns it's dst
@NonCPS
def getDst(prefix, jsonResp) {
    def json = new groovy.json.JsonSlurper().parseText(jsonResp)
    def dtab = json.dtab
    for (dentry in dtab) {
      if (dentry.prefix == prefix) {
        return dentry.dst
      }
    }
}
