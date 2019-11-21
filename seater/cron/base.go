package cron

import (
	"sync"
	"time"

	"seater/models"

	"github.com/astaxie/beego/orm"
	"github.com/juju/errors"
)

type jobOptions struct {
	autoBeginTran bool
}

type jobOption func(*jobOptions)

func withTran(autoBegin bool) jobOption {
	return func(o *jobOptions) {
		o.autoBeginTran = autoBegin
	}
}

// Job defines cron j interface
type Job interface {
	init() error
	// cronExpression specific the job running time using cron syntax
	// if the expression starts with @every, this job is not passed to cron manager,
	// but managed by itself.
	cronExpression() string
	setJobOptions([]jobOption)
	jobOptions() []jobOption
	prepare() error
	run() error
	finish() error
	// concurrent means the next cron job of the same job instance
	// will be run even if the fore one is running
	concurrent() bool
	model() *models.SeaterModel
	renew() error
	lock()
	unlock()
	// name is the required method for every job
	name() string
	setName(string)
	// setFailed if the cron job failed
	setFailed(bool)
	// failed checks if cron job failed
	failed() bool
	// enabled checks if cron job is enabled to run
	enabled() bool
	setInterval(duration time.Duration)
	// interval return the running interval	if cronExpression start with @every
	// otherwise return 0s
	interval() time.Duration
	// setStartTime sets started time of current running
	setStartTime(time.Time)
	// run registered defered funcs
	runDefers()
	// reset init defers
	reset()
}

type baseJob struct {
	m                    *models.SeaterModel
	o                    orm.Ormer
	l                    sync.Mutex
	f                    bool
	jm                   *JobManager
	intervalDur          time.Duration
	startTime            time.Time
	expectedFinishedTime time.Time
	jobName              string

	opts           *jobOptions

	deferList []func()
}

func (j *baseJob) init() (err error) {
	j.opts = &jobOptions{
		autoBeginTran: true,
	}
	err = j.renew()
	if err != nil {
		err = errors.Trace(err)
		return
	}
	j.jm = jobManager
	j.f = false
	return
}

func (j *baseJob) reset() {
	j.deferList = make([]func(), 0, 1)
}

func (j *baseJob) runDefers() {
	// run defered funcs in last-in-first-out order
	for i := len(j.deferList) - 1; i >= 0; i-- {
		j.deferList[i]()
	}
}

func (j *baseJob) registerDefer(f func()) {
	j.deferList = append(j.deferList, f)
}

func (j *baseJob) renew() (err error) {
	m, err := models.NewModel()
	if err != nil {
		err = errors.Annotate(err, "failed to init models")
		return
	}
	j.o = m.O()
	j.m = m
	return nil
}

// name should be overrided
func (j *baseJob) name() string {
	return j.jobName
}

func (j *baseJob) setName(name string) {
	j.jobName = name
}

// run should be overrided
func (j *baseJob) run() error {
	return nil
}

func (j *baseJob) model() *models.SeaterModel {
	return j.m
}

func (j *baseJob) setJobOptions(opts []jobOption) {
	for _, opt := range opts {
		opt(j.opts)
	}
}

func (j *baseJob) jobOptions() []jobOption {
	return nil
}

func (j *baseJob) beginTransaction() (err error) {
	if err = j.m.Begin(); err != nil {
		return errors.Annotate(err, "failed to begin transaction")
	}
	return
}

func (j *baseJob) finishTransaction() (err error) {
	if !j.failed() {
		if err = j.m.Commit(); err != nil {
			j.setFailed(true)
		} else {
			return nil
		}
	}
	if err = j.m.Rollback(); err != nil {
		j.setFailed(true)
		return errors.Annotatef(err, "failed to rollback transaction")
	}
	return nil
}

func (j *baseJob) prepare() (err error) {
	name := j.name()
	if err = j.renew(); err != nil {
		return errors.Annotatef(err, "failed to renew %s job", name)
	}
	if j.opts.autoBeginTran {
		if err = j.beginTransaction(); err != nil {
			return errors.Trace(err)
		}
	}
	return
}

func (j *baseJob) finish() (err error) {
	if !j.opts.autoBeginTran {
		return
	}
	if err = j.finishTransaction(); err != nil {
		return errors.Trace(err)
	}
	return
}

func (j *baseJob) inTransaction(f func() error) (err error) {
	if err = j.beginTransaction(); err != nil {
		return errors.Trace(err)
	}
	runerr := f()
	if runerr != nil {
		j.setFailed(true)
	}
	if err = j.finishTransaction(); err != nil {
		return errors.Trace(err)
	}
	return errors.Trace(runerr)
}

func (j *baseJob) concurrent() bool {
	return false
}

func (j *baseJob) lock() {
	j.l.Lock()
}

func (j *baseJob) unlock() {
	j.l.Unlock()
}

func (j *baseJob) setFailed(failed bool) {
	j.f = failed
}

func (j *baseJob) failed() bool {
	return j.f
}

func (j *baseJob) enabled() bool {
	return true
}

func (j *baseJob) setStartTime(t time.Time) {
	j.startTime = t
	if !j.startTime.IsZero() && j.intervalDur > 0 {
		j.expectedFinishedTime = j.startTime.Add(j.intervalDur)
	}
}

func (j *baseJob) interval() time.Duration {
	return j.intervalDur
}

func (j *baseJob) setInterval(i time.Duration) {
	j.intervalDur = i
}
