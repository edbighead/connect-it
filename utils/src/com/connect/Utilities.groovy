package com.connect

class Utilities implements Serializable{
    def script
    def scanner = "/mnt/disks/jenkins/sonar-scanner/bin/sonar-scanner"

    Utilities(script) {this.script = script}
 
    def go(goal) {
            script.sh "go ${goal}"
    }
    
    def checkVersion() {
        script.sh "go version"
        script.sh "gcloud version"
        script.sh "kubectl version"
        script.sh "helm version"
        script.sh "docker version"
    }
    
    def checkTools() {
            script.sh "gcloud auth activate-service-account --key-file /path/to/your/key.json"
            script.sh "gcloud container clusters get-credentials cluster --zone cluster-zone --project your-project-name"
            script.sh "kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default || true"
    }
    
    def sonar(project) {
    script.sh "${scanner} \
        -Dsonar.projectKey=${project} \
        -Dsonar.organization=connect-it \
        -Dsonar.sources=. \
        -Dsonar.host.url=https://sonarcloud.io \
        -Dsonar.login=sonar-login-id \
        -Dsonar.exclusions=**/*_test.go,**/*.html,src/** \
        -Dsonar.go.coverage.reportPaths=coverage.out"
    }

    def qualityCheck(repo) {
        script.sh("sleep 5")
        def response = script.httpRequest "https://sonarcloud.io/api/measures/component?component=${repo}&metricKeys=alert_status"
        script.echo("Status: "+response.content)
        def content = response.content
        def check = content.contains("\"value\":\"OK\"")
        if (!check){
            script.error("Quality gate failed")
        }
    }

    def envsubst(repo,pattern){
        script.sh "envsubst < \"chart/${repo}/Chart-${pattern}\" > \"chart/${repo}/Chart.yaml\""
        script.sh "envsubst < \"chart/${repo}/values-${pattern}\" > \"chart/${repo}/values.yaml\""
        script.sh "rm -rf chart/${repo}/*-template.yaml"
    }

    def integrationTest(url){
        script.sh "go run integration.go -url=${url}"
    }
}