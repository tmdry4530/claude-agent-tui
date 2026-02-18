package collector

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
	"github.com/fsnotify/fsnotify"
)

// Collector는 이벤트 수집기 인터페이스입니다.
type Collector interface {
	Start(ctx context.Context) error
	Events() <-chan schema.RawEvent
	Stop()
}

// FileCollector는 파일 기반 이벤트 수집기 구현체입니다.
// fsnotify를 사용하여 JSONL 파일을 tail-follow 방식으로 읽습니다.
type FileCollector struct {
	watchPath string
	events    chan schema.RawEvent
	stopOnce  sync.Once
	wg        sync.WaitGroup
	cancel    context.CancelFunc

	// circuit-breaker 상태
	mu            sync.Mutex
	failCount     int
	backoffUntil  time.Time
	backoffLevels []time.Duration
}

// NewFileCollector는 새로운 FileCollector를 생성합니다.
func NewFileCollector(watchPath string) *FileCollector {
	return &FileCollector{
		watchPath: watchPath,
		events:    make(chan schema.RawEvent, 1000),
		backoffLevels: []time.Duration{
			10 * time.Second,
			30 * time.Second,
			60 * time.Second,
		},
	}
}

// Start는 파일 감시를 시작합니다.
func (fc *FileCollector) Start(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("fsnotify watcher 생성 실패: %w", err)
	}

	err = watcher.Add(fc.watchPath)
	if err != nil {
		watcher.Close()
		return fmt.Errorf("경로 감시 추가 실패 (%s): %w", fc.watchPath, err)
	}

	// 내부 context 생성 (Stop()에서 취소 가능)
	internalCtx, cancel := context.WithCancel(ctx)
	fc.cancel = cancel

	// 기존 파일들의 현재 위치를 추적
	filePositions := make(map[string]int64)

	fc.wg.Add(1)
	go func() {
		defer fc.wg.Done()
		defer watcher.Close()

		for {
			select {
			case <-internalCtx.Done():
				return

			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// WRITE 또는 CREATE 이벤트만 처리
				if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
					continue
				}

				// circuit-breaker 확인
				if fc.isBackedOff() {
					continue
				}

				// 파일 읽기 시도
				if err := fc.readNewLines(event.Name, filePositions); err != nil {
					fc.recordFailure()
				} else {
					fc.recordSuccess()
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fc.recordFailure()
				// 에러는 무시하고 계속 진행 (장애 격리)
				_ = err
			}
		}
	}()

	return nil
}

// Events는 수집된 이벤트 채널을 반환합니다.
func (fc *FileCollector) Events() <-chan schema.RawEvent {
	return fc.events
}

// Stop은 수집기를 중지합니다. 최대 5초 대기 후 강제 종료합니다.
func (fc *FileCollector) Stop() {
	fc.stopOnce.Do(func() {
		// context 취소로 goroutine 종료 유도
		if fc.cancel != nil {
			fc.cancel()
		}
		// goroutine 종료 대기 (타임아웃 포함)
		done := make(chan struct{})
		go func() {
			fc.wg.Wait()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
		}
		// 채널 닫기
		close(fc.events)
	})
}

// isBackedOff는 현재 backoff 상태인지 확인합니다.
func (fc *FileCollector) isBackedOff() bool {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if time.Now().Before(fc.backoffUntil) {
		return true
	}
	return false
}

// recordFailure는 실패를 기록하고 필요시 backoff를 설정합니다.
func (fc *FileCollector) recordFailure() {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.failCount++
	if fc.failCount >= 3 {
		level := fc.failCount - 3
		if level >= len(fc.backoffLevels) {
			level = len(fc.backoffLevels) - 1
		}
		fc.backoffUntil = time.Now().Add(fc.backoffLevels[level])
	}
}

// recordSuccess는 성공을 기록하고 실패 카운터를 초기화합니다.
func (fc *FileCollector) recordSuccess() {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.failCount = 0
	fc.backoffUntil = time.Time{}
}

// readNewLines는 파일의 새로운 라인들을 읽어 이벤트로 변환합니다.
func (fc *FileCollector) readNewLines(filePath string, positions map[string]int64) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("파일 열기 실패: %w", err)
	}
	defer file.Close()

	// 이전 위치로 이동
	startPos := positions[filePath]
	if _, err := file.Seek(startPos, io.SeekStart); err != nil {
		return fmt.Errorf("파일 위치 이동 실패: %w", err)
	}

	scanner := bufio.NewScanner(file)
	var bytesProcessed int64

	for scanner.Scan() {
		line := scanner.Bytes()
		// 처리된 바이트 수 추적 (라인 내용 + 개행 문자)
		bytesProcessed += int64(len(line)) + 1

		if len(line) == 0 {
			continue
		}

		// JSON 파싱 시도
		var data json.RawMessage
		if err := json.Unmarshal(line, &data); err != nil {
			// 잘못된 JSON은 건너뛰고 에러 카운트 증가
			fc.recordFailure()
			continue
		}

		// RawEvent 생성 및 전송
		event := schema.RawEvent{
			Source:   filePath,
			Data:     data,
			Received: time.Now(),
		}

		// buffered channel이 가득 찬 경우 이벤트 드롭
		select {
		case fc.events <- event:
		default:
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("파일 읽기 중 오류: %w", err)
	}

	// 처리된 완전한 라인 기준으로 위치 저장 (불완전한 라인 재읽기 보장)
	positions[filePath] = startPos + bytesProcessed

	return nil
}
