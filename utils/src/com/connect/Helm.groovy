package com.connect

class Helm implements Serializable{
    
    def script
    
    def stagingChartRepo="your-charts-staging" 
    def prodChartRepo="your-charts"
    def chartRepo=""
    def folder="environment"
    def chartFolder="${folder}/charts"
    def envName="staging"

    Helm(script) {this.script = script}
    
    def init() {
        script.sh "helm init --upgrade"
        script.sh "helm init --client-only"
        script.sh "helm repo add ${stagingChartRepo} https://${stagingChartRepo}.storage.googleapis.com"
        script.sh "helm repo add ${prodChartRepo} https://${prodChartRepo}.storage.googleapis.com"
        script.sh "helm repo update"
    }

    def pushNewChart(repo,env) {
        if (env=="master") {
            chartRepo=prodChartRepo
        } else {
            chartRepo=stagingChartRepo
        }
        script.sh "mkdir env"
        script.sh "gsutil cp gs://${chartRepo}/index.yaml index_from_bucket.yaml"
        script.sh "helm package chart/${repo}"
        script.sh "find . -regex \".*tgz\" -exec mv -i {} -t env \\;"
        script.sh "helm repo index env/ --merge index_from_bucket.yaml --url https://${chartRepo}.storage.googleapis.com"
        script.sh "gsutil cp env/*.* gs://${chartRepo}"
    }

    def createEnv(repo,version) {
        script.sh "helm create ${folder} && rm -rf ${folder}/templates/*.yaml ${folder}/values.yaml ${folder}/templates/*.txt"
        script.sh "gsutil cp gs://${prodChartRepo}/*.tgz ${chartFolder}"
        script.sh "cd ${chartFolder} && find . -maxdepth 1 -name \"*.tgz\" -exec tar -xvzf {} \\;"
        script.sh "cp -a chart/${repo}/. ${chartFolder}/${repo} && rm -rf ${chartFolder}/*.tgz"
        script.sh "sed -i 's/0.1.0/${version}/g' ${folder}/Chart.yaml"
    }

    def transformCommitMessage(){
        script.sh ''' 
            git log -1 --pretty=%B  | head -n 1 > commit 
            cat commit
            sed -i 's/(//g' commit
            sed -i 's/)//g' commit
            sed -i 's/#/pr./g' commit
            cat commit
            '''
    }

    def getCandidateName(repo,commit){
        def version = commit.split(" ").last().trim()
        def branchFrom = commit.split(" ")[2].replaceAll("/","-")
        def chartPackage = "${repo}-${version}.tgz"
        def fullChartPackage = "gs://${stagingChartRepo}/${chartPackage}"
        
        return [fullChartPackage,version,branchFrom]
    }

    def promoteChart(repo,commit,env) {
        def (downloadUrl,version,branchFrom) = getCandidateName(repo,commit)
        script.sh "rm -rf chart/*"
        script.sh "gsutil cp ${downloadUrl} ."
        script.sh "gsutil rm ${downloadUrl}"
        script.sh "tar -xvzf *.tgz -C chart && rm *.tgz"
        script.sh "sed -i 's/${version}/master/g' chart/${repo}/Chart.yaml"
        script.sh "echo ${branchFrom}"
        pushNewChart(repo,env)
    }

    def cleanup(repo,commit){
        def (downloadUrl,version,releaseName) = getCandidateName(repo,commit)
        script.sh "helm delete --purge ${releaseName}"
        script.sh "kubectl delete namespace ${releaseName}"
    }

    def installChart(name,chart,namespace){
        script.sh "helm upgrade --install --debug --wait ${name} ${chart} --namespace ${namespace}"
    }

}
