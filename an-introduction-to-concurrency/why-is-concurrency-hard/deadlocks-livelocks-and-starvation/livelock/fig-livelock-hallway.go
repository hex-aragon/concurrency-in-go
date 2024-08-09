package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	cadence := sync.NewCond(&sync.Mutex{})	// 1. 조건 변수를 생성하고 뮤텍스로 보호
	go func() { // 2. 고루틴을 시작하여 매 1밀리초마다 Broadcast 신호를 보냄
		for range time.Tick(1 * time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	takeStep := func() { // 3. 신호가 올 때까지 기다렸다가 뮤텍스를 해제하는 함수
		cadence.L.Lock() // 4. 뮤텍스를 잠금
		cadence.Wait()   // 5. 신호를 기다림
		cadence.L.Unlock() // 6. 신호가 오면 뮤텍스를 해제
	}

	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool { // 7. 한 방향으로 움직이려고 시도하는 함수
		fmt.Fprintf(out, " %v", dirName) // 8. 시도하려는 방향을 출력
		atomic.AddInt32(dir, 1)          // 9. 해당 방향으로 시도한 것을 원자적으로 증가
		takeStep()                       // 10. 신호가 올 때까지 대기
		if atomic.LoadInt32(dir) == 1 {  // 11. 해당 방향으로 첫 번째 시도인지 확인
			fmt.Fprint(out, ". Success!") // 12. 성공 메시지 출력
			return true                   // 13. 성공하면 true 반환
		}
		takeStep()                    // 14. 다시 신호가 올 때까지 대기
		atomic.AddInt32(dir, -1)      // 15. 실패한 시도를 원자적으로 감소
		return false                  // 16. 실패하면 false 반환
	}

	var left, right int32 // 17. 왼쪽과 오른쪽 시도를 카운트하는 변수 선언
	tryLeft := func(out *bytes.Buffer) bool { return tryDir("left", &left, out) }   // 18. 왼쪽으로 이동하려고 시도하는 함수
	tryRight := func(out *bytes.Buffer) bool { return tryDir("right", &right, out) } // 19. 오른쪽으로 이동하려고 시도하는 함수
	walk := func(walking *sync.WaitGroup, name string) { // 20. 사람이 이동을 시도하는 함수
		var out bytes.Buffer // 21. 결과 출력을 저장할 버퍼
		defer func() { fmt.Println(out.String()) }() // 22. 함수 종료 시 결과를 출력
		defer walking.Done() // 23. 함수 종료 시 WaitGroup에 완료 신호를 보냄
		fmt.Fprintf(&out, "%v is trying to scoot:", name) // 24. 시도 시작 메시지 출력
		for i := 0; i < 5; i++ { // 25. 최대 5번 시도
			if tryLeft(&out) || tryRight(&out) { // 26. 왼쪽 또는 오른쪽으로 이동 시도
				return // 27. 성공하면 함수 종료
			}
		}
		fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name) // 28. 모든 시도가 실패했을 때 메시지 출력
	}

	var peopleInHallway sync.WaitGroup // 29. 두 명의 사람이 움직임을 완료할 때까지 대기할 WaitGroup 생성
	peopleInHallway.Add(2) // 30. 두 명의 사람이 움직임을 시도할 것이므로 WaitGroup에 2를 추가
	go walk(&peopleInHallway, "Alice") // 31. Alice가 이동을 시도하는 고루틴 시작
	go walk(&peopleInHallway, "Barbara") // 32. Barbara가 이동을 시도하는 고루틴 시작
	peopleInHallway.Wait() // 33. 두 사람의 이동이 끝날 때까지 대기
}
