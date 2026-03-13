# 바이브 코딩 전환 계획서

## 현재 상태 분석

| 항목 | 초기 상태 | 최종 상태 |
|------|----------|----------|
| 스펙 문서 | 없음 | OpenSpec specs/ **5개** |
| 커스텀 스킬 | 없음 | **7개** (controller, template, integration-test, e2e-test, model, debug-htmx, add-feature) |
| 단위 테스트 | 0% | **84.9%** 커버리지 (controller 89.4%, htmx 83.7%, model 85.2%, internal 100%) |
| 통합 테스트 | 없음 | **65+개** Fiber app.Test() 기반 |
| E2E 테스트 | 없음 | **16개** Playwright (인증, 기사 CRUD, 피드) |
| 린트 | 없음 | golangci-lint v2 **0 issues** |
| CI/CD (로컬) | 없음 | lefthook pre-commit + pre-push |
| CI/CD (원격) | 없음 | GitHub Actions (build, test, lint, e2e) |
| 빌드 자동화 | 없음 | Makefile **11개** 타겟 |
| README | 기본 설치 안내만 | 팀 온보딩 가이드 (환경 설정, 구조, Hook, CI) |

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

### Week 1: 문서화 기반 구축 ✅

**목표:** OpenSpec 스펙 문서화 + 커스텀 스킬 생성

| 태스크 | 산출물 | 상태 |
|--------|--------|------|
| OpenSpec config 확장 | `openspec/config.yaml` | ✅ |
| 아키텍처 스펙 | `openspec/specs/architecture.md` | ✅ |
| 데이터 모델 스펙 | `openspec/specs/data-model.md` | ✅ |
| 라우트 스펙 | `openspec/specs/routes.md` | ✅ |
| 인증 스펙 | `openspec/specs/authentication.md` | ✅ |
| 템플릿 스펙 | `openspec/specs/templates.md` | ✅ |
| Controller 스킬 | `.claude/skills/fiber-htmx-controller/` | ✅ |
| Template 스킬 | `.claude/skills/fiber-htmx-template/` | ✅ |

### Week 2: 테스트 인프라 + 단위 테스트 ✅

**목표:** 테스트 프레임워크 구축 + 모델/유틸리티 단위 테스트

| 태스크 | 산출물 | 상태 |
|--------|--------|------|
| testify 의존성 추가 | `go.mod` 업데이트 | ✅ |
| 테스트 헬퍼 | `internal/testutil/apptest.go` | ✅ |
| User 모델 테스트 | `cmd/web/model/user_test.go` | ✅ |
| Article 모델 테스트 | `cmd/web/model/article_test.go` | ✅ |
| Validator 테스트 | `internal/validator_test.go` | ✅ |
| ErrorMessage 테스트 | `internal/errormessage_test.go` | ✅ |
| Makefile 생성 | `Makefile` | ✅ (11개 타겟) |
| golangci-lint 설정 | `.golangci.yml` | ✅ |

### Week 3: 통합 테스트 + 린트 정리 ✅

**목표:** HTTP 핸들러 통합 테스트 + 기존 코드 린트 통과

| 태스크 | 산출물 | 상태 |
|--------|--------|------|
| 페이지 컨트롤러 테스트 (7개 파일) | `cmd/web/controller/*_test.go` | ✅ 89.4% |
| HTMX 컨트롤러 테스트 (9개 파일) | `cmd/web/controller/htmx/*_test.go` | ✅ 83.7% |
| 린트 수정 (errcheck, goimports) | 기존 소스 코드 18건 수정 | ✅ 0 issues |
| 커버리지 확인 | 패키지별 80%+ | ✅ 84.9% |

### Week 4: CI/CD + E2E 테스트 ✅

**목표:** Git hook 설정 + Playwright E2E 테스트

| 태스크 | 산출물 | 상태 |
|--------|--------|------|
| lefthook 설치 | `lefthook.yml` | ✅ |
| pre-commit hook | `make check` (lint + test) | ✅ |
| pre-push hook | build + coverage-check + e2e | ✅ |
| Playwright 설치 | `e2e/` 디렉토리 | ✅ |
| 인증 E2E | `e2e/tests/auth.spec.ts` (6개) | ✅ |
| 기사 E2E | `e2e/tests/article.spec.ts` (5개) | ✅ |
| 피드 E2E | `e2e/tests/feed.spec.ts` (5개) | ✅ |
| GitHub Actions CI | `.github/workflows/ci.yml` | ✅ |

### Week 5: 마무리 + 검증 ✅

**목표:** 전환 계획서 완성 + 전체 검증

| 태스크 | 산출물 | 상태 |
|--------|--------|------|
| 추가 커스텀 스킬 5개 | `.claude/skills/` | ✅ |
| GitHub Issues 등록 | #1~#5 | ✅ (전체 클로즈) |
| 전환 계획서 최종 업데이트 | `migration-plan.md` | ✅ |
| CLAUDE.md 최종 업데이트 | `CLAUDE.md` | ✅ |
| README 온보딩 가이드 | `README.md` | ✅ |
| 바이브 코딩 튜토리얼 | `docs/vibe-coding-migration-tutorial.md` | ✅ |

## 위험 평가

### High Risk
| 위험 | 영향 | 대응 방안 | 결과 |
|------|------|----------|------|
| SQLite 동시성 제한 | 통합 테스트 간 DB 충돌 | in-memory DB 격리 (testutil) | ✅ 해결 |
| HTMX fragment 테스트 어려움 | E2E 커버리지 부족 | Playwright로 실제 브라우저 테스트 | ✅ 16개 통과 |
| 기존 코드 린트 위반 다수 | Week 3 지연 | 직접 수정 (nolint 없이) | ✅ 18건 수정 |

### Medium Risk
| 위험 | 영향 | 대응 방안 | 결과 |
|------|------|----------|------|
| golangci-lint 버전 호환성 | CI 불안정 | v2 형식 대응 (3차 수정) | ✅ 해결 |
| Playwright + Go 서버 통합 | E2E 셋업 복잡 | CI에서 서버 빌드→기동→테스트 | ✅ 해결 |
| lefthook 팀원 미설치 | hook 우회 | README에 설치 안내 + GitHub Actions CI | ✅ 이중 안전망 |

### Low Risk
| 위험 | 영향 | 대응 방안 | 결과 |
|------|------|----------|------|
| OpenSpec 학습 곡선 | 초기 생산성 저하 | explore → propose 단계별 도입 | ✅ 스킬로 자동화 |
| 스킬 패턴 변경 필요 | 스킬 업데이트 | SKILL.md 버전 관리 | ✅ 7개 스킬 안정 |

## 성공 기준 체크리스트

- [x] OpenSpec `specs/` 에 5개 이상 문서 — **5개** (architecture, data-model, routes, authentication, templates)
- [x] 커스텀 스킬 2개+ 동작 확인 — **7개** (controller, template + 5개 추가)
- [x] `make test` — 전체 단위 + 통합 테스트 통과 — **87개+ 테스트 PASS**
- [x] `make coverage-check` — 80%+ 커버리지 — **84.9%** (패키지별 전체 80%+)
- [x] `make lint` — golangci-lint 통과 — **0 issues**
- [x] `make build` — 빌드 성공
- [x] `lefthook run pre-commit` — 정상 동작 (`make check`)
- [x] `lefthook run pre-push` — 정상 동작 (build + coverage-check + e2e)
- [x] Playwright E2E 핵심 시나리오 3개 이상 통과 — **16개 통과**
- [x] GitHub Actions CI — **4개 job** (build, test, lint, e2e)
- [x] `migration-plan.md` 완성
- [x] `CLAUDE.md` 최종 업데이트
- [x] `README.md` 팀 온보딩 가이드

## 바이브 코딩 워크플로우 (도입 후)

도입이 완료되면 새 기능 개발 시 다음 워크플로우를 따릅니다:

1. **탐색** — `/opsx explore` 로 아이디어 탐색, 기존 코드 조사
2. **제안** — `/opsx propose` 로 변경 제안서 생성 (proposal → design → tasks)
3. **구현** — `/opsx apply` 로 태스크별 TDD 구현
   - 스킬 활용: `/add-feature`, `/gen-integration-test`, `/gen-e2e-test`, `/add-model`
   - 문제 발생 시: `/debug-htmx`
   - 테스트 먼저 작성 (RED → GREEN → REFACTOR)
4. **검증** — `make check` 로 품질 확인
5. **커밋** — lefthook이 pre-commit (lint + test) 자동 검증
6. **푸시** — lefthook이 pre-push (build + coverage + e2e) 자동 검증
7. **CI** — GitHub Actions가 원격에서 이중 검증
8. **아카이브** — `/opsx archive` 로 완료된 변경 정리
