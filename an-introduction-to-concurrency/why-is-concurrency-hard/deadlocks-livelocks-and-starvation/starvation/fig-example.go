package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup        // 1. 두 개의 고루틴이 완료될 때까지 기다리기 위한 WaitGroup 선언
	var sharedLock sync.Mutex    // 2. 두 고루틴이 공유하는 뮤텍스 선언
	const runtime = 1 * time.Second // 3. 고루틴이 실행될 시간(1초) 상수로 선언

	greedyWorker := func() { // 4. "탐욕스러운" 작업자 함수 정의
		defer wg.Done() // 5. 함수 종료 시 WaitGroup에 완료 신호를 보냄

		var count int // 6. 작업자가 몇 번 반복했는지를 기록할 카운터 변수
		for begin := time.Now(); time.Since(begin) <= runtime; { // 7. 시작 시간부터 지정된 시간(runtime)까지 반복
			sharedLock.Lock()         // 8. 뮤텍스를 잠금
			time.Sleep(3 * time.Nanosecond) // 9. 3나노초 동안 작업(잠금 유지)
			sharedLock.Unlock()       // 10. 뮤텍스 잠금 해제
			count++                   // 11. 작업 반복 횟수를 증가
		}

		fmt.Printf("Greedy worker was able to execute %v work loops\n", count) // 12. 탐욕스러운 작업자의 작업 횟수 출력
	}

	politeWorker := func() { // 13. "예의 바른" 작업자 함수 정의
		defer wg.Done() // 14. 함수 종료 시 WaitGroup에 완료 신호를 보냄

		var count int // 15. 작업자가 몇 번 반복했는지를 기록할 카운터 변수
		for begin := time.Now(); time.Since(begin) <= runtime; { // 16. 시작 시간부터 지정된 시간(runtime)까지 반복
			sharedLock.Lock()         // 17. 뮤텍스를 잠금
			time.Sleep(1 * time.Nanosecond) // 18. 1나노초 동안 작업(잠금 유지)
			sharedLock.Unlock()       // 19. 뮤텍스 잠금 해제

			sharedLock.Lock()         // 20. 다시 뮤텍스를 잠금
			time.Sleep(1 * time.Nanosecond) // 21. 1나노초 동안 작업(잠금 유지)
			sharedLock.Unlock()       // 22. 뮤텍스 잠금 해제

			sharedLock.Lock()         // 23. 다시 뮤텍스를 잠금
			time.Sleep(1 * time.Nanosecond) // 24. 1나노초 동안 작업(잠금 유지)
			sharedLock.Unlock()       // 25. 뮤텍스 잠금 해제

			count++                   // 26. 작업 반복 횟수를 증가
		}

		fmt.Printf("Polite worker was able to execute %v work loops.\n", count) // 27. 예의 바른 작업자의 작업 횟수 출력
	}

	wg.Add(2)           // 28. 두 개의 고루틴을 기다리기 위해 WaitGroup 카운트 증가
	go greedyWorker()   // 29. 탐욕스러운 작업자를 고루틴으로 실행
	go politeWorker()   // 30. 예의 바른 작업자를 고루틴으로 실행

	wg.Wait()           // 31. 두 고루틴이 모두 종료될 때까지 대기
}
