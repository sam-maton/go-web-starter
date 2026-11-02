# Go Web Starter: learning and implementation plan

## Goal

Build `gws`, a Go CLI that asks a few questions and generates a working web
application from a bundled, known-working version of the baseline.

This is a learning project. Implement it yourself, one task at a time. The plan
provides outcomes and checks, not implementation solutions. Keep the architecture
small and introduce abstractions only when a task actually needs them.

## Agreed first-release scope

| Area | Decision |
| --- | --- |
| Questions | Project directory name, Go module path, include authentication |
| Interaction | Interactive prompts first; flags are a later milestone |
| Generator dependencies | Prefer the Go standard library |
| Template source | Bundle a tested baseline snapshot inside the executable |
| Authentication enabled | Sign up, log in, log out; authentication protects the shared todo list |
| Authentication disabled | Anyone can view and modify the shared todo list; no user accounts |
| Todo ownership | One shared list in both variants, not a list per user |
| Database | SQLite in both variants |
| Sessions | Keep sessions and session-backed success messages in both variants |
| Setup automation | Generate files and print instructions; do not run setup commands |
| Baseline repairs | Fix and verify the baseline separately before importing a snapshot |
| Task style | Small tasks with learning objectives and completion checks; no implementation code |

Generating files must work without network access or a checkout of the baseline.
Installing the CLI or resolving the generated app's dependencies may still need
network access. The generator does not promise that the web app's external font
resources work offline.

### Deliberately outside v1

- PostgreSQL or selectable databases.
- Optional todo example, per-user todos, or additional authentication providers.
- Automatic dependency installation, migrations, formatting, Git initialisation,
  or server startup.
- Downloading templates at generation time.
- Updating an already-generated project or overwriting an existing directory.
- A terminal UI framework, plugin system, or general-purpose scaffolding engine.
- Claims of production readiness or tested cross-platform support.

## Starting point

- The current [CLI entry point](cmd/gws/main.go) collects three answers using
  `fmt.Scan`, then prints the answers. It does not generate files yet.
- The [CLI module](go.mod) declares Go `1.27.0`; the
  [baseline module](../go-web-starter-basline/go.mod) declares Go `1.25.5`.
  Establish a tested toolchain policy rather than changing versions arbitrarily.
- The baseline directory is spelled
  [`go-web-starter-basline`](../go-web-starter-basline), while its Go module path
  uses `go-web-starter-baseline`. Use the real directory name when working locally.
- The baseline's [routes](../go-web-starter-basline/cmd/web/routes.go) protect
  todos, but its [todo model](../go-web-starter-basline/internal/models/todos.go)
  and [schema](../go-web-starter-basline/sql/schema/00003_create_todos_table.sql)
  have no user ownership. Preserve that shared-list behavior.
- Authentication spans startup, handlers, middleware, helpers, templates,
  dependencies, and migrations. An authentication toggle cannot just omit the
  login page.

## Working method

- Work through the milestones in order, following each task's dependencies.
- Each checkbox is a task, not a requirement to finish everything in one sitting.
- Add focused tests alongside the behavior being introduced. The later test
  milestone connects those checks into end-to-end coverage.
- When a task is complete, tick it and make a small, meaningful commit.
- If a task grows too large, split it at an observable outcome rather than by
  arbitrary file boundaries.
- Do not add a package just to practise package organisation. A small prompt layer
  and generator package are enough to begin.
- Change this plan explicitly if a requirement changes; do not silently expand v1.

## Milestone 0: establish a working source application

These tasks belong to the baseline repository, apart from recording decisions in
this plan. Do not mix baseline repairs with the CLI implementation.

### T01 - Establish prerequisites and reproduce current behavior
- [ ] Complete
- **Depends on:** None.
- **Learn:** Go toolchains, modules, CGO, and repeatable local setup.
- **Work:** Run both projects with a compatible Go toolchain. Check the C compiler
  required by the baseline's SQLite driver. Set up Goose using the baseline's
  documented migration workflow. Record which Go versions the CLI and generated
  apps will support and test.
- **Done when:** You can run the existing prompts, migrate a fresh baseline
  database, and start the baseline using documented commands.

### T02 - Make authentication actually stop unauthenticated requests
- [ ] Complete
- **Depends on:** T01.
- **Learn:** HTTP middleware, redirects, and `net/http/httptest`.
- **Work:** Repair the baseline's authentication middleware so it does not call
  the protected handler after redirecting. Verify that session cookies work under
  the documented local-development setup without weakening production guidance.
- **Done when:** A test proves an unauthenticated request is redirected without
  invoking the next handler; a signed-in request reaches it. Signup, login, and
  logout also work in the documented local setup.

### T03 - Make todo deletion real and verify the todo lifecycle
- [ ] Complete
- **Depends on:** T01.
- **Learn:** SQL mutations, model errors, and handler/model tests.
- **Work:** Add the missing model-backed delete behavior and connect it to the
  Delete action. Ensure invalid IDs return immediately and missing records receive
  an appropriate response.
- **Done when:** Create, list, complete, and delete work against a fresh database.
  Tests show deletion changes stored data and invalid/missing IDs do not report
  success.

### T04 - Verify the baseline as a usable starter
- [ ] Complete
- **Depends on:** T02, T03.
- **Learn:** Smoke testing and distinguishing a compiling app from a working app.
- **Work:** Check signup/login/logout, the shared todo lifecycle, flash messages,
  fresh migrations, and local static assets. Resolve missing local asset
  references in the base page by supplying the intended asset or removing an
  obsolete reference. Keep unrelated redesigns out of this work.
- **Done when:** A fresh checkout can follow its README to a working application.
  Relevant tests pass, local asset requests succeed, and the snapshot source is a
  specific baseline commit.

**Milestone check:** There is a baseline worth copying. Do not use the generator
to conceal application bugs.

## Milestone 1: make the questionnaire reliable

### T05 - Define configuration and v1 input rules
- [ ] Complete
- **Depends on:** None.
- **Learn:** Structs, zero values, validation, and separating data from interaction.
- **Work:** Define one configuration value for the directory, module path, and
  authentication choice. Document defaults and a conservative supported input
  policy. For v1, accept a single new directory name under the current working
  directory, not arbitrary destination paths. Use it as the app's display name.
- **Done when:** Each option has an unambiguous meaning. Empty names, path
  separators, `.`/`..`, control characters, and unsupported module paths have
  explicit validation outcomes. Inputs cannot inject source or template syntax.

### T06 - Replace token scanning with line-based prompts
- [ ] Complete
- **Depends on:** T05.
- **Learn:** Buffered input, readers/writers, and testable functions.
- **Work:** Read complete answers and make the prompt code accept input/output
  streams so tests do not need a real terminal. Retain the current simple welcome
  and question flow.
- **Done when:** Scripted input can exercise all three questions. Pressing Enter
  selects a documented default; read errors and premature EOF are not treated as
  successful answers. A final answer without a newline is handled deliberately.

### T07 - Validate answers and handle yes/no choices
- [ ] Complete
- **Depends on:** T05, T06.
- **Learn:** Table-driven tests, string handling, and user-facing errors.
- **Work:** Trim surrounding whitespace, validate directory/module answers, and
  accept documented case-insensitive yes/no forms. Re-prompt on invalid answers.
  Derive the suggested module path from the chosen directory.
- **Done when:** Tests cover defaults, valid values, invalid names/module paths,
  all accepted yes/no forms, unrecognised answers, and EOF during a retry.

### T08 - Separate orchestration from process exit
- [ ] Complete
- **Depends on:** T06, T07.
- **Learn:** Error propagation, standard output/error, and small program boundaries.
- **Work:** Keep process exit decisions at the entry point; let the command flow
  return errors. Report failures clearly and print success only after generation
  has actually completed. Retain a readable plain-text experience.
- **Done when:** Tests distinguish successful completion from input failures
  without terminating the test process. Diagnostic and normal output have
  consistent destinations.

**Milestone check:** The questionnaire produces validated configuration, but
still does not write a project.

## Milestone 2: generate one small project safely

Start with a tiny fixture, not the entire baseline. This keeps filesystem and
rendering problems separate from application-variation problems.

### T09 - Decide which files are rendered and which are copied
- [ ] Complete
- **Depends on:** T05.
- **Learn:** `text/template`, literal file copying, and template boundaries.
- **Work:** Pick a clear suffix for generation templates. Store generated Go source
  using that suffix so Go tooling does not try to build it as CLI source. Render
  only files needing generator substitutions and copy other assets unchanged.
  Preserve the generated app's runtime `html/template` expressions.
- **Done when:** A small fixture renders a module path and copies an HTML template
  containing runtime template actions without interpreting them. If a file needs
  both kinds of templating, distinct delimiters or another explicit boundary are
  tested. Missing generation values cause an error.

### T10 - Bundle and enumerate the fixture
- [ ] Complete
- **Depends on:** T09.
- **Learn:** `embed`, `io/fs`, and packaged resources.
- **Work:** Embed the fixture and enumerate its intended output files. Account
  explicitly for dotfiles: ordinary embed directory patterns omit hidden files.
  Define how template suffixes map to destination names.
- **Done when:** Tests find all expected files, including a representative
  dotfile, with no template suffix in output names. A built executable can read
  its resources when run outside the repository.

### T11 - Create a destination without overwriting anything
- [ ] Complete
- **Depends on:** T05, T10.
- **Learn:** `os`, `path/filepath`, exclusive directory creation, and filesystem errors.
- **Work:** Create a new destination under the current working directory. Reject
  an existing file, directory, or symlink at that destination, including an empty
  directory. Ensure every planned output stays inside the destination.
- **Done when:** Tests use temporary directories and prove existing contents are
  unchanged. Unsafe names and fixture paths that escape the destination are
  rejected. Directory-creation failures are reported.

### T12 - Write the fixture and report failures safely
- [ ] Complete
- **Depends on:** T09, T10, T11.
- **Learn:** Rendering errors, file writes, error context, and cleanup ownership.
- **Work:** Render/copy files, creating needed subdirectories with appropriate
  permissions. Use a simple failure policy: remove only output created by this
  invocation; never delete a pre-existing path. If cleanup also fails, report it
  alongside the original error.
- **Done when:** Tests compare generated contents and exercise a failed render
  and write. Failure produces no success message, preserves unrelated files, and
  follows the documented cleanup policy.

### T13 - Connect prompts to fixture generation
- [ ] Complete
- **Depends on:** T08, T12.
- **Learn:** Composing small functions and testing a complete command flow.
- **Work:** Pass validated answers into the generator and produce the fixture.
  Keep the generator independent of terminal prompts.
- **Done when:** A test feeds answers into the command flow and verifies the
  resulting directory and files. The flow does not execute external commands.

**Milestone check:** You have a small, working scaffolder before introducing the
full web application.

## Milestone 3: generate the authentication-enabled baseline

### T14 - Import a clean, traceable baseline snapshot
- [ ] Complete
- **Depends on:** T04, T10.
- **Learn:** Reproducible inputs and maintaining copied source.
- **Work:** Import only intended source files, UI assets, migration SQL, module
  metadata, useful development configuration, and documentation. Record the source
  commit and manual refresh procedure. Keep the baseline repository unchanged by
  generation.
- **Done when:** The bundled snapshot contains no live database, credentials,
  environment files, repository metadata, build binaries, or temporary logs.
  Required dotfiles are present, generated Go files use the chosen template
  suffix, and an empty database location can be created from the migrations.

### T15 - Substitute the module path and app name
- [ ] Complete
- **Depends on:** T09, T14.
- **Learn:** Targeted substitutions and preserving valid Go/HTML.
- **Work:** Parameterise the generated module declaration, baseline-module imports,
  and intended app-title placeholder. Avoid blanket replacements across unrelated
  text or binary assets.
- **Done when:** Two different valid module paths generate correctly. Generated
  source imports reference the new module, the page title uses the project name,
  and runtime HTML template actions remain intact.

### T16 - Generate and compile the full authenticated app
- [ ] Complete
- **Depends on:** T12, T15.
- **Learn:** Generated-project integration tests and Go module boundaries.
- **Work:** Generate the authenticated variant into a temporary directory. In your
  development/test workflow, resolve its dependencies and run formatting checks,
  tests, vet, and a build. Repair generation-specific failures.
- **Done when:** The generated app passes those checks without hand-editing its
  source. The CLI itself still executes no Go, Git, migration, or shell commands
  during generation.

### T17 - Verify an authenticated generated app from scratch
- [ ] Complete
- **Depends on:** T16.
- **Learn:** Acceptance testing beyond compilation.
- **Work:** Follow the generated app's setup workflow with a new database. Exercise
  signup, login, shared todo operations, flash messages, and logout. Verify a
  second signed-in user sees the same list.
- **Done when:** Protected requests do not execute without authentication, all
  intended operations work, and the original baseline checkout is not needed.

**Milestone check:** One complete variant works before you add conditional output.

## Milestone 4: add the public-todo variant

### T18 - Map the authentication toggle across the app
- [ ] Complete
- **Depends on:** T17.
- **Learn:** Following dependencies through a real application.
- **Work:** Record a short inclusion/variation checklist near the bundled-template
  maintenance documentation. Cover user models, authentication forms/handlers,
  startup fields, middleware, helpers, template data, routes, navigation, HTML
  pages, migrations, constants, and module dependencies.
- **Done when:** Each affected surface has a clear auth-on/auth-off outcome.
  SQLite, sessions, flash messages, and shared todos are explicitly retained in
  both variants. No third project option is introduced.

### T19 - Select the right files and migrations
- [ ] Complete
- **Depends on:** T18.
- **Learn:** Explicit file selection and coherent generated project structure.
- **Work:** Exclude user-account models, login/signup pages, and the users
  migration when authentication is disabled. Preserve the sessions and todos
  migrations. Keep a single generation pipeline with a small, explicit selection
  rule rather than duplicating the whole starter.
- **Done when:** Tests assert the exact expected inclusion/exclusion of files for
  both variants. Fresh Goose migrations work for each variant despite any retained
  migration-number gaps.

### T20 - Vary startup, handlers, routes, and template data
- [ ] Complete
- **Depends on:** T18, T19.
- **Learn:** Conditional generation without leaving invalid references.
- **Work:** Remove authentication-only wiring from public output. Allow public
  read/create/complete/delete requests, while keeping session middleware around
  handlers using flash messages. Make template data and helpers consistent with
  the remaining application fields.
- **Done when:** Both variants compile; public requests do not redirect to login.
  No public-variant code references absent user models, auth handlers, or
  authentication-only helpers. Anonymous flash messages still work.

### T21 - Vary navigation and the todo page
- [ ] Complete
- **Depends on:** T20.
- **Learn:** Coordinating server behavior and rendered UI.
- **Work:** Remove account navigation from public output and show the todo-create
  form to anonymous visitors. Preserve authenticated navigation in the protected
  variant. Do not merely hide links while leaving unwanted routes active.
- **Done when:** Both rendered interfaces match their access rules. In the public
  variant, signup/login/logout routes are absent, todos can be created and changed,
  and success messages display.

### T22 - Maintain correct dependencies for both variants
- [ ] Complete
- **Depends on:** T19, T20.
- **Learn:** Direct/transitive dependencies and module checksums.
- **Work:** Prepare tested module metadata for each variant. Keep session, SQLite,
  and todo form dependencies; omit unused authentication-only requirements.
  Resolve/tidy dependencies during template maintenance and testing, not when a
  user runs the generator.
- **Done when:** Both variants have usable module metadata and pass tests/builds.
  Tidying does not reveal stale required dependencies or missing requirements.

### T23 - Verify the public app from a fresh database
- [ ] Complete
- **Depends on:** T19, T20, T21, T22.
- **Learn:** Testing feature combinations and documented behavior.
- **Work:** Generate public output, migrate it, and exercise the complete todo
  lifecycle without signing in. Use separate browser sessions or clients to
  confirm the list is shared.
- **Done when:** Anyone can view and modify the same todos. There is no users
  table, users are never asked to authenticate, and session-backed success
  messages continue to work. Documentation clearly warns that the public demo
  intentionally gives all visitors write access.

**Milestone check:** Answering yes or no produces a different, complete, working
application, not one application with broken optional pieces.

## Milestone 5: make generation dependable and useful

### T24 - Add fast generator regression tests
- [ ] Complete
- **Depends on:** T13, T22.
- **Learn:** Focused fixtures, table-driven cases, and maintainable assertions.
- **Work:** Consolidate fast coverage for configuration, prompts, file selection,
  rendering, dotfiles, nested assets, import rewriting, path handling, and failure
  cleanup. Test representative contents rather than enormous brittle snapshots.
- **Done when:** `go test ./...` tests the CLI without downloads, a database
  server, a baseline checkout, or other external tools. Both authentication
  answers and different project/module names are covered.

### T25 - Make generated-project validation repeatable
- [ ] Complete
- **Depends on:** T17, T23, T24.
- **Learn:** Integration boundaries, subprocess testing, and explicit prerequisites.
- **Work:** Add a separately invoked integration workflow that generates both
  variants, resolves dependencies, checks formatting, runs tests/vet, and builds.
  Include repeatable fresh-migration and HTTP smoke checks for the behavior matrix
  below. Reuse suitable tests imported from the repaired baseline.
- **Done when:** One documented invocation validates both variants. Missing Go,
  CGO/compiler, Goose, or network prerequisites are reported clearly. Normal CLI
  use still never runs these commands. Temporary servers/databases are cleaned up
  and no checks modify a developer's real database.

### T26 - Print accurate, variant-aware next steps
- [ ] Complete
- **Depends on:** T17, T23.
- **Learn:** Designing useful CLI output and executable documentation.
- **Work:** On success, report the destination and selected variant, then print the
  exact commands to enter it, resolve dependencies, apply Goose migrations, and
  run the app. Explain prerequisites and local cookie configuration. Generate a
  matching app README.
- **Done when:** Following only those instructions starts either variant from
  scratch. Commands are printed, not executed; success is never printed after a
  generation failure.

### T27 - Document installation, usage, and template maintenance
- [ ] Complete
- **Depends on:** T08, T14, T26.
- **Learn:** Writing documentation for someone without your local environment.
- **Work:** Add the CLI README with install/run instructions, default answers,
  destination rules, examples of both variants, and troubleshooting. Document how
  to refresh the pinned baseline snapshot and repeat both-variant validation.
  Note that authentication protects a shared list, not user-owned data.
- **Done when:** Someone can build/install and use `gws` without either developer's
  checkout paths. The documentation explains which output is intentionally
  omitted, which setup commands are manual, and the supported test environment.

### T28 - Verify the packaged executable outside the repository
- [ ] Complete
- **Depends on:** T24, T25, T26, T27.
- **Learn:** Build artifacts, working directories, and release acceptance.
- **Work:** Build/install the CLI and invoke it from an unrelated temporary working
  directory with no baseline checkout. Verify offline generation for both answers.
  Recheck existing-directory refusal and run the complete validation workflow.
- **Done when:** The executable has everything needed to generate files, does
  not look up templates relative to the source repository, and both generated
  variants pass acceptance checks.

**Milestone check:** v1 is complete. Resist adding more options until this point.

## Generated-application acceptance matrix

| Check | Authentication enabled | Authentication disabled |
| --- | --- | --- |
| Fresh migrations | Users, sessions, todos | Sessions, todos; no users |
| Anonymous home request | Redirect; protected handler does not run | Public todo page |
| Signup/login/logout | Present and working | Routes and pages absent |
| Shared todo list | Shared by signed-in users | Shared by all visitors |
| Create/complete/delete | Requires authentication | Publicly available |
| Success messages | Session-backed | Session-backed, anonymous sessions |
| App navigation | Correct account actions | No account actions |
| Local assets | Referenced assets load | Referenced assets load |
| Module/source checks | Format, tidy, test, vet, build pass | Same checks pass |
| Baseline checkout needed | No | No |

## Definition of done for v1

- Every v1 task above is checked off with its completion condition verified.
- The generator itself uses the standard library unless a documented, justified
  need changed that decision.
- Both outputs work from fresh directories and databases.
- Generated projects have the chosen module path and display name, with no stale
  baseline imports, template suffixes, or unresolved generator placeholders.
- Runtime HTML template expressions are preserved.
- Existing paths are never overwritten, failures are explicit, and cleanup does
  not touch unrelated data.
- Normal generation does not execute setup commands or access the network.
- The built executable works without the CLI source or baseline repositories.
- Installation, manual setup, variant behavior, and template refresh are documented.

## Later milestone: flags, only after v1

These tasks are intentionally outside the first-release checklist.

### L01 - Define explicit non-interactive behavior
- [ ] Complete
- **Depends on:** T28.
- **Learn:** CLI contracts and configuration precedence.
- **Work:** Decide how flags interact with prompts and defaults, how to express
  authentication both on and off, and what happens when input is incomplete
  without an interactive terminal.
- **Done when:** These rules are documented before implementing flags.

### L02 - Add standard-library flags using the same generator
- [ ] Complete
- **Depends on:** L01.
- **Learn:** `flag`, reusable validation, and testing command arguments.
- **Work:** Add directory/module/auth options without duplicating validation or
  generation logic. Keep the original questionnaire as the default experience.
- **Done when:** Scripted commands generate both variants without prompts, while
  interactive behavior and all existing checks continue to pass.

### L03 - Add small help/version polish
- [ ] Complete
- **Depends on:** L02.
- **Learn:** Command discovery and release identification.
- **Work:** Add clear help, version reporting, and tested examples. Consider
  release automation or additional platforms only when there is a real need.
- **Done when:** Users can discover supported options and identify the installed
  CLI version without generating a project.

## First task to pick up

Start with **T01**, then repair and verify the baseline in **T02-T04**. This gives
you a trustworthy target. After that, **T05-T08** improve the small CLI you already
have without introducing file generation and template variation all at once.
