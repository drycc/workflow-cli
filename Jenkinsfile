def workpath_linux = "/src/github.com/deis/workflow-cli"
def keyfile = "tmp/key.json"

def getBasePath = { String filepath ->
	def filename = filepath.lastIndexOf(File.separator)
	return filepath.substring(0, filename)
}

def upload_artifacts = { String filepath ->
	withCredentials([[$class: 'FileBinding', credentialsId: 'e80fd033-dd76-4d96-be79-6c272726fb82', variable: 'GCSKEY']]) {
		sh "mkdir -p ${getBasePath(filepath)}"
		sh "cat \"\${GCSKEY}\" > ${filepath}"
		sh "make upload-gcs"
	}
}

def gopath_linux = {
	def gopath = pwd() + "/gopath"
	env.GOPATH = gopath
	gopath
}

def workdir_linux = { String gopath ->
	gopath + workpath_linux
}

properties([[$class: 'ParametersDefinitionProperty',
	parameterDefinitions: [
		[$class: 'StringParameterDefinition',
			defaultValue: '',
			description: 'controller-sdk-go commit sha to use in glide.yaml (not updated if empty)',
			name : 'SDK_SHA'],
		[$class: 'StringParameterDefinition',
			defaultValue: '',
			description: 'controller-sdk-go repo to use in glide.yaml (not updated if empty)',
			name: 'SDK_GO_REPO']
		]]])

node('windows') {
	def gopath = pwd() + "\\gopath"
	env.GOPATH = gopath
	def workdir = gopath + "\\src\\github.com\\deis\\workflow-cli"

	def pscmd = { String cmd ->
		"powershell -NoProfile -ExecutionPolicy Bypass -Command \"${cmd}\""
	}

	dir(workdir) {
		stage 'Checkout Windows'
			checkout scm
		stage 'Install Windows'
			bat pscmd('.\\make bootstrap')
		stage 'Test Windows'
			bat pscmd('.\\make test')
	}
}

node('linux') {
	def gopath = gopath_linux()
	def workdir = workdir_linux(gopath)

	dir(workdir) {
		stage 'Checkout Linux'
			checkout scm
		stage 'Install Linux'
			sh 'make bootstrap'
		stage 'Test Linux'
			sh 'make test'
	}
}

def git_commit = ''
def git_branch = ''

stage 'Git Info'
node('linux') {

	def gopath = gopath_linux()
	def workdir = workdir_linux(gopath)

	dir(workdir) {
		checkout scm

		// HACK: Recommended approach for getting command output is writing to and then reading a file.
		sh 'mkdir -p tmp'
		sh 'git describe --all > tmp/GIT_BRANCH'
		sh 'git rev-parse HEAD > tmp/GIT_COMMIT'
		git_branch = readFile('tmp/GIT_BRANCH').trim()
		git_commit = readFile('tmp/GIT_COMMIT').trim()
	}
}

stage 'Build and Upload Artifacts'

parallel(
	revision: {
		node('linux') {
			def gopath = gopath_linux()
			def workdir = workdir_linux(gopath)

			dir(workdir) {
					checkout scm

					if (git_branch != "remotes/origin/master") {
						echo "Skipping build of 386 binaries to shorten CI for Pull Requests"
						env.BUILD_ARCH = "amd64"
					}

					sh 'make bootstrap'

					if (SDK_GO_REPO && SDK_SHA) {
						echo "Updating local glide.yaml with controller-sdk-go repo '${SDK_GO_REPO}' and version '${SDK_SHA}'"

						def pattern = "github\\.com\\/deis\\/controller-sdk-go\\n\\s+version:\\s+[a-f0-9]+"
						def replacement = "${SDK_GO_REPO.replace("/", "\\/")}\\n  version: ${SDK_SHA}"
						sh "perl -i -0pe 's/${pattern}/${replacement}/' glide.yaml"

						def glideYaml = readFile('glide.yaml')
						echo "Updated glide.yaml:\n${glideYaml}"

						sh 'make glideup'
					}

					sh 'make build-revision'

					upload_artifacts(keyfile)
			}
		}
	},
	latest: {
		node('linux') {
			def gopath = gopath_linux()
			def workdir = workdir_linux(gopath)

			dir(workdir) {
					checkout scm

					if (git_branch == "remotes/origin/master") {
						sh 'make bootstrap'
						sh 'make build-latest'

						upload_artifacts(keyfile)
					} else {
						echo "Skipping build of latest artifacts because this build is not on the master branch (branch: ${git_branch})"
					}
			}
		}
	}
)

stage 'Trigger e2e tests'

// If build is on master, trigger workflow-test, otherwise, assume build is a PR and trigger workflow-test-pr
waitUntil {
	try {
		if (git_branch == "remotes/origin/master") {
			build job: '/workflow-test', parameters: [[$class: 'StringParameterValue', name: 'WORKFLOW_CLI_SHA', value: git_commit]]
		} else {
			build job: '/workflow-test-pr', parameters: [[$class: 'StringParameterValue', name: 'WORKFLOW_CLI_SHA', value: git_commit]]
		}
		true
	} catch(error) {
		 input "Retry the e2e tests?"
		 false
	}
}
