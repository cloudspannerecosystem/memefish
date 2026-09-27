-- Example from https://cloud.google.com/spanner/docs/reference/standard-sql/operators#new_operator
NEW Universe {
  name: "Sol"
  closest_planets: ["Mercury", "Venus", "Earth" ]
  star {
    radius_miles: 432690
    age: 4603000000
  }
  constellations: [{
    name: "Libra"
    index: 0
  }, {
    name: "Scorpio"
    index: 1
  }]
  all_planets: (SELECT planets FROM SolTable)
}