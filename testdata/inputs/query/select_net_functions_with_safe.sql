-- original: https://cloud.google.com/spanner/docs/reference/standard-sql/net_functions#nethost
SELECT
  FORMAT("%T", input) AS input,
  description,
  FORMAT("%T", NET.HOST(input)) AS host,
  FORMAT("%T", NET.PUBLIC_SUFFIX(input)) AS suffix,
  FORMAT("%T", NET.REG_DOMAIN(input)) AS domain,
  FORMAT("%T", SAFE.NET.HOST(input)) AS safe_host,
  FORMAT("%T", SAFE.NET.PUBLIC_SUFFIX(input)) AS safe_suffix,
  FORMAT("%T", SAFE.NET.REG_DOMAIN(input)) AS safe_domain
FROM (
    SELECT "" AS input, "invalid input" AS description
    UNION ALL SELECT "http://abc.xyz", "standard URL"
    UNION ALL SELECT "//user:password@a.b:80/path?query",
    "standard URL with relative scheme, port, path and query, but no public suffix"
    UNION ALL SELECT "https://[::1]:80", "standard URL with IPv6 host"
    UNION ALL SELECT "http://例子.卷筒纸.中国", "standard URL with internationalized domain name"
    UNION ALL SELECT "    www.Example.Co.UK    ",
    "non-standard URL with spaces, upper case letters, and without scheme"
    UNION ALL SELECT "mailto:?to=&subject=&body=", "URI rather than URL--unsupported"
)
