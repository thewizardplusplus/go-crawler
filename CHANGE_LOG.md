# Change Log

## [v1.4.1](https://github.com/thewizardplusplus/go-crawler/tree/v1.4.1) (2020-11-22)

- extract model of a sourced link:
  - use it in the `crawler.LinkChecker` interface;
  - use it in the `crawler.LinkHandler` interface.

## [v1.4](https://github.com/thewizardplusplus/go-crawler/tree/v1.4) (2020-10-02)

- crawling of all relative links for specified ones:
  - delayed extracting of relative links (optional):
    - reducing of a delay time by the time elapsed since the last request;
    - using of individual delays for each thread;
- refactoring:
  - extracting utility entities for syncing to the [single](https://github.com/thewizardplusplus/go-sync-utils) package;
  - passing a thread ID to an extractor;
  - `extractors.RepeatingExtractor` structure:
    - passing an abstract sleeper as a parameter;
    - checking of call count in tests.

## [v1.3](https://github.com/thewizardplusplus/go-crawler/tree/v1.3) (2020-09-22)

- calling of an outer handler for an each found link:
  - handling only of unique links (optional):
    - supporting of sanitizing of a link before checking of uniqueness (optional);
- refactoring:
  - extract the `sanitizing` package;
  - extract the `register.LinkRegister` structure;
  - add a function that makes it easier:
    - to specify initial links;
    - to wait for completion.

## [v1.2](https://github.com/thewizardplusplus/go-crawler/tree/v1.2) (2020-09-03)

- calling of an outer handler for an each found link:
  - handling of links immediately after they have been extracted;
  - passing of the source link in the outer handler;
- custom filtering of considered links:
  - by uniqueness of an extracted link (optional):
    - supporting of sanitizing of a link before checking of uniqueness (optional);
  - supporting of grouping of link filters:
    - result of group filtering is successful only when all filters are successful.

## [v1.1](https://github.com/thewizardplusplus/go-crawler/tree/v1.1) (2020-08-26)

- repeated extracting of relative links on error (optional):
  - only specified repeat count;
  - supporting of delay between repeats;
- parallelization possibilities:
  - crawling of relative links in parallel;
  - simulate an unbounded channel of links to avoid a deadlock.

## [v1.0](https://github.com/thewizardplusplus/go-crawler/tree/v1.0) (2020-08-19)
