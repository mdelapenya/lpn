import net.sf.json.JSONArray;
import net.sf.json.JSONObject;

String cronString = BRANCH_NAME == "master" ? "59 23 * * 1-5" : "";

pipeline {
  agent {
    node {
      label 'api'
    }
  }
  triggers { cron(cronString) }
  stages {
    stage('Install Dependencies') {
      steps {
        buildStep('Install Dependencies') {
          //
        }
      }
    }
    stage('Build') {
      steps {
        buildStep('Build') {
          //
        }
      }
    }
    stage('Unit Tests') {
      steps {
        buildStep('Unit Tests') {
          //
        }
      }
    }
    stage('Integration Tests') {
      steps {
        buildStep('Integration Tests') {
          //
        }
      }
    }
    stage('Push image to dockerhub') {
      when {
        branch 'master'
      }
      steps {
        buildStep('Push image to dockerhub') {
          //
        }
      }
    }
  }
  post {
    always {
      junit(allowEmptyResults: true, testResults: '**/test-results/TEST*.xml')
    }
    changed {
      handleCurrentResultChange();
    }
  }
}

void handleCurrentResultChange() {
  switch(currentBuild.currentResult) {
    case 'SUCCESS':
      pushSuccessToSlack();
    break
  }
}

JSONArray buildAttachments(String pretext, String text, String fallback, String title, String color) {
  JSONArray attachments = new JSONArray();

  attachment = new JSONObject();
  attachment.put('pretext', pretext);
  attachment.put('text', text);
  attachment.put('fallback', fallback);
  attachment.put('color', color);
  attachment.put('author_name', getGitAuthor());
  attachment.put('title', title);
  attachment.put('title_link', env.BUILD_URL);
  attachment.put('footer', 'WeDeploy CI Team');
  attachment.put('footer_icon', 'https://a.slack-edge.com/7bf4/img/services/jenkins-ci_48.png')

  JSONArray attachmentFields = new JSONArray();

  lastCommitField = new JSONObject();
  lastCommitField.put('title', 'Last Commit');
  lastCommitField.put('value', getLastCommitMessage());
  lastCommitField.put('short', false);

  attachmentFields.add(lastCommitField);

  attachment.put('fields', attachmentFields);

  attachments.add(attachment);

  return attachments;
}

void buildStep(String message, Closure closure) {
  try {
    setBuildStatus(message, "PENDING");
    closure();
    setBuildStatus(message, "SUCCESS");
  }
  catch (Exception e) {
    setBuildStatus(message, "FAILURE");
    pushFailureToSlack(message);
    throw e
  }
}

String getGitAuthor() {
  def commit = sh(returnStdout: true, script: 'git rev-parse HEAD')
  return sh(returnStdout: true, script: "git --no-pager show -s --format='%an' ${commit}").trim()
}

String getLastCommitMessage() {
  return sh(returnStdout: true, script: 'git log -1 --pretty=%B').trim()
}

String getRandom(String[] array) {
    int rnd = new Random().nextInt(array.length);
    return array[rnd];
}

void pushFailureToSlack(step) {
  String[] errorMessages = [
    'Hey, Vader seems to be mad at you.',
    'I find your your lack of quality disturbing',
    'Please! Don\'t break the CI ;/',
    'Houston, we have a problem'
  ];

  String title = "FAILED: Job ${env.JOB_NAME} - ${env.BUILD_NUMBER}";
  String text = getRandom(errorMessages);

  JSONArray attachments = buildAttachments(
    "BUILD FAILED: ${step} - mdelapenya/lpn",
    text,
    'CI BUILD FAILED',
    title,
    '#ff0000'
  );

  slackSend (color: '#ff0000', attachments: attachments.toString());
}

void pushSuccessToSlack() {
  String[] successMessages = [
    'Howdy, we\'re back on track.',
    'YAY!',
    'The force is strong with this one.'
  ];

  String title = "BUILD FIXED: Job ${env.JOB_NAME} - ${env.BUILD_NUMBER}";
  String text = getRandom(successMessages);

  JSONArray attachments = buildAttachments(
    'BUILD FIXED - mdelapenya/lpn',
    text,
    'CI BUILD FIXED',
    title,
    '#5fba7d'
  );

  slackSend (color: '#5fba7d', attachments: attachments.toString());
}

void setBuildStatus(String message, String state) {
  step([
      $class: "GitHubCommitStatusSetter",
      reposSource: [$class: "ManuallyEnteredRepositorySource", url: "https://github.com/mdelapenya/lpn"],
      contextSource: [$class: "ManuallyEnteredCommitContextSource", context: "ci/jenkins/build-status"],
      errorHandlers: [[$class: "ChangingBuildStatusErrorHandler", result: "UNSTABLE"]],
      statusResultSource: [ $class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: message, state: state]] ]
  ]);
}
