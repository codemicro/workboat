package core

type dockerJob struct {
	RepoOwner     string
	RepoName      string
	ManifestEntry *workflowManifestEntry
}

var jobQueue = make(chan *dockerJob, 512)

func enqueueDockerJob(job *dockerJob) {
	jobQueue <- job
}

func runDockerJob(job *dockerJob) error {
	// Spin up Docker container
	// Setup repository within container
	//     Clone repository within the container
	//     Checkout
	// Run user script
	return nil
}
