create index VenuesByCityAndName on Venues (
  (json_value(VenueData.address.city)) asc,
  VenueName desc,
  (Rating * 2)
)
