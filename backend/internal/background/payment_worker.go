package background

import (
	"context"
	"go-expense-management-system/internal/model"
	"time"

	"github.com/sirupsen/logrus"
)

type PaymentProcessorFunc func(context.Context, model.PaymentJob) error

type PaymentWorker struct {
	jobs       chan model.PaymentJob
	log        *logrus.Logger
	retryCount int
	retryDelay time.Duration
	timeout    time.Duration
	processFn  PaymentProcessorFunc
}

func NewPaymentWorker(
	buffer int,
	retryCount int,
	retryDelay time.Duration,
	timeout time.Duration,
	log *logrus.Logger,
	processFn PaymentProcessorFunc,
) *PaymentWorker {
	if buffer <= 0 {
		buffer = 100
	}
	if retryCount <= 0 {
		retryCount = 3
	}
	if retryDelay <= 0 {
		retryDelay = time.Second
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &PaymentWorker{
		jobs:       make(chan model.PaymentJob, buffer),
		log:        log,
		retryCount: retryCount,
		retryDelay: retryDelay,
		timeout:    timeout,
		processFn:  processFn,
	}
}

func (w *PaymentWorker) Start() {
	go func() {
		for job := range w.jobs {
			w.handleJob(job)
		}
	}()
}

func (w *PaymentWorker) Enqueue(job model.PaymentJob) bool {
	select {
	case w.jobs <- job:
		return true
	default:
		if w.log != nil {
			w.log.Warnf("Payment queue full, dropping job for expense %s", job.ExpenseID)
		}
		return false
	}
}

func (w *PaymentWorker) handleJob(job model.PaymentJob) {
	for attempt := 1; attempt <= w.retryCount; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
		err := w.processFn(ctx, job)
		cancel()

		if err == nil {
			return
		}

		if w.log != nil {
			w.log.Warnf("Payment job failed (attempt %d/%d) for %s: %+v", attempt, w.retryCount, job.ExpenseID, err)
		}

		if attempt < w.retryCount {
			time.Sleep(w.retryDelay * time.Duration(attempt))
		}
	}
}
