# Change Log

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
