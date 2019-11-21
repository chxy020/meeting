package cron

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"seater/models"

	"github.com/astaxie/beego/orm"
	"github.com/juju/errors"
	cron "gopkg.in/robfig/cron.v2"
)

// JobManager is the manager for cron jobs
type JobManager struct {
	mode           string
	lock           sync.Mutex
	m              *models.SeaterModel
	o              orm.Ormer
	registeredJobs map[string]Job
	c              *cron.Cron
	stop           chan bool
}

var jobManager = new(JobManager)

// Init initialize a JobManager
func (manager *JobManager) Init() (err error) {
	manager.stop = make(chan bool)
	m, err := models.NewModel()
	if err != nil {
		err = errors.Trace(err)
		return
	}
	manager.m = m
	manager.o = m.O()
	manager.c = cron.New()
	return
}

func (manager *JobManager) registerJob(jobs ...Job) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if manager.registeredJobs == nil {
		manager.registeredJobs = make(map[string]Job)
	}
	for _, job := range jobs {
		name := job.name()
		if name == "" {
			panic(fmt.Sprintf("job %T name is empty", job))
		}
		_, ok := manager.registeredJobs[name]
		if ok {
			panic(fmt.Sprintf("cron job %s already registered", name))
		}
		manager.registeredJobs[name] = job
		job.setName(name)
	}
}

// GetJobManager returns the job manager
func GetJobManager() *JobManager {
	return jobManager
}

func (manager *JobManager) jobNames() []string {
	names := make([]string, 0, len(manager.registeredJobs))
	for k := range manager.registeredJobs {
		names = append(names, k)
	}
	return names
}

// Start start the cron jobs
func (manager *JobManager) Start() {
	for _, job := range manager.registeredJobs {
		if err := job.init(); err != nil {
			continue
		}
		manager.addCronJob(job, manager.stop)
	}
	manager.c.Start()
}

// Stop stops the cron jobs
func (manager *JobManager) Stop() {
	manager.c.Stop()
	close(manager.stop)
}

func (manager *JobManager) addCronJob(job Job, stop chan bool) {
	cronfunc := func() {
		if !job.enabled() {
			return
		}
		defer func() {
			_ = recover()
		}()

		var err error
		_ = job.name()
		// TODO(heha37): not supoort concurrent yet, because all cron jobs share the same job instance,
		// we should create new job instance for each job running
		if !job.concurrent() {
			job.lock()
			defer job.unlock()
		}

		job.reset()
		defer job.runDefers()
		job.setJobOptions(job.jobOptions())
		if err = job.prepare(); err != nil {
			job.setFailed(true)
			return
		}
		defer job.finish()
		if err = job.run(); err != nil {
			job.setFailed(true)
		} else {
			job.setFailed(false)
		}
	}

	cronExpression := job.cronExpression()
	if strings.HasPrefix(cronExpression, "@every") {
		a := strings.SplitN(cronExpression, " ", 2)
		interval, err := time.ParseDuration(strings.TrimSpace(a[1]))
		if err != nil {
			return
		}
		job.setInterval(interval)

		go func() {
			ticker := time.NewTicker(interval)
			for {
				select {
				case t := <-ticker.C:
					job.setStartTime(t)
					cronfunc()
				case <-stop:
					ticker.Stop()
					return
				}
			}
		}()
	} else {
		if _, err := manager.c.AddFunc(cronExpression, cronfunc); err != nil {
			return
		}
	}
}
