# prbot

## How to use with Kubernetes CronJob

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: bundle-update
spec:
  schedule: "0 2 * * WED"
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: Never
          containers:
          - image: library/ruby
            name: ruby
            command: [ "/bin/sh" ]
            args:
            - "-c"
            - |2
              set -e
              wget https://github.com/satococoa/prbot/releases/download/v${PRBOT_VERSION}/prbot_${PRBOT_VERSION}_linux_amd64.tar.gz
              tar -xvzf prbot_${PRBOT_VERSION}_linux_amd64.tar.gz
              mv prbot /bin/prbot
              chmod +x /bin/prbot
              gem install bundler
              prbot
            env:
            - name: PRBOT_VERSION
              value: "0.1.2"
            - name: GITHUB_REPOSITORY
              value: "org/repository"
            - name: GITHUB_ACCESS_TOKEN
              value: YOUR_GITHUB_PERSONAL_ACCESS_TOKEN
            - name: BASE_BRANCH
              value: "master"
            - name: COMMAND
              value: "bin/bundle update"
            - name: TITLE
              value: "title of pull request"
            - name: AUTHOR_NAME
              value: "prbot"
            - name: AUTHOR_EMAIL
              value: "prbot@example.com"
```
