# Change Log

## [v1.9.1](https://github.com/thewizardplusplus/go-crawler/tree/v1.9.1) (2021-06-24)

- perform the refactoring:
  - move the interfaces of the `models` package to the separate file;
  - replace the `registers.LinkGenerator` interface with `models.LinkExtractor`:
    - replace the `sitemap.GeneratorGroup` type with `extractors.ExtractorGroup`;
  - pass a thread ID to the `registers.SitemapRegister.RegisterSitemap()` method;
  - rename the `sanitizing` package to `urlutils`;
  - add the `urlutils.GenerateHierarchicalLinks()` function:
    - use in the `registers.RobotsTXTRegister` structure;
    - use in the `sitemap.HierarchicalGenerator` structure:
      - replace the `sitemap.SimpleGenerator` structure with `sitemap.HierarchicalGenerator`;
- simplify the examples.

## [v1.9](https://github.com/thewizardplusplus/go-crawler/tree/v1.9) (2021-06-18)

- crawling of all relative links for specified ones:
  - extracting links from a `sitemap.xml` file (optional):
    - supporting of a gzip compression of a `sitemap.xml` file.

## [v1.8](https://github.com/thewizardplusplus/go-crawler/tree/v1.8) (2021-06-13)

- crawling of all relative links for specified ones:
  - extracting links from a `sitemap.xml` file (optional):
    - supporting of few `sitemap.xml` files for a single link:
      - supporting of an outer generator for `sitemap.xml` links:
        - generators:
          - simple generator (it returns the `sitemap.xml` file in the site root);
          - hierarchical generator (it returns the suitable `sitemap.xml` file for each part of the URL path);
          - generator based on the `robots.txt` file;
        - supporting of grouping of generators:
          - result of group generating is merged results of each generator in the group;
          - generating concurrently:
            - processing of each generator is done in a separate goroutine.

## [v1.7.1](https://github.com/thewizardplusplus/go-crawler/tree/v1.7.1) (2021-05-28)

- crawling of all relative links for specified ones:
  - extracting links from a `sitemap.xml` file (optional):
    - ignoring of the error on loading of the `sitemap.xml` file:
      - logging of the received error;
      - returning of an empty Sitemap instead;
    - supporting of few `sitemap.xml` files for a single link:
      - processing of each `sitemap.xml` file is done in a separate goroutine;
  - supporting of grouping of link extractors:
    - extracting links concurrently:
      - processing of each link extractor is done in a separate goroutine.

## [v1.7](https://github.com/thewizardplusplus/go-crawler/tree/v1.7) (2021-05-02)

- crawling of all relative links for specified ones:
  - extracting links from a `sitemap.xml` file (optional):
    - supporting of few `sitemap.xml` files for a single link;
    - supporting of a Sitemap index file;
    - supporting of a delay before loading of a specific `sitemap.xml` file;
  - supporting of grouping of link extractors:
    - result of group extracting is merged results of each extractor in the group.

## [v1.6](https://github.com/thewizardplusplus/go-crawler/tree/v1.6) (2021-02-20)

- calling of an outer handler for an each found link:
  - handling links concurrently (optional);
- refactoring:
  - extracting the `models` package:
    - `SourcedLink` structure;
    - `LinkExtractor` interface;
    - `LinkChecker` interface;
    - `LinkHandler` interface;
  - crawling of all relative links for specified ones:
    - unioning concurrency parameters in the `crawler.ConcurrencyConfig` structure;
    - removing a leak of goroutines from the `crawler.Crawl()` function;
    - adding the `crawler.CrawlByConcurrentHandler()` function that triggers crawling using a concurrent handler.

## [v1.5.1](https://github.com/thewizardplusplus/go-crawler/tree/v1.5.1) (2021-02-10)

- crawling of all relative links for specified ones:
  - use the `httputils.HTTPClient` interface from the [github.com/thewizardplusplus/go-http-utils](https://github.com/thewizardplusplus/go-http-utils) package;
- calling of an outer handler for an each found link:
  - passing a context to the `crawler.LinkHandler` interface;
  - handling links filtered by a custom link filter (optional):
    - removing the `handlers.UniqueHandler` structure;
    - removing the `handlers.RobotsTXTHandler` structure;
- custom filtering of considered links:
  - passing a context to the `crawler.LinkChecker` interface;
  - by a `robots.txt` file (optional):
    - use the `httputils.HTTPClient` interface from the [github.com/thewizardplusplus/go-http-utils](https://github.com/thewizardplusplus/go-http-utils) package.

## [v1.5](https://github.com/thewizardplusplus/go-crawler/tree/v1.5) (2021-01-29)

- calling of an outer handler for an each found link:
  - handling only of links allowed by a `robots.txt` file (optional):
    - customized user agent;
- custom filtering of considered links:
  - by a `robots.txt` file (optional):
    - customized user agent.

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
