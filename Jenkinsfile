def windows = 'windows'
def linux = 'linux'
def git_commit = ''
def git_branch = ''
def git_tag = ''

def getBasePath = { String filepath ->
	def filename = filepath.lastIndexOf(File.separator)
	return filepath.substring(0, filename)
}

def sh = { cmd ->
	wrap([$class: 'AnsiColorBuildWrapper', 'colorMapName': 'XTerm']) {
		sh cmd
	}
}

def pscmd = { String cmd ->
	"powershell -NoProfile -ExecutionPolicy Bypass -Command \"${cmd}\""
}

def bootstrap = { String node ->
	bootstrapCmd = node == windows ? bat(pscmd('.\\make bootstrap')) : make('bootstrap')

	try {
		bootstrapCmd
	} catch(error) {
		echo "bootstrap failed; wiping 'vendor' directory and trying again..."
		retry(1) {
			dir('vendor') { deleteDir() }
			bootstrapCmd
		}
	}
}

// Disabling until we have a more sustainable Windows Jenkins Agent plan
// See https://github.com/drycc/jenkins-jobs/issues/351
// node(windows) {
// 	def gopath = pwd() + "\\gopath"
// 	env.GOPATH = gopath
// 	def workdir = gopath + "\\src\\github.com\\drycc\\workflow-cli"
//
// 	dir(workdir) {
// 		stage 'Checkout Windows'
// 			checkout scm
// 		stage 'Install Windows'
// 			bootstrap windows
// 		stage 'Test Windows'
// 			bat pscmd('.\\make test')
// 	}
// }

stage 'Git Info'
node(linux) {
	checkout scm

	git_branch = env.BRANCH_NAME
	git_commit = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
	git_tag = sh(returnStdout: true, script: 'git describe --abbrev=0 --tags').trim()

	if (git_branch != "master") {
		// Determine actual PR commit, if necessary
		merge_commit_parents= sh(returnStdout: true, script: 'git rev-parse HEAD | git log --pretty=%P -n 1 --date-order').trim()
		if (merge_commit_parents.length() > 40) {
			echo 'More than one merge commit parent signifies that the merge commit is not the PR commit'
			echo "Changing git_commit from '${git_commit}' to '${merge_commit_parents.take(40)}'"
			git_commit = merge_commit_parents.take(40)
		} else {
			echo 'Only one merge commit parent signifies that the merge commit is also the PR commit'
			echo "Keeping git_commit as '${git_commit}'"
		}
	}
}

def test_image = "quay.io/drycc/workflow-cli-dev:${git_commit.take(7)}"
def mutable_image = 'quay.io/drycc/workflow-cli-dev:latest'

node(linux) {
		stage 'Build and push test container'
			checkout scm
			def quayUsername = "drycc+jenkins"
			def quayEmail = "drycc+jenkins@drycc.com"
			withCredentials([[$class: 'StringBinding',
												credentialsId: 'c67dc0a1-c8c4-4568-a73d-53ad8530ceeb',
									 			variable: 'QUAY_PASSWORD']]) {

				sh """
					docker login -e="${quayEmail}" -u="${quayUsername}" -p="\${QUAY_PASSWORD}" quay.io
					docker build \${DOCKER_BUILD_FLAGS} -t ${test_image} .
					docker push ${test_image}
				"""

				if (git_branch == "master") {
					sh """
						docker tag ${test_image} ${mutable_image}
						docker push ${mutable_image}
					"""
				}
			}
}


stage 'Lint and test container'
// TODO: re-parallelize these tasks when race condition is fixed.
		node(linux) {
			sh "docker run --rm ${test_image} lint"
		}
		node(linux) {
			withCredentials([[$class: 'StringBinding',
												credentialsId: '995d99a7-466b-4beb-bf75-f3ba91cbbc18',
												variable: 'CODECOV_TOKEN']]) {
				def codecov = "codecov -Z -C ${git_commit} -B ${branch_name}"
				if (branch_name != "master") {
					def branch_index = branch_name.indexOf('-')
					def pr = branch_name.substring(branch_index+1, branch_name.length())
					codecov += " -P ${pr}"
				}
				sh "docker run -e CODECOV_TOKEN=\${CODECOV_TOKEN} --rm ${test_image} sh -c 'test-cover.sh &&  ${codecov}'"
			}
		}

stage 'Build and Upload Artifacts'

def old_gcs_bucket = "gs://workflow-cli"
def pr_gcs_bucket = "gs://workflow-cli-pr"
def master_gcs_bucket = "gs://workflow-cli-master"

def upload_artifacts = { String dist_dir, String auth_id, String bucket, boolean cache ->
	def headers  = "-h 'x-goog-meta-git-branch:${git_branch}' "
	headers += "-h 'x-goog-meta-git-sha:${git_commit}' "
	headers += "-h 'x-goog-meta-ci-job:${env.JOB_NAME}' "
	headers += "-h 'x-goog-meta-ci-number:${env.BUILD_NUMBER}' "
	headers += "-h 'x-goog-meta-ci-url:${env.BUILD_URL}'"
	if(!cache) {
		headers += ' -h "Cache-Control:no-cache"'
	}

	def script = "sh -c 'echo \${GCS_KEY_JSON} | base64 -d - > /tmp/key.json "
	script += "&& gcloud auth activate-service-account -q --key-file /tmp/key.json "
	script += "&& gsutil -mq ${headers} cp -a public-read -r /upload/* ${bucket}'"

	withCredentials([[$class: 'StringBinding',
										credentialsId: auth_id,
										variable: 'GCSKEY']]) {
		sh "docker run ${dist_dir} -e GCS_KEY_JSON=\"\${GCSKEY}\" --rm ${test_image} ${script}"
	}
}

def mktmp = {
	// Create tmp directory to store files
	sh 'mktemp -d > tmp_dir'
	def tmp = readFile('tmp_dir').trim()
	echo "Storing binaries in ${tmp}"
	sh 'rm tmp_dir'
	return tmp
}

def version_flags = "-e REVISION=${git_commit.take(7)} -e GIT_TAG=${git_tag}"

parallel(
	revision: {
		node(linux) {

			def flags = ""
			if (git_branch != "master") {
				echo "Skipping build of 386 binaries to shorten CI for Pull Requests"
				flags += "-e BUILD_ARCH=amd64"
			}

			def tmp_dir = mktmp()
			def dist_dir = "-e DIST_DIR=/upload -v ${tmp_dir}:/upload"
			sh "docker run ${flags} ${version_flags} ${dist_dir} --rm ${test_image} make build-revision"

			if (git_branch == "master") {
				upload_artifacts(dist_dir, '6029cf4e-eaa3-4a8e-9dc7-744d118ebe6a', master_gcs_bucket, true)
			} else {
				upload_artifacts(dist_dir, '6029cf4e-eaa3-4a8e-9dc7-744d118ebe6a', pr_gcs_bucket, true)
			}
			sh "docker run ${dist_dir} --rm ${test_image} sh -c 'rm -rf /upload/*'"
			sh "rm -rf ${tmp_dir}"
		}
	},
	latest: {
		node(linux) {
			if (git_branch == "master") {
				def tmp_dir = mktmp()
				def dist_dir = "-e DIST_DIR=/upload -v ${tmp_dir}:/upload"
				sh "docker run ${dist_dir} ${version_flags} --rm ${test_image} make build-latest"

				upload_artifacts(dist_dir, '6029cf4e-eaa3-4a8e-9dc7-744d118ebe6a', master_gcs_bucket, false)
				sh "docker run ${dist_dir} --rm ${test_image} sh -c 'rm -rf /upload/*'"
				sh "rm -rf ${tmp_dir}"
			} else {
				echo "Skipping build of latest artifacts because this build is not on the master branch (branch: ${git_branch})"
			}
		}
	}
)

stage 'Trigger e2e tests'

waitUntil {
	try {
		def chartRepoType = git_branch == "master" ? 'dev' : 'pr'
		build job: 'workflow-chart-e2e', parameters: [
			[$class: 'StringParameterValue', name: 'WORKFLOW_CLI_SHA', value: git_commit],
			[$class: 'StringParameterValue', name: 'ACTUAL_COMMIT', value: git_commit],
			[$class: 'StringParameterValue', name: 'COMPONENT_REPO', value: 'workflow-cli'],
			[$class: 'StringParameterValue', name: 'CHART_REPO_TYPE', value: chartRepoType],
			[$class: 'StringParameterValue', name: 'UPSTREAM_SLACK_CHANNEL', value: '#controller']]
		true
	} catch(error) {
		if (git_branch == "master") {
			throw error
		}

		node(linux) {
			withCredentials([[$class: 'StringBinding', credentialsId: '8a727911-596f-4057-97c2-b9e23de5268d', variable: 'SLACKEMAIL']]) {
				mail body: """<!DOCTYPE html>
<html>
<head>
<meta content='text/html; charset=UTF-8' http-equiv='Content-Type' />
</head>
<body>
<div>Author: ${env.CHANGE_AUTHOR}<br/>
Branch: ${env.BRANCH_NAME}<br/>
Commit: ${env.CHANGE_TITLE}<br/>
<a href="${env.BUILD_URL}console">Click here</a> to view logs.</p>
<a href="${env.BUILD_URL}input/">Click here</a> to restart e2e.</p>
</div>
</html>
""", from: 'jenkins@ci.drycc.cc', subject: 'Workflow CLI E2E Test Failure', to: env.SLACKEMAIL, mimeType: 'text/html'
			}
			input "Retry the e2e tests?"
		}
		false
	}
}
