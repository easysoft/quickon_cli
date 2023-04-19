pipeline {
  agent {
    kubernetes {
      inheritFrom "code-scan xuanim"
    }
  }

  stages {
    stage("checkout code") {
      steps {
        echo "checkout code success"
      }
    }

    stage('Sonar Scanner') {
      parallel {
        stage('SonarQube') {
          steps {
            container('sonar') {
              withSonarQubeEnv('sonarqube') {
                sh 'make fix-local-version'
                sh 'git config --global --add safe.directory $(pwd)'
                sh 'sonar-scanner -Dsonar.inclusions=$(git diff --name-only HEAD~1|tr "\\n" ",") -Dsonar.analysis.user=$(git show -s --format=%an)'
              }
            }
          }

          post {
            success {
              container('xuanimbot') {
                sh 'git config --global --add safe.directory $(pwd)'
                sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%an)" --title "sonar scanner" --url "https://sonar.qc.oop.cc/dashboard?id=quickon_cli&branch=${GIT_BRANCH}" --content "sonar scanner quickon_cli success" --debug --custom'
              }
            }
            failure {
              container('xuanimbot') {
                sh 'git config --global --add safe.directory $(pwd)'
                sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%an)" --title "sonar scanner" --url "https://sonar.qc.oop.cc/dashboard?id=quickon_cli&branch=${GIT_BRANCH}" --content "sonar scanner quickon_cli failure" --debug --custom'
              }
            }
          }
        }
      }
    }
  }
}
