REPLACE_FIELDS(
  NEW Book(
    "The Hummingbird" AS title,
    NEW BookDetails(10 AS chapters) AS details),
  "The Hummingbird II" AS title,
  11 AS details.chapters)