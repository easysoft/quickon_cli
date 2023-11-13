pipeline {
  agent {
    kubernetes {
      inheritFrom "code-scan build-go xuanim"
    }
  }

  stages {
    stage("checkout code") {
      steps {
        echo "checkout code success"
      }
    }

    stage("build") {
      environment {
        GOPROXY = "https://goproxy.cn,direct"
      }

      steps {
        container('golang') {
          sh "sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories"
          sh "apk --no-cache add make"
          sh 'go mod download'
          sh 'go install github.com/go-task/task/v3/cmd/task@latest'
          sh 'task fmt'
          sh 'task linux'
        }
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
                sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%ce)" --title "sonar scanner" --url "https://sonar.qc.oop.cc/dashboard?id=quickon_cli&branch=${GIT_BRANCH}" --content "sonar scanner quickon_cli success" --debug --custom'
              }
            }
            failure {
              container('xuanimbot') {
                sh 'git config --global --add safe.directory $(pwd)'
                sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%ce)" --title "sonar scanner" --url "https://sonar.qc.oop.cc/dashboard?id=quickon_cli&branch=${GIT_BRANCH}" --content "sonar scanner quickon_cli failure" --debug --custom'
              }
            }
          }
        }
      }
    }
  }
}
