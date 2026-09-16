package task

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
)

var (
	mu      sync.Mutex
	enabled = map[uint]models.Task{}
	lastRun = map[uint]time.Time{}
	ticker  *time.Ticker
	started bool
)

// handlers is the registry of built-in function tasks, keyed by handler name.
var handlers = map[string]func() (string, bool){
	"heartbeat": func() (string, bool) {
		return "ok", true
	},
	"clear_operation_logs": func() (string, bool) {
		global.DB.Where("created_at < ?", time.Now().AddDate(0, 0, -30)).Delete(&models.OperationLog{})
		return "cleared operation logs older than 30 days", true
	},
}

// Init loads enabled tasks from the database and starts the scheduler.
func Init() {
	loadFromDB()
	start()
}

func loadFromDB() {
	var tasks []models.Task
	global.DB.Where("status = ?", 1).Find(&tasks)
	mu.Lock()
	defer mu.Unlock()
	enabled = make(map[uint]models.Task, len(tasks))
	for _, t := range tasks {
		enabled[t.ID] = t
	}
}

func start() {
	mu.Lock()
	if started {
		mu.Unlock()
		return
	}
	started = true
	mu.Unlock()

	ticker = time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			tick()
		}
	}()
}

func tick() {
	now := time.Now()

	mu.Lock()
	tasks := make([]models.Task, 0, len(enabled))
	for _, t := range enabled {
		tasks = append(tasks, t)
	}
	mu.Unlock()

	for _, t := range tasks {
		if !Match(t.Spec, now) {
			continue
		}
		mu.Lock()
		last, ok := lastRun[t.ID]
		mu.Unlock()
		if ok && last.Truncate(time.Minute).Equal(now.Truncate(time.Minute)) {
			continue // already fired this minute
		}
		mu.Lock()
		lastRun[t.ID] = now
		mu.Unlock()
		go execute(t)
	}
}

// AddTask registers or updates a task in the live scheduler.
func AddTask(t models.Task) {
	mu.Lock()
	defer mu.Unlock()
	if t.Status == 1 {
		enabled[t.ID] = t
	} else {
		delete(enabled, t.ID)
	}
}

// RemoveTask stops and forgets a task.
func RemoveTask(id uint) {
	mu.Lock()
	defer mu.Unlock()
	delete(enabled, id)
	delete(lastRun, id)
}

// SetEnabled toggles a task on/off in the live scheduler.
func SetEnabled(t models.Task) { AddTask(t) }

// RunNow executes a task immediately, regardless of its schedule.
func RunNow(t models.Task) { go execute(t) }

func execute(t models.Task) {
	start := time.Now()
	output, ok := runJob(t)
	dur := time.Since(start).Milliseconds()
	now := time.Now()

	global.DB.Model(&models.Task{}).Where("id = ?", t.ID).Update("last_run_at", now)
	global.DB.Create(&models.TaskLog{
		TaskID:     t.ID,
		TaskName:   t.Name,
		Status:     boolToInt(ok),
		Output:     output,
		DurationMs: dur,
	})
}

func runJob(t models.Task) (string, bool) {
	switch t.Type {
	case 1:
		return runHTTP(t)
	case 2:
		return runFunc(t)
	default:
		return "unknown task type", false
	}
}

func runHTTP(t models.Task) (string, bool) {
	method := t.Method
	if method == "" {
		method = http.MethodGet
	}
	req, err := http.NewRequest(method, t.Url, nil)
	if err != nil {
		return err.Error(), false
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err.Error(), false
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)), resp.StatusCode < 400
}

func runFunc(t models.Task) (string, bool) {
	h, ok := handlers[t.Handler]
	if !ok {
		return "no such handler: " + t.Handler, false
	}
	return h()
}

// Match reports whether the 5-field cron spec fires at time t.
// Fields: minute hour day-of-month month day-of-week. Supports *, ?, ,, -, /.
func Match(spec string, t time.Time) bool {
	fields := strings.Fields(spec)
	if len(fields) != 5 {
		return false
	}
	if !matchField(fields[0], t.Minute(), 0, 59) {
		return false
	}
	if !matchField(fields[1], t.Hour(), 0, 23) {
		return false
	}
	if !matchField(fields[2], t.Day(), 1, 31) {
		return false
	}
	if !matchField(fields[3], int(t.Month()), 1, 12) {
		return false
	}
	if !matchField(fields[4], int(t.Weekday()), 0, 6) {
		return false
	}
	return true
}

func matchField(field string, val, min, max int) bool {
	if field == "*" || field == "?" {
		return true
	}
	for _, part := range strings.Split(field, ",") {
		if matchPart(part, val, min, max) {
			return true
		}
	}
	return false
}

func matchPart(part string, val, min, max int) bool {
	step := 1
	if i := strings.Index(part, "/"); i >= 0 {
		fmt.Sscanf(part[i+1:], "%d", &step)
		part = part[:i]
		if step < 1 {
			step = 1
		}
	}
	if part == "*" {
		return (val-min)%step == 0
	}
	if strings.Contains(part, "-") {
		bounds := strings.SplitN(part, "-", 2)
		a, _ := strconv.Atoi(bounds[0])
		b, _ := strconv.Atoi(bounds[1])
		if b < a {
			a, b = b, a
		}
		if val < a || val > b {
			return false
		}
		return (val-a)%step == 0
	}
	n, err := strconv.Atoi(part)
	if err != nil {
		return false
	}
	return n == val
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
