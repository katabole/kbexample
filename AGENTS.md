# AGENTS.md

This file provides guidance to AI agents (Claude, Codex, Cursor, etc.) when working in this repository.

Read the README.md for high level information about the project.

## Architecture Overivew

This project uses the [Katabole server framework](https://github.com/katabole/katabole).

## Code Style

### Comments

- Comments should explain *why*, not *what*.
- Focus on non-obvious details, relevant design, intentional trade-offs, or important context.
- Do not add comments that just summarize or repeat the code.
- If the code is self-explanatory, no comment is needed.

### Variables

- Inline single use variables at their use sites if the expression is straight forward.
- Intermediate variables are useful when they clarify intent or when the expression is complex.
- Prefer defining and documenting constants over inline use of constant literals. The exception is for simple, trivial constants (e.g., `0`) if their use is self-explanatory.

## Tests

Always add or update tests when appropriate for a change.

The majority of your tests should be unit tests covering representative, realistic flows—do not litter the test files with every possible test, we are not fuzz testing here (though if you believe fuzz testing is necessary for a particularly security-sensitive component, you may recommend it)—as well as representative and important edge cases.

Tests should be human-readable and easy to follow along, and easily understood by a reader to be correct by visual inspection. So don't add too many abstractions or too much indirection in your tests.

For test cases covering one logical concept but with multiple potential sub-cases to be covered, follow the [subtest pattern](https://go.dev/blog/subtests) to parameterize a test case.

Follow the "arrange, act, assert" pattern.

After unit tests, the minority of your tests should be integration tests.

## Security

This is a server framework and template. Consequently, any issues here will be propagated and induce bugs in downstream dependents, so care needs to be taken with respect to correctness and especially security.

Always check for potential security issues, e.g., IDOR, SSRF, other confused deputy or abuse of server ambient authority issues, code injection or code execution, both preexisting and in new code.

Adhere to best practices and defense-in-depth. When planning, verify there is a clear authn and authz model and you understand it, and that it's consistently used throughout.

## Workflow Rules

Make sure you explore the code base and understand it before jumping into planning or coding.

Claude: Prefer using the Go LSP (gopls) plugin to navigate rather than `grep`ing  or `find`ing symbols across the Go code.

Explore, interview and ask clarifying questions (if necessary), *then* plan. If there are open questions and multiple approaches, make a recommendation and explain tradeoffs.

### Making changes

- **Never commit directly to `main`.** Always create a feature branch and open a PR.
- Always `go fmt` changed files.

### Submitting changes

Your commit messages and PR descriptions must be **concise**, **high level** summaries, and contain only the most salient information a reader should know.

Provide context, describe the problem or goal and motivation, summarize any design decisions made and any principles involved, if relevant, and summarize the key changes (existing behavior if relevant, new behavior), any outstanding TODOs or followups for the future. For each of these sections, only use them if they're relevant to the change.

DO NOT go into specific lines of code, specific files or line numbers. Keep it high level, focused on high level concepts, semantics, systems involved.

Don't add extra fluff about "Changes added: full test coverage" or "Test plan: add unit tests, all tests passed." It is taken for granted every change comes with relevant tests and they all pass. Only call out changes to tests if it's really important to know, e.g., fixing broken or incorrect tests, fixing gaps, etc.
