# 바이브 코딩 전환 계획서

## 현재 상태 분석

| 항목 | 현재 상태 | 목표 상태 |
|------|----------|----------|
| 스펙 문서 | 없음 | OpenSpec specs/ 5개 이상 |
| 커스텀 스킬 | 없음 | 2개 (controller, template) |
| 단위 테스트 | 0% | 80%+ 커버리지 |
| 통합 테스트 | 없음 | Fiber app.Test() 기반 |
| E2E 테스트 | 없음 | Playwright 핵심 시나리오 3개+ |
| 린트 | 없음 | golangci-lint 통과 |
| CI/CD | 없음 | lefthook pre-commit/pre-push |
| 빌드 자동화 | 없음 | Makefile 9개 타겟 |

## 의존성 그래프

```
Epic 1 (문서화) ──────────────────→ Epic 4 (전환 계획서)
                                ↗
Epic 2 (테스트) ──→ Epic 3 (CI/CD) ──→ Epic 4
```

- Epic 1과 Epic 2는 병렬 진행 가능
- Epic 3은 Epic 2의 테스트/린트 인프라에 의존
- Epic 4는 모든 에픽 결과를 종합

## 주 단위 타임라인

### Week 1: 문서화 기반 구축

**목표:** OpenSpec 스펙 문서화 + 커스텀 스킬 생성

| 태스크 | 산출물 | 완료 기준 |
|--------|--------|----------|
| OpenSpec config 확장 | `openspec/config.yaml` | project context, artifact rules 정의 |
| 아키텍처 스펙 | `openspec/specs/architecture.md` | 듀얼 컨트롤러 패턴, 요청 흐름 문서화 |
| 데이터 모델 스펙 | `openspec/specs/data-model.md` | 전체 Entity 스키마, 관계 문서화 |
| 라우트 스펙 | `openspec/specs/routes.md` | 페이지 + HTMX 전체 라우트 목록 |
| 인증 스펙 | `openspec/specs/authentication.md` | 세션 기반 인증 흐름 문서화 |
| 템플릿 스펙 | `openspec/specs/templates.md` | 레이아웃, 커스텀 함수, HTMX 패턴 |
| Controller 스킬 | `.claude/skills/fiber-htmx-controller/` | 듀얼 컨트롤러 생성 패턴 |
| Template 스킬 | `.claude/skills/fiber-htmx-template/` | 템플릿 생성 패턴 |

### Week 2: 테스트 인프라 + 단위 테스트

**목표:** 테스트 프레임워크 구축 + 모델/유틸리티 단위 테스트

| 태스크 | 산출물 | 완료 기준 |
|--------|--------|----------|
| testify 의존성 추가 | `go.mod` 업데이트 | `go mod tidy` 성공 |
| 테스트 헬퍼 | `internal/testutil/testutil.go` | in-memory SQLite DB, fixture 생성 |
| User 모델 테스트 | `cmd/web/model/user_test.go` | HashPassword, CheckPassword, FollowedBy 등 |
| Article 모델 테스트 | `cmd/web/model/article_test.go` | GetFavoriteCount, FavoritedBy, Tags 등 |
| Validator 테스트 | `internal/validator_test.go` | 유효/무효 입력 검증 |
| ErrorMessage 테스트 | `internal/errormessage_test.go` | required, email, unknown 타입 |
| Makefile 생성 | `Makefile` | 9개 타겟 동작 확인 |
| golangci-lint 설정 | `.golangci.yml` | lint 통과 |

### Week 3: 통합 테스트 + 린트 정리

**목표:** HTTP 핸들러 통합 테스트 + 기존 코드 린트 통과

| 태스크 | 산출물 | 완료 기준 |
|--------|--------|----------|
| 홈페이지 통합 테스트 | `cmd/web/controller/home_test.go` | Fiber app.Test() 기반 |
| 로그인 통합 테스트 | `cmd/web/controller/htmx/sign_in_test.go` | 로그인 성공/실패 시나리오 |
| 기사 CRUD 통합 테스트 | `cmd/web/controller/htmx/article_test.go` | 생성/조회/수정/삭제 |
| 댓글 통합 테스트 | `cmd/web/controller/htmx/comment_test.go` | 작성/목록 조회 |
| 린트 수정 | 기존 소스 코드 | golangci-lint run 통과 |
| 커버리지 확인 | `coverage.html` | 80%+ 달성 |

### Week 4: CI/CD + E2E 테스트

**목표:** Git hook 설정 + Playwright E2E 테스트

| 태스크 | 산출물 | 완료 기준 |
|--------|--------|----------|
| lefthook 설치 | `lefthook.yml` | pre-commit, pre-push hook 동작 |
| pre-commit hook | gofmt + go vet + lint | staged .go 파일 대상 실행 |
| pre-push hook | build + test | 전체 빌드/테스트 통과 확인 |
| Playwright 설치 | `e2e/` 디렉토리 | npm init playwright |
| 인증 E2E | `e2e/auth.spec.ts` | 회원가입/로그인/로그아웃 |
| 기사 E2E | `e2e/article.spec.ts` | 작성/조회/수정/삭제 |
| 피드 E2E | `e2e/feed.spec.ts` | 글로벌/내 피드/태그 피드 |
| .gitignore 정리 | `.gitignore` | 빌드 아티팩트, IDE 파일 제외 |

### Week 5: 마무리 + 검증

**목표:** 전환 계획서 완성 + 전체 검증

| 태스크 | 산출물 | 완료 기준 |
|--------|--------|----------|
| 전환 계획서 | `migration-plan.md` | 전체 로드맵 문서화 |
| CLAUDE.md 최종 업데이트 | `CLAUDE.md` | 스킬, OpenSpec, 테스트 명령어 추가 |
| 전체 테스트 실행 | CI 파이프라인 | `make test` + `make lint` 통과 |
| E2E 테스트 실행 | Playwright | `make e2e` 통과 |
| lefthook 검증 | Git hooks | pre-commit/pre-push 동작 확인 |
| 팀 온보딩 가이드 | README 업데이트 | 개발 환경 설정 안내 |

## 위험 평가

### High Risk
| 위험 | 영향 | 대응 방안 |
|------|------|----------|
| SQLite 동시성 제한 | 통합 테스트 간 DB 충돌 | in-memory DB 격리 (testutil) |
| HTMX fragment 테스트 어려움 | E2E 커버리지 부족 | Playwright로 실제 브라우저 테스트 |
| 기존 코드 린트 위반 다수 | Week 3 지연 | nolint 주석으로 점진적 해결 |

### Medium Risk
| 위험 | 영향 | 대응 방안 |
|------|------|----------|
| golangci-lint 버전 호환성 | CI 불안정 | .golangci.yml에 Go 버전 고정 |
| Playwright + Go 서버 통합 | E2E 셋업 복잡 | Makefile에 서버 시작/종료 스크립트 |
| lefthook 팀원 미설치 | hook 우회 | README에 설치 안내, CI에서도 검증 |

### Low Risk
| 위험 | 영향 | 대응 방안 |
|------|------|----------|
| OpenSpec 학습 곡선 | 초기 생산성 저하 | explore → propose 단계별 도입 |
| 스킬 패턴 변경 필요 | 스킬 업데이트 | SKILL.md 버전 관리 |

## 성공 기준 체크리스트

- [ ] OpenSpec `specs/` 에 5개 이상 문서
- [ ] 커스텀 스킬 2개 동작 확인 (controller, template)
- [ ] `make test` — 전체 단위 + 통합 테스트 통과
- [ ] `make test-coverage` — 80%+ 커버리지
- [ ] `make lint` — golangci-lint 통과
- [ ] `make build` — 빌드 성공
- [ ] `lefthook run pre-commit` — 정상 동작
- [ ] `lefthook run pre-push` — 정상 동작
- [ ] Playwright E2E 핵심 시나리오 3개 이상 통과
- [ ] `migration-plan.md` 완성
- [ ] `CLAUDE.md` 최종 업데이트

## 바이브 코딩 워크플로우 (도입 후)

도입이 완료되면 새 기능 개발 시 다음 워크플로우를 따릅니다:

1. **탐색** — `/opsx explore` 로 아이디어 탐색, 기존 코드 조사
2. **제안** — `/opsx propose` 로 변경 제안서 생성 (proposal → design → tasks)
3. **구현** — `/opsx apply` 로 태스크별 TDD 구현
   - 스킬 활용: `fiber-htmx-controller`, `fiber-htmx-template`
   - 테스트 먼저 작성 (RED → GREEN → REFACTOR)
4. **검증** — `make test && make lint` 로 품질 확인
5. **커밋** — lefthook이 pre-commit/pre-push 자동 검증
6. **아카이브** — `/opsx archive` 로 완료된 변경 정리
