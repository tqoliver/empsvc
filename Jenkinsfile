#!groovy
  
def imageTag = "timothyoliver/employeeservices"
def prevImageTag = ''
def prevBuildNum = ''
def firstDeploy = false

node {
stage 'Build in QA'
openshiftBuild(namespace:'ccdemoqa', buildConfig: 'employeeservices', showBuildLogs: 'true')
stage 'Deploy to QA'
openshiftDeploy(namespace: 'ccdemoqa', deploymentConfig: 'employeeservices')
openshiftScale(namespace: 'ccdemoqa', deploymentConfig: 'employeeservices', replicaCount: '1')
}
